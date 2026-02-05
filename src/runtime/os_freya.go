// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freya

package runtime

import (
	"internal/abi"
	"internal/runtime/atomic"
	"unsafe"
)

// sigPerThreadSyscall is the same signal (SIGSETXID) used by glibc for
// per-thread syscalls. We use it for the same purpose in non-cgo binaries.
const sigPerThreadSyscall = _SIGRTMIN + 1

type mOS struct {
	// waitsema is used for parking on locks.
	waitsema uint32
}

// Freya futex support.
//
//	futexsleep(uint32 *addr, uint32 val)
//	futexwakeup(uint32 *addr)
//
// Futexsleep atomically checks if *addr == val and if so, sleeps on addr.
// Futexwakeup wakes up threads sleeping on addr.
// Futexsleep is allowed to wake up spuriously.

const (
	_FUTEX_PRIVATE_FLAG = 128
	_FUTEX_WAIT_PRIVATE = 0 | _FUTEX_PRIVATE_FLAG
	_FUTEX_WAKE_PRIVATE = 1 | _FUTEX_PRIVATE_FLAG
)

// Atomically,
//
//	if(*addr == val) sleep
//
// Might be woken up spuriously; that's allowed.
// Don't sleep longer than ns; ns < 0 means forever.
//
//go:nosplit
func futexsleep(addr *uint32, val uint32, ns int64) {
	if ns < 0 {
		futex(unsafe.Pointer(addr), _FUTEX_WAIT_PRIVATE, val, nil, nil, 0)
		return
	}

	var ts timespec
	ts.setNsec(ns)
	futex(unsafe.Pointer(addr), _FUTEX_WAIT_PRIVATE, val, unsafe.Pointer(&ts), nil, 0)
}

// If any procs are sleeping on addr, wake up at most cnt.
//
//go:nosplit
func futexwakeup(addr *uint32, cnt uint32) {
	ret := futex(unsafe.Pointer(addr), _FUTEX_WAKE_PRIVATE, cnt, nil, nil, 0)
	if ret >= 0 {
		return
	}

	systemstack(func() {
		print("futexwakeup addr=", addr, " returned ", ret, "\n")
	})

	*(*int32)(unsafe.Pointer(uintptr(0x1006))) = 0x1006
}

//go:noescape
func futex(addr unsafe.Pointer, op int32, val uint32, ts, addr2 unsafe.Pointer, val3 uint32) int32

//go:noescape
func clone(flags int32, stk, mp, gp, fn unsafe.Pointer) int32

// Thread creation flags for Freya.
const (
	_CLONE_VM     = 0x100
	_CLONE_FS     = 0x200
	_CLONE_FILES  = 0x400
	_CLONE_SIGHAND = 0x800
	_CLONE_THREAD = 0x10000
	_CLONE_SYSVSEM = 0x40000

	cloneFlags = _CLONE_VM | /* share memory */
		_CLONE_FS | /* share cwd, etc */
		_CLONE_FILES | /* share fd table */
		_CLONE_SIGHAND | /* share sig handler table */
		_CLONE_SYSVSEM | /* share SysV semaphore undo lists */
		_CLONE_THREAD /* revisit - okay for now */
)

// May run with m.p==nil, so write barriers are not allowed.
//
//go:nowritebarrier
func newosproc(mp *m) {
	stk := unsafe.Pointer(mp.g0.stack.hi)
	if false {
		print("newosproc stk=", stk, " m=", mp, " g=", mp.g0, " clone=", abi.FuncPCABI0(clone), " id=", mp.id, " ostk=", &mp, "\n")
	}

	// Disable signals during clone, so that the new thread starts
	// with signals disabled. It will enable them in minit.
	var oset sigset
	sigprocmask(_SIG_SETMASK, &sigset_all, &oset)
	ret := retryOnEAGAIN(func() int32 {
		r := clone(cloneFlags, stk, unsafe.Pointer(mp), unsafe.Pointer(mp.g0), unsafe.Pointer(abi.FuncPCABI0(mstart)))
		// clone returns positive TID, negative errno.
		// We don't care about the TID.
		if r >= 0 {
			return 0
		}
		return -r
	})
	sigprocmask(_SIG_SETMASK, &oset, nil)

	if ret != 0 {
		print("runtime: failed to create new OS thread (have ", mcount(), " already; errno=", ret, ")\n")
		if ret == _EAGAIN {
			println("runtime: may need to increase max user processes (ulimit -u)")
		}
		throw("newosproc")
	}
}

// Version of newosproc that doesn't require a valid G.
//
//go:nosplit
func newosproc0(stacksize uintptr, fn unsafe.Pointer) {
	stack := sysAlloc(stacksize, &memstats.stacks_sys, "OS thread stack")
	if stack == nil {
		writeErrStr(failallocatestack)
		exit(1)
	}
	ret := clone(cloneFlags, unsafe.Pointer(uintptr(stack)+stacksize), nil, nil, fn)
	if ret < 0 {
		writeErrStr(failthreadcreate)
		exit(1)
	}
}

const (
	_AT_NULL   = 0 // End of vector
	_AT_PAGESZ = 6 // System physical page size
	_AT_RANDOM = 25
)

func sysargs(argc int32, argv **byte) {
	// On Freya, we parse the auxiliary vector similarly to Linux.
	n := argc + 1

	// skip over argv, envp to get to auxv
	for argv_index(argv, n) != nil {
		n++
	}

	// skip NULL separator
	n++

	// now argv+n is auxv
	auxvp := (*[1 << 28]uintptr)(add(unsafe.Pointer(argv), uintptr(n)*unsafe.Sizeof(argv)))
	for i := 0; auxvp[i] != _AT_NULL; i += 2 {
		tag, val := auxvp[i], auxvp[i+1]
		switch tag {
		case _AT_PAGESZ:
			physPageSize = val
		case _AT_RANDOM:
			startupRand = (*[16]byte)(unsafe.Pointer(val))[:]
		}
	}
}

func osinit() {
	numCPUStartup = getCPUCount()
	if physPageSize == 0 {
		physPageSize = 4096
	}
}

// Called to initialize a new m (including the bootstrap m).
// Called on the parent thread (main thread in case of bootstrap), can allocate memory.
func mpreinit(mp *m) {
	mp.gsignal = malg(32 * 1024) // Freya wants >= 2K
	mp.gsignal.m = mp
}

func gettid() uint32

// Called to initialize a new m (including the bootstrap m).
// Called on the new thread, cannot allocate memory.
func minit() {
	minitSignals()

	// Cgo-created threads and the bootstrap m are missing a
	// procid. We need this for asynchronous preemption and it's
	// useful in debuggers.
	getg().m.procid = uint64(gettid())
}

// Called from dropm to undo the effect of an minit.
//
//go:nosplit
func unminit() {
	unminitSignals()
	getg().m.procid = 0
}

// Called from mexit, but not from dropm, to undo the effect of thread-owned
// resources in minit, semacreate, or elsewhere. Do not take locks after calling this.
//
// This always runs without a P, so //go:nowritebarrierrec is required.
//
//go:nowritebarrierrec
func mdestroy(mp *m) {
}

func sigreturn__sigaction()
func sigtramp() // Called via C ABI
func cgoSigtramp()

//go:noescape
func sigaltstack(new, old *stackt)

//go:noescape
func setitimer(mode int32, new, old *itimerval)

//go:noescape
func rtsigprocmask(how int32, new, old *sigset, size int32)

//go:nosplit
//go:nowritebarrierrec
func sigprocmask(how int32, new, old *sigset) {
	rtsigprocmask(how, new, old, int32(unsafe.Sizeof(*new)))
}

func raise(sig uint32)
func raiseproc(sig uint32)

//go:noescape
func sched_getaffinity(pid, len uintptr, buf *byte) int32
func osyield()

//go:nosplit
func osyield_no_g() {
	osyield()
}

func pipe2(flags int32) (r, w int32, errno int32)

const (
	_si_max_size    = 128
	_sigev_max_size = 64
)

//go:nosplit
//go:nowritebarrierrec
func setsig(i uint32, fn uintptr) {
	var sa sigactiont
	sa.sa_flags = _SA_SIGINFO | _SA_ONSTACK | _SA_RESTORER | _SA_RESTART
	sigfillset(&sa.sa_mask)
	if GOARCH == "386" || GOARCH == "amd64" {
		sa.sa_restorer = abi.FuncPCABI0(sigreturn__sigaction)
	}
	if fn == abi.FuncPCABIInternal(sighandler) {
		if iscgo {
			fn = abi.FuncPCABI0(cgoSigtramp)
		} else {
			fn = abi.FuncPCABI0(sigtramp)
		}
	}
	sa.sa_handler = fn
	sigaction(i, &sa, nil)
}

//go:nosplit
//go:nowritebarrierrec
func setsigstack(i uint32) {
	var sa sigactiont
	sigaction(i, nil, &sa)
	if sa.sa_flags&_SA_ONSTACK != 0 {
		return
	}
	sa.sa_flags |= _SA_ONSTACK
	sigaction(i, &sa, nil)
}

//go:nosplit
//go:nowritebarrierrec
func getsig(i uint32) uintptr {
	var sa sigactiont
	sigaction(i, nil, &sa)
	return sa.sa_handler
}

// setSignalstackSP sets the ss_sp field of a stackt.
//
//go:nosplit
func setSignalstackSP(s *stackt, sp uintptr) {
	*(*uintptr)(unsafe.Pointer(&s.ss_sp)) = sp
}

//go:nosplit
func (c *sigctxt) fixsigcode(sig uint32) {
}

// sysSigaction calls the rt_sigaction system call.
//
//go:nosplit
func sysSigaction(sig uint32, new, old *sigactiont) {
	if rt_sigaction(uintptr(sig), new, old, unsafe.Sizeof(sigactiont{}.sa_mask)) != 0 {
		// Ignore errors for signals 32, 33, 64 (like Linux does for QEMU).
		if sig != 32 && sig != 33 && sig != 64 {
			systemstack(func() {
				throw("sigaction failed")
			})
		}
	}
}

// rt_sigaction is implemented in assembly.
//
//go:noescape
func rt_sigaction(sig uintptr, new, old *sigactiont, size uintptr) int32

func getpid() int
func tgkill(tgid, tid, sig int)

// signalM sends a signal to mp.
func signalM(mp *m, sig int) {
	tgkill(getpid(), int(mp.procid), sig)
}

// validSIGPROF compares this signal delivery's code against the signal sources
// that the profiler uses, returning whether the delivery should be processed.
//
//go:nosplit
func validSIGPROF(mp *m, c *sigctxt) bool {
	code := int32(c.sigcode())
	setitimer_sig := code == _SI_KERNEL
	if !setitimer_sig {
		return true
	}
	return setitimer_sig
}

func setProcessCPUProfiler(hz int32) {
	setProcessCPUProfilerTimer(hz)
}

func setThreadCPUProfiler(hz int32) {
	mp := getg().m
	mp.profilehz = hz
}

const (
	_SI_USER  = 0
	_SI_TKILL = -6
)

func readRandom(r []byte) int {
	// Freya: try reading from /dev/urandom
	var urandom_dev = []byte("/dev/urandom\x00")
	fd := open(&urandom_dev[0], 0 /* O_RDONLY */, 0)
	if fd < 0 {
		return 0
	}
	n := read(fd, unsafe.Pointer(&r[0]), int32(len(r)))
	closefd(fd)
	return int(n)
}

func goenvs() {
	goenvs_unix()
}

// Called to do synchronous initialization of Go code built with
// -buildmode=c-archive or -buildmode=c-shared.
// None of the Go runtime is initialized.
//
//go:nosplit
//go:nowritebarrierrec
func libpreinit() {
	initsig(true)
}

//go:nosplit
func fixSigactionForCgo(new *sigactiont) {
	if GOARCH == "386" && new != nil {
		new.sa_flags &^= _SA_RESTORER
		new.sa_restorer = 0
	}
}

//go:nosplit
func mprotect(addr unsafe.Pointer, n uintptr, prot int32) (ret int32, errno int32) {
	// Freya: stub for now
	return 0, 0
}

// getCPUCount returns the number of logical CPUs usable by the current process.
// On Freya, we don't have sched_getaffinity, so we return 1.
func getCPUCount() int32 {
	return 1
}

// fcntl is implemented in sys_freya_*.s
//
//go:noescape
func fcntl(fd, cmd, arg int32) (ret int32, errno int32)

//go:nosplit
func runPerThreadSyscall() {
	throw("runPerThreadSyscall only valid on linux")
}

// dummy declarations to satisfy the linker for unused references
var _ atomic.Uint8

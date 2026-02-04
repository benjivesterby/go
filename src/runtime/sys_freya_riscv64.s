// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// System calls and other sys.stuff for riscv64, Freya
//

#include "textflag.h"
#include "go_asm.h"

#define CLOCK_REALTIME 0
#define CLOCK_MONOTONIC 1

// Freya syscall numbers
#define SYS_exit		10
#define SYS_getpid		13
#define SYS_thread_yield	15
#define SYS_thread_exit		16
#define SYS_thread_create	17
#define SYS_thread_sleep	18
#define SYS_get_time		19
#define SYS_mmap		20
#define SYS_munmap		21
#define SYS_mprotect		22
#define SYS_process_kill	27
#define SYS_nanosleep		36
#define SYS_clock_gettime	37
#define SYS_futex		44
#define SYS_dup			54
#define SYS_pipe_create		55
#define SYS_dup2		86
#define SYS_rt_sigaction	94
#define SYS_rt_sigprocmask	95
#define SYS_rt_sigreturn	97
#define SYS_open		100
#define SYS_read		101
#define SYS_write		102
#define SYS_close		103
#define SYS_seek		105
#define SYS_stat		106

// func exit(code int32)
TEXT runtime·exit(SB),NOSPLIT|NOFRAME,$0-4
	MOVW	code+0(FP), A0
	MOV	$SYS_exit, A7
	ECALL
	RET

// func exitThread(wait *atomic.Uint32)
TEXT runtime·exitThread(SB),NOSPLIT|NOFRAME,$0-8
	MOV	wait+0(FP), A0
	// We're done using the stack.
	FENCE
	MOVW	ZERO, (A0)
	FENCE
	MOV	$0, A0	// exit code
	MOV	$SYS_thread_exit, A7
	ECALL
	JMP	0(PC)

// func open(name *byte, mode, perm int32) int32
TEXT runtime·open(SB),NOSPLIT|NOFRAME,$0-20
	MOV	name+0(FP), A0
	MOVW	mode+8(FP), A1
	MOVW	perm+12(FP), A2
	MOV	$SYS_open, A7
	ECALL
	MOV	$-4096, T0
	BGEU	T0, A0, 2(PC)
	MOV	$-1, A0
	MOVW	A0, ret+16(FP)
	RET

// func closefd(fd int32) int32
TEXT runtime·closefd(SB),NOSPLIT|NOFRAME,$0-12
	MOVW	fd+0(FP), A0
	MOV	$SYS_close, A7
	ECALL
	MOV	$-4096, T0
	BGEU	T0, A0, 2(PC)
	MOV	$-1, A0
	MOVW	A0, ret+8(FP)
	RET

// func write1(fd uintptr, p unsafe.Pointer, n int32) int32
TEXT runtime·write1(SB),NOSPLIT|NOFRAME,$0-28
	MOV	fd+0(FP), A0
	MOV	p+8(FP), A1
	MOVW	n+16(FP), A2
	MOV	$SYS_write, A7
	ECALL
	MOVW	A0, ret+24(FP)
	RET

// func read(fd int32, p unsafe.Pointer, n int32) int32
TEXT runtime·read(SB),NOSPLIT|NOFRAME,$0-28
	MOVW	fd+0(FP), A0
	MOV	p+8(FP), A1
	MOVW	n+16(FP), A2
	MOV	$SYS_read, A7
	ECALL
	MOVW	A0, ret+24(FP)
	RET

// func pipe2(flags int32) (r, w int32, errno int32)
TEXT runtime·pipe2(SB),NOSPLIT|NOFRAME,$0-20
	MOV	$r+8(FP), A0
	MOVW	flags+0(FP), A1
	MOV	$SYS_pipe_create, A7
	ECALL
	MOVW	A0, errno+16(FP)
	RET

// func usleep(usec uint32)
TEXT runtime·usleep(SB),NOSPLIT,$24-4
	MOVWU	usec+0(FP), A0
	MOV	$1000, A1
	MUL	A1, A0, A0
	MOV	$1000000000, A1
	DIV	A1, A0, A2
	MOV	A2, 8(X2)
	REM	A1, A0, A3
	MOV	A3, 16(X2)
	ADD	$8, X2, A0
	MOV	ZERO, A1
	MOV	$SYS_nanosleep, A7
	ECALL
	RET

// func gettid() uint32
TEXT runtime·gettid(SB),NOSPLIT,$0-4
	MOV	$SYS_getpid, A7
	ECALL
	MOVW	A0, ret+0(FP)
	RET

// func raise(sig uint32)
TEXT runtime·raise(SB),NOSPLIT|NOFRAME,$0
	MOV	$SYS_getpid, A7
	ECALL
	// arg 1 pid - already in A0
	MOVW	sig+0(FP), A1	// arg 2
	MOV	$SYS_process_kill, A7
	ECALL
	RET

// func raiseproc(sig uint32)
TEXT runtime·raiseproc(SB),NOSPLIT|NOFRAME,$0
	MOV	$SYS_getpid, A7
	ECALL
	// arg 1 pid - already in A0
	MOVW	sig+0(FP), A1	// arg 2
	MOV	$SYS_process_kill, A7
	ECALL
	RET

// func getpid() int
TEXT ·getpid(SB),NOSPLIT,$0-8
	MOV	$SYS_getpid, A7
	ECALL
	MOV	A0, ret+0(FP)
	RET

// func tgkill(tgid, tid, sig int)
TEXT ·tgkill(SB),NOSPLIT,$0-24
	MOV	tgid+0(FP), A0
	MOV	tid+8(FP), A1
	MOV	sig+16(FP), A2
	MOV	$SYS_process_kill, A7
	ECALL
	RET

// func walltime() (sec int64, nsec int32)
TEXT runtime·walltime(SB),NOSPLIT,$24-12
	MOV	$CLOCK_REALTIME, A0
	ADD	$8, X2, A1
	MOV	$SYS_clock_gettime, A7
	ECALL
	MOV	8(X2), T0	// sec
	MOV	16(X2), T1	// nsec
	MOV	T0, sec+0(FP)
	MOVW	T1, nsec+8(FP)
	RET

// func nanotime1() int64
TEXT runtime·nanotime1(SB),NOSPLIT,$24-8
	MOV	$CLOCK_MONOTONIC, A0
	ADD	$8, X2, A1
	MOV	$SYS_clock_gettime, A7
	ECALL
	MOV	8(X2), T0	// sec
	MOV	16(X2), T1	// nsec
	MOV	$1000000000, T2
	MUL	T2, T0
	ADD	T1, T0
	MOV	T0, ret+0(FP)
	RET

// func rtsigprocmask(how int32, new, old *sigset, size int32)
TEXT runtime·rtsigprocmask(SB),NOSPLIT|NOFRAME,$0-28
	MOVW	how+0(FP), A0
	MOV	new+8(FP), A1
	MOV	old+16(FP), A2
	MOVW	size+24(FP), A3
	MOV	$SYS_rt_sigprocmask, A7
	ECALL
	MOV	$-4096, T0
	BLTU	A0, T0, 2(PC)
	WORD	$0	// crash
	RET

// func rt_sigaction(sig uintptr, new, old *sigactiont, size uintptr) int32
TEXT runtime·rt_sigaction(SB),NOSPLIT|NOFRAME,$0-36
	MOV	sig+0(FP), A0
	MOV	new+8(FP), A1
	MOV	old+16(FP), A2
	MOV	size+24(FP), A3
	MOV	$SYS_rt_sigaction, A7
	ECALL
	MOVW	A0, ret+32(FP)
	RET

// func sigfwd(fn uintptr, sig uint32, info *siginfo, ctx unsafe.Pointer)
TEXT runtime·sigfwd(SB),NOSPLIT,$0-32
	MOVW	sig+8(FP), A0
	MOV	info+16(FP), A1
	MOV	ctx+24(FP), A2
	MOV	fn+0(FP), T1
	JALR	RA, T1
	RET

// func sigtramp(signo, ureg, ctxt unsafe.Pointer)
TEXT runtime·sigtramp(SB),NOSPLIT|TOPFRAME,$64
	MOVW	A0, 8(X2)
	MOV	A1, 16(X2)
	MOV	A2, 24(X2)

	// this might be called in external code context,
	// where g is not set.
	MOVBU	runtime·iscgo(SB), A0
	BEQ	A0, ZERO, 2(PC)
	CALL	runtime·load_g(SB)

	MOV	$runtime·sigtrampgo(SB), A0
	JALR	RA, A0
	RET

// func cgoSigtramp()
TEXT runtime·cgoSigtramp(SB),NOSPLIT,$0
	MOV	$runtime·sigtramp(SB), T1
	JALR	ZERO, T1

// func mmap(addr unsafe.Pointer, n uintptr, prot, flags, fd int32, off uint32) (p unsafe.Pointer, err int)
TEXT runtime·mmap(SB),NOSPLIT|NOFRAME,$0
	MOV	addr+0(FP), A0
	MOV	n+8(FP), A1
	MOVW	prot+16(FP), A2
	MOVW	flags+20(FP), A3
	MOVW	fd+24(FP), A4
	MOVW	off+28(FP), A5
	MOV	$SYS_mmap, A7
	ECALL
	MOV	$-4096, T0
	BGEU	T0, A0, 5(PC)
	SUB	A0, ZERO, A0
	MOV	ZERO, p+32(FP)
	MOV	A0, err+40(FP)
	RET
ok:
	MOV	A0, p+32(FP)
	MOV	ZERO, err+40(FP)
	RET

// func munmap(addr unsafe.Pointer, n uintptr)
TEXT runtime·munmap(SB),NOSPLIT|NOFRAME,$0
	MOV	addr+0(FP), A0
	MOV	n+8(FP), A1
	MOV	$SYS_munmap, A7
	ECALL
	MOV	$-4096, T0
	BLTU	A0, T0, 2(PC)
	WORD	$0	// crash
	RET

// func futex(addr unsafe.Pointer, op int32, val uint32, ts, addr2 unsafe.Pointer, val3 uint32) int32
TEXT runtime·futex(SB),NOSPLIT|NOFRAME,$0-48
	MOV	addr+0(FP), A0
	MOVW	op+8(FP), A1
	MOVW	val+12(FP), A2
	MOV	ts+16(FP), A3
	MOV	addr2+24(FP), A4
	MOVW	val3+32(FP), A5
	MOV	$SYS_futex, A7
	ECALL
	MOVW	A0, ret+40(FP)
	RET

// func clone(flags int32, stk, mp, gp, fn unsafe.Pointer) int32
TEXT runtime·clone(SB),NOSPLIT|NOFRAME,$0
	MOVW	flags+0(FP), A0
	MOV	stk+8(FP), A1

	// Copy mp, gp, fn off parent stack for use by child.
	MOV	mp+16(FP), T0
	MOV	gp+24(FP), T1
	MOV	fn+32(FP), T2

	MOV	T0, -8(A1)
	MOV	T1, -16(A1)
	MOV	T2, -24(A1)
	MOV	$1234, T0
	MOV	T0, -32(A1)

	MOV	$SYS_thread_create, A7
	ECALL

	// In parent, return.
	BEQ	ZERO, A0, child
	MOVW	ZERO, ret+40(FP)
	RET

child:
	// In child, on new stack.
	MOV	-32(X2), T0
	MOV	$1234, A0
	BEQ	A0, T0, good
	WORD	$0	// crash

good:
	// Initialize m->procid to tid
	MOV	$SYS_getpid, A7
	ECALL

	MOV	-24(X2), T2	// fn
	MOV	-16(X2), T1	// g
	MOV	-8(X2), T0	// m

	BEQ	ZERO, T0, nog
	BEQ	ZERO, T1, nog

	MOV	A0, m_procid(T0)

	// In child, set up new stack
	MOV	T0, g_m(T1)
	MOV	T1, g

nog:
	// Call fn
	JALR	RA, T2

	// It shouldn't return. If it does, exit this thread.
	MOV	$111, A0
	MOV	$SYS_thread_exit, A7
	ECALL
	JMP	-3(PC)	// keep exiting

// func sigaltstack(new, old *stackt)
TEXT runtime·sigaltstack(SB),NOSPLIT|NOFRAME,$0
	// Freya: sigaltstack is a no-op for now since the kernel
	// handles signal stacks internally. Keep the interface for compatibility.
	RET

// func osyield()
TEXT runtime·osyield(SB),NOSPLIT,$0
	MOV	$SYS_thread_yield, A7
	ECALL
	RET

// func sched_getaffinity(pid, len uintptr, buf *uintptr) int32
TEXT runtime·sched_getaffinity(SB),NOSPLIT,$0-28
	// Freya doesn't have sched_getaffinity; return 0 (will use default ncpu=1)
	MOVW	ZERO, ret+24(FP)
	RET

// func madvise(addr unsafe.Pointer, n uintptr, flags int32) int32
TEXT runtime·madvise(SB),NOSPLIT,$0-28
	// Freya doesn't have madvise; return 0
	MOVW	ZERO, ret+24(FP)
	RET

// sbrk0 is provided by stubs_nonlinux.go for non-linux targets.

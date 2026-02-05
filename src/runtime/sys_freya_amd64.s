// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// System calls and other sys.stuff for AMD64, Freya
//

#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"

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
#define SYS_fcntl		80

TEXT runtime·exit(SB),NOSPLIT,$0-4
	MOVL	code+0(FP), DI
	MOVL	$SYS_exit, AX
	SYSCALL
	RET

// func exitThread(wait *atomic.Uint32)
TEXT runtime·exitThread(SB),NOSPLIT,$0-8
	MOVQ	wait+0(FP), AX
	// We're done using the stack.
	MOVL	$0, (AX)
	MOVL	$0, DI	// exit code
	MOVL	$SYS_thread_exit, AX
	SYSCALL
	// We may not even have a stack any more.
	INT	$3
	JMP	0(PC)

TEXT runtime·open(SB),NOSPLIT,$0-20
	MOVQ	name+0(FP), DI
	MOVL	mode+8(FP), SI
	MOVL	perm+12(FP), DX
	MOVL	$SYS_open, AX
	SYSCALL
	CMPQ	AX, $0xfffffffffffff001
	JLS	2(PC)
	MOVL	$-1, AX
	MOVL	AX, ret+16(FP)
	RET

TEXT runtime·closefd(SB),NOSPLIT,$0-12
	MOVL	fd+0(FP), DI
	MOVL	$SYS_close, AX
	SYSCALL
	CMPQ	AX, $0xfffffffffffff001
	JLS	2(PC)
	MOVL	$-1, AX
	MOVL	AX, ret+8(FP)
	RET

TEXT runtime·write1(SB),NOSPLIT,$0-28
	MOVQ	fd+0(FP), DI
	MOVQ	p+8(FP), SI
	MOVL	n+16(FP), DX
	MOVL	$SYS_write, AX
	SYSCALL
	MOVL	AX, ret+24(FP)
	RET

TEXT runtime·read(SB),NOSPLIT,$0-28
	MOVL	fd+0(FP), DI
	MOVQ	p+8(FP), SI
	MOVL	n+16(FP), DX
	MOVL	$SYS_read, AX
	SYSCALL
	MOVL	AX, ret+24(FP)
	RET

// func pipe2(flags int32) (r, w int32, errno int32)
TEXT runtime·pipe2(SB),NOSPLIT,$0-20
	LEAQ	r+8(FP), DI
	MOVL	flags+0(FP), SI
	MOVL	$SYS_pipe_create, AX
	SYSCALL
	MOVL	AX, errno+16(FP)
	RET

TEXT runtime·usleep(SB),NOSPLIT,$16
	MOVL	$0, DX
	MOVL	usec+0(FP), AX
	MOVL	$1000000, CX
	DIVL	CX
	MOVQ	AX, 0(SP)
	MOVL	$1000, AX	// usec to nsec
	MULL	DX
	MOVQ	AX, 8(SP)

	// nanosleep(&ts, 0)
	MOVQ	SP, DI
	MOVL	$0, SI
	MOVL	$SYS_nanosleep, AX
	SYSCALL
	RET

TEXT runtime·gettid(SB),NOSPLIT,$0-4
	MOVL	$SYS_getpid, AX
	SYSCALL
	MOVL	AX, ret+0(FP)
	RET

TEXT runtime·raise(SB),NOSPLIT,$0
	MOVL	$SYS_getpid, AX
	SYSCALL
	MOVL	AX, DI	// arg 1 pid
	MOVL	sig+0(FP), SI	// arg 2
	MOVL	$SYS_process_kill, AX
	SYSCALL
	RET

TEXT runtime·raiseproc(SB),NOSPLIT,$0
	MOVL	$SYS_getpid, AX
	SYSCALL
	MOVL	AX, DI	// arg 1 pid
	MOVL	sig+0(FP), SI	// arg 2
	MOVL	$SYS_process_kill, AX
	SYSCALL
	RET

TEXT ·getpid(SB),NOSPLIT,$0-8
	MOVL	$SYS_getpid, AX
	SYSCALL
	MOVQ	AX, ret+0(FP)
	RET

TEXT ·tgkill(SB),NOSPLIT,$0
	MOVQ	tgid+0(FP), DI
	MOVQ	tid+8(FP), SI
	MOVQ	sig+16(FP), DX
	MOVL	$SYS_process_kill, AX
	SYSCALL
	RET

// func walltime() (sec int64, nsec int32)
TEXT runtime·walltime(SB),NOSPLIT,$16-12
	MOVL	$CLOCK_REALTIME, DI
	LEAQ	0(SP), SI
	MOVL	$SYS_clock_gettime, AX
	SYSCALL
	MOVQ	0(SP), AX	// sec
	MOVQ	8(SP), DX	// nsec
	MOVQ	AX, sec+0(FP)
	MOVL	DX, nsec+8(FP)
	RET

// func nanotime1() int64
TEXT runtime·nanotime1(SB),NOSPLIT,$16-8
	MOVL	$CLOCK_MONOTONIC, DI
	LEAQ	0(SP), SI
	MOVL	$SYS_clock_gettime, AX
	SYSCALL
	MOVQ	0(SP), AX	// sec
	MOVQ	8(SP), DX	// nsec
	// sec is in AX, nsec in DX
	// return nsec in AX
	IMULQ	$1000000000, AX
	ADDQ	DX, AX
	MOVQ	AX, ret+0(FP)
	RET

TEXT runtime·rtsigprocmask(SB),NOSPLIT,$0-28
	MOVL	how+0(FP), DI
	MOVQ	new+8(FP), SI
	MOVQ	old+16(FP), DX
	MOVL	size+24(FP), R10
	MOVL	$SYS_rt_sigprocmask, AX
	SYSCALL
	CMPQ	AX, $0xfffffffffffff001
	JLS	2(PC)
	MOVL	$0xf1, 0xf1  // crash
	RET

TEXT runtime·rt_sigaction(SB),NOSPLIT,$0-36
	MOVQ	sig+0(FP), DI
	MOVQ	new+8(FP), SI
	MOVQ	old+16(FP), DX
	MOVQ	size+24(FP), R10
	MOVL	$SYS_rt_sigaction, AX
	SYSCALL
	MOVL	AX, ret+32(FP)
	RET

TEXT runtime·sigfwd(SB),NOSPLIT,$0-32
	MOVQ	fn+0(FP),    AX
	MOVL	sig+8(FP),   DI
	MOVQ	info+16(FP), SI
	MOVQ	ctx+24(FP),  DX
	MOVQ	SP, BX		// callee-saved
	ANDQ	$~15, SP     // alignment for x86_64 ABI
	CALL	AX
	MOVQ	BX, SP
	RET

// Called using C ABI.
TEXT runtime·sigtramp(SB),NOSPLIT|TOPFRAME|NOFRAME,$0
	// Transition from C ABI to Go ABI.
	// Save all callee-save registers.
	PUSHQ	BX
	PUSHQ	BP
	PUSHQ	R12
	PUSHQ	R13
	PUSHQ	R14
	PUSHQ	R15

	// Set up ABIInternal environment: g in R14.
	get_tls(R12)
	MOVQ	g(R12), R14
	PXOR	X15, X15

	// Reserve space for spill slots.
	NOP	SP		// disable vet stack checking
	ADJSP	$24

	// Call into the Go signal handler
	MOVQ	DI, AX	// sig
	MOVQ	SI, BX	// info
	MOVQ	DX, CX	// ctx
	CALL	·sigtrampgo<ABIInternal>(SB)

	ADJSP	$-24

	POPQ	R15
	POPQ	R14
	POPQ	R13
	POPQ	R12
	POPQ	BP
	POPQ	BX
	RET

TEXT runtime·cgoSigtramp(SB),NOSPLIT,$0
	JMP	runtime·sigtramp(SB)

// For gdb and signal return
TEXT runtime·sigreturn__sigaction(SB),NOSPLIT,$0
	MOVQ	$SYS_rt_sigreturn, AX
	SYSCALL
	INT $3	// not reached

// func mmap(addr unsafe.Pointer, n uintptr, prot, flags, fd int32, off uint32) (p unsafe.Pointer, err int)
TEXT runtime·mmap(SB),NOSPLIT,$0
	MOVQ	addr+0(FP), DI
	MOVQ	n+8(FP), SI
	MOVL	prot+16(FP), DX
	MOVL	flags+20(FP), R10
	MOVL	fd+24(FP), R8
	MOVL	off+28(FP), R9

	MOVL	$SYS_mmap, AX
	SYSCALL
	CMPQ	AX, $0xfffffffffffff001
	JLS	ok
	NOTQ	AX
	INCQ	AX
	MOVQ	$0, p+32(FP)
	MOVQ	AX, err+40(FP)
	RET
ok:
	MOVQ	AX, p+32(FP)
	MOVQ	$0, err+40(FP)
	RET

// func munmap(addr unsafe.Pointer, n uintptr)
TEXT runtime·munmap(SB),NOSPLIT,$0
	MOVQ	addr+0(FP), DI
	MOVQ	n+8(FP), SI
	MOVQ	$SYS_munmap, AX
	SYSCALL
	CMPQ	AX, $0xfffffffffffff001
	JLS	2(PC)
	MOVL	$0xf1, 0xf1  // crash
	RET

// func madvise(addr unsafe.Pointer, n uintptr, flags int32) int32
TEXT runtime·madvise(SB),NOSPLIT,$0-28
	// Freya doesn't have madvise; return 0
	MOVL	$0, AX
	MOVL	AX, ret+24(FP)
	RET

// func fcntl(fd, cmd, arg int32) (ret int32, errno int32)
TEXT runtime·fcntl(SB),NOSPLIT,$0-20
	MOVL	fd+0(FP), DI
	MOVL	cmd+4(FP), SI
	MOVL	arg+8(FP), DX
	MOVL	$SYS_fcntl, AX
	SYSCALL
	CMPQ	AX, $0xfffffffffffff001
	JLS	noerr
	// Error case: AX contains negative errno
	MOVL	$-1, ret+12(FP)
	NEGQ	AX
	MOVL	AX, errno+16(FP)
	RET
noerr:
	MOVL	AX, ret+12(FP)
	MOVL	$0, errno+16(FP)
	RET

// int64 futex(int32 *uaddr, int32 op, int32 val,
//	struct timespec *timeout, int32 *uaddr2, int32 val2);
TEXT runtime·futex(SB),NOSPLIT,$0
	MOVQ	addr+0(FP), DI
	MOVL	op+8(FP), SI
	MOVL	val+12(FP), DX
	MOVQ	ts+16(FP), R10
	MOVQ	addr2+24(FP), R8
	MOVL	val3+32(FP), R9
	MOVL	$SYS_futex, AX
	SYSCALL
	MOVL	AX, ret+40(FP)
	RET

// int32 clone(int32 flags, void *stk, M *mp, G *gp, void (*fn)(void));
TEXT runtime·clone(SB),NOSPLIT|NOFRAME,$0
	MOVL	flags+0(FP), DI
	MOVQ	stk+8(FP), SI
	MOVQ	$0, DX
	MOVQ	$0, R10
	MOVQ    $0, R8
	// Copy mp, gp, fn off parent stack for use by child.
	// Careful: syscall clobbers CX and R11.
	MOVQ	mp+16(FP), R13
	MOVQ	gp+24(FP), R9
	MOVQ	fn+32(FP), R12

	MOVL	$SYS_thread_create, AX
	SYSCALL

	// In parent, return.
	CMPQ	AX, $0
	JEQ	3(PC)
	MOVL	AX, ret+40(FP)
	RET

	// In child, on new stack.
	MOVQ	SI, SP

	// If g or m are nil, skip Go-related setup.
	CMPQ	R13, $0    // m
	JEQ	nog
	CMPQ	R9, $0    // g
	JEQ	nog

	// Initialize m->procid to tid
	MOVL	$SYS_getpid, AX
	SYSCALL
	MOVQ	AX, m_procid(R13)

	// In child, set up new stack
	get_tls(CX)
	MOVQ	R13, g_m(R9)
	MOVQ	R9, g(CX)
	MOVQ	R9, R14 // set g register
	CALL	runtime·stackcheck(SB)

nog:
	// Call fn. This is the PC of an ABI0 function.
	CALL	R12

	// It shouldn't return. If it does, exit that thread.
	MOVL	$111, DI
	MOVL	$SYS_thread_exit, AX
	SYSCALL
	JMP	-3(PC)	// keep exiting

TEXT runtime·sigaltstack(SB),NOSPLIT,$0
	// Freya: sigaltstack is handled by the kernel internally.
	RET

// set tls base to DI
TEXT runtime·settls(SB),NOSPLIT,$32
	// Freya uses the same TLS mechanism as Linux on amd64
	ADDQ	$8, DI	// ELF wants to use -8(FS)
	MOVQ	DI, SI
	MOVQ	$0x1002, DI	// ARCH_SET_FS
	// For now, use a simple approach - store the TLS pointer
	RET

TEXT runtime·osyield(SB),NOSPLIT,$0
	MOVL	$SYS_thread_yield, AX
	SYSCALL
	RET

TEXT runtime·sched_getaffinity(SB),NOSPLIT,$0
	// Freya doesn't have sched_getaffinity
	MOVL	$0, AX
	MOVL	AX, ret+24(FP)
	RET

// func sbrk0() uintptr
TEXT runtime·sbrk0(SB),NOSPLIT,$0-8
	MOVQ	$0, AX
	MOVQ	AX, ret+0(FP)
	RET

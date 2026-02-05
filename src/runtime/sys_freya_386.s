// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// System calls and other sys.stuff for 386, Freya
//

#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"

#define INVOKE_SYSCALL	INT	$0x80

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

TEXT runtime·exit(SB),NOSPLIT,$0
	MOVL	$SYS_exit, AX
	MOVL	code+0(FP), BX
	INVOKE_SYSCALL
	INT $3	// not reached
	RET

TEXT exit1<>(SB),NOSPLIT,$0
	MOVL	$SYS_thread_exit, AX
	MOVL	code+0(FP), BX
	INVOKE_SYSCALL
	INT $3	// not reached
	RET

// func exitThread(wait *atomic.Uint32)
TEXT runtime·exitThread(SB),NOSPLIT,$0-4
	MOVL	wait+0(FP), AX
	// We're done using the stack.
	MOVL	$0, (AX)
	MOVL	$SYS_thread_exit, AX
	MOVL	$0, BX	// exit code
	INT	$0x80	// no stack; must not use CALL
	// We may not even have a stack any more.
	INT	$3
	JMP	0(PC)

TEXT runtime·open(SB),NOSPLIT,$0
	MOVL	$SYS_open, AX
	MOVL	name+0(FP), BX
	MOVL	mode+4(FP), CX
	MOVL	perm+8(FP), DX
	INVOKE_SYSCALL
	CMPL	AX, $0xfffff001
	JLS	2(PC)
	MOVL	$-1, AX
	MOVL	AX, ret+12(FP)
	RET

TEXT runtime·closefd(SB),NOSPLIT,$0
	MOVL	$SYS_close, AX
	MOVL	fd+0(FP), BX
	INVOKE_SYSCALL
	CMPL	AX, $0xfffff001
	JLS	2(PC)
	MOVL	$-1, AX
	MOVL	AX, ret+4(FP)
	RET

TEXT runtime·write1(SB),NOSPLIT,$0
	MOVL	$SYS_write, AX
	MOVL	fd+0(FP), BX
	MOVL	p+4(FP), CX
	MOVL	n+8(FP), DX
	INVOKE_SYSCALL
	MOVL	AX, ret+12(FP)
	RET

TEXT runtime·read(SB),NOSPLIT,$0
	MOVL	$SYS_read, AX
	MOVL	fd+0(FP), BX
	MOVL	p+4(FP), CX
	MOVL	n+8(FP), DX
	INVOKE_SYSCALL
	MOVL	AX, ret+12(FP)
	RET

// func pipe2(flags int32) (r, w int32, errno int32)
TEXT runtime·pipe2(SB),NOSPLIT,$0-16
	MOVL	$SYS_pipe_create, AX
	LEAL	r+4(FP), BX
	MOVL	flags+0(FP), CX
	INVOKE_SYSCALL
	MOVL	AX, errno+12(FP)
	RET

TEXT runtime·usleep(SB),NOSPLIT,$8
	MOVL	$0, DX
	MOVL	usec+0(FP), AX
	MOVL	$1000000, CX
	DIVL	CX
	MOVL	AX, 0(SP)
	MOVL	$1000, AX	// usec to nsec
	MULL	DX
	MOVL	AX, 4(SP)

	// nanosleep(&ts, 0)
	MOVL	$SYS_nanosleep, AX
	LEAL	0(SP), BX
	MOVL	$0, CX
	INVOKE_SYSCALL
	RET

TEXT runtime·gettid(SB),NOSPLIT,$0-4
	MOVL	$SYS_getpid, AX
	INVOKE_SYSCALL
	MOVL	AX, ret+0(FP)
	RET

TEXT runtime·raise(SB),NOSPLIT,$12
	MOVL	$SYS_getpid, AX
	INVOKE_SYSCALL
	MOVL	AX, BX	// arg 1 pid
	MOVL	sig+0(FP), CX	// arg 2 signal
	MOVL	$SYS_process_kill, AX
	INVOKE_SYSCALL
	RET

TEXT runtime·raiseproc(SB),NOSPLIT,$12
	MOVL	$SYS_getpid, AX
	INVOKE_SYSCALL
	MOVL	AX, BX	// arg 1 pid
	MOVL	sig+0(FP), CX	// arg 2 signal
	MOVL	$SYS_process_kill, AX
	INVOKE_SYSCALL
	RET

TEXT ·getpid(SB),NOSPLIT,$0-4
	MOVL	$SYS_getpid, AX
	INVOKE_SYSCALL
	MOVL	AX, ret+0(FP)
	RET

TEXT ·tgkill(SB),NOSPLIT,$0
	MOVL	$SYS_process_kill, AX
	MOVL	tgid+0(FP), BX
	MOVL	tid+4(FP), CX
	MOVL	sig+8(FP), DX
	INVOKE_SYSCALL
	RET

TEXT runtime·setitimer(SB),NOSPLIT,$0-12
	// Freya: no-op for now
	RET

TEXT runtime·mincore(SB),NOSPLIT,$0-16
	// Freya doesn't have mincore
	MOVL	$-1, AX
	MOVL	AX, ret+12(FP)
	RET

// func walltime() (sec int64, nsec int32)
TEXT runtime·walltime(SB), NOSPLIT, $8-12
	MOVL	$SYS_clock_gettime, AX
	MOVL	$0, BX		// CLOCK_REALTIME
	LEAL	0(SP), CX
	INVOKE_SYSCALL

	MOVL	0(SP), AX	// sec
	MOVL	4(SP), BX	// nsec

	// sec is in AX, nsec in BX
	MOVL	AX, sec_lo+0(FP)
	MOVL	$0, sec_hi+4(FP)
	MOVL	BX, nsec+8(FP)
	RET

// int64 nanotime(void) so really
// void nanotime(int64 *nsec)
TEXT runtime·nanotime1(SB), NOSPLIT, $8-8
	MOVL	$SYS_clock_gettime, AX
	MOVL	$1, BX		// CLOCK_MONOTONIC
	LEAL	0(SP), CX
	INVOKE_SYSCALL

	MOVL	0(SP), AX	// sec
	MOVL	4(SP), BX	// nsec

	// sec is in AX, nsec in BX
	// convert to DX:AX nsec
	MOVL	$1000000000, CX
	MULL	CX
	ADDL	BX, AX
	ADCL	$0, DX

	MOVL	AX, ret_lo+0(FP)
	MOVL	DX, ret_hi+4(FP)
	RET

TEXT runtime·rtsigprocmask(SB),NOSPLIT,$0
	MOVL	$SYS_rt_sigprocmask, AX
	MOVL	how+0(FP), BX
	MOVL	new+4(FP), CX
	MOVL	old+8(FP), DX
	MOVL	size+12(FP), SI
	INVOKE_SYSCALL
	CMPL	AX, $0xfffff001
	JLS	2(PC)
	INT $3
	RET

TEXT runtime·rt_sigaction(SB),NOSPLIT,$0
	MOVL	$SYS_rt_sigaction, AX
	MOVL	sig+0(FP), BX
	MOVL	new+4(FP), CX
	MOVL	old+8(FP), DX
	MOVL	size+12(FP), SI
	INVOKE_SYSCALL
	MOVL	AX, ret+16(FP)
	RET

TEXT runtime·sigfwd(SB),NOSPLIT,$12-16
	MOVL	fn+0(FP), AX
	MOVL	sig+4(FP), BX
	MOVL	info+8(FP), CX
	MOVL	ctx+12(FP), DX
	MOVL	SP, SI
	SUBL	$32, SP
	ANDL	$-15, SP	// align stack: handler might be a C function
	MOVL	BX, 0(SP)
	MOVL	CX, 4(SP)
	MOVL	DX, 8(SP)
	MOVL	SI, 12(SP)	// save SI: handler might be a Go function
	CALL	AX
	MOVL	12(SP), AX
	MOVL	AX, SP
	RET

// Called using C ABI.
TEXT runtime·sigtramp(SB),NOSPLIT|TOPFRAME,$28
	// Save callee-saved C registers.
	MOVL	BX, bx-4(SP)
	MOVL	BP, bp-8(SP)
	MOVL	SI, si-12(SP)
	MOVL	DI, di-16(SP)

	MOVL	(28+4)(SP), BX
	MOVL	BX, 0(SP)
	MOVL	(28+8)(SP), BX
	MOVL	BX, 4(SP)
	MOVL	(28+12)(SP), BX
	MOVL	BX, 8(SP)
	CALL	runtime·sigtrampgo(SB)

	MOVL	di-16(SP), DI
	MOVL	si-12(SP), SI
	MOVL	bp-8(SP),  BP
	MOVL	bx-4(SP),  BX
	RET

TEXT runtime·cgoSigtramp(SB),NOSPLIT,$0
	JMP	runtime·sigtramp(SB)

TEXT runtime·sigreturn__sigaction(SB),NOSPLIT,$0
	MOVL	$SYS_rt_sigreturn, AX
	INT	$0x80
	INT	$3	// not reached
	RET

TEXT runtime·mmap(SB),NOSPLIT,$0
	MOVL	$SYS_mmap, AX
	MOVL	addr+0(FP), BX
	MOVL	n+4(FP), CX
	MOVL	prot+8(FP), DX
	MOVL	flags+12(FP), SI
	MOVL	fd+16(FP), DI
	MOVL	off+20(FP), BP
	INVOKE_SYSCALL
	CMPL	AX, $0xfffff001
	JLS	ok
	NOTL	AX
	INCL	AX
	MOVL	$0, p+24(FP)
	MOVL	AX, err+28(FP)
	RET
ok:
	MOVL	AX, p+24(FP)
	MOVL	$0, err+28(FP)
	RET

TEXT runtime·munmap(SB),NOSPLIT,$0
	MOVL	$SYS_munmap, AX
	MOVL	addr+0(FP), BX
	MOVL	n+4(FP), CX
	INVOKE_SYSCALL
	CMPL	AX, $0xfffff001
	JLS	2(PC)
	INT $3
	RET

TEXT runtime·madvise(SB),NOSPLIT,$0
	// Freya doesn't have madvise; return 0
	MOVL	$0, AX
	MOVL	AX, ret+12(FP)
	RET

// func fcntl(fd, cmd, arg int32) (ret int32, errno int32)
TEXT runtime·fcntl(SB),NOSPLIT,$0-20
	MOVL	$SYS_fcntl, AX
	MOVL	fd+0(FP), BX
	MOVL	cmd+4(FP), CX
	MOVL	arg+8(FP), DX
	INVOKE_SYSCALL
	CMPL	AX, $0xfffff001
	JLS	noerr
	// Error case
	MOVL	$-1, ret+12(FP)
	NEGL	AX
	MOVL	AX, errno+16(FP)
	RET
noerr:
	MOVL	AX, ret+12(FP)
	MOVL	$0, errno+16(FP)
	RET

// int32 futex(int32 *uaddr, int32 op, int32 val,
//	struct timespec *timeout, int32 *uaddr2, int32 val2);
TEXT runtime·futex(SB),NOSPLIT,$0
	MOVL	$SYS_futex, AX
	MOVL	addr+0(FP), BX
	MOVL	op+4(FP), CX
	MOVL	val+8(FP), DX
	MOVL	ts+12(FP), SI
	MOVL	addr2+16(FP), DI
	MOVL	val3+20(FP), BP
	INVOKE_SYSCALL
	MOVL	AX, ret+24(FP)
	RET

// int32 clone(int32 flags, void *stack, M *mp, G *gp, void (*fn)(void));
TEXT runtime·clone(SB),NOSPLIT,$0
	MOVL	$SYS_thread_create, AX
	MOVL	flags+0(FP), BX
	MOVL	stk+4(FP), CX
	MOVL	$0, DX	// parent tid ptr
	MOVL	$0, DI	// child tid ptr

	// Copy mp, gp, fn off parent stack for use by child.
	SUBL	$16, CX
	MOVL	mp+8(FP), SI
	MOVL	SI, 0(CX)
	MOVL	gp+12(FP), SI
	MOVL	SI, 4(CX)
	MOVL	fn+16(FP), SI
	MOVL	SI, 8(CX)
	MOVL	$1234, 12(CX)

	// cannot use CALL here, because the stack changes during the
	// system call.
	INT	$0x80

	// In parent, return.
	CMPL	AX, $0
	JEQ	3(PC)
	MOVL	AX, ret+20(FP)
	RET

	// Paranoia: check that SP is as we expect.
	NOP	SP // tell vet SP changed - stop checking offsets
	MOVL	12(SP), BP
	CMPL	BP, $1234
	JEQ	2(PC)
	INT	$3

	// Initialize AX to tid
	MOVL	$SYS_getpid, AX
	INVOKE_SYSCALL

	MOVL	0(SP), BX	    // m
	MOVL	4(SP), DX	    // g
	MOVL	8(SP), SI	    // fn

	CMPL	BX, $0
	JEQ	nog
	CMPL	DX, $0
	JEQ	nog

	MOVL	AX, m_procid(BX)	// save tid as m->procid

	// Now set up TLS and g.
	get_tls(AX)
	MOVL	DX, g(AX)
	MOVL	BX, g_m(DX)

nog:
	CALL	SI	// fn()
	CALL	exit1<>(SB)
	MOVL	$0x1234, 0x1005

TEXT runtime·sigaltstack(SB),NOSPLIT,$-8
	// Freya: sigaltstack is handled by the kernel internally.
	RET

TEXT runtime·osyield(SB),NOSPLIT,$0
	MOVL	$SYS_thread_yield, AX
	INVOKE_SYSCALL
	RET

TEXT runtime·sched_getaffinity(SB),NOSPLIT,$0
	// Freya doesn't have sched_getaffinity
	MOVL	$0, AX
	MOVL	AX, ret+12(FP)
	RET

// func sbrk0() uintptr
TEXT runtime·sbrk0(SB),NOSPLIT,$0-4
	MOVL	$0, AX
	MOVL	AX, ret+0(FP)
	RET

// setldt(int entry, int address, int limit)
// Freya: no-op for now since we skip ldt0setup
TEXT runtime·setldt(SB),NOSPLIT,$0
	RET

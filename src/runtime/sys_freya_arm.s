// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// System calls and other sys.stuff for arm, Freya
//

#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"

#define CLOCK_REALTIME	0
#define CLOCK_MONOTONIC	1

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

TEXT runtime·open(SB),NOSPLIT,$0
	MOVW	name+0(FP), R0
	MOVW	mode+4(FP), R1
	MOVW	perm+8(FP), R2
	MOVW	$SYS_open, R7
	SWI	$0
	MOVW	$0xfffff001, R1
	CMP	R1, R0
	MOVW.HI	$-1, R0
	MOVW	R0, ret+12(FP)
	RET

TEXT runtime·closefd(SB),NOSPLIT,$0
	MOVW	fd+0(FP), R0
	MOVW	$SYS_close, R7
	SWI	$0
	MOVW	$0xfffff001, R1
	CMP	R1, R0
	MOVW.HI	$-1, R0
	MOVW	R0, ret+4(FP)
	RET

TEXT runtime·write1(SB),NOSPLIT,$0
	MOVW	fd+0(FP), R0
	MOVW	p+4(FP), R1
	MOVW	n+8(FP), R2
	MOVW	$SYS_write, R7
	SWI	$0
	MOVW	R0, ret+12(FP)
	RET

TEXT runtime·read(SB),NOSPLIT,$0
	MOVW	fd+0(FP), R0
	MOVW	p+4(FP), R1
	MOVW	n+8(FP), R2
	MOVW	$SYS_read, R7
	SWI	$0
	MOVW	R0, ret+12(FP)
	RET

// func pipe2(flags int32) (r, w int32, errno int32)
TEXT runtime·pipe2(SB),NOSPLIT,$0-16
	MOVW	$r+4(FP), R0
	MOVW	flags+0(FP), R1
	MOVW	$SYS_pipe_create, R7
	SWI	$0
	MOVW	R0, errno+12(FP)
	RET

TEXT runtime·exit(SB),NOSPLIT|NOFRAME,$0
	MOVW	code+0(FP), R0
	MOVW	$SYS_exit, R7
	SWI	$0
	MOVW	$1234, R0
	MOVW	$1002, R1
	MOVW	R0, (R1)	// fail hard

TEXT exit1<>(SB),NOSPLIT|NOFRAME,$0
	MOVW	code+0(FP), R0
	MOVW	$SYS_thread_exit, R7
	SWI	$0
	MOVW	$1234, R0
	MOVW	$1003, R1
	MOVW	R0, (R1)	// fail hard

// func exitThread(wait *atomic.Uint32)
TEXT runtime·exitThread(SB),NOSPLIT|NOFRAME,$0-4
	MOVW	wait+0(FP), R0
	// We're done using the stack.
	MOVW	$0, R1
	MOVW	R1, (R0)
	MOVW	$0, R0	// exit code
	MOVW	$SYS_thread_exit, R7
	SWI	$0
	MOVW	$1234, R0
	MOVW	$1004, R1
	MOVW	R0, (R1)	// fail hard
	JMP	0(PC)

TEXT runtime·gettid(SB),NOSPLIT,$0-4
	MOVW	$SYS_getpid, R7
	SWI	$0
	MOVW	R0, ret+0(FP)
	RET

TEXT	runtime·raise(SB),NOSPLIT|NOFRAME,$0
	MOVW	$SYS_getpid, R7
	SWI	$0
	// arg 1 pid already in R0
	MOVW	sig+0(FP), R1	// arg 2 - signal
	MOVW	$SYS_process_kill, R7
	SWI	$0
	RET

TEXT	runtime·raiseproc(SB),NOSPLIT|NOFRAME,$0
	MOVW	$SYS_getpid, R7
	SWI	$0
	// arg 1 pid already in R0
	MOVW	sig+0(FP), R1	// arg 2 - signal
	MOVW	$SYS_process_kill, R7
	SWI	$0
	RET

TEXT ·getpid(SB),NOSPLIT,$0-4
	MOVW	$SYS_getpid, R7
	SWI	$0
	MOVW	R0, ret+0(FP)
	RET

TEXT ·tgkill(SB),NOSPLIT,$0-12
	MOVW	tgid+0(FP), R0
	MOVW	tid+4(FP), R1
	MOVW	sig+8(FP), R2
	MOVW	$SYS_process_kill, R7
	SWI	$0
	RET

TEXT runtime·mmap(SB),NOSPLIT,$0
	MOVW	addr+0(FP), R0
	MOVW	n+4(FP), R1
	MOVW	prot+8(FP), R2
	MOVW	flags+12(FP), R3
	MOVW	fd+16(FP), R4
	MOVW	off+20(FP), R5
	MOVW	$SYS_mmap, R7
	SWI	$0
	MOVW	$0xfffff001, R6
	CMP		R6, R0
	MOVW	$0, R1
	RSB.HI	$0, R0
	MOVW.HI	R0, R1		// if error, put in R1
	MOVW.HI	$0, R0
	MOVW	R0, p+24(FP)
	MOVW	R1, err+28(FP)
	RET

TEXT runtime·munmap(SB),NOSPLIT,$0
	MOVW	addr+0(FP), R0
	MOVW	n+4(FP), R1
	MOVW	$SYS_munmap, R7
	SWI	$0
	MOVW	$0xfffff001, R6
	CMP 	R6, R0
	MOVW.HI	$0, R8  // crash on syscall failure
	MOVW.HI	R8, (R8)
	RET

TEXT runtime·madvise(SB),NOSPLIT,$0
	MOVW	addr+0(FP), R0
	MOVW	n+4(FP), R1
	MOVW	flags+8(FP), R2
	// Freya doesn't have madvise; just return 0
	MOVW	$0, R0
	MOVW	R0, ret+12(FP)
	RET

// func fcntl(fd, cmd, arg int32) (ret int32, errno int32)
TEXT runtime·fcntl(SB),NOSPLIT,$0-20
	MOVW	fd+0(FP), R0
	MOVW	cmd+4(FP), R1
	MOVW	arg+8(FP), R2
	MOVW	$SYS_fcntl, R7
	SWI	$0
	MOVW	$0xfffff001, R1
	CMP	R1, R0
	BLS	noerr
	// Error case
	MOVW	$-1, R1
	MOVW	R1, ret+12(FP)
	RSB	$0, R0, R0
	MOVW	R0, errno+16(FP)
	RET
noerr:
	MOVW	R0, ret+12(FP)
	MOVW	$0, R1
	MOVW	R1, errno+16(FP)
	RET

TEXT runtime·setitimer(SB),NOSPLIT,$0
	// Freya: no-op for now
	RET

TEXT runtime·mincore(SB),NOSPLIT,$0
	// Freya doesn't have mincore
	MOVW	$-1, R0
	MOVW	R0, ret+12(FP)
	RET

// func walltime() (sec int64, nsec int32)
TEXT runtime·walltime(SB),NOSPLIT,$12-12
	MOVW	$CLOCK_REALTIME, R0
	MOVW	$spec-12(SP), R1	// timespec

	MOVW	$SYS_clock_gettime, R7
	SWI	$0

	MOVW	sec-12(SP), R0  // sec
	MOVW	nsec-8(SP), R2  // nsec

	MOVW	R0, sec_lo+0(FP)
	MOVW	$0, R1
	MOVW	R1, sec_hi+4(FP)
	MOVW	R2, nsec+8(FP)
	RET

// func nanotime1() int64
TEXT runtime·nanotime1(SB),NOSPLIT,$12-8
	MOVW	$CLOCK_MONOTONIC, R0
	MOVW	$spec-12(SP), R1	// timespec

	MOVW	$SYS_clock_gettime, R7
	SWI	$0

	MOVW	sec-12(SP), R0  // sec
	MOVW	nsec-8(SP), R2  // nsec

	MOVW	$1000000000, R3
	MULLU	R0, R3, (R1, R0)
	ADD.S	R2, R0
	ADC	$0, R1	// Add carry bit to upper half.

	MOVW	R0, ret_lo+0(FP)
	MOVW	R1, ret_hi+4(FP)

	RET

// int32 futex(int32 *uaddr, int32 op, int32 val,
//	struct timespec *timeout, int32 *uaddr2, int32 val2);
TEXT runtime·futex(SB),NOSPLIT,$0
	MOVW    addr+0(FP), R0
	MOVW    op+4(FP), R1
	MOVW    val+8(FP), R2
	MOVW    ts+12(FP), R3
	MOVW    addr2+16(FP), R4
	MOVW    val3+20(FP), R5
	MOVW	$SYS_futex, R7
	SWI	$0
	MOVW	R0, ret+24(FP)
	RET

// int32 clone(int32 flags, void *stack, M *mp, G *gp, void (*fn)(void));
TEXT runtime·clone(SB),NOSPLIT,$0
	MOVW	flags+0(FP), R0
	MOVW	stk+4(FP), R1
	MOVW	$0, R2	// parent tid ptr
	MOVW	$0, R3	// tls_val
	MOVW	$0, R4	// child tid ptr
	MOVW	$0, R5

	// Copy mp, gp, fn off parent stack for use by child.
	MOVW	$-16(R1), R1
	MOVW	mp+8(FP), R6
	MOVW	R6, 0(R1)
	MOVW	gp+12(FP), R6
	MOVW	R6, 4(R1)
	MOVW	fn+16(FP), R6
	MOVW	R6, 8(R1)
	MOVW	$1234, R6
	MOVW	R6, 12(R1)

	MOVW	$SYS_thread_create, R7
	SWI	$0

	// In parent, return.
	CMP	$0, R0
	BEQ	3(PC)
	MOVW	R0, ret+20(FP)
	RET

	// Paranoia: check that SP is as we expect. Use R13 to avoid linker 'fixup'
	NOP	R13	// tell vet SP/R13 changed - stop checking offsets
	MOVW	12(R13), R0
	MOVW	$1234, R1
	CMP	R0, R1
	BEQ	2(PC)
	BL	runtime·abort(SB)

	MOVW	0(R13), R8    // m
	MOVW	4(R13), R0    // g

	CMP	$0, R8
	BEQ	nog
	CMP	$0, R0
	BEQ	nog

	MOVW	R0, g
	MOVW	R8, g_m(g)

	// paranoia; check they are not nil
	MOVW	0(R8), R0
	MOVW	0(g), R0

	BL	runtime·emptyfunc(SB)	// fault if stack check is wrong

	// Initialize m->procid to tid
	MOVW	$SYS_getpid, R7
	SWI	$0
	MOVW	g_m(g), R8
	MOVW	R0, m_procid(R8)

nog:
	// Call fn
	MOVW	8(R13), R0
	MOVW	$16(R13), R13
	BL	(R0)

	// It shouldn't return. If it does, exit that thread.
	SUB	$16, R13 // restore the stack pointer to avoid memory corruption
	MOVW	$0, R0
	MOVW	R0, 4(R13)
	BL	exit1<>(SB)

	MOVW	$1234, R0
	MOVW	$1005, R1
	MOVW	R0, (R1)

TEXT runtime·sigaltstack(SB),NOSPLIT,$0
	// Freya: sigaltstack is handled by the kernel internally.
	RET

TEXT runtime·sigfwd(SB),NOSPLIT,$0-16
	MOVW	sig+4(FP), R0
	MOVW	info+8(FP), R1
	MOVW	ctx+12(FP), R2
	MOVW	fn+0(FP), R11
	MOVW	R13, R4
	SUB	$24, R13
	BIC	$0x7, R13 // alignment for ELF ABI
	BL	(R11)
	MOVW	R4, R13
	RET

TEXT runtime·sigtramp(SB),NOSPLIT|TOPFRAME,$0
	// Reserve space for callee-save registers and arguments.
	MOVM.DB.W [R4-R11], (R13)
	SUB	$16, R13

	// this might be called in external code context,
	// where g is not set.
	// first save R0, because runtime·load_g will clobber it
	MOVW	R0, 4(R13)
	MOVB	runtime·iscgo(SB), R0
	CMP 	$0, R0
	BL.NE	runtime·load_g(SB)

	MOVW	R1, 8(R13)
	MOVW	R2, 12(R13)
	MOVW  	$runtime·sigtrampgo(SB), R11
	BL	(R11)

	// Restore callee-save registers.
	ADD	$16, R13
	MOVM.IA.W (R13), [R4-R11]

	RET

TEXT runtime·cgoSigtramp(SB),NOSPLIT,$0
	MOVW  	$runtime·sigtramp(SB), R11
	B	(R11)

TEXT runtime·rtsigprocmask(SB),NOSPLIT,$0
	MOVW	how+0(FP), R0
	MOVW	new+4(FP), R1
	MOVW	old+8(FP), R2
	MOVW	size+12(FP), R3
	MOVW	$SYS_rt_sigprocmask, R7
	SWI	$0
	RET

TEXT runtime·rt_sigaction(SB),NOSPLIT,$0
	MOVW	sig+0(FP), R0
	MOVW	new+4(FP), R1
	MOVW	old+8(FP), R2
	MOVW	size+12(FP), R3
	MOVW	$SYS_rt_sigaction, R7
	SWI	$0
	MOVW	R0, ret+16(FP)
	RET

TEXT runtime·usleep(SB),NOSPLIT,$12
	MOVW	usec+0(FP), R0
	CALL	runtime·usplitR0(SB)
	MOVW	R0, 4(R13)
	MOVW	$1000, R0	// usec to nsec
	MUL	R0, R1
	MOVW	R1, 8(R13)
	MOVW	$4(R13), R0
	MOVW	$0, R1
	MOVW	$SYS_nanosleep, R7
	SWI	$0
	RET

TEXT runtime·osyield(SB),NOSPLIT,$0
	MOVW	$SYS_thread_yield, R7
	SWI	$0
	RET

TEXT runtime·sched_getaffinity(SB),NOSPLIT,$0
	// Freya doesn't have sched_getaffinity
	MOVW	$0, R0
	MOVW	R0, ret+12(FP)
	RET

// func sbrk0() uintptr
TEXT runtime·sbrk0(SB),NOSPLIT,$0-4
	MOVW	$0, R0
	MOVW	R0, ret+0(FP)
	RET

TEXT ·publicationBarrier(SB),NOSPLIT|NOFRAME,$0-0
	B	runtime·armPublicationBarrier(SB)

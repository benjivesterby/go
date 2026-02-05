// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// System calls and other sys.stuff for arm64, Freya
//

#include "go_asm.h"
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

TEXT runtime·exit(SB),NOSPLIT|NOFRAME,$0-4
	MOVW	code+0(FP), R0
	MOVD	$SYS_exit, R8
	SVC
	RET

// func exitThread(wait *atomic.Uint32)
TEXT runtime·exitThread(SB),NOSPLIT|NOFRAME,$0-8
	MOVD	wait+0(FP), R0
	// We're done using the stack.
	MOVW	$0, R1
	STLRW	R1, (R0)
	MOVW	$0, R0	// exit code
	MOVD	$SYS_thread_exit, R8
	SVC
	JMP	0(PC)

TEXT runtime·open(SB),NOSPLIT|NOFRAME,$0-20
	MOVD	name+0(FP), R0
	MOVW	mode+8(FP), R1
	MOVW	perm+12(FP), R2
	MOVD	$SYS_open, R8
	SVC
	CMN	$4095, R0
	BCC	done
	MOVW	$-1, R0
done:
	MOVW	R0, ret+16(FP)
	RET

TEXT runtime·closefd(SB),NOSPLIT|NOFRAME,$0-12
	MOVW	fd+0(FP), R0
	MOVD	$SYS_close, R8
	SVC
	CMN	$4095, R0
	BCC	done
	MOVW	$-1, R0
done:
	MOVW	R0, ret+8(FP)
	RET

TEXT runtime·write1(SB),NOSPLIT|NOFRAME,$0-28
	MOVD	fd+0(FP), R0
	MOVD	p+8(FP), R1
	MOVW	n+16(FP), R2
	MOVD	$SYS_write, R8
	SVC
	MOVW	R0, ret+24(FP)
	RET

TEXT runtime·read(SB),NOSPLIT|NOFRAME,$0-28
	MOVW	fd+0(FP), R0
	MOVD	p+8(FP), R1
	MOVW	n+16(FP), R2
	MOVD	$SYS_read, R8
	SVC
	MOVW	R0, ret+24(FP)
	RET

// func pipe2(flags int32) (r, w int32, errno int32)
TEXT runtime·pipe2(SB),NOSPLIT|NOFRAME,$0-20
	MOVD	$r+8(FP), R0
	MOVW	flags+0(FP), R1
	MOVW	$SYS_pipe_create, R8
	SVC
	MOVW	R0, errno+16(FP)
	RET

TEXT runtime·usleep(SB),NOSPLIT,$24-4
	MOVWU	usec+0(FP), R3
	MOVD	R3, R5
	MOVW	$1000000, R4
	UDIV	R4, R3
	MOVD	R3, 8(RSP)
	MUL	R3, R4
	SUB	R4, R5
	MOVW	$1000, R4
	MUL	R4, R5
	MOVD	R5, 16(RSP)

	// nanosleep(&ts, 0)
	ADD	$8, RSP, R0
	MOVD	$0, R1
	MOVD	$SYS_nanosleep, R8
	SVC
	RET

TEXT runtime·gettid(SB),NOSPLIT,$0-4
	MOVD	$SYS_getpid, R8
	SVC
	MOVW	R0, ret+0(FP)
	RET

TEXT runtime·raise(SB),NOSPLIT|NOFRAME,$0
	MOVD	$SYS_getpid, R8
	SVC
	MOVW	R0, R0		// arg 1 pid
	MOVW	sig+0(FP), R1	// arg 2
	MOVD	$SYS_process_kill, R8
	SVC
	RET

TEXT runtime·raiseproc(SB),NOSPLIT|NOFRAME,$0
	MOVD	$SYS_getpid, R8
	SVC
	MOVW	R0, R0		// arg 1 pid
	MOVW	sig+0(FP), R1	// arg 2
	MOVD	$SYS_process_kill, R8
	SVC
	RET

TEXT ·getpid(SB),NOSPLIT|NOFRAME,$0-8
	MOVD	$SYS_getpid, R8
	SVC
	MOVD	R0, ret+0(FP)
	RET

TEXT ·tgkill(SB),NOSPLIT,$0-24
	MOVD	tgid+0(FP), R0
	MOVD	tid+8(FP), R1
	MOVD	sig+16(FP), R2
	MOVD	$SYS_process_kill, R8
	SVC
	RET

// func walltime() (sec int64, nsec int32)
TEXT runtime·walltime(SB),NOSPLIT,$24-12
	MOVW	$CLOCK_REALTIME, R0
	ADD	$8, RSP, R1
	MOVD	$SYS_clock_gettime, R8
	SVC
	MOVD	8(RSP), R3	// sec
	MOVD	16(RSP), R5	// nsec
	MOVD	R3, sec+0(FP)
	MOVW	R5, nsec+8(FP)
	RET

TEXT runtime·nanotime1(SB),NOSPLIT,$24-8
	MOVW	$CLOCK_MONOTONIC, R0
	ADD	$8, RSP, R1
	MOVD	$SYS_clock_gettime, R8
	SVC
	MOVD	8(RSP), R3	// sec
	MOVD	16(RSP), R5	// nsec
	// sec is in R3, nsec in R5
	// return nsec in R3
	MOVD	$1000000000, R4
	MUL	R4, R3
	ADD	R5, R3
	MOVD	R3, ret+0(FP)
	RET

TEXT runtime·rtsigprocmask(SB),NOSPLIT|NOFRAME,$0-28
	MOVW	how+0(FP), R0
	MOVD	new+8(FP), R1
	MOVD	old+16(FP), R2
	MOVW	size+24(FP), R3
	MOVD	$SYS_rt_sigprocmask, R8
	SVC
	CMN	$4095, R0
	BCC	done
	MOVD	$0, R0
	MOVD	R0, (R0)	// crash
done:
	RET

TEXT runtime·rt_sigaction(SB),NOSPLIT|NOFRAME,$0-36
	MOVD	sig+0(FP), R0
	MOVD	new+8(FP), R1
	MOVD	old+16(FP), R2
	MOVD	size+24(FP), R3
	MOVD	$SYS_rt_sigaction, R8
	SVC
	MOVW	R0, ret+32(FP)
	RET

TEXT runtime·sigfwd(SB),NOSPLIT,$0-32
	MOVW	sig+8(FP), R0
	MOVD	info+16(FP), R1
	MOVD	ctx+24(FP), R2
	MOVD	fn+0(FP), R11
	BL	(R11)
	RET

TEXT runtime·sigtramp(SB),NOSPLIT|TOPFRAME,$176
	// Save callee-save registers in the case of signal forwarding.
	STP	(R19, R20), 4*8(RSP)
	STP	(R21, R22), 6*8(RSP)
	STP	(R23, R24), 8*8(RSP)
	STP	(R25, R26), 10*8(RSP)
	STP	(R27, g), 12*8(RSP)

	// this might be called in external code context,
	// where g is not set.
	// first save R0, because runtime·load_g will clobber it
	MOVW	R0, 8(RSP)
	MOVBU	runtime·iscgo(SB), R0
	CBZ	R0, 2(PC)
	BL	runtime·load_g(SB)

	// Restore signum to R0.
	MOVW	8(RSP), R0
	// R1 and R2 already contain info and ctx, respectively.
	MOVD	$runtime·sigtrampgo<ABIInternal>(SB), R3
	BL	(R3)

	// Restore callee-save registers.
	LDP	4*8(RSP), (R19, R20)
	LDP	6*8(RSP), (R21, R22)
	LDP	8*8(RSP), (R23, R24)
	LDP	10*8(RSP), (R25, R26)
	LDP	12*8(RSP), (R27, g)

	RET

TEXT runtime·cgoSigtramp(SB),NOSPLIT|NOFRAME,$0
	B	runtime·sigtramp(SB)

// func mmap(addr unsafe.Pointer, n uintptr, prot, flags, fd int32, off uint32) (p unsafe.Pointer, err int)
TEXT runtime·mmap(SB),NOSPLIT|NOFRAME,$0
	MOVD	addr+0(FP), R0
	MOVD	n+8(FP), R1
	MOVW	prot+16(FP), R2
	MOVW	flags+20(FP), R3
	MOVW	fd+24(FP), R4
	MOVW	off+28(FP), R5

	MOVD	$SYS_mmap, R8
	SVC
	CMN	$4095, R0
	BCC	ok
	NEG	R0,R0
	MOVD	$0, p+32(FP)
	MOVD	R0, err+40(FP)
	RET
ok:
	MOVD	R0, p+32(FP)
	MOVD	$0, err+40(FP)
	RET

// func munmap(addr unsafe.Pointer, n uintptr)
TEXT runtime·munmap(SB),NOSPLIT|NOFRAME,$0
	MOVD	addr+0(FP), R0
	MOVD	n+8(FP), R1
	MOVD	$SYS_munmap, R8
	SVC
	CMN	$4095, R0
	BCC	cool
	MOVD	R0, 0xf0(R0)
cool:
	RET

// func madvise(addr unsafe.Pointer, n uintptr, flags int32) int32
TEXT runtime·madvise(SB),NOSPLIT|NOFRAME,$0-28
	// Freya doesn't have madvise; return 0
	MOVW	$0, R0
	MOVW	R0, ret+24(FP)
	RET

// func fcntl(fd, cmd, arg int32) (ret int32, errno int32)
TEXT runtime·fcntl(SB),NOSPLIT|NOFRAME,$0-20
	MOVW	fd+0(FP), R0
	MOVW	cmd+4(FP), R1
	MOVW	arg+8(FP), R2
	MOVD	$SYS_fcntl, R8
	SVC
	CMN	$4095, R0
	BCC	noerr
	// Error case
	MOVW	$-1, R1
	MOVW	R1, ret+12(FP)
	NEG	R0, R0
	MOVW	R0, errno+16(FP)
	RET
noerr:
	MOVW	R0, ret+12(FP)
	MOVW	$0, R1
	MOVW	R1, errno+16(FP)
	RET

// int64 futex(int32 *uaddr, int32 op, int32 val,
//	struct timespec *timeout, int32 *uaddr2, int32 val3);
TEXT runtime·futex(SB),NOSPLIT|NOFRAME,$0
	MOVD	addr+0(FP), R0
	MOVW	op+8(FP), R1
	MOVW	val+12(FP), R2
	MOVD	ts+16(FP), R3
	MOVD	addr2+24(FP), R4
	MOVW	val3+32(FP), R5
	MOVD	$SYS_futex, R8
	SVC
	MOVW	R0, ret+40(FP)
	RET

// int64 clone(int32 flags, void *stk, M *mp, G *gp, void (*fn)(void));
TEXT runtime·clone(SB),NOSPLIT|NOFRAME,$0
	MOVW	flags+0(FP), R0
	MOVD	stk+8(FP), R1

	// Copy mp, gp, fn off parent stack for use by child.
	MOVD	mp+16(FP), R10
	MOVD	gp+24(FP), R11
	MOVD	fn+32(FP), R12

	MOVD	R10, -8(R1)
	MOVD	R11, -16(R1)
	MOVD	R12, -24(R1)
	MOVD	$1234, R10
	MOVD	R10, -32(R1)

	MOVD	$SYS_thread_create, R8
	SVC

	// In parent, return.
	CMP	ZR, R0
	BEQ	child
	MOVW	R0, ret+40(FP)
	RET
child:

	// In child, on new stack.
	MOVD	-32(RSP), R10
	MOVD	$1234, R0
	CMP	R0, R10
	BEQ	good
	MOVD	$0, R0
	MOVD	R0, (R0)	// crash

good:
	// Initialize m->procid to tid
	MOVD	$SYS_getpid, R8
	SVC

	MOVD	-24(RSP), R12     // fn
	MOVD	-16(RSP), R11     // g
	MOVD	-8(RSP), R10      // m

	CMP	$0, R10
	BEQ	nog
	CMP	$0, R11
	BEQ	nog

	MOVD	R0, m_procid(R10)

	// In child, set up new stack
	MOVD	R10, g_m(R11)
	MOVD	R11, g

nog:
	// Call fn
	MOVD	R12, R0
	BL	(R0)

	// It shouldn't return. If it does, exit that thread.
	MOVW	$111, R0
again:
	MOVD	$SYS_thread_exit, R8
	SVC
	B	again	// keep exiting

TEXT runtime·sigaltstack(SB),NOSPLIT|NOFRAME,$0
	// Freya: sigaltstack is handled by the kernel internally.
	RET

TEXT runtime·osyield(SB),NOSPLIT|NOFRAME,$0
	MOVD	$SYS_thread_yield, R8
	SVC
	RET

TEXT runtime·sched_getaffinity(SB),NOSPLIT|NOFRAME,$0
	// Freya doesn't have sched_getaffinity
	MOVW	$0, R0
	MOVW	R0, ret+24(FP)
	RET

// func sbrk0() uintptr
TEXT runtime·sbrk0(SB),NOSPLIT,$0-8
	MOVD	$0, R0
	MOVD	R0, ret+0(FP)
	RET

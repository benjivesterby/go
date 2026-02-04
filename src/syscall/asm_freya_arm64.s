// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

//
// System calls for arm64, Freya
//

// func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno)
TEXT ·RawSyscall6(SB),NOSPLIT,$0-80
	MOVD	a1+8(FP), R0
	MOVD	a2+16(FP), R1
	MOVD	a3+24(FP), R2
	MOVD	a4+32(FP), R3
	MOVD	a5+40(FP), R4
	MOVD	a6+48(FP), R5
	MOVD	trap+0(FP), R8	// syscall entry
	SVC
	CMN	$4095, R0
	BCC	ok
	MOVD	$-1, R4
	MOVD	R4, r1+56(FP)	// r1
	MOVD	ZR, r2+64(FP)	// r2
	NEG	R0, R0
	MOVD	R0, err+72(FP)	// errno
	RET
ok:
	MOVD	R0, r1+56(FP)	// r1
	MOVD	R1, r2+64(FP)	// r2
	MOVD	ZR, err+72(FP)	// errno
	RET

// func rawSyscallNoError(trap, a1, a2, a3 uintptr) (r1, r2 uintptr)
TEXT ·rawSyscallNoError(SB),NOSPLIT,$0-48
	MOVD	a1+8(FP), R0
	MOVD	a2+16(FP), R1
	MOVD	a3+24(FP), R2
	MOVD	$0, R3
	MOVD	$0, R4
	MOVD	$0, R5
	MOVD	trap+0(FP), R8	// syscall entry
	SVC
	MOVD	R0, r1+32(FP)
	MOVD	R1, r2+40(FP)
	RET

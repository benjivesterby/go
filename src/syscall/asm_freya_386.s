// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"
#include "funcdata.h"

//
// System call support for 386, Freya
// Uses Linux-style error handling (negative values are errors)
//

// func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, errno uintptr)
TEXT	·RawSyscall6(SB),NOSPLIT,$0-40
	MOVL	trap+0(FP), AX // syscall number
	MOVL	a1+4(FP), BX   // a1
	MOVL	a2+8(FP), CX   // a2
	MOVL	a3+12(FP), DX  // a3
	MOVL	a4+16(FP), SI  // a4
	MOVL	a5+20(FP), DI  // a5
	MOVL	a6+24(FP), BP  // a6
	INT	$0x80
	CMPL	AX, $0xfffff001
	JLS	ok
	MOVL	$-1, r1+28(FP)
	MOVL	$0, r2+32(FP)
	NEGL	AX
	MOVL	AX, err+36(FP)
	RET
ok:
	MOVL	AX, r1+28(FP)
	MOVL	DX, r2+32(FP)
	MOVL	$0, err+36(FP)
	RET

// func rawSyscallNoError(trap, a1, a2, a3 uintptr) (r1, r2 uintptr)
TEXT	·rawSyscallNoError(SB),NOSPLIT,$0-24
	MOVL	trap+0(FP), AX // syscall number
	MOVL	a1+4(FP), BX   // a1
	MOVL	a2+8(FP), CX   // a2
	MOVL	a3+12(FP), DX  // a3
	INT	$0x80
	MOVL	AX, r1+16(FP)
	MOVL	DX, r2+20(FP)
	RET

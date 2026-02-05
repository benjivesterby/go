// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"
#include "funcdata.h"

//
// System call support for ARM, Freya
// Uses Linux-style error handling (negative values are errors)
//

// func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, errno uintptr)
TEXT	·RawSyscall6(SB),NOSPLIT,$0-40
	MOVW trap+0(FP), R7 // syscall number
	MOVW a1+4(FP), R0 // a1
	MOVW a2+8(FP), R1 // a2
	MOVW a3+12(FP), R2 // a3
	MOVW a4+16(FP), R3 // a4
	MOVW a5+20(FP), R4 // a5
	MOVW a6+24(FP), R5 // a6
	SWI $0 // syscall
	MOVW $0xfffff001, R4
	CMP R4, R0
	BLS ok
	MOVW $-1, R3
	MOVW R3, r1+28(FP) // r1
	MOVW $0, R3
	MOVW R3, r2+32(FP) // r2
	RSB $0, R0, R0
	MOVW R0, err+36(FP) // errno
	RET
ok:
	MOVW R0, r1+28(FP) // r1
	MOVW R1, r2+32(FP) // r2
	MOVW $0, R3
	MOVW R3, err+36(FP) // errno
	RET

// func rawSyscallNoError(trap, a1, a2, a3 uintptr) (r1, r2 uintptr)
TEXT	·rawSyscallNoError(SB),NOSPLIT,$0-24
	MOVW trap+0(FP), R7 // syscall number
	MOVW a1+4(FP), R0 // a1
	MOVW a2+8(FP), R1 // a2
	MOVW a3+12(FP), R2 // a3
	SWI $0 // syscall
	MOVW R0, r1+16(FP) // r1
	MOVW R1, r2+20(FP) // r2
	RET

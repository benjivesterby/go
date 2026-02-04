// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

TEXT _rt0_arm_freya(SB),NOSPLIT|NOFRAME,$0
	MOVW	(R13), R0	// argc
	MOVW	$4(R13), R1	// argv
	B	runtime·rt0_go(SB)

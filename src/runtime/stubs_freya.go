// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freya

package runtime

// sbrk0 returns the current process brk, or 0 if not implemented.
// Freya doesn't implement brk/sbrk, so this always returns 0.
func sbrk0() uintptr

// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freya

package runtime

// secureMode is only ever mutated in schedinit, so we don't need to worry about
// synchronization primitives.
var secureMode bool

func initSecureMode() {
	// Freya does not currently support setuid/setgid binaries.
	secureMode = false
}

func isSecureMode() bool {
	return secureMode
}

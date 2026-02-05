// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freya

package sysrand

// Freya uses /dev/urandom for random data.
func read(b []byte) error {
	return urandomRead(b)
}

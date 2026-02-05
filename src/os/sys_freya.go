// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import "syscall"

func hostname() (name string, err error) {
	var un syscall.Utsname
	err = syscall.Uname(&un)
	if err != nil {
		return "", NewSyscallError("uname", err)
	}

	var buf [512]byte
	for i, b := range un.Nodename[:] {
		buf[i] = uint8(b)
		if b == 0 {
			name = string(buf[:i])
			break
		}
	}
	return name, nil
}

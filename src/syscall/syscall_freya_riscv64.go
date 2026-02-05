// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freya && riscv64

package syscall

func (cmsg *Cmsghdr) SetLen(length int) {
	cmsg.Len = uint64(length)
}

func (iov *Iovec) SetLen(length int) {
	iov.Len = uint64(length)
}

func (msghdr *Msghdr) SetControllen(length int) {
	msghdr.Controllen = uint64(length)
}

func setTimespec(sec, nsec int64) Timespec {
	return Timespec{Sec: sec, Nsec: nsec}
}

func setTimeval(sec, usec int64) Timeval {
	return Timeval{Sec: sec, Usec: usec}
}

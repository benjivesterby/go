// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Freya system calls.

//go:build freya

package syscall

import (
	"unsafe"
)

// Pull in entersyscall/exitsyscall for Syscall/Syscall6.
//go:linkname runtime_entersyscall runtime.entersyscall
func runtime_entersyscall()

//go:linkname runtime_exitsyscall runtime.exitsyscall
func runtime_exitsyscall()

// Assembly-defined functions
func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno)
func rawSyscallNoError(trap, a1, a2, a3 uintptr) (r1, r2 uintptr)

//go:uintptrkeepalive
//go:nosplit
//go:norace
//go:linkname RawSyscall
func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno) {
	return RawSyscall6(trap, a1, a2, a3, 0, 0, 0)
}

//go:uintptrkeepalive
//go:nosplit
//go:linkname Syscall
func Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno) {
	runtime_entersyscall()
	r1, r2, err = RawSyscall6(trap, a1, a2, a3, 0, 0, 0)
	runtime_exitsyscall()
	return
}

//go:uintptrkeepalive
//go:nosplit
//go:linkname Syscall6
func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno) {
	runtime_entersyscall()
	r1, r2, err = RawSyscall6(trap, a1, a2, a3, a4, a5, a6)
	runtime_exitsyscall()
	return
}

// Dirent helpers

func direntIno(buf []byte) (uint64, bool) {
	return readInt(buf, unsafe.Offsetof(Dirent{}.Ino), unsafe.Sizeof(Dirent{}.Ino))
}

func direntReclen(buf []byte) (uint64, bool) {
	return readInt(buf, unsafe.Offsetof(Dirent{}.Reclen), unsafe.Sizeof(Dirent{}.Reclen))
}

func direntNamlen(buf []byte) (uint64, bool) {
	reclen, ok := direntReclen(buf)
	if !ok {
		return 0, false
	}
	return reclen - uint64(unsafe.Offsetof(Dirent{}.Name)), true
}

// WaitStatus represents the status of a process
type WaitStatus uint32

func (w WaitStatus) Exited() bool       { return w&0x7f == 0 }
func (w WaitStatus) Signaled() bool     { return w&0x7f != 0x7f && w&0x7f != 0 }
func (w WaitStatus) Stopped() bool      { return w&0xff == 0x7f }
func (w WaitStatus) Continued() bool    { return w == 0xffff }
func (w WaitStatus) CoreDump() bool     { return w.Signaled() && w&0x80 != 0 }
func (w WaitStatus) ExitStatus() int    { return int(w >> 8 & 0xff) }
func (w WaitStatus) Signal() Signal     { return Signal(w & 0x7f) }
func (w WaitStatus) StopSignal() Signal { return Signal(w >> 8 & 0xff) }
func (w WaitStatus) TrapCause() int     { return 0 }

// Wait4 waits for process state changes
func Wait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (wpid int, err error) {
	var status _C_int
	r1, _, e1 := Syscall6(SYS_WAIT4, uintptr(pid), uintptr(unsafe.Pointer(&status)), uintptr(options), uintptr(unsafe.Pointer(rusage)), 0, 0)
	wpid = int(r1)
	if wstatus != nil {
		*wstatus = WaitStatus(status)
	}
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// readlen reads into a buffer
func readlen(fd int, buf *byte, nbuf int) (n int, err error) {
	r0, _, e1 := Syscall(SYS_READ, uintptr(fd), uintptr(unsafe.Pointer(buf)), uintptr(nbuf))
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// forkAndExecFailureCleanup cleans up after fork/exec failure
func forkAndExecFailureCleanup(attr *ProcAttr, sys *SysProcAttr) {
	// Nothing to do for Freya.
}

// Getrlimit gets resource limits
func Getrlimit(resource int, rlim *Rlimit) (err error) {
	_, _, e1 := RawSyscall(SYS_GETRLIMIT, uintptr(resource), uintptr(unsafe.Pointer(rlim)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// setrlimit sets resource limits (called by Setrlimit in rlimit.go)
func setrlimit(resource int, rlim *Rlimit) (err error) {
	_, _, e1 := RawSyscall(SYS_SETRLIMIT, uintptr(resource), uintptr(unsafe.Pointer(rlim)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// fcntl wrapper
func fcntl(fd int, cmd int, arg int) (val int, err error) {
	r0, _, e1 := Syscall(SYS_FCNTL, uintptr(fd), uintptr(cmd), uintptr(arg))
	val = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

func fcntlPtr(fd int, cmd int, arg unsafe.Pointer) (val int, err error) {
	r0, _, e1 := Syscall(SYS_FCNTL, uintptr(fd), uintptr(cmd), uintptr(arg))
	val = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Low-level *at syscall wrappers that take int dirfd parameter

func mkdirat(dirfd int, path string, mode uint32) (err error) {
	p, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	_, _, e1 := Syscall(SYS_MKDIRAT, uintptr(dirfd), uintptr(unsafe.Pointer(p)), uintptr(mode))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

func unlinkat(dirfd int, path string, flags int) (err error) {
	p, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	_, _, e1 := Syscall(SYS_UNLINKAT, uintptr(dirfd), uintptr(unsafe.Pointer(p)), uintptr(flags))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

func renameat(olddirfd int, oldpath string, newdirfd int, newpath string) (err error) {
	p1, err := BytePtrFromString(oldpath)
	if err != nil {
		return err
	}
	p2, err := BytePtrFromString(newpath)
	if err != nil {
		return err
	}
	_, _, e1 := Syscall6(SYS_RENAMEAT, uintptr(olddirfd), uintptr(unsafe.Pointer(p1)), uintptr(newdirfd), uintptr(unsafe.Pointer(p2)), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

func linkat(olddirfd int, oldpath string, newdirfd int, newpath string, flags int) (err error) {
	p1, err := BytePtrFromString(oldpath)
	if err != nil {
		return err
	}
	p2, err := BytePtrFromString(newpath)
	if err != nil {
		return err
	}
	_, _, e1 := Syscall6(SYS_LINKAT, uintptr(olddirfd), uintptr(unsafe.Pointer(p1)), uintptr(newdirfd), uintptr(unsafe.Pointer(p2)), uintptr(flags), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

func symlinkat(oldpath string, newdirfd int, newpath string) (err error) {
	p1, err := BytePtrFromString(oldpath)
	if err != nil {
		return err
	}
	p2, err := BytePtrFromString(newpath)
	if err != nil {
		return err
	}
	_, _, e1 := Syscall(SYS_SYMLINKAT, uintptr(unsafe.Pointer(p1)), uintptr(newdirfd), uintptr(unsafe.Pointer(p2)))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

func readlinkat(dirfd int, path string, buf []byte) (n int, err error) {
	p, err := BytePtrFromString(path)
	if err != nil {
		return 0, err
	}
	var _p0 unsafe.Pointer
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
	}
	r0, _, e1 := Syscall6(SYS_READLINKAT, uintptr(dirfd), uintptr(unsafe.Pointer(p)), uintptr(_p0), uintptr(len(buf)), 0, 0)
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

func fchmodat(dirfd int, path string, mode uint32, flags int) (err error) {
	p, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	_, _, e1 := Syscall6(SYS_FCHMODAT, uintptr(dirfd), uintptr(unsafe.Pointer(p)), uintptr(mode), uintptr(flags), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

func fchownat(dirfd int, path string, uid int, gid int, flags int) (err error) {
	p, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	_, _, e1 := Syscall6(SYS_FCHOWNAT, uintptr(dirfd), uintptr(unsafe.Pointer(p)), uintptr(uid), uintptr(gid), uintptr(flags), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

func faccessat(dirfd int, path string, mode uint32) (err error) {
	p, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	_, _, e1 := Syscall(SYS_FACCESSAT, uintptr(dirfd), uintptr(unsafe.Pointer(p)), uintptr(mode))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Faccessat checks file accessibility with flags support
func Faccessat(dirfd int, path string, mode uint32, flags int) (err error) {
	if flags == 0 {
		return faccessat(dirfd, path, mode)
	}
	// Freya doesn't support faccessat2 yet, so we ignore AT_EACCESS
	// and AT_SYMLINK_NOFOLLOW for now
	return faccessat(dirfd, path, mode)
}

func openat(dirfd int, path string, flags int, mode uint32) (fd int, err error) {
	p, err := BytePtrFromString(path)
	if err != nil {
		return 0, err
	}
	r0, _, e1 := Syscall6(SYS_OPENAT, uintptr(dirfd), uintptr(unsafe.Pointer(p)), uintptr(flags), uintptr(mode), 0, 0)
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

func fstatat(dirfd int, path string, stat *Stat_t, flags int) (err error) {
	p, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	_, _, e1 := Syscall6(SYS_FSTATAT, uintptr(dirfd), uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(stat)), uintptr(flags), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Public file system functions

// Fstat stats a file by fd
func Fstat(fd int, stat *Stat_t) (err error) {
	_, _, e1 := Syscall(SYS_FSTAT, uintptr(fd), uintptr(unsafe.Pointer(stat)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Stat stats a file by path
func Stat(path string, stat *Stat_t) (err error) {
	return fstatat(_AT_FDCWD, path, stat, 0)
}

// Lstat stats a symlink
func Lstat(path string, stat *Stat_t) (err error) {
	return fstatat(_AT_FDCWD, path, stat, AT_SYMLINK_NOFOLLOW)
}

// Open opens a file
func Open(path string, mode int, perm uint32) (fd int, err error) {
	return openat(_AT_FDCWD, path, mode, perm)
}

// Pipe creates a pipe
func Pipe(p []int) (err error) {
	if len(p) != 2 {
		return EINVAL
	}
	var pp [2]int32
	_, _, e1 := RawSyscall(SYS_PIPE2, uintptr(unsafe.Pointer(&pp)), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	p[0] = int(pp[0])
	p[1] = int(pp[1])
	return
}

func Pipe2(p []int, flags int) (err error) {
	if len(p) != 2 {
		return EINVAL
	}
	var pp [2]int32
	_, _, e1 := RawSyscall(SYS_PIPE2, uintptr(unsafe.Pointer(&pp)), uintptr(flags), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	p[0] = int(pp[0])
	p[1] = int(pp[1])
	return
}

// Dup duplicates a file descriptor
func Dup(oldfd int) (fd int, err error) {
	r1, _, e1 := Syscall(SYS_DUP, uintptr(oldfd), 0, 0)
	fd = int(r1)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

func Dup2(oldfd, newfd int) (err error) {
	_, _, e1 := Syscall(SYS_DUP2, uintptr(oldfd), uintptr(newfd), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

func Dup3(oldfd, newfd, flags int) (err error) {
	_, _, e1 := Syscall(SYS_DUP3, uintptr(oldfd), uintptr(newfd), uintptr(flags))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Mkdir creates a directory
func Mkdir(path string, mode uint32) (err error) {
	return mkdirat(_AT_FDCWD, path, mode)
}

// Rmdir removes a directory
func Rmdir(path string) (err error) {
	return unlinkat(_AT_FDCWD, path, _AT_REMOVEDIR)
}

// Unlink removes a file
func Unlink(path string) (err error) {
	return unlinkat(_AT_FDCWD, path, 0)
}

// Rename renames a file
func Rename(from, to string) (err error) {
	return renameat(_AT_FDCWD, from, _AT_FDCWD, to)
}

// Link creates a hard link
func Link(oldpath, newpath string) (err error) {
	return linkat(_AT_FDCWD, oldpath, _AT_FDCWD, newpath, 0)
}

// Symlink creates a symbolic link
func Symlink(oldpath, newpath string) (err error) {
	return symlinkat(oldpath, _AT_FDCWD, newpath)
}

// Readlink reads a symbolic link
func Readlink(path string, buf []byte) (n int, err error) {
	return readlinkat(_AT_FDCWD, path, buf)
}

// Chmod changes file mode
func Chmod(path string, mode uint32) (err error) {
	return fchmodat(_AT_FDCWD, path, mode, 0)
}

// Chown changes file owner
func Chown(path string, uid, gid int) (err error) {
	return fchownat(_AT_FDCWD, path, uid, gid, 0)
}

// Lchown changes symlink owner
func Lchown(path string, uid, gid int) (err error) {
	return fchownat(_AT_FDCWD, path, uid, gid, AT_SYMLINK_NOFOLLOW)
}

// Fchown changes file owner by fd
func Fchown(fd int, uid, gid int) (err error) {
	_, _, e1 := Syscall(SYS_FCHOWN, uintptr(fd), uintptr(uid), uintptr(gid))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Fchmod changes file mode by fd
func Fchmod(fd int, mode uint32) (err error) {
	_, _, e1 := Syscall(SYS_FCHMOD, uintptr(fd), uintptr(mode), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Fchdir changes to dir by fd
func Fchdir(fd int) (err error) {
	_, _, e1 := Syscall(SYS_CHDIR, uintptr(fd), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Truncate truncates a file
func Truncate(path string, length int64) (err error) {
	p, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	_, _, e1 := Syscall(SYS_TRUNCATE, uintptr(unsafe.Pointer(p)), uintptr(length), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Ftruncate truncates a file by fd
func Ftruncate(fd int, length int64) (err error) {
	_, _, e1 := Syscall(SYS_FTRUNCATE, uintptr(fd), uintptr(length), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Fsync syncs a file
func Fsync(fd int) (err error) {
	_, _, e1 := Syscall(SYS_FSYNC, uintptr(fd), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Fdatasync is an alias for Fsync
func Fdatasync(fd int) (err error) {
	return Fsync(fd)
}

// Sync syncs all filesystems
func Sync() {
	Syscall(SYS_SYNC, 0, 0, 0)
}

// Chdir changes current directory
func Chdir(path string) (err error) {
	p, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	_, _, e1 := Syscall(SYS_CHDIR, uintptr(unsafe.Pointer(p)), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Getcwd gets current directory
func Getcwd(buf []byte) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
	}
	r0, _, e1 := Syscall(SYS_GETCWD, uintptr(_p0), uintptr(len(buf)), 0)
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// ImplementsGetwd indicates whether Getwd is implemented
const ImplementsGetwd = true

// Getwd returns the current working directory
func Getwd() (wd string, err error) {
	var buf [PathMax]byte
	n, err := Getcwd(buf[0:])
	if err != nil {
		return "", err
	}
	// Getcwd returns the number of bytes written to buf, including the NUL.
	if n < 1 || n > len(buf) || buf[n-1] != 0 {
		return "", EINVAL
	}
	return string(buf[:n-1]), nil
}

// utimensat changes file timestamps with nanosecond precision
func utimensat(dirfd int, path string, times *[2]Timespec, flag int) (err error) {
	p, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	_, _, e1 := Syscall6(SYS_UTIMENSAT, uintptr(dirfd), uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(times)), uintptr(flag), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// UtimesNano changes file timestamps with nanosecond precision
func UtimesNano(path string, ts []Timespec) (err error) {
	if len(ts) != 2 {
		return EINVAL
	}
	return utimensat(_AT_FDCWD, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), 0)
}

// Umask sets file creation mask
func Umask(mask int) (oldmask int) {
	r1, _, _ := RawSyscall(SYS_UMASK, uintptr(mask), 0, 0)
	return int(r1)
}

// Getdents reads directory entries
func Getdents(fd int, buf []byte) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
	}
	r0, _, e1 := Syscall(SYS_GETDENTS64, uintptr(fd), uintptr(_p0), uintptr(len(buf)))
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// ReadDirent reads directory entries from a directory file descriptor
func ReadDirent(fd int, buf []byte) (n int, err error) {
	return Getdents(fd, buf)
}

// Uname gets system name
func Uname(buf *Utsname) (err error) {
	_, _, e1 := RawSyscall(SYS_UNAME, uintptr(unsafe.Pointer(buf)), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Getpgid gets process group
func Getpgid(pid int) (pgid int, err error) {
	r1, _, e1 := RawSyscall(SYS_GETPGID, uintptr(pid), 0, 0)
	pgid = int(r1)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Getsid gets session ID
func Getsid(pid int) (sid int, err error) {
	r1, _, e1 := RawSyscall(SYS_GETSID, uintptr(pid), 0, 0)
	sid = int(r1)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Access checks file access
func Access(path string, mode uint32) (err error) {
	return faccessat(_AT_FDCWD, path, mode)
}

// mmap is the low-level mmap implementation required by //go:linkname
func mmap(addr uintptr, length uintptr, prot int, flags int, fd int, offset int64) (xaddr uintptr, err error) {
	r0, _, e1 := Syscall6(SYS_MMAP, addr, length, uintptr(prot), uintptr(flags), uintptr(fd), uintptr(offset))
	xaddr = r0
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Mmap maps files or devices into memory
func Mmap(fd int, offset int64, length, prot, flags int) (data []byte, err error) {
	addr, err := mmap(0, uintptr(length), prot, flags, fd, offset)
	if err != nil {
		return nil, err
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(addr)), length), nil
}

// Munmap unmaps memory
func Munmap(b []byte) (err error) {
	_, _, e1 := Syscall(SYS_MUNMAP, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Mprotect changes memory protection
func Mprotect(b []byte, prot int) (err error) {
	_, _, e1 := Syscall(SYS_MPROTECT, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)), uintptr(prot))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Gettimeofday gets current time
func Gettimeofday(tv *Timeval) (err error) {
	_, _, e1 := RawSyscall(SYS_GETTIMEOFDAY, uintptr(unsafe.Pointer(tv)), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Nanosleep sleeps for a duration
func Nanosleep(time *Timespec, leftover *Timespec) (err error) {
	_, _, e1 := Syscall(SYS_NANOSLEEP, uintptr(unsafe.Pointer(time)), uintptr(unsafe.Pointer(leftover)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Ioctl performs an ioctl
func Ioctl(fd int, req uint, arg uintptr) (err error) {
	_, _, e1 := Syscall(SYS_IOCTL, uintptr(fd), uintptr(req), arg)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// ioctlPtr is like Ioctl but passes a pointer
func ioctlPtr(fd int, req uint, arg unsafe.Pointer) (err error) {
	_, _, e1 := Syscall(SYS_IOCTL, uintptr(fd), uintptr(req), uintptr(arg))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Close closes a file descriptor
func Close(fd int) (err error) {
	_, _, e1 := Syscall(SYS_CLOSE, uintptr(fd), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// Socket address conversion

func anyToSockaddr(rsa *RawSockaddrAny) (Sockaddr, error) {
	switch rsa.Addr.Family {
	case AF_UNIX:
		pp := (*RawSockaddrUnix)(unsafe.Pointer(rsa))
		sa := new(SockaddrUnix)
		if pp.Path[0] == 0 {
			// "Abstract" Unix domain socket.
			pp.Path[0] = '@'
		}
		n := 0
		for n < len(pp.Path) && pp.Path[n] != 0 {
			n++
		}
		sa.Name = string(unsafe.Slice((*byte)(unsafe.Pointer(&pp.Path[0])), n))
		return sa, nil

	case AF_INET:
		pp := (*RawSockaddrInet4)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet4)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		sa.Addr = pp.Addr
		return sa, nil

	case AF_INET6:
		pp := (*RawSockaddrInet6)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet6)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		sa.ZoneId = pp.Scope_id
		sa.Addr = pp.Addr
		return sa, nil
	}
	return nil, EAFNOSUPPORT
}

func Accept(fd int) (nfd int, sa Sockaddr, err error) {
	return Accept4(fd, 0)
}

func Accept4(fd int, flags int) (nfd int, sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	nfd, err = accept4(fd, &rsa, &len, flags)
	if err != nil {
		return
	}
	if len > SizeofSockaddrAny {
		panic("RawSockaddrAny too small")
	}
	sa, err = anyToSockaddr(&rsa)
	if err != nil {
		Close(nfd)
		nfd = 0
	}
	return
}

func recvmsgRaw(fd int, p, oob []byte, flags int, rsa *RawSockaddrAny) (n, oobn int, recvflags int, err error) {
	var msg Msghdr
	msg.Name = (*byte)(unsafe.Pointer(rsa))
	msg.Namelen = uint32(SizeofSockaddrAny)
	var iov Iovec
	if len(p) > 0 {
		iov.Base = &p[0]
		iov.SetLen(len(p))
	}
	var dummy byte
	if len(oob) > 0 {
		if len(p) == 0 {
			var sockType int
			sockType, err = GetsockoptInt(fd, SOL_SOCKET, SO_TYPE)
			if err != nil {
				return
			}
			// receive at least one normal byte
			if sockType != SOCK_DGRAM {
				iov.Base = &dummy
				iov.SetLen(1)
			}
		}
		msg.Control = &oob[0]
		msg.SetControllen(len(oob))
	}
	msg.Iov = &iov
	msg.Iovlen = 1
	if n, err = recvmsg(fd, &msg, flags); err != nil {
		return
	}
	oobn = int(msg.Controllen)
	recvflags = int(msg.Flags)
	return
}

func sendmsgN(fd int, p, oob []byte, ptr unsafe.Pointer, salen _Socklen, flags int) (n int, err error) {
	var msg Msghdr
	msg.Name = (*byte)(ptr)
	msg.Namelen = uint32(salen)
	var iov Iovec
	if len(p) > 0 {
		iov.Base = &p[0]
		iov.SetLen(len(p))
	}
	var dummy byte
	if len(oob) > 0 {
		if len(p) == 0 {
			var sockType int
			sockType, err = GetsockoptInt(fd, SOL_SOCKET, SO_TYPE)
			if err != nil {
				return 0, err
			}
			// send at least one normal byte
			if sockType != SOCK_DGRAM {
				iov.Base = &dummy
				iov.SetLen(1)
			}
		}
		msg.Control = &oob[0]
		msg.SetControllen(len(oob))
	}
	msg.Iov = &iov
	msg.Iovlen = 1
	if n, err = sendmsg(fd, &msg, flags); err != nil {
		return 0, err
	}
	if len(oob) > 0 && len(p) == 0 {
		n = 0
	}
	return n, nil
}

// Sockaddr methods

func (sa *SockaddrInet4) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_INET
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	sa.raw.Addr = sa.Addr
	return unsafe.Pointer(&sa.raw), SizeofSockaddrInet4, nil
}

func (sa *SockaddrInet6) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_INET6
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	sa.raw.Flowinfo = 0
	sa.raw.Scope_id = sa.ZoneId
	sa.raw.Addr = sa.Addr
	return unsafe.Pointer(&sa.raw), SizeofSockaddrInet6, nil
}

func (sa *SockaddrUnix) sockaddr() (unsafe.Pointer, _Socklen, error) {
	name := sa.Name
	n := len(name)
	if n > len(sa.raw.Path) {
		return nil, 0, EINVAL
	}
	isAbstract := n > 0 && (name[0] == '@' || name[0] == '\x00')
	if n == len(sa.raw.Path) && !isAbstract {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_UNIX
	for i := 0; i < n; i++ {
		sa.raw.Path[i] = int8(name[i])
	}
	// Length is family + name (+ NUL if non-abstract).
	sl := _Socklen(2)
	if isAbstract {
		sl += _Socklen(n)
		sa.raw.Path[0] = 0
	} else {
		sl += _Socklen(n + 1)
	}
	return unsafe.Pointer(&sa.raw), sl, nil
}

// sendfile - Freya doesn't have a sendfile syscall, so implement using read/write
func sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) {
	// If offset is provided, seek to it first
	if offset != nil {
		_, err = Seek(infd, *offset, 0) // SEEK_SET
		if err != nil {
			return 0, err
		}
	}

	// Use a buffer to copy data
	buf := make([]byte, 32*1024) // 32KB buffer
	remaining := count
	for remaining > 0 {
		toRead := len(buf)
		if toRead > remaining {
			toRead = remaining
		}
		n, readErr := read(infd, buf[:toRead])
		if n > 0 {
			nw, writeErr := write(outfd, buf[:n])
			written += nw
			if writeErr != nil {
				err = writeErr
				break
			}
			if nw < n {
				err = EAGAIN
				break
			}
		}
		if readErr != nil {
			if readErr != EAGAIN && readErr != EINTR {
				err = readErr
			}
			break
		}
		if n == 0 {
			// EOF
			break
		}
		remaining -= n
	}

	// Update offset if provided
	if offset != nil {
		*offset += int64(written)
	}
	return
}

// Getgroups returns the supplementary group IDs of the calling process
func Getgroups() (gids []int, err error) {
	n, err := getgroups(0, nil)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}

	// Sanity check group count
	if n < 0 || n > 1<<20 {
		return nil, EINVAL
	}

	a := make([]_Gid_t, n)
	n, err = getgroups(n, &a[0])
	if err != nil {
		return nil, err
	}
	gids = make([]int, n)
	for i, v := range a[0:n] {
		gids[i] = int(v)
	}
	return
}

// Getrusage returns resource usage
func Getrusage(who int, rusage *Rusage) (err error) {
	_, _, e1 := RawSyscall(SYS_GETRUSAGE, uintptr(who), uintptr(unsafe.Pointer(rusage)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// RUSAGE constants
const (
	RUSAGE_SELF     = 0
	RUSAGE_CHILDREN = -1
	RUSAGE_THREAD   = 1
)

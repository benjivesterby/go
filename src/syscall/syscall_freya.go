// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Freya system calls.

//go:build freya

package syscall

import (
	"runtime"
	"sync"
	"unsafe"
)

// Freya syscall numbers.
const (
	// IPC
	SYS_SEND    = 1
	SYS_RECEIVE = 2
	SYS_CALL    = 3
	SYS_REPLY   = 4
	SYS_NOTIFY  = 5
	SYS_WAIT    = 6

	// Process management
	SYS_PROCESS_EXIT   = 10
	SYS_PROCESS_WAIT   = 12
	SYS_GETPID         = 13
	SYS_GETTID         = 14
	SYS_GETPPID        = 35
	SYS_PROCESS_CREATE = 25
	SYS_PROCESS_EXEC   = 26
	SYS_PROCESS_KILL   = 27

	// Thread management
	SYS_THREAD_YIELD  = 15
	SYS_THREAD_EXIT   = 16
	SYS_THREAD_CREATE = 17
	SYS_THREAD_SLEEP  = 18

	// Time
	SYS_GETTIME       = 19
	SYS_NANOSLEEP     = 36
	SYS_CLOCK_GETTIME = 37
	SYS_FUTEX         = 44

	// Memory management
	SYS_MEMORY_MAP     = 20
	SYS_MEMORY_UNMAP   = 21
	SYS_MEMORY_PROTECT = 22

	// User/Group
	SYS_GETUID = 28
	SYS_SETUID = 29
	SYS_GETGID = 33
	SYS_SETGID = 34

	// File descriptors
	SYS_FILE_DUP    = 54
	SYS_PIPE_CREATE = 55
	SYS_FILE_CHDIR  = 56
	SYS_FILE_GETCWD = 57
	SYS_FILE_CHMOD  = 58
	SYS_FILE_FCNTL  = 80
	SYS_FILE_IOCTL  = 81
	SYS_FILE_DUP2   = 86
	SYS_FILE_POLL   = 87
	SYS_FILE_READV  = 88
	SYS_FILE_WRITEV = 89
	SYS_FILE_PREAD  = 90
	SYS_FILE_PWRITE = 91

	// Signals
	SYS_RT_SIGACTION   = 94
	SYS_RT_SIGPROCMASK = 95
	SYS_RT_SIGRETURN   = 97

	// File operations
	SYS_FILE_OPEN     = 100
	SYS_FILE_READ     = 101
	SYS_FILE_WRITE    = 102
	SYS_FILE_CLOSE    = 103
	SYS_FILE_READDIR  = 104
	SYS_FILE_SEEK     = 105
	SYS_FILE_STAT     = 106
	SYS_FILE_FSTAT    = 107
	SYS_FILE_MKDIR    = 108
	SYS_FILE_UNLINK   = 109
	SYS_FILE_TRUNCATE = 110
	SYS_FILE_RMDIR    = 111
	SYS_FILE_RENAME   = 112
	SYS_FILE_FSYNC    = 113

	// Links
	SYS_FILE_LINK     = 67
	SYS_FILE_SYMLINK  = 68
	SYS_FILE_READLINK = 69

	// Sockets
	SYS_SOCKET_CREATE  = 70
	SYS_SOCKET_BIND    = 71
	SYS_SOCKET_CONNECT = 72
	SYS_SOCKET_SEND    = 73
	SYS_SOCKET_RECV    = 74
	SYS_SOCKET_CLOSE   = 77
	SYS_SOCKET_LISTEN  = 78
	SYS_SOCKET_ACCEPT  = 79
)

// File open flags.
const (
	O_RDONLY   = 0x0
	O_WRONLY   = 0x1
	O_RDWR     = 0x2
	O_CREAT    = 0x40
	O_EXCL     = 0x80
	O_TRUNC    = 0x200
	O_APPEND   = 0x400
	O_NONBLOCK = 0x800
	O_CLOEXEC  = 0x80000
)

// Seek whence values.
const (
	SEEK_SET = 0
	SEEK_CUR = 1
	SEEK_END = 2
)

// Memory protection flags.
const (
	PROT_NONE  = 0x0
	PROT_READ  = 0x1
	PROT_WRITE = 0x2
	PROT_EXEC  = 0x4
)

// Memory mapping flags.
const (
	MAP_SHARED    = 0x1
	MAP_PRIVATE   = 0x2
	MAP_ANONYMOUS = 0x20
	MAP_ANON      = MAP_ANONYMOUS
)

// File mode bits.
const (
	S_IFMT   = 0xf000
	S_IFIFO  = 0x1000
	S_IFCHR  = 0x2000
	S_IFDIR  = 0x4000
	S_IFBLK  = 0x6000
	S_IFREG  = 0x8000
	S_IFLNK  = 0xa000
	S_IFSOCK = 0xc000
	S_ISUID  = 0x800
	S_ISGID  = 0x400
	S_ISVTX  = 0x200
)

// Wait flags.
const (
	WNOHANG    = 0x1
	WUNTRACED  = 0x2
	WCONTINUED = 0x8
)

// Socket address families.
const (
	AF_UNSPEC = 0x0
	AF_UNIX   = 0x1
	AF_INET   = 0x2
	AF_INET6  = 0xa
)

// Path limits.
const (
	PathMax = 0x1000
)

// Size constants.
const (
	sizeofPtr      = 0x8
	sizeofShort    = 0x2
	sizeofInt      = 0x4
	sizeofLong     = 0x8
	sizeofLongLong = 0x8
)

// C type aliases.
type (
	_C_short     int16
	_C_int       int32
	_C_long      int64
	_C_long_long int64
)

// Errno values (POSIX-compatible).
const (
	EPERM           Errno = 0x1
	ENOENT          Errno = 0x2
	ESRCH           Errno = 0x3
	EINTR           Errno = 0x4
	EIO             Errno = 0x5
	ENXIO           Errno = 0x6
	E2BIG           Errno = 0x7
	ENOEXEC         Errno = 0x8
	EBADF           Errno = 0x9
	ECHILD          Errno = 0xa
	EAGAIN          Errno = 0xb
	ENOMEM          Errno = 0xc
	EACCES          Errno = 0xd
	EFAULT          Errno = 0xe
	EBUSY           Errno = 0x10
	EEXIST          Errno = 0x11
	EXDEV           Errno = 0x12
	ENODEV          Errno = 0x13
	ENOTDIR         Errno = 0x14
	EISDIR          Errno = 0x15
	EINVAL          Errno = 0x16
	ENFILE          Errno = 0x17
	EMFILE          Errno = 0x18
	ENOTTY          Errno = 0x19
	EFBIG           Errno = 0x1b
	ENOSPC          Errno = 0x1c
	ESPIPE          Errno = 0x1d
	EROFS           Errno = 0x1e
	EPIPE           Errno = 0x20
	ERANGE          Errno = 0x22
	ENAMETOOLONG    Errno = 0x24
	ENOSYS          Errno = 0x26
	ENOTEMPTY       Errno = 0x27
	ELOOP           Errno = 0x28
	ENOTSUP         Errno = 0x5f
	EOPNOTSUPP      Errno = 0x5f
	EAFNOSUPPORT    Errno = 0x61
	EADDRINUSE      Errno = 0x62
	ECONNREFUSED    Errno = 0x6f
	ETIMEDOUT       Errno = 0x6e
	EWOULDBLOCK     Errno = EAGAIN
	ECONNRESET      Errno = 0x68
	ECONNABORTED    Errno = 0x67
	ENETUNREACH     Errno = 0x65
	EHOSTUNREACH    Errno = 0x71
	EALREADY        Errno = 0x72
	EINPROGRESS     Errno = 0x73
	EDESTADDRREQ    Errno = 0x59
	EMSGSIZE        Errno = 0x5a
	ENOBUFS         Errno = 0x69
	EISCONN         Errno = 0x6a
	ENOTCONN        Errno = 0x6b
	ENOTSOCK        Errno = 0x58
	EPROTONOSUPPORT Errno = 0x5d
)

// Signal numbers (POSIX-compatible).
const (
	SIGHUP    = Signal(0x1)
	SIGINT    = Signal(0x2)
	SIGQUIT   = Signal(0x3)
	SIGILL    = Signal(0x4)
	SIGTRAP   = Signal(0x5)
	SIGABRT   = Signal(0x6)
	SIGBUS    = Signal(0x7)
	SIGFPE    = Signal(0x8)
	SIGKILL   = Signal(0x9)
	SIGUSR1   = Signal(0xa)
	SIGSEGV   = Signal(0xb)
	SIGUSR2   = Signal(0xc)
	SIGPIPE   = Signal(0xd)
	SIGALRM   = Signal(0xe)
	SIGTERM   = Signal(0xf)
	SIGCHLD   = Signal(0x11)
	SIGCONT   = Signal(0x12)
	SIGSTOP   = Signal(0x13)
	SIGTSTP   = Signal(0x14)
	SIGTTIN   = Signal(0x15)
	SIGTTOU   = Signal(0x16)
	SIGURG    = Signal(0x17)
	SIGXCPU   = Signal(0x18)
	SIGXFSZ   = Signal(0x19)
	SIGVTALRM = Signal(0x1a)
	SIGPROF   = Signal(0x1b)
	SIGWINCH  = Signal(0x1c)
	SIGIO     = Signal(0x1d)
	SIGSYS    = Signal(0x1f)
)

// Error strings table.
var errors = [...]string{
	1:   "operation not permitted",
	2:   "no such file or directory",
	3:   "no such process",
	4:   "interrupted system call",
	5:   "input/output error",
	6:   "no such device or address",
	7:   "argument list too long",
	8:   "exec format error",
	9:   "bad file descriptor",
	10:  "no child processes",
	11:  "resource temporarily unavailable",
	12:  "cannot allocate memory",
	13:  "permission denied",
	14:  "bad address",
	16:  "device or resource busy",
	17:  "file exists",
	18:  "invalid cross-device link",
	19:  "no such device",
	20:  "not a directory",
	21:  "is a directory",
	22:  "invalid argument",
	23:  "too many open files in system",
	24:  "too many open files",
	25:  "inappropriate ioctl for device",
	27:  "file too large",
	28:  "no space left on device",
	29:  "illegal seek",
	30:  "read-only file system",
	32:  "broken pipe",
	34:  "numerical result out of range",
	36:  "file name too long",
	38:  "function not implemented",
	39:  "directory not empty",
	40:  "too many levels of symbolic links",
	88:  "socket operation on non-socket",
	89:  "destination address required",
	90:  "message too long",
	93:  "protocol not supported",
	95:  "operation not supported",
	97:  "address family not supported by protocol",
	98:  "address already in use",
	99:  "cannot assign requested address",
	101: "network is unreachable",
	103: "software caused connection abort",
	104: "connection reset by peer",
	105: "no buffer space available",
	106: "transport endpoint is already connected",
	107: "transport endpoint is not connected",
	110: "connection timed out",
	111: "connection refused",
	113: "no route to host",
	114: "operation already in progress",
	115: "operation now in progress",
}

// Signal strings table.
var signals = [...]string{
	1:  "hangup",
	2:  "interrupt",
	3:  "quit",
	4:  "illegal instruction",
	5:  "trace/breakpoint trap",
	6:  "aborted",
	7:  "bus error",
	8:  "floating point exception",
	9:  "killed",
	10: "user defined signal 1",
	11: "segmentation fault",
	12: "user defined signal 2",
	13: "broken pipe",
	14: "alarm clock",
	15: "terminated",
	17: "child exited",
	18: "continued",
	19: "stopped (signal)",
	20: "stopped",
	21: "stopped (tty input)",
	22: "stopped (tty output)",
	23: "urgent I/O condition",
	24: "CPU time limit exceeded",
	25: "file size limit exceeded",
	26: "virtual timer expired",
	27: "profiling timer expired",
	28: "window changed",
	29: "I/O possible",
	31: "bad system call",
}

// Timespec represents a time with nanosecond precision.
type Timespec struct {
	Sec  int64
	Nsec int64
}

// Timeval represents a time with microsecond precision.
type Timeval struct {
	Sec  int64
	Usec int64
}

// Time_t is the type for time values.
type Time_t int64

// Stat_t is the file status structure.
type Stat_t struct {
	Dev     uint64
	Ino     uint64
	Mode    uint32
	Nlink   uint32
	Uid     uint32
	Gid     uint32
	Rdev    uint64
	_       uint64
	Size    int64
	Blksize int32
	_       int32
	Blocks  int64
	Atim    Timespec
	Mtim    Timespec
	Ctim    Timespec
	_       [2]int32
}

// Dirent is the directory entry structure.
type Dirent struct {
	Ino    uint64
	Off    int64
	Reclen uint16
	Type   uint8
	Name   [256]uint8
	_      [5]byte
}

// Rusage is the resource usage structure.
type Rusage struct {
	Utime    Timeval
	Stime    Timeval
	Maxrss   int64
	Ixrss    int64
	Idrss    int64
	Isrss    int64
	Minflt   int64
	Majflt   int64
	Nswap    int64
	Inblock  int64
	Oublock  int64
	Msgsnd   int64
	Msgrcv   int64
	Nsignals int64
	Nvcsw    int64
	Nivcsw   int64
}

// Rlimit is the resource limit structure.
type Rlimit struct {
	Cur uint64
	Max uint64
}

// WaitStatus represents the status of a waited-for process.
type WaitStatus uint32

const (
	mask    = 0x7F
	core    = 0x80
	exited  = 0x00
	stopped = 0x7F
	shift   = 8
)

func (w WaitStatus) Exited() bool    { return w&mask == exited }
func (w WaitStatus) Signaled() bool  { return w&mask != stopped && w&mask != exited }
func (w WaitStatus) Stopped() bool   { return w&0xFF == stopped }
func (w WaitStatus) Continued() bool { return w == 0xFFFF }
func (w WaitStatus) CoreDump() bool  { return w.Signaled() && w&core != 0 }

func (w WaitStatus) ExitStatus() int {
	if !w.Exited() {
		return -1
	}
	return int(w>>shift) & 0xFF
}

func (w WaitStatus) Signal() Signal {
	if !w.Signaled() {
		return -1
	}
	return Signal(w & mask)
}

func (w WaitStatus) StopSignal() Signal {
	if !w.Stopped() {
		return -1
	}
	return Signal(w>>shift) & 0xFF
}

func (w WaitStatus) TrapCause() int {
	if w.StopSignal() != SIGTRAP {
		return -1
	}
	return int(w>>shift) >> 8
}

// Standard file descriptors.
var (
	Stdin  = 0
	Stdout = 1
	Stderr = 2
)

const ImplementsGetwd = true

// For testing: clients can set this flag to force
// creation of IPv6 sockets to return EAFNOSUPPORT.
var SocketDisableIPv6 bool

// Pull in entersyscall/exitsyscall for Syscall/Syscall6.
//
//go:linkname runtime_entersyscall runtime.entersyscall
func runtime_entersyscall()

//go:linkname runtime_exitsyscall runtime.exitsyscall
func runtime_exitsyscall()

// Syscall performs a system call, notifying the scheduler.
//
//go:uintptrkeepalive
//go:nosplit
//go:linkname Syscall
func Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno) {
	runtime_entersyscall()
	r1, r2, err = RawSyscall6(trap, a1, a2, a3, 0, 0, 0)
	runtime_exitsyscall()
	return
}

// Syscall6 performs a system call with 6 arguments, notifying the scheduler.
//
//go:uintptrkeepalive
//go:nosplit
//go:linkname Syscall6
func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno) {
	runtime_entersyscall()
	r1, r2, err = RawSyscall6(trap, a1, a2, a3, a4, a5, a6)
	runtime_exitsyscall()
	return
}

// RawSyscall performs a raw system call without scheduler notification.
//
//go:uintptrkeepalive
//go:nosplit
//go:norace
//go:linkname RawSyscall
func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno) {
	return RawSyscall6(trap, a1, a2, a3, 0, 0, 0)
}

// RawSyscall6 is implemented in assembly.
//
//go:uintptrkeepalive
//go:nosplit
//go:norace
//go:linkname RawSyscall6
func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno)

// Do the interface allocations only once for common
// Errno values.
var (
	errEAGAIN error = EAGAIN
	errEINVAL error = EINVAL
	errENOENT error = ENOENT
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e Errno) error {
	switch e {
	case 0:
		return nil
	case EAGAIN:
		return errEAGAIN
	case EINVAL:
		return errEINVAL
	case ENOENT:
		return errENOENT
	}
	return e
}

/*
 * Wrapped system calls
 */

func Open(path string, mode int, perm uint32) (fd int, err error) {
	p, err := BytePtrFromString(path)
	if err != nil {
		return -1, err
	}
	r, _, e := Syscall(SYS_FILE_OPEN, uintptr(unsafe.Pointer(p)), uintptr(mode), uintptr(perm))
	if e != 0 {
		return -1, errnoErr(e)
	}
	return int(r), nil
}

func read(fd int, p []byte) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	}
	r, _, e := Syscall(SYS_FILE_READ, uintptr(fd), uintptr(_p0), uintptr(len(p)))
	if e != 0 {
		return 0, errnoErr(e)
	}
	return int(r), nil
}

func write(fd int, p []byte) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	}
	r, _, e := Syscall(SYS_FILE_WRITE, uintptr(fd), uintptr(_p0), uintptr(len(p)))
	if e != 0 {
		return 0, errnoErr(e)
	}
	return int(r), nil
}

func Close(fd int) (err error) {
	_, _, e := Syscall(SYS_FILE_CLOSE, uintptr(fd), 0, 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Seek(fd int, offset int64, whence int) (newoffset int64, err error) {
	r, _, e := Syscall(SYS_FILE_SEEK, uintptr(fd), uintptr(offset), uintptr(whence))
	if e != 0 {
		return -1, errnoErr(e)
	}
	return int64(r), nil
}

func Stat(path string, stat *Stat_t) (err error) {
	p, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	_, _, e := Syscall(SYS_FILE_STAT, uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(stat)), 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Fstat(fd int, stat *Stat_t) (err error) {
	_, _, e := Syscall(SYS_FILE_FSTAT, uintptr(fd), uintptr(unsafe.Pointer(stat)), 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Lstat(path string, stat *Stat_t) (err error) {
	// Freya uses FILE_STAT; symlink-aware stat may be added later.
	return Stat(path, stat)
}

func Mkdir(path string, mode uint32) (err error) {
	p, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	_, _, e := Syscall(SYS_FILE_MKDIR, uintptr(unsafe.Pointer(p)), uintptr(mode), 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Unlink(path string) error {
	p, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	_, _, e := Syscall(SYS_FILE_UNLINK, uintptr(unsafe.Pointer(p)), 0, 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Rmdir(path string) error {
	p, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	_, _, e := Syscall(SYS_FILE_RMDIR, uintptr(unsafe.Pointer(p)), 0, 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Rename(oldpath string, newpath string) (err error) {
	oldp, err := BytePtrFromString(oldpath)
	if err != nil {
		return err
	}
	newp, err := BytePtrFromString(newpath)
	if err != nil {
		return err
	}
	_, _, e := Syscall(SYS_FILE_RENAME, uintptr(unsafe.Pointer(oldp)), uintptr(unsafe.Pointer(newp)), 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Truncate(path string, length int64) (err error) {
	p, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	_, _, e := Syscall(SYS_FILE_TRUNCATE, uintptr(unsafe.Pointer(p)), uintptr(length), 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Fsync(fd int) (err error) {
	_, _, e := Syscall(SYS_FILE_FSYNC, uintptr(fd), 0, 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Link(oldpath string, newpath string) (err error) {
	oldp, err := BytePtrFromString(oldpath)
	if err != nil {
		return err
	}
	newp, err := BytePtrFromString(newpath)
	if err != nil {
		return err
	}
	_, _, e := Syscall(SYS_FILE_LINK, uintptr(unsafe.Pointer(oldp)), uintptr(unsafe.Pointer(newp)), 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Symlink(oldpath string, newpath string) (err error) {
	oldp, err := BytePtrFromString(oldpath)
	if err != nil {
		return err
	}
	newp, err := BytePtrFromString(newpath)
	if err != nil {
		return err
	}
	_, _, e := Syscall(SYS_FILE_SYMLINK, uintptr(unsafe.Pointer(oldp)), uintptr(unsafe.Pointer(newp)), 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Readlink(path string, buf []byte) (n int, err error) {
	p, err := BytePtrFromString(path)
	if err != nil {
		return 0, err
	}
	var _p0 unsafe.Pointer
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
	}
	r, _, e := Syscall(SYS_FILE_READLINK, uintptr(unsafe.Pointer(p)), uintptr(_p0), uintptr(len(buf)))
	if e != 0 {
		return 0, errnoErr(e)
	}
	return int(r), nil
}

func Exit(code int) {
	RawSyscall(SYS_PROCESS_EXIT, uintptr(code), 0, 0)
}

func Getpid() (pid int) {
	r, _, _ := RawSyscall(SYS_GETPID, 0, 0, 0)
	return int(r)
}

func Getppid() (ppid int) {
	r, _, _ := RawSyscall(SYS_GETPPID, 0, 0, 0)
	return int(r)
}

func Gettid() (tid int) {
	r, _, _ := RawSyscall(SYS_GETTID, 0, 0, 0)
	return int(r)
}

func Kill(pid int, sig Signal) (err error) {
	_, _, e := RawSyscall(SYS_PROCESS_KILL, uintptr(pid), uintptr(sig), 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Getuid() (uid int) {
	r, _, _ := RawSyscall(SYS_GETUID, 0, 0, 0)
	return int(r)
}

func Getgid() (gid int) {
	r, _, _ := RawSyscall(SYS_GETGID, 0, 0, 0)
	return int(r)
}

func Geteuid() (euid int) {
	// Freya does not distinguish effective vs real uid.
	return Getuid()
}

func Getegid() (egid int) {
	// Freya does not distinguish effective vs real gid.
	return Getgid()
}

func Setuid(uid int) (err error) {
	_, _, e := RawSyscall(SYS_SETUID, uintptr(uid), 0, 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Setgid(gid int) (err error) {
	_, _, e := RawSyscall(SYS_SETGID, uintptr(gid), 0, 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Getgroups() (gids []int, err error) {
	return make([]int, 0), nil
}

func Setgroups(gids []int) (err error) {
	return ENOSYS
}

func Pipe(p []int) (err error) {
	if len(p) != 2 {
		return EINVAL
	}
	var pp [2]_C_int
	_, _, e := RawSyscall(SYS_PIPE_CREATE, uintptr(unsafe.Pointer(&pp)), 0, 0)
	if e != 0 {
		return errnoErr(e)
	}
	p[0] = int(pp[0])
	p[1] = int(pp[1])
	return nil
}

func Dup(oldfd int) (fd int, err error) {
	r, _, e := Syscall(SYS_FILE_DUP, uintptr(oldfd), 0, 0)
	if e != 0 {
		return -1, errnoErr(e)
	}
	return int(r), nil
}

func Dup2(oldfd int, newfd int) (err error) {
	_, _, e := Syscall(SYS_FILE_DUP2, uintptr(oldfd), uintptr(newfd), 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Chdir(path string) (err error) {
	p, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	_, _, e := Syscall(SYS_FILE_CHDIR, uintptr(unsafe.Pointer(p)), 0, 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Fchdir(fd int) (err error) {
	_, _, e := Syscall(SYS_FILE_CHDIR, uintptr(fd), 0, 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Getcwd(buf []byte) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
	}
	r, _, e := Syscall(SYS_FILE_GETCWD, uintptr(_p0), uintptr(len(buf)), 0)
	if e != 0 {
		return 0, errnoErr(e)
	}
	return int(r), nil
}

func Getwd() (wd string, err error) {
	var buf [PathMax]byte
	n, err := Getcwd(buf[0:])
	if err != nil {
		return "", err
	}
	if n < 1 || n > len(buf) || buf[n-1] != 0 {
		return "", EINVAL
	}
	if buf[0] != '/' {
		return "", ENOENT
	}
	return string(buf[0 : n-1]), nil
}

func Chmod(path string, mode uint32) (err error) {
	p, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	_, _, e := Syscall(SYS_FILE_CHMOD, uintptr(unsafe.Pointer(p)), uintptr(mode), 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Fchmod(fd int, mode uint32) (err error) {
	_, _, e := Syscall(SYS_FILE_CHMOD, uintptr(fd), uintptr(mode), 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Nanosleep(time *Timespec, leftover *Timespec) (err error) {
	_, _, e := Syscall(SYS_NANOSLEEP, uintptr(unsafe.Pointer(time)), uintptr(unsafe.Pointer(leftover)), 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

func Gettimeofday(tv *Timeval) error {
	var ts Timespec
	_, _, e := RawSyscall(SYS_CLOCK_GETTIME, 0, uintptr(unsafe.Pointer(&ts)), 0)
	if e != 0 {
		return errnoErr(e)
	}
	tv.Sec = ts.Sec
	tv.Usec = ts.Nsec / 1000
	return nil
}

func NsecToTimeval(nsec int64) (tv Timeval) {
	nsec += 999 // round up to microsecond
	tv.Usec = int64(nsec % 1e9 / 1e3)
	tv.Sec = int64(nsec / 1e9)
	return
}

func NsecToTimespec(nsec int64) (ts Timespec) {
	ts.Sec = int64(nsec / 1e9)
	ts.Nsec = int64(nsec % 1e9)
	return
}

func TimespecToNsec(ts Timespec) int64 {
	return ts.Sec*1e9 + ts.Nsec
}

func TimevalToNsec(tv Timeval) int64 {
	return tv.Sec*1e9 + tv.Usec*1000
}

// Wait4 waits for a process to change state.
func Wait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (wpid int, err error) {
	var status _C_int
	r, _, e := Syscall6(SYS_PROCESS_WAIT, uintptr(pid), uintptr(unsafe.Pointer(&status)), uintptr(options), uintptr(unsafe.Pointer(rusage)), 0, 0)
	if e != 0 {
		return -1, errnoErr(e)
	}
	if wstatus != nil {
		*wstatus = WaitStatus(status)
	}
	return int(r), nil
}

// ReadDirent reads directory entries from fd into buf.
func ReadDirent(fd int, buf []byte) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
	}
	r, _, e := Syscall(SYS_FILE_READDIR, uintptr(fd), uintptr(_p0), uintptr(len(buf)))
	if e != 0 {
		return 0, errnoErr(e)
	}
	return int(r), nil
}

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

func Pread(fd int, p []byte, offset int64) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	}
	r, _, e := Syscall6(SYS_FILE_PREAD, uintptr(fd), uintptr(_p0), uintptr(len(p)), uintptr(offset), 0, 0)
	if e != 0 {
		return 0, errnoErr(e)
	}
	return int(r), nil
}

func Pwrite(fd int, p []byte, offset int64) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	}
	r, _, e := Syscall6(SYS_FILE_PWRITE, uintptr(fd), uintptr(_p0), uintptr(len(p)), uintptr(offset), 0, 0)
	if e != 0 {
		return 0, errnoErr(e)
	}
	return int(r), nil
}

func Fcntl(fd int, cmd int, arg int) (val int, err error) {
	r, _, e := Syscall(SYS_FILE_FCNTL, uintptr(fd), uintptr(cmd), uintptr(arg))
	if e != 0 {
		return -1, errnoErr(e)
	}
	return int(r), nil
}

// mmap wraps the Freya memory map syscall.
func mmap(addr uintptr, length uintptr, prot int, flags int, fd int, offset int64) (xaddr uintptr, err error) {
	r, _, e := Syscall6(SYS_MEMORY_MAP, addr, length, uintptr(prot), uintptr(flags), uintptr(fd), uintptr(offset))
	if e != 0 {
		return 0, errnoErr(e)
	}
	return r, nil
}

// munmap wraps the Freya memory unmap syscall.
func munmap(addr uintptr, length uintptr) error {
	_, _, e := Syscall(SYS_MEMORY_UNMAP, addr, length, 0)
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

var mapper = &mmapper{
	active: make(map[*byte][]byte),
	mmap:   mmap,
	munmap: munmap,
}

func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
	return mapper.Mmap(fd, offset, length, prot, flags)
}

func Munmap(b []byte) (err error) {
	return mapper.Munmap(b)
}

func Mprotect(b []byte, prot int) (err error) {
	var _p0 unsafe.Pointer
	if len(b) > 0 {
		_p0 = unsafe.Pointer(&b[0])
	}
	_, _, e := Syscall(SYS_MEMORY_PROTECT, uintptr(_p0), uintptr(len(b)), uintptr(prot))
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}

// Mmap manager, for use by operating system-specific implementations.
type mmapper struct {
	sync.Mutex
	active map[*byte][]byte
	mmap   func(addr, length uintptr, prot, flags, fd int, offset int64) (uintptr, error)
	munmap func(addr uintptr, length uintptr) error
}

func (m *mmapper) Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
	if length <= 0 {
		return nil, EINVAL
	}

	addr, errno := m.mmap(0, uintptr(length), prot, flags, fd, offset)
	if errno != nil {
		return nil, errno
	}

	b := unsafe.Slice((*byte)(unsafe.Pointer(addr)), length)

	p := &b[cap(b)-1]
	m.Lock()
	defer m.Unlock()
	m.active[p] = b
	return b, nil
}

func (m *mmapper) Munmap(data []byte) (err error) {
	if len(data) == 0 || len(data) != cap(data) {
		return EINVAL
	}

	p := &data[cap(data)-1]
	m.Lock()
	defer m.Unlock()
	b := m.active[p]
	if b == nil || &b[0] != &data[0] {
		return EINVAL
	}

	if errno := m.munmap(uintptr(unsafe.Pointer(&b[0])), uintptr(len(b))); errno != nil {
		return errno
	}
	delete(m.active, p)
	return nil
}

// Socket operations.

func Socket(domain, typ, proto int) (fd int, err error) {
	if domain == AF_INET6 && SocketDisableIPv6 {
		return -1, EAFNOSUPPORT
	}
	r, _, e := Syscall(SYS_SOCKET_CREATE, uintptr(domain), uintptr(typ), uintptr(proto))
	if e != 0 {
		return -1, errnoErr(e)
	}
	return int(r), nil
}

// Utility functions to satisfy the runtime's needs.

// Getrlimit is a stub for Freya.
func Getrlimit(resource int, rlim *Rlimit) (err error) {
	return ENOSYS
}

// Setrlimit is a stub for Freya.
func Setrlimit(resource int, rlim *Rlimit) (err error) {
	return ENOSYS
}

// Umask sets the file mode creation mask.
func Umask(mask int) (oldmask int) {
	// Freya may not support umask directly; return 0.
	_ = mask
	return 0
}

var _ = runtime.GOOS // ensure runtime is imported

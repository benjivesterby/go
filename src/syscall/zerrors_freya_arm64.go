// Freya OS error codes and constants for riscv64.

//go:build freya && arm64

package syscall

// Error numbers (POSIX compatible)
const (
	E2BIG           = Errno(0x7)
	EACCES          = Errno(0xd)
	EADDRINUSE      = Errno(0x62)
	EADDRNOTAVAIL   = Errno(0x63)
	EAFNOSUPPORT    = Errno(0x61)
	EAGAIN          = Errno(0xb)
	EALREADY        = Errno(0x72)
	EBADF           = Errno(0x9)
	EBADMSG         = Errno(0x4a)
	EBUSY           = Errno(0x10)
	ECANCELED       = Errno(0x7d)
	ECHILD          = Errno(0xa)
	ECONNABORTED    = Errno(0x67)
	ECONNREFUSED    = Errno(0x6f)
	ECONNRESET      = Errno(0x68)
	EDEADLK         = Errno(0x23)
	EDESTADDRREQ    = Errno(0x59)
	EDOM            = Errno(0x21)
	EDQUOT          = Errno(0x7a)
	EEXIST          = Errno(0x11)
	EFAULT          = Errno(0xe)
	EFBIG           = Errno(0x1b)
	EHOSTDOWN       = Errno(0x70)
	EHOSTUNREACH    = Errno(0x71)
	EIDRM           = Errno(0x2b)
	EILSEQ          = Errno(0x54)
	EINPROGRESS     = Errno(0x73)
	EINTR           = Errno(0x4)
	EINVAL          = Errno(0x16)
	EIO             = Errno(0x5)
	EISCONN         = Errno(0x6a)
	EISDIR          = Errno(0x15)
	ELOOP           = Errno(0x28)
	EMFILE          = Errno(0x18)
	EMLINK          = Errno(0x1f)
	EMSGSIZE        = Errno(0x5a)
	EMULTIHOP       = Errno(0x48)
	ENAMETOOLONG    = Errno(0x24)
	ENETDOWN        = Errno(0x64)
	ENETRESET       = Errno(0x66)
	ENETUNREACH     = Errno(0x65)
	ENFILE          = Errno(0x17)
	ENOBUFS         = Errno(0x69)
	ENODATA         = Errno(0x3d)
	ENODEV          = Errno(0x13)
	ENOENT          = Errno(0x2)
	ENOEXEC         = Errno(0x8)
	ENOLCK          = Errno(0x25)
	ENOLINK         = Errno(0x43)
	ENOMEM          = Errno(0xc)
	ENOMSG          = Errno(0x2a)
	ENOPROTOOPT     = Errno(0x5c)
	ENOSPC          = Errno(0x1c)
	ENOSR           = Errno(0x3f)
	ENOSTR          = Errno(0x3c)
	ENOSYS          = Errno(0x26)
	ENOTBLK         = Errno(0xf)
	ENOTCONN        = Errno(0x6b)
	ENOTDIR         = Errno(0x14)
	ENOTEMPTY       = Errno(0x27)
	ENOTRECOVERABLE = Errno(0x83)
	ENOTSOCK        = Errno(0x58)
	ENOTSUP         = Errno(0x5f)
	ENOTTY          = Errno(0x19)
	ENXIO           = Errno(0x6)
	EOPNOTSUPP      = Errno(0x5f)
	EOVERFLOW       = Errno(0x4b)
	EOWNERDEAD      = Errno(0x82)
	EPERM           = Errno(0x1)
	EPFNOSUPPORT    = Errno(0x60)
	EPIPE           = Errno(0x20)
	EPROTO          = Errno(0x47)
	EPROTONOSUPPORT = Errno(0x5d)
	EPROTOTYPE      = Errno(0x5b)
	ERANGE          = Errno(0x22)
	EROFS           = Errno(0x1e)
	ESPIPE          = Errno(0x1d)
	ESRCH           = Errno(0x3)
	ESTALE          = Errno(0x74)
	ETIME           = Errno(0x3e)
	ETIMEDOUT       = Errno(0x6e)
	ETXTBSY         = Errno(0x1a)
	EWOULDBLOCK     = Errno(0xb)
	EXDEV           = Errno(0x12)
)

// Signals
const (
	SIGABRT   = Signal(0x6)
	SIGALRM   = Signal(0xe)
	SIGBUS    = Signal(0x7)
	SIGCHLD   = Signal(0x11)
	SIGCONT   = Signal(0x12)
	SIGFPE    = Signal(0x8)
	SIGHUP    = Signal(0x1)
	SIGILL    = Signal(0x4)
	SIGINT    = Signal(0x2)
	SIGIO     = Signal(0x1d)
	SIGIOT    = Signal(0x6)
	SIGKILL   = Signal(0x9)
	SIGPIPE   = Signal(0xd)
	SIGPOLL   = Signal(0x1d)
	SIGPROF   = Signal(0x1b)
	SIGPWR    = Signal(0x1e)
	SIGQUIT   = Signal(0x3)
	SIGSEGV   = Signal(0xb)
	SIGSTKFLT = Signal(0x10)
	SIGSTOP   = Signal(0x13)
	SIGSYS    = Signal(0x1f)
	SIGTERM   = Signal(0xf)
	SIGTRAP   = Signal(0x5)
	SIGTSTP   = Signal(0x14)
	SIGTTIN   = Signal(0x15)
	SIGTTOU   = Signal(0x16)
	SIGURG    = Signal(0x17)
	SIGUSR1   = Signal(0xa)
	SIGUSR2   = Signal(0xc)
	SIGVTALRM = Signal(0x1a)
	SIGWINCH  = Signal(0x1c)
	SIGXCPU   = Signal(0x18)
	SIGXFSZ   = Signal(0x19)
)

// Address families
const (
	AF_INET     = 0x2
	AF_INET6    = 0xa
	AF_UNIX     = 0x1
	AF_LOCAL    = 0x1
	AF_UNSPEC   = 0x0
)

// Socket types
const (
	SOCK_STREAM    = 0x1
	SOCK_DGRAM     = 0x2
	SOCK_RAW       = 0x3
	SOCK_SEQPACKET = 0x5
	SOCK_CLOEXEC   = 0x80000
	SOCK_NONBLOCK  = 0x800
)

// Protocol numbers
const (
	IPPROTO_IP   = 0x0
	IPPROTO_ICMP = 0x1
	IPPROTO_TCP  = 0x6
	IPPROTO_UDP  = 0x11
	IPPROTO_IPV6 = 0x29
	IPPROTO_RAW  = 0xff
)

// Socket options
const (
	SOL_SOCKET   = 0x1
	SO_REUSEADDR = 0x2
	SO_TYPE      = 0x3
	SO_ERROR     = 0x4
	SO_BROADCAST = 0x6
	SO_SNDBUF    = 0x7
	SO_RCVBUF    = 0x8
	SO_KEEPALIVE = 0x9
	SO_LINGER    = 0xd
	SO_RCVTIMEO  = 0x14
	SO_SNDTIMEO  = 0x15
	SCM_RIGHTS   = 0x1
)

// File open flags
const (
	O_RDONLY    = 0x0
	O_WRONLY    = 0x1
	O_RDWR      = 0x2
	O_APPEND    = 0x400
	O_CREAT     = 0x40
	O_EXCL      = 0x80
	O_SYNC      = 0x101000
	O_TRUNC     = 0x200
	O_NONBLOCK  = 0x800
	O_NOCTTY    = 0x100
	O_CLOEXEC   = 0x80000
	O_DIRECTORY = 0x10000
	O_NOFOLLOW  = 0x20000
)

// File mode bits
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
	S_IRUSR  = 0x100
	S_IWUSR  = 0x80
	S_IXUSR  = 0x40
	S_IRGRP  = 0x20
	S_IWGRP  = 0x10
	S_IXGRP  = 0x8
	S_IROTH  = 0x4
	S_IWOTH  = 0x2
	S_IXOTH  = 0x1
)

// Fcntl commands
const (
	F_DUPFD         = 0x0
	F_GETFD         = 0x1
	F_SETFD         = 0x2
	F_GETFL         = 0x3
	F_SETFL         = 0x4
	F_GETLK         = 0x5
	F_SETLK         = 0x6
	F_SETLKW        = 0x7
	F_SETOWN        = 0x8
	F_GETOWN        = 0x9
	F_DUPFD_CLOEXEC = 0x406
	FD_CLOEXEC      = 0x1
)

// Mmap flags
const (
	PROT_NONE  = 0x0
	PROT_READ  = 0x1
	PROT_WRITE = 0x2
	PROT_EXEC  = 0x4
	MAP_SHARED    = 0x1
	MAP_PRIVATE   = 0x2
	MAP_FIXED     = 0x10
	MAP_ANONYMOUS = 0x20
	MAP_ANON      = 0x20
)

// Wait flags
const (
	WNOHANG    = 0x1
	WUNTRACED  = 0x2
	WCONTINUED = 0x8
)

// Poll events
const (
	POLLIN   = 0x1
	POLLPRI  = 0x2
	POLLOUT  = 0x4
	POLLERR  = 0x8
	POLLHUP  = 0x10
	POLLNVAL = 0x20
)

// Rlimit constants
const (
	RLIMIT_NOFILE = 0x7
	RLIM_INFINITY = 0xffffffffffffffff
)

// Ptrace
const (
	PTRACE_TRACEME = 0x0
)

// TTY ioctls
const (
	TIOCNOTTY = 0x5422
	TIOCSCTTY = 0x540e
	TIOCGPGRP = 0x540f
	TIOCSPGRP = 0x5410
	TIOCGWINSZ = 0x5413
	TIOCSWINSZ = 0x5414
)

// Shutdown how
const (
	SHUT_RD   = 0x0
	SHUT_WR   = 0x1
	SHUT_RDWR = 0x2
)

// Clock IDs
const (
	CLOCK_REALTIME  = 0x0
	CLOCK_MONOTONIC = 0x1
)

// AT_ constants for *at syscalls
const (
	AT_FDCWD            = -0x64
	AT_REMOVEDIR        = 0x200
	AT_SYMLINK_NOFOLLOW = 0x100
	AT_EACCESS          = 0x200
	AT_EMPTY_PATH       = 0x1000
)

// Futex operations
const (
	FUTEX_WAIT            = 0x0
	FUTEX_WAKE            = 0x1
	FUTEX_PRIVATE_FLAG    = 0x80
	FUTEX_WAIT_PRIVATE    = 0x80
	FUTEX_WAKE_PRIVATE    = 0x81
)

// Clone flags (for compatibility - Freya uses spawn)
const (
	CLONE_VM             = 0x100
	CLONE_FS             = 0x200
	CLONE_FILES          = 0x400
	CLONE_SIGHAND        = 0x800
	CLONE_THREAD         = 0x10000
)

// Signal action flags
const (
	SA_NOCLDSTOP = 0x1
	SA_NOCLDWAIT = 0x2
	SA_SIGINFO   = 0x4
	SA_ONSTACK   = 0x8000000
	SA_RESTART   = 0x10000000
	SA_NODEFER   = 0x40000000
	SA_RESETHAND = 0x80000000
)

// Seek whence
const (
	SEEK_SET = 0x0
	SEEK_CUR = 0x1
	SEEK_END = 0x2
)

// Dirent types
const (
	DT_UNKNOWN = 0x0
	DT_FIFO    = 0x1
	DT_CHR     = 0x2
	DT_DIR     = 0x4
	DT_BLK     = 0x6
	DT_REG     = 0x8
	DT_LNK     = 0xa
	DT_SOCK    = 0xc
	DT_WHT     = 0xe
)

// Misc
const (
	PATH_MAX = 0x1000
	NAME_MAX = 0xff
)

// Error table
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
	15:  "block device required",
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
	26:  "text file busy",
	27:  "file too large",
	28:  "no space left on device",
	29:  "illegal seek",
	30:  "read-only file system",
	31:  "too many links",
	32:  "broken pipe",
	33:  "numerical argument out of domain",
	34:  "numerical result out of range",
	35:  "resource deadlock avoided",
	36:  "file name too long",
	37:  "no locks available",
	38:  "function not implemented",
	39:  "directory not empty",
	40:  "too many levels of symbolic links",
	42:  "no message of desired type",
	43:  "identifier removed",
	60:  "device not a stream",
	61:  "no data available",
	62:  "timer expired",
	63:  "out of streams resources",
	67:  "link has been severed",
	71:  "protocol error",
	72:  "multihop attempted",
	74:  "bad message",
	75:  "value too large for defined data type",
	84:  "invalid or incomplete multibyte or wide character",
	88:  "socket operation on non-socket",
	89:  "destination address required",
	90:  "message too long",
	91:  "protocol wrong type for socket",
	92:  "protocol not available",
	93:  "protocol not supported",
	95:  "operation not supported",
	96:  "protocol family not supported",
	97:  "address family not supported by protocol",
	98:  "address already in use",
	99:  "cannot assign requested address",
	100: "network is down",
	101: "network is unreachable",
	102: "network dropped connection on reset",
	103: "software caused connection abort",
	104: "connection reset by peer",
	105: "no buffer space available",
	106: "transport endpoint is already connected",
	107: "transport endpoint is not connected",
	110: "connection timed out",
	111: "connection refused",
	112: "host is down",
	113: "no route to host",
	114: "operation already in progress",
	115: "operation now in progress",
	116: "stale file handle",
	122: "disk quota exceeded",
	125: "operation canceled",
	130: "owner died",
	131: "state not recoverable",
}

// Signal table
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
	16: "stack fault",
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
	30: "power failure",
	31: "bad system call",
}

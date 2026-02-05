// Freya OS syscall numbers for 386.
// These map to Freya's actual syscall numbers from kernel/syscall/syscall.zig

//go:build freya && 386

package syscall

const (
	// File I/O
	SYS_READ          = 101  // FileRead
	SYS_WRITE         = 102  // FileWrite
	SYS_OPENAT        = 100  // FileOpen
	SYS_OPEN          = 100  // FileOpen (compat)
	SYS_CLOSE         = 103  // FileClose
	SYS_LSEEK         = 105  // FileSeek
	SYS_FSTATAT       = 106  // FileStat
	SYS_FSTATAT64     = 106  // FileStat
	SYS_STAT          = 106  // FileStat (compat)
	SYS_STAT64        = 106  // FileStat
	SYS_LSTAT         = 106  // FileStat (compat)
	SYS_LSTAT64       = 106  // FileStat
	SYS_FSTAT         = 107  // FileFStat
	SYS_FSTAT64       = 107  // FileFStat
	SYS_DUP           = 54   // FileDup
	SYS_DUP2          = 86   // FileDup2
	SYS_DUP3          = 86   // Use Dup2
	SYS_FCNTL         = 80   // FileFcntl
	SYS_FCNTL64       = 80   // FileFcntl
	SYS_IOCTL         = 81   // FileIoctl
	SYS_PIPE2         = 55   // PipeCreate
	SYS_PIPE          = 55   // PipeCreate (compat)
	SYS_READV         = 88   // FileReadv
	SYS_WRITEV        = 89   // FileWritev
	SYS_PREAD64       = 90   // FilePread
	SYS_PWRITE64      = 91   // FilePwrite
	SYS_FSYNC         = 113  // FileFsync
	SYS_SYNC          = 92   // FileSync
	SYS_PPOLL         = 87   // FilePoll
	SYS_POLL          = 87   // FilePoll
	SYS_PSELECT6      = 93   // FileSelect
	SYS_SELECT        = 93   // FileSelect (compat)
	SYS__NEWSELECT    = 93   // FileSelect

	// Directory operations
	SYS_MKDIR         = 108  // FileMkdir
	SYS_MKDIRAT       = 108  // FileMkdir
	SYS_UNLINK        = 109  // FileUnlink
	SYS_UNLINKAT      = 109  // FileUnlink
	SYS_RMDIR         = 111  // FileRmdir
	SYS_LINKAT        = 67   // FileLink
	SYS_LINK          = 67   // FileLink
	SYS_SYMLINKAT     = 68   // FileSymlink
	SYS_SYMLINK       = 68   // FileSymlink
	SYS_READLINKAT    = 69   // FileReadlink
	SYS_READLINK      = 69   // FileReadlink
	SYS_GETDENTS64    = 104  // FileReadDir
	SYS_GETDENTS      = 104  // FileReadDir
	SYS_CHDIR         = 56   // FileChdir
	SYS_GETCWD        = 57   // FileGetcwd
	SYS_FCHMOD        = 58   // FileChmod
	SYS_CHMOD         = 58   // FileChmod
	SYS_FCHMODAT      = 58   // FileChmod
	SYS_FACCESSAT     = 85   // FileAccess
	SYS_ACCESS        = 85   // FileAccess
	SYS_TRUNCATE      = 110  // FileTruncate
	SYS_TRUNCATE64    = 110  // FileTruncate
	SYS_FTRUNCATE     = 110  // FileTruncate
	SYS_FTRUNCATE64   = 110  // FileTruncate
	SYS_UMASK         = 84   // FileUmask
	SYS_UTIMENSAT     = 59   // FileUtime
	SYS_UTIME         = 59   // FileUtime
	SYS_RENAMEAT      = 112  // FileRename
	SYS_RENAME        = 112  // FileRename

	// Mount operations
	SYS_MOUNT         = 82   // FileMount
	SYS_UMOUNT2       = 83   // FileUnmount
	SYS_UMOUNT        = 83   // FileUnmount

	// Memory management
	SYS_MMAP          = 20   // MemoryMap
	SYS_MMAP2         = 20   // MemoryMap
	SYS_MUNMAP        = 21   // MemoryUnmap
	SYS_MPROTECT      = 22   // MemoryProtect
	SYS_BRK           = 20   // Use mmap

	// Process management
	SYS_EXIT          = 10   // ProcessExit
	SYS_EXIT_GROUP    = 10   // ProcessExit
	SYS_GETPID        = 13   // GetPid
	SYS_GETTID        = 14   // GetTid
	SYS_GETPPID       = 35   // GetPPid
	SYS_WAIT4         = 12   // ProcessWait
	SYS_WAITPID       = 12   // ProcessWait
	SYS_WAITID        = 12   // ProcessWait
	SYS_KILL          = 27   // ProcessKill
	SYS_TGKILL        = 27   // ProcessKill (use same)
	SYS_TKILL         = 27   // ProcessKill

	// Process creation (Freya uses spawn, not fork)
	SYS_CLONE         = 250  // ProcessSpawn (not real clone)
	SYS_FORK          = 250  // ProcessSpawn (not real fork)
	SYS_VFORK         = 250  // ProcessSpawn
	SYS_EXECVE        = 26   // ProcessExec

	// User/Group
	SYS_GETUID        = 28   // GetUid
	SYS_GETUID32      = 28   // GetUid
	SYS_SETUID        = 29   // SetUid
	SYS_SETUID32      = 29   // SetUid
	SYS_GETGID        = 33   // GetGid
	SYS_GETGID32      = 33   // GetGid
	SYS_SETGID        = 34   // SetGid
	SYS_SETGID32      = 34   // SetGid
	SYS_GETEUID       = 28   // GetUid (no euid distinction)
	SYS_GETEUID32     = 28   // GetUid
	SYS_GETEGID       = 33   // GetGid (no egid distinction)
	SYS_GETEGID32     = 33   // GetGid
	SYS_SETREUID      = 29   // SetUid
	SYS_SETREUID32    = 29   // SetUid
	SYS_SETREGID      = 34   // SetGid
	SYS_SETREGID32    = 34   // SetGid
	SYS_SETGROUPS     = 34   // SetGid (simplified)
	SYS_SETGROUPS32   = 34   // SetGid
	SYS_GETGROUPS     = 33   // GetGid (simplified)
	SYS_GETGROUPS32   = 33   // GetGid

	// Process groups/sessions
	SYS_SETPGID       = 228  // SetPgid
	SYS_GETPGID       = 227  // GetPgid
	SYS_GETPGRP       = 229  // GetPgrp
	SYS_SETSID        = 232  // SetSid
	SYS_GETSID        = 231  // GetSid

	// Time
	SYS_NANOSLEEP         = 36   // Nanosleep
	SYS_CLOCK_GETTIME     = 37   // ClockGetTime
	SYS_CLOCK_GETTIME32   = 37   // ClockGetTime
	SYS_CLOCK_GETTIME64   = 37   // ClockGetTime
	SYS_CLOCK_SETTIME     = 38   // ClockSetTime
	SYS_CLOCK_SETTIME32   = 38   // ClockSetTime
	SYS_CLOCK_SETTIME64   = 38   // ClockSetTime
	SYS_CLOCK_GETRES      = 37   // ClockGetTime (same for now)
	SYS_CLOCK_GETRES_TIME64 = 37 // ClockGetTime
	SYS_GETTIMEOFDAY      = 42   // GetTimeOfDay
	SYS_SETTIMEOFDAY      = 38   // ClockSetTime
	SYS_TIME              = 42   // GetTimeOfDay
	SYS_TIMER_CREATE      = 222  // TimerCreate
	SYS_TIMER_SETTIME     = 223  // TimerSettime
	SYS_TIMER_SETTIME64   = 223  // TimerSettime
	SYS_TIMER_GETTIME     = 224  // TimerGettime
	SYS_TIMER_GETTIME64   = 224  // TimerGettime
	SYS_TIMER_DELETE      = 225  // TimerDelete
	SYS_TIMER_GETOVERRUN  = 226  // TimerGetoverrun
	SYS_GETITIMER         = 42   // GetTimeOfDay (stub)
	SYS_SETITIMER         = 38   // ClockSetTime (stub)

	// Resource limits
	SYS_GETRLIMIT     = 40   // GetRLimit
	SYS_SETRLIMIT     = 41   // SetRLimit
	SYS_PRLIMIT64     = 40   // GetRLimit (use same)
	SYS_GETRUSAGE     = 39   // GetRusage

	// Synchronization
	SYS_FUTEX         = 44   // Futex
	SYS_FUTEX_TIME64  = 44   // Futex

	// Signals
	SYS_RT_SIGACTION    = 94   // RtSigaction
	SYS_RT_SIGPROCMASK  = 95   // RtSigprocmask
	SYS_RT_SIGPENDING   = 96   // RtSigpending
	SYS_RT_SIGRETURN    = 97   // RtSigreturn
	SYS_RT_SIGSUSPEND   = 94   // RtSigaction (stub)
	SYS_SIGALTSTACK     = 94   // RtSigaction (stub)
	SYS_SIGNAL          = 94   // RtSigaction (compat)
	SYS_SIGACTION       = 94   // RtSigaction (compat)
	SYS_SIGPROCMASK     = 95   // RtSigprocmask (compat)
	SYS_SIGPENDING      = 96   // RtSigpending (compat)
	SYS_SIGSUSPEND      = 94   // RtSigaction (compat)

	// Network/Socket
	SYS_SOCKET        = 70   // SocketCreate
	SYS_BIND          = 71   // SocketBind
	SYS_CONNECT       = 72   // SocketConnect
	SYS_LISTEN        = 78   // SocketListen
	SYS_ACCEPT        = 79   // SocketAccept
	SYS_ACCEPT4       = 79   // SocketAccept
	SYS_SENDTO        = 75   // SocketSendTo
	SYS_RECVFROM      = 76   // SocketRecvFrom
	SYS_SEND          = 73   // SocketSend
	SYS_RECV          = 74   // SocketRecv
	SYS_SENDMSG       = 73   // SocketSend
	SYS_RECVMSG       = 74   // SocketRecv
	SYS_SHUTDOWN      = 77   // SocketClose
	SYS_SETSOCKOPT    = 70   // SocketCreate (stub)
	SYS_GETSOCKOPT    = 70   // SocketCreate (stub)
	SYS_GETSOCKNAME   = 70   // SocketCreate (stub)
	SYS_GETPEERNAME   = 70   // SocketCreate (stub)
	SYS_SOCKETPAIR    = 70   // SocketCreate (stub)
	SYS_SOCKETCALL    = 70   // SocketCreate (multiplexer)

	// System info
	SYS_SYSINFO       = 99   // SysInfo
	SYS_UNAME         = 99   // SysInfo

	// Ptrace
	SYS_PTRACE        = 117  // Linux compat number (not implemented)

	// Misc
	SYS_CHROOT        = 56   // FileChdir (no real chroot)
	SYS_FCHOWN        = 58   // FileChmod (simplified)
	SYS_FCHOWN32      = 58   // FileChmod
	SYS_CHOWN         = 58   // FileChmod
	SYS_CHOWN32       = 58   // FileChmod
	SYS_LCHOWN        = 58   // FileChmod
	SYS_LCHOWN32      = 58   // FileChmod
	SYS_FCHOWNAT      = 58   // FileChmod
	SYS_GETRANDOM     = 120  // DebugPutChar (stub - needs /dev/urandom)
	SYS_EPOLL_CREATE1 = 87   // FilePoll (stub)
	SYS_EPOLL_CREATE  = 87   // FilePoll (stub)
	SYS_EPOLL_CTL     = 87   // FilePoll (stub)
	SYS_EPOLL_PWAIT   = 87   // FilePoll (stub)
	SYS_EPOLL_WAIT    = 87   // FilePoll (stub)
	SYS_EVENTFD2      = 55   // PipeCreate (stub)
	SYS_EVENTFD       = 55   // PipeCreate (stub)

	// Thread local storage
	SYS_SET_TID_ADDRESS = 14   // GetTid (stub)
	SYS_SET_THREAD_AREA = 14   // GetTid (stub)
	SYS_GET_THREAD_AREA = 14   // GetTid (stub)

	// Prctl
	SYS_PRCTL         = 99   // SysInfo (stub)

	// Scheduling
	SYS_SCHED_YIELD           = 15  // ThreadYield
	SYS_SCHED_GETAFFINITY     = 99  // SysInfo (stub)
	SYS_SCHED_SETAFFINITY     = 99  // SysInfo (stub)
	SYS_SCHED_GETSCHEDULER    = 99  // SysInfo (stub)
	SYS_SCHED_SETSCHEDULER    = 99  // SysInfo (stub)
	SYS_SCHED_GETPARAM        = 99  // SysInfo (stub)
	SYS_SCHED_SETPARAM        = 99  // SysInfo (stub)
	SYS_SCHED_GET_PRIORITY_MAX = 99 // SysInfo (stub)
	SYS_SCHED_GET_PRIORITY_MIN = 99 // SysInfo (stub)

	// Memory locking (stubs)
	SYS_MLOCK         = 22   // MemoryProtect (stub)
	SYS_MUNLOCK       = 22   // MemoryProtect (stub)
	SYS_MLOCKALL      = 22   // MemoryProtect (stub)
	SYS_MUNLOCKALL    = 22   // MemoryProtect (stub)
	SYS_MADVISE       = 22   // MemoryProtect (stub)
	SYS_MINCORE       = 22   // MemoryProtect (stub)
	SYS_MSYNC         = 92   // FileSync (stub)

	// Statfs
	SYS_STATFS        = 106  // FileStat (stub)
	SYS_STATFS64      = 106  // FileStat (stub)
	SYS_FSTATFS       = 107  // FileFStat (stub)
	SYS_FSTATFS64     = 107  // FileFStat (stub)

	// IPC (stubs - Freya uses capability-based IPC)
	SYS_IPC           = 99   // SysInfo (stub) - multiplexer
	SYS_SEMGET        = 99   // SysInfo (stub)
	SYS_SEMOP         = 99   // SysInfo (stub)
	SYS_SEMCTL        = 99   // SysInfo (stub)
	SYS_SEMTIMEDOP    = 99   // SysInfo (stub)
	SYS_SHMGET        = 99   // SysInfo (stub)
	SYS_SHMAT         = 99   // SysInfo (stub)
	SYS_SHMDT         = 99   // SysInfo (stub)
	SYS_SHMCTL        = 99   // SysInfo (stub)
	SYS_MSGGET        = 99   // SysInfo (stub)
	SYS_MSGSND        = 99   // SysInfo (stub)
	SYS_MSGRCV        = 99   // SysInfo (stub)
	SYS_MSGCTL        = 99   // SysInfo (stub)

	// Debug
	SYS_ARCH_SPECIFIC_SYSCALL = 120  // DebugPutChar
)

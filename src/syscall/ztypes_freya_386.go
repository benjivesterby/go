// Freya OS type definitions for 386.

//go:build freya && 386

package syscall

const (
	sizeofPtr      = 0x4
	sizeofShort    = 0x2
	sizeofInt      = 0x4
	sizeofLong     = 0x4
	sizeofLongLong = 0x8
	PathMax        = 0x1000
)

type (
	_C_short     int16
	_C_int       int32
	_C_long      int32
	_C_long_long int64
)

type Timespec struct {
	Sec  int32
	Nsec int32
}

type Timeval struct {
	Sec  int32
	Usec int32
}

type Time_t int32

type Rusage struct {
	Utime    Timeval
	Stime    Timeval
	Maxrss   int32
	Ixrss    int32
	Idrss    int32
	Isrss    int32
	Minflt   int32
	Majflt   int32
	Nswap    int32
	Inblock  int32
	Oublock  int32
	Msgsnd   int32
	Msgrcv   int32
	Nsignals int32
	Nvcsw    int32
	Nivcsw   int32
}

type Rlimit struct {
	Cur uint64
	Max uint64
}

type _Gid_t uint32

type Stat_t struct {
	Dev     uint64
	_       uint16
	_       uint16
	Ino     uint64
	Mode    uint32
	Nlink   uint32
	Uid     uint32
	Gid     uint32
	Rdev    uint64
	_       uint16
	_       [4]byte
	Size    int64
	Blksize int32
	_       [4]byte
	Blocks  int64
	Atim    Timespec
	Mtim    Timespec
	Ctim    Timespec
	_       uint64
}

type Statfs_t struct {
	Type    int32
	Bsize   int32
	Blocks  uint64
	Bfree   uint64
	Bavail  uint64
	Files   uint64
	Ffree   uint64
	Fsid    Fsid
	Namelen int32
	Frsize  int32
	Flags   int32
	Spare   [4]int32
	_       [4]byte
}

type Dirent struct {
	Ino    uint64
	Off    int64
	Reclen uint16
	Type   uint8
	Name   [256]uint8
	_      [1]byte
}

type Fsid struct {
	Val [2]int32
}

type Flock_t struct {
	Type   int16
	Whence int16
	Start  int64
	Len    int64
	Pid    int32
}

type RawSockaddrInet4 struct {
	Family uint16
	Port   uint16
	Addr   [4]byte /* in_addr */
	Zero   [8]uint8
}

type RawSockaddrInet6 struct {
	Family   uint16
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte /* in6_addr */
	Scope_id uint32
}

type RawSockaddrUnix struct {
	Family uint16
	Path   [108]int8
}

type RawSockaddr struct {
	Family uint16
	Data   [14]int8
}

type RawSockaddrAny struct {
	Addr RawSockaddr
	Pad  [96]int8
}

type _Socklen uint32

type Linger struct {
	Onoff  int32
	Linger int32
}

type Iovec struct {
	Base *byte
	Len  uint32
}

type IPMreq struct {
	Multiaddr [4]byte /* in_addr */
	Interface [4]byte /* in_addr */
}

type IPMreqn struct {
	Multiaddr [4]byte /* in_addr */
	Address   [4]byte /* in_addr */
	Ifindex   int32
}

type IPv6Mreq struct {
	Multiaddr [16]byte /* in6_addr */
	Interface uint32
}

type Msghdr struct {
	Name       *byte
	Namelen    uint32
	Iov        *Iovec
	Iovlen     uint32
	Control    *byte
	Controllen uint32
	Flags      int32
}

type Cmsghdr struct {
	Len   uint32
	Level int32
	Type  int32
}

type Inet4Pktinfo struct {
	Ifindex  int32
	Spec_dst [4]byte /* in_addr */
	Addr     [4]byte /* in_addr */
}

type Inet6Pktinfo struct {
	Addr    [16]byte /* in6_addr */
	Ifindex uint32
}

type IPv6MTUInfo struct {
	Addr RawSockaddrInet6
	Mtu  uint32
}

type ICMPv6Filter struct {
	Data [8]uint32
}

type Ucred struct {
	Pid int32
	Uid uint32
	Gid uint32
}

const (
	SizeofSockaddrInet4  = 0x10
	SizeofSockaddrInet6  = 0x1c
	SizeofSockaddrAny    = 0x70
	SizeofSockaddrUnix   = 0x6e
	SizeofLinger         = 0x8
	SizeofIPMreq         = 0x8
	SizeofIPMreqn        = 0xc
	SizeofIPv6Mreq       = 0x14
	SizeofMsghdr         = 0x1c
	SizeofCmsghdr        = 0xc
	SizeofInet4Pktinfo   = 0xc
	SizeofInet6Pktinfo   = 0x14
	SizeofIPv6MTUInfo    = 0x20
	SizeofICMPv6Filter   = 0x20
	SizeofUcred          = 0xc
)

type FdSet struct {
	Bits [32]int32
}

type Utsname struct {
	Sysname    [65]int8
	Nodename   [65]int8
	Release    [65]int8
	Version    [65]int8
	Machine    [65]int8
	Domainname [65]int8
}

const (
	_AT_FDCWD            = -0x64
	_AT_REMOVEDIR        = 0x200
	_AT_SYMLINK_NOFOLLOW = 0x100
	_AT_EACCESS          = 0x200
	_AT_EMPTY_PATH       = 0x1000
)

type PollFd struct {
	Fd      int32
	Events  int16
	Revents int16
}

type Termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Line   uint8
	Cc     [32]uint8
	_      [3]byte
	Ispeed uint32
	Ospeed uint32
}

const (
	IUCLC  = 0x200
	OLCUC  = 0x2
	TCGETS = 0x5401
	TCSETS = 0x5402
	XCASE  = 0x4
)

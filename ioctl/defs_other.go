//go:build !cgo && linux && (mips64 || mips64le || mips || mipsle || ppc64 || ppc64le || ppc || sparc64)

package ioctl

const (
	_IOC_NRBITS   = 8
	_IOC_TYPEBITS = 8
	_IOC_SIZEBITS = 13
	_IOC_DIRBITS  = 3

	_IOC_NRSHIFT   = 0
	_IOC_TYPESHIFT = _IOC_NRSHIFT + _IOC_NRBITS
	_IOC_SIZESHIFT = _IOC_TYPESHIFT + _IOC_TYPEBITS
	_IOC_DIRSHIFT  = _IOC_SIZESHIFT + _IOC_SIZEBITS

	_IOC_NONE  = 1
	_IOC_READ  = 2
	_IOC_WRITE = 4
)

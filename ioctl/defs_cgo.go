//go:build cgo

package ioctl

/*
#include <linux/ioctl.h>
*/
import "C"

const (
	_IOC_DIRSHIFT  = C._IOC_DIRSHIFT
	_IOC_TYPESHIFT = C._IOC_TYPESHIFT
	_IOC_NRSHIFT   = C._IOC_NRSHIFT
	_IOC_SIZESHIFT = C._IOC_SIZESHIFT
	_IOC_NONE      = C._IOC_NONE
	_IOC_READ      = C._IOC_READ
	_IOC_WRITE     = C._IOC_WRITE
)

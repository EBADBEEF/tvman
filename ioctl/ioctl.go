package ioctl

func IOC(dir uint32, itype uint32, nr uint32, size uintptr) uint32 {
	return (dir << _IOC_DIRSHIFT) | (itype << _IOC_TYPESHIFT) | (nr << _IOC_NRSHIFT) | (uint32(size) << _IOC_SIZESHIFT)
}

func IO(itype uint32, nr uint32) uint32 {
	return IOC(_IOC_NONE, itype, nr, 0)
}

func IOR(itype uint32, nr uint32, size uintptr) uint32 {
	return IOC(_IOC_READ, itype, nr, size)
}

func IOW(itype uint32, nr uint32, size uintptr) uint32 {
	return IOC(_IOC_WRITE, itype, nr, size)
}

func IOWR(itype uint32, nr uint32, size uintptr) uint32 {
	return IOC(_IOC_READ|_IOC_WRITE, itype, nr, size)
}

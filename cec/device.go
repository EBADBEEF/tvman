package cec

import (
	"os"
	"syscall"
	"unsafe"
)

type Device struct {
	Rx      chan *CecMsg
	Event   chan *CecEvent
	Tx      chan *CecMsg
	TxReady chan bool
	Error   chan error
	file    *os.File
}

// Open a cec device, sets the mode (based on CEC_MODE flags), and reads
// current settings. Empty flags means the device is an initiator (one who
// transmits independently), but not a follower (one who replies to messages).
// A device can be both an initiator and follower but if the device is a
// monitor it cannot be either initiator or follower.
//
// example modeFlags: CEC_MODE_MONITOR_ALL
// example modeFlags: CEC_MODE_INITIATOR | CEC_MODE_FOLLOWER
func Open(pathname string, modeFlags uint32) (*Device, error) {
	d := Device{}

	file, err := os.OpenFile(pathname, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), CEC_S_MODE, uintptr(unsafe.Pointer(&modeFlags)))
	if errno != 0 {
		file.Close()
		return nil, errno
	}

	d.file = file
	d.Error = make(chan error)

	return &d, nil
}

func (d *Device) StartReceiver() error {
	d.Rx = make(chan *CecMsg, 10)
	d.Event = make(chan *CecEvent, 10)
	go func() {
		for {
			cm := CecMsg{}
			if err := d.Receive(&cm); err != nil {
				d.Error <- err
			} else {
				d.Rx <- &cm
			}
		}
	}()
	go func() {
		for {
			cm := CecEvent{}
			if err := d.DqEvent(&cm); err != nil {
				d.Error <- err
			} else {
				d.Event <- &cm
			}
		}
	}()
	return nil
}

func (d *Device) StartTransmitter() error {
	d.Tx = make(chan *CecMsg)
	d.TxReady = make(chan bool)
	go func() {
		for {
			d.TxReady <- true
			cm, open := <-d.Tx
			if !open {
				return
			}
			if err := d.Transmit(cm); err != nil {
				d.Error <- err
			}
		}
	}()
	return nil
}

func (d *Device) Close() error {
	if d.file != nil {
		return d.file.Close()
	} else {
		return nil
	}
}

// Wait for an event
func (d *Device) DqEvent(event *CecEvent) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, d.file.Fd(), CEC_DQEVENT, uintptr(unsafe.Pointer(event)))
	if errno == 0 {
		return nil
	} else {
		return errno
	}
}

func (d *Device) Transmit(cm *CecMsg) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, d.file.Fd(), CEC_TRANSMIT, uintptr(unsafe.Pointer(cm)))
	if errno == 0 {
		return nil
	} else {
		return errno
	}
}

func (d *Device) Receive(cm *CecMsg) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, d.file.Fd(), CEC_RECEIVE, uintptr(unsafe.Pointer(cm)))
	if errno == 0 {
		return nil
	} else {
		return errno
	}
}

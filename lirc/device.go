package lirc

import (
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	delayBetweenWrites = time.Duration(10 * time.Millisecond)
)

type Device struct {
	Rx      chan any
	Tx      chan any
	TxReady chan bool
	Error   chan error
	file    *os.File
	table   map[Scancode]any
}

func getLircDevice(uevent string) (string, string, error) {
	const rcSysfs = "/sys/class/rc/"
	lircDevice := ""
	protocolsPath := ""

	devices, err := os.ReadDir(rcSysfs)
	if err != nil {
		return lircDevice, protocolsPath, errors.New("could not read " + rcSysfs)
	}

	for _, rc := range devices {
		// check to see if uevent matches
		if uevent != "" {
			file, err := os.Open(filepath.Join(rcSysfs, rc.Name(), "uevent"))
			if err != nil {
				continue
			}
			buf := make([]byte, 256)
			if _, err := file.Read(buf); err != nil {
				continue
			}
			if !strings.Contains(string(buf), uevent) {
				continue
			}
		}
		// check to see if lirc device present
		rd, _ := os.ReadDir(filepath.Join(rcSysfs, rc.Name()))
		for _, file := range rd {
			if strings.HasPrefix(file.Name(), "lirc") {
				lircDevice = filepath.Join("/dev", file.Name())
				protocolsPath = filepath.Join(rcSysfs, rc.Name(), "protocols")
				return lircDevice, protocolsPath, nil
			}
		}
	}

	return lircDevice, protocolsPath, errors.New("could not find lirc device")
}

func Open(uevent string, protocols string) (*Device, error) {

	lircDevice, protocolsPath, err := getLircDevice(uevent)
	if err != nil {
		return nil, err
	}

	if protocols == "" {
		protocols = "all"
	}

	setProtosBytes := []byte(strings.Join(strings.Fields("none "+protocols), "\n+"))
	if err := os.WriteFile(protocolsPath, setProtosBytes, 0644); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(lircDevice, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	d := Device{
		file:  file,
		table: make(map[Scancode]any),
	}

	return &d, nil
}

func (d *Device) Write(data any) error {
	return binary.Write(d.file, binary.NativeEndian, data)
}

func (d *Device) Ioctl(req uintptr, arg uint32) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, d.file.Fd(), req, uintptr(unsafe.Pointer(&arg)))
	if errno == 0 {
		return nil
	} else {
		return errno
	}
}

func (d *Device) SetReceiveMode(flags uint32) error {
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL, d.file.Fd(), LIRC_SET_REC_MODE, uintptr(unsafe.Pointer(&flags)))
	if errno != 0 {
		return errno
	}
	return nil
}

func (d *Device) SetSendMode(flags uint32) error {
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL, d.file.Fd(), LIRC_SET_SEND_MODE, uintptr(unsafe.Pointer(&flags)))
	if errno != 0 {
		return errno
	}
	return nil
}

func (d *Device) RegisterScancode(sc *Scancode, val any) error {
	if d.table == nil {
		return errors.New("device not open")
	}
	key := Scancode{
		RcProto:  sc.RcProto,
		Scancode: sc.Scancode,
	}
	if _, found := d.table[key]; found {
		return errors.New("duplicate key")
	}
	d.table[key] = val
	return nil
}

func (d *Device) StartScancodeTransmitter() error {
	if err := d.SetSendMode(LIRC_MODE_SCANCODE); err != nil {
		return err
	}
	d.Tx = make(chan any)
	d.TxReady = make(chan bool)
	go func() {
		for {
			d.TxReady <- true
			sc, open := <-d.Tx
			if !open {
				return
			}
			if _, ok := sc.(Scancode); !ok {
				d.Error <- errors.New("not a scancode")
			}
			// Put this before write because this transmit could have been
			// triggered by a real IR remote
			time.Sleep(delayBetweenWrites)
			if err := d.Write(sc); err != nil {
				d.Error <- err
			}
		}
	}()
	return nil
}

type ScancodeReaderEvent struct {
	Scancode Scancode
	User     any
}

func (d *Device) StartScancodeReader() error {
	if err := d.SetReceiveMode(LIRC_MODE_SCANCODE); err != nil {
		return err
	}

	d.Rx = make(chan any, 10)
	go func() {
		for {
			sc := Scancode{}
			if err := d.Receive(&sc); err != nil {
				d.Error <- err
				return
			}
			val, found := d.table[Scancode{RcProto: sc.RcProto, Scancode: sc.Scancode}]
			if found {
				d.Rx <- ScancodeReaderEvent{Scancode: sc, User: val}
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

func (d *Device) Receive(sc *Scancode) error {
	return binary.Read(d.file, binary.NativeEndian, sc)

}

package lgtv

import (
	"errors"
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

const (
	protectedDuration = time.Duration(10 * time.Second)
)

var (
	ProtectedError = errors.New("Blocked by CrashProtect")
)

type Device struct {
	Rx        chan string
	Tx        chan string
	TxReady   chan bool
	Error     chan error
	file      *os.File
	protected time.Time
	mu        sync.Mutex
}

func Open(pathname string) (*Device, error) {
	d := Device{}

	file, err := os.OpenFile(pathname, os.O_RDWR|syscall.O_NOCTTY|syscall.O_NDELAY, 0)
	if err != nil {
		return nil, err
	}

	termios := syscall.Termios{
		Iflag: syscall.IGNPAR,
		Cflag: syscall.CS8 | syscall.CREAD | syscall.CLOCAL | syscall.B9600,
		// see termios(3): noncanonical read "MIN > 0, TIME > 0 (read with interbyte timeout)"
		// VTIME is in deciseconds
		Cc:     [32]uint8{syscall.VMIN: 1, syscall.VTIME: 1},
		Ispeed: syscall.B9600,
		Ospeed: syscall.B9600,
	}

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), syscall.TCSETS, uintptr(unsafe.Pointer(&termios)))
	if errno != 0 {
		file.Close()
		return nil, errno
	}

	d.file = file
	d.Error = make(chan error)

	return &d, nil
}

func (d *Device) StartReceiver() error {
	d.Rx = make(chan string, 10)
	go func() {
		for {
			buf := make([]uint8, 256)
			n, err := d.Receive(buf)
			str := string(buf[:n])
			if err != nil {
				d.Error <- err
			} else {
				d.Rx <- str
			}
			// TODO: report discarded reads?
		}
	}()
	return nil
}

func (d *Device) StartTransmitter() error {
	d.Tx = make(chan string)
	d.TxReady = make(chan bool)
	go func() {
		for {
			d.TxReady <- true
			item, open := <-d.Tx
			if !open {
				return
			}
			if err := d.Transmit(item); err != nil {
				d.Error <- err
			}
		}
	}()
	return nil
}

func (d *Device) Receive(b []uint8) (int, error) {
	n, err := d.file.Read(b)
	return n, err
}

func (d *Device) CrashProtect() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.protected = time.Now()
}

func (d *Device) isProtected() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return time.Since(d.protected) < protectedDuration
}

type tvMessage struct {
	buf     string
	protect bool
}

var commands = map[string]tvMessage{
	"on":         {buf: "ka 0 01", protect: true},
	"off":        {buf: "ka 0 00", protect: true},
	"check":      {buf: "ka 0 ff", protect: true},
	"remote-off": {buf: "km 0 01"},
	"remote-on":  {buf: "km 0 00"},
	"osd-off":    {buf: "kl 0 00"},
	"osd-on":     {buf: "kl 0 01"},
	"mute":       {buf: "kf 0 00"},
	"vol":        {buf: "kf 0 a0"},
	"blank":      {buf: "kd 0 01"},
	"unblank":    {buf: "kd 0 00"},

	"input-tv":           {buf: "xb 0 00"},
	"input-cable":        {buf: "xb 0 01"},
	"input-analog":       {buf: "xb 0 10"},
	"input-analog-cable": {buf: "xb 0 11"},
	"input-av1":          {buf: "xb 0 20"},
	"input-av2":          {buf: "xb 0 21"},
	"input-component1":   {buf: "xb 0 40"},
	"input-component2":   {buf: "xb 0 41"},
	"input-component3":   {buf: "xb 0 42"},
	"input-rgb-pc":       {buf: "xb 0 60"},
	"input-hdmi1":        {buf: "xb 0 90"},
	"input-hdmi2":        {buf: "xb 0 91"},
	"input-hdmi3":        {buf: "xb 0 92"},
	"input-hdmi4":        {buf: "xb 0 93"},

	"3d-off":       {buf: "xt 0 01 00 00 00"},
	"3d-on-sbs-rl": {buf: "xt 0 00 01 00 00"},
	"3d-on-sbs-lr": {buf: "xt 0 00 01 01 00"},
	"3d-on-tb-rl":  {buf: "xt 0 00 00 00 00"},
	"3d-on-tb-lr":  {buf: "xt 0 00 00 01 00"},
	"3d-on-fp-rl":  {buf: "xt 0 00 03 00 00"},
	"3d-on-fp-lr":  {buf: "xt 0 00 03 01 00"},
	"3d-lr":        {buf: "xv 0 00 01"},
	"3d-rl":        {buf: "xv 0 00 00"},

	"rc-channel-up":   {buf: "mc 0 00"},
	"rc-channel-down": {buf: "mc 0 01"},
	"rc-volume-up":    {buf: "mc 0 02"},
	"rc-volume-down":  {buf: "mc 0 03"},
	"rc-right":        {buf: "mc 0 06"},
	"rc-left":         {buf: "mc 0 07"},
	"rc-power":        {buf: "mc 0 08"},
	"rc-mute":         {buf: "mc 0 09"},
	"rc-input":        {buf: "mc 0 0b"},
	"rc-sleep":        {buf: "mc 0 0e"},
	"rc-tv/rad":       {buf: "mc 0 0f"},
	"rc-num-0":        {buf: "mc 0 10"},
	"rc-num-1":        {buf: "mc 0 11"},
	"rc-num-2":        {buf: "mc 0 12"},
	"rc-num-3":        {buf: "mc 0 13"},
	"rc-num-4":        {buf: "mc 0 14"},
	"rc-num-5":        {buf: "mc 0 15"},
	"rc-num-6":        {buf: "mc 0 16"},
	"rc-num-7":        {buf: "mc 0 17"},
	"rc-num-8":        {buf: "mc 0 18"},
	"rc-num-9":        {buf: "mc 0 19"},
	"rc-flashback":    {buf: "mc 0 1a"},
	"rc-favorite":     {buf: "mc 0 1e"},
	"rc-teletext":     {buf: "mc 0 20"},
	"rc-teletext-opt": {buf: "mc 0 21"},
	"rc-back":         {buf: "mc 0 28"},
	"rc-av-mode":      {buf: "mc 0 30"},
	"rc-subtitle":     {buf: "mc 0 39"},
	"rc-up":           {buf: "mc 0 40"},
	"rc-down":         {buf: "mc 0 41"},
	"rc-apps":         {buf: "mc 0 42"},
	"rc-settings":     {buf: "mc 0 43"},
	"rc-select":       {buf: "mc 0 44"},
	"rc-menu":         {buf: "mc 0 45"},
	"rc-minus":        {buf: "mc 0 4c"},
	"rc-picture":      {buf: "mc 0 4d"},
	"rc-sound":        {buf: "mc 0 52"},
	"rc-list":         {buf: "mc 0 53"},
	"rc-exit":         {buf: "mc 0 5b"},
	"rc-pip":          {buf: "mc 0 60"},
	"rc-blue":         {buf: "mc 0 61"},
	"rc-yellow":       {buf: "mc 0 63"},
	"rc-green":        {buf: "mc 0 71"},
	"rc-red":          {buf: "mc 0 72"},
	"rc-aspect-ratio": {buf: "mc 0 79"},
	"rc-audio-desc":   {buf: "mc 0 91"},
	"rc-live-menu":    {buf: "mc 0 9e"},
	"rc-userguide":    {buf: "mc 0 7a"},
	"rc-home":         {buf: "mc 0 7c"},
	"rc-simplink":     {buf: "mc 0 7e"},
	"rc-ff":           {buf: "mc 0 8e"},
	"rc-rewind":       {buf: "mc 0 8f"},
	"rc-info":         {buf: "mc 0 aa"},
	"rc-guide":        {buf: "mc 0 ab"},
	"rc-play":         {buf: "mc 0 b0"},
	"rc-stop":         {buf: "mc 0 b1"},
	"rc-recent":       {buf: "mc 0 b5"},
	"rc-pause":        {buf: "mc 0 ba"},
	"rc-record":       {buf: "mc 0 bd"},
	"rc-3d":           {buf: "mc 0 dc"},
}

func (d *Device) Transmit(s string) error {
	msg, found := commands[s]
	if !found {
		return errors.New("invalid command " + s)
	}
	if msg.protect && d.isProtected() {
		return ProtectedError
	}
	d.file.WriteString("\r\n")
	d.file.WriteString(msg.buf)
	d.file.WriteString("\r\n")
	/*
		nbToWrite := len(msg.buf)
		nbWritten := 0
		for nbWritten < nbToWrite {
			n, err := d.file.WriteString(msg.buf[nbWritten:])
			if n == 0 {
				return nil
			} else if n < 0 {
				return err
			}
			nbWritten += n
		}
	*/
	return nil
}

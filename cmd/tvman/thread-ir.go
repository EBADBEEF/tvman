package main

import (
	"container/list"
	"time"
	"tvman/lirc"
	"tvman/msg"
)

type action map[string]any

type irKey struct {
	name         string
	crashProtect bool
	state        *state
	actions      action
}

const (
	TV_ON         = (1 << 0)
	TV_OFF        = (1 << 1)
	SHIELD_ON     = (1 << 2)
	SHIELD_OFF    = (1 << 3)
	BLURAY_ON     = (1 << 4)
	BLURAY_OFF    = (1 << 5)
	RECEIVER_ON   = (1 << 6)
	RECEIVER_OFF  = (1 << 7)
	TV_MASK       = (TV_ON | TV_OFF)
	SHIELD_MASK   = (SHIELD_ON | SHIELD_OFF)
	BLURAY_MASK   = (BLURAY_ON | BLURAY_OFF)
	RECEIVER_MASK = (RECEIVER_ON | RECEIVER_OFF)
)

var txTable = make(map[string]lirc.Scancode)

func registerKeys(dev *lirc.Device) {
	sc := lirc.Scancode{}
	reg := func(proto uint16, code uint64, key irKey) {
		sc.RcProto = proto
		sc.Scancode = code
		txTable[key.name] = sc
		if err := dev.RegisterScancode(&sc, key); err != nil {
			msg.Error("ir: could not register key (%v)", err)
		}
	}
	reg(lirc.RC_PROTO_NEC, 0x4c4, irKey{name: "lg-tv-on", crashProtect: true})
	reg(lirc.RC_PROTO_NEC, 0x4c5, irKey{name: "lg-tv-off", state: STATE_MASTER_OFF, crashProtect: true})
	reg(lirc.RC_PROTO_NEC, 0x456, irKey{name: "lg-tv-red", actions: action{"live_tv": tvCommand("rc-subtitle")}})
	reg(lirc.RC_PROTO_NEC, 0x45c, irKey{name: "lg-tv-blue"})
	reg(lirc.RC_PROTO_NECX, 0xd26d04, irKey{name: "onkyo-on"})
	reg(lirc.RC_PROTO_NECX, 0xd26c47, irKey{name: "onkyo-off", state: STATE_MASTER_OFF})
	reg(lirc.RC_PROTO_NECX, 0x2d2d75, irKey{name: "shield-on", state: STATE_SHIELD})
	reg(lirc.RC_PROTO_NECX, 0x2d2d76, irKey{name: "shield-off"})
	reg(lirc.RC_PROTO_SONY20, 0x1ae22e, irKey{name: "bluray-on", state: STATE_BLURAY})
	reg(lirc.RC_PROTO_SONY20, 0x1ae22f, irKey{name: "bluray-off"})
	reg(lirc.RC_PROTO_NECX, 0xd26ccb, irKey{name: "onkyo-power-toggle", state: STATE_MASTER_OFF})
	reg(lirc.RC_PROTO_NECX, 0xd26d03, irKey{name: "onkyo-volume-down"})
	reg(lirc.RC_PROTO_NECX, 0xd26d02, irKey{name: "onkyo-volume-up"})
	reg(lirc.RC_PROTO_NECX, 0xd26c8c, irKey{name: "onkyo-input-bd-dvd", state: STATE_BLURAY})
	reg(lirc.RC_PROTO_NECX, 0xd26d0e, irKey{name: "onkyo-input-cbl-sat"})
	reg(lirc.RC_PROTO_NECX, 0xd26d0c, irKey{name: "onkyo-input-strm-box", state: STATE_SHIELD})
	reg(lirc.RC_PROTO_NECX, 0xd26d9c, irKey{name: "onkyo-input-pc", state: STATE_PC})
	reg(lirc.RC_PROTO_NECX, 0xd26d0d, irKey{name: "onkyo-input-game", state: STATE_GAME})
	reg(lirc.RC_PROTO_NECX, 0xd26d9f, irKey{name: "onkyo-input-aux"})
	reg(lirc.RC_PROTO_NECX, 0xd26d09, irKey{name: "onkyo-input-cd"})
	reg(lirc.RC_PROTO_NECX, 0xd26d0a, irKey{name: "onkyo-input-phono", state: STATE_PORTAL})
	reg(lirc.RC_PROTO_NECX, 0xd26d48, irKey{name: "onkyo-input-tv", state: STATE_LIVE_TV})
	reg(lirc.RC_PROTO_NECX, 0xd26d0b, irKey{name: "onkyo-input-tuner"})
	reg(lirc.RC_PROTO_NEC32, 0xd20287, irKey{name: "onkyo-input-net"})
	reg(lirc.RC_PROTO_NEC32, 0x1ed20e50, irKey{name: "onkyo-input-bluetooth"})
	reg(lirc.RC_PROTO_NEC32, 0x1fd25000, irKey{name: "onkyo-bddvd-up"})
	reg(lirc.RC_PROTO_NEC32, 0x1fd25100, irKey{name: "onkyo-bddvd-down"})
	reg(lirc.RC_PROTO_NEC32, 0x1fd25200, irKey{name: "onkyo-bddvd-left"})
	reg(lirc.RC_PROTO_NEC32, 0x1fd25300, irKey{name: "onkyo-bddvd-right"})
	reg(lirc.RC_PROTO_NEC32, 0x1fd20800, irKey{name: "onkyo-bddvd-enter"})
	reg(lirc.RC_PROTO_NEC32, 0x1fd24d00, irKey{name: "onkyo-bddvd-setup"})
	reg(lirc.RC_PROTO_NEC32, 0x1fd20900, irKey{name: "onkyo-bddvd-return"})

	reg(lirc.RC_PROTO_NECX, 0x860561, irKey{name: "insignia-dvd-on", state: STATE_MENU})
	reg(lirc.RC_PROTO_NECX, 0x860560, irKey{name: "insignia-dvd-off"})
	reg(lirc.RC_PROTO_NECX, 0x860516, irKey{name: "insignia-dvd-left", actions: action{
		"live_tv": tvCommand("rc-left"), "menu": menuCommand("left"), "portal": cecCommand("left")}})
	reg(lirc.RC_PROTO_NECX, 0x860515, irKey{name: "insignia-dvd-right", actions: action{
		"live_tv": tvCommand("rc-right"), "menu": menuCommand("right"), "portal": cecCommand("right")}})
	reg(lirc.RC_PROTO_NECX, 0x860542, irKey{name: "insignia-dvd-up", actions: action{
		"live_tv": tvCommand("rc-up"), "menu": menuCommand("up"), "portal": cecCommand("up")}})
	reg(lirc.RC_PROTO_NECX, 0x860543, irKey{name: "insignia-dvd-down", actions: action{
		"live_tv": tvCommand("rc-down"), "menu": menuCommand("down"), "portal": cecCommand("down")}})
	reg(lirc.RC_PROTO_NECX, 0x860518, irKey{name: "insignia-dvd-select", actions: action{
		"live_tv": tvCommand("rc-select"), "menu": menuCommand("select"), "portal": cecCommand("select")}})
	reg(lirc.RC_PROTO_NECX, 0x86057b, irKey{name: "insignia-dvd-back", actions: action{
		"live_tv": tvCommand("rc-back"), "menu": menuCommand("back"), "portal": cecCommand("back")}})
	reg(lirc.RC_PROTO_NECX, 0x8605b9, irKey{name: "insignia-dvd-home", actions: action{
		"live_tv": tvCommand("rc-apps"), "menu": menuCommand("home"), "portal": cecCommand("home")}})
	reg(lirc.RC_PROTO_NECX, 0x86051b, irKey{name: "insignia-dvd-cancel", actions: action{"live_tv": tvCommand("rc-exit")}})
	reg(lirc.RC_PROTO_NECX, 0x86050a, irKey{name: "insignia-dvd-channel-up", actions: action{"live_tv": tvCommand("rc-channel-up")}})
	reg(lirc.RC_PROTO_NECX, 0x86050b, irKey{name: "insignia-dvd-channel-down", actions: action{"live_tv": tvCommand("rc-channel-down")}})
	reg(lirc.RC_PROTO_NECX, 0x860574, irKey{name: "insignia-dvd-list", actions: action{"live_tv": tvCommand("rc-list")}})
	reg(lirc.RC_PROTO_NECX, 0x860570, irKey{name: "insignia-dvd-play", actions: action{"live_tv": tvCommand("rc-play")}})
	reg(lirc.RC_PROTO_NECX, 0x86057d, irKey{name: "insignia-dvd-fast-forward", actions: action{"live_tv": tvCommand("rc-ff")}})
	reg(lirc.RC_PROTO_NECX, 0x860581, irKey{name: "insignia-dvd-pause", actions: action{"live_tv": tvCommand("rc-pause")}})
	reg(lirc.RC_PROTO_NECX, 0x86057e, irKey{name: "insignia-dvd-rewind", actions: action{"live_tv": tvCommand("rc-rewind")}})
	reg(lirc.RC_PROTO_NECX, 0x860575, irKey{name: "insignia-dvd-settings", actions: action{"live_tv": tvCommand("rc-menu")}})
	reg(lirc.RC_PROTO_NECX, 0x860571, irKey{name: "insignia-dvd-stop", actions: action{"live_tv": tvCommand("rc-stop")}})
	reg(lirc.RC_PROTO_NECX, 0x860545, irKey{name: "insignia-dvd-bookmarks", actions: action{"live_tv": tvCommand("rc-subtitle")}})
	reg(lirc.RC_PROTO_NECX, 0x8605b7, irKey{name: "insignia-dvd-record", actions: action{"live_tv": tvCommand("rc-record")}})
	reg(lirc.RC_PROTO_NECX, 0x860517, irKey{name: "insignia-dvd-info", actions: action{"live_tv": tvCommand("rc-info")}})
	reg(lirc.RC_PROTO_NECX, 0x860500, irKey{name: "insignia-dvd-1", actions: action{"live_tv": tvCommand("rc-num-1")}})
	reg(lirc.RC_PROTO_NECX, 0x860501, irKey{name: "insignia-dvd-2", actions: action{"live_tv": tvCommand("rc-num-2")}})
	reg(lirc.RC_PROTO_NECX, 0x860502, irKey{name: "insignia-dvd-3", actions: action{"live_tv": tvCommand("rc-num-3")}})
	reg(lirc.RC_PROTO_NECX, 0x860503, irKey{name: "insignia-dvd-4", actions: action{"live_tv": tvCommand("rc-num-4")}})
	reg(lirc.RC_PROTO_NECX, 0x860504, irKey{name: "insignia-dvd-5", actions: action{"live_tv": tvCommand("rc-num-5")}})
	reg(lirc.RC_PROTO_NECX, 0x860505, irKey{name: "insignia-dvd-6", actions: action{"live_tv": tvCommand("rc-num-6")}})
	reg(lirc.RC_PROTO_NECX, 0x860506, irKey{name: "insignia-dvd-7", actions: action{"live_tv": tvCommand("rc-num-7")}})
	reg(lirc.RC_PROTO_NECX, 0x860507, irKey{name: "insignia-dvd-8", actions: action{"live_tv": tvCommand("rc-num-8")}})
	reg(lirc.RC_PROTO_NECX, 0x860508, irKey{name: "insignia-dvd-9", actions: action{"live_tv": tvCommand("rc-num-9")}})
	reg(lirc.RC_PROTO_NECX, 0x860509, irKey{name: "insignia-dvd-0", actions: action{"live_tv": tvCommand("rc-num-0")}})
	//missing reg(lirc.RC_PROTO_NECX, 0x860517, irKey{name: "insignia-dvd-minus", tv: "rc-minus" })
}

func irThread(wakeMain chan any, wake chan any) {
	var curState string

	dev, err := lirc.Open(global.rcUevent, global.rcProtos)
	if err != nil {
		msg.Fatal("ir: failed to open device matching \"%s\" (%v)", global.rcUevent, err)
	}

	registerKeys(dev)

	if err := dev.StartScancodeReader(); err != nil {
		msg.Fatal("ir: failed to start scancode reader (%v)", err)
	}

	if err := dev.StartScancodeTransmitter(); err != nil {
		msg.Fatal("ir: failed to start scancode reader (%v)", err)
	}

	pauseTimer := time.NewTimer(time.Duration(0))
	pauseTimer.Stop()
	txPaused := false
	canWrite := false
	receiverIsOn := false

	cmds := list.New()

	msg.Info("ir: started")

	for {
		keepProcessing := true
		for cmds.Front() != nil && keepProcessing {
			elem := cmds.Front()
			switch item := elem.Value.(type) {
			case txPause:
				msg.Verbose2("ir: wait for %s", item.duration)
				pauseTimer.Reset(item.duration)
				txPaused = true
			case string:
				scanCode, found := txTable[item]
				// if we have a known key, submit it to the driver then wait
				// for event to wake us up when we can safely transmit another
				keepProcessing = !found
				if found && !txPaused && canWrite {
					msg.Verbose("ir: submit \"%s\"", item)
					cmds.Remove(elem)
					dev.Tx <- scanCode
					canWrite = false
					if item == "onkyo-on" {
						if !receiverIsOn {
							cmds.PushFront(txPause{duration: receiverOnTime})
						}
						receiverIsOn = true
					} else if item == "onkyo-off" {
						if receiverIsOn {
							cmds.PushFront(txPause{duration: receiverOffTime})
						}
						receiverIsOn = false
					}
				}
			default:
				msg.Fatal("ir: unknown cmd %v", item)
			}
			if keepProcessing {
				cmds.Remove(elem)
			}
		}

		select {
		case event := <-dev.Rx:
			switch item := event.(type) {
			case lirc.ScancodeReaderEvent:
				key := item.User.(irKey)
				//TODO: can I use item.Scancode.Timestamp (nanosecond
				//CLOCK_MONOTONIC) to do the timer reset? maybe I can save the
				//first timestamp and use that to determine what the timeout
				//should be updated to.
				msg.Verbose("ir: got key %v", key.name)
				if key.state != nil {
					wakeMain <- key.state
				}
				if key.crashProtect {
					wakeMain <- crashProtect{}
				}
				if (item.Scancode.Flags & lirc.LIRC_SCANCODE_FLAG_REPEAT) == 0 {
					cmd, ok := key.actions[curState]
					if ok {
						switch cmd.(type) {
						case tvCommand:
							wakeMain <- cmd
						case cecCommand:
							wakeMain <- cmd
						case menuCommand:
							wakeMain <- cmd
						}
					}
				}
			case error:
				msg.Error("ir: scancode reader failed (%v)", item)
				break
			}
		case wakeCommand := <-wake:
			switch item := wakeCommand.(type) {
			case *state:
				msg.Verbose("ir: new commands %v", item.irCommands)
				cmds = list.New()
				for _, cmd := range item.irCommands {
					cmds.PushBack(cmd)
				}
				curState = item.name
			case irCommand:
				msg.Verbose2("ir: add command %v", item)
				cmds.PushBack(string(item))
				cmds.PushBack(txPause{duration: irKeyCooldownTime})
			default:
				msg.Fatal("ir: unknown wake %v", item)
			}
		case <-dev.TxReady:
			canWrite = true
		case <-pauseTimer.C:
			msg.Verbose("ir: done waiting")
			txPaused = false
		}
	}
}

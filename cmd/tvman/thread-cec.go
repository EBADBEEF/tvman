package main

import (
	"container/list"
	"time"
	"tvman/cec"
	"tvman/msg"
	"tvman/tristate"
)

// TODO: it would be nice to move this into the cec package and expose an API
// TODO: how to detect if logical address changes?
var cecDeviceStatus [16]struct {
	Present      tristate.Value
	Power        tristate.Value
	PhysicalAddr uint16
}

var cecOsdName = map[string]uint8{
	"playback1":    4,
	"audio-system": 5,
	"playback2":    8,
	"playback3":    11,
}

// Save the last turned on device, used for idle check and ui cmds. If the
// device turns off we should broadcast master off.
var cecActiveDevice struct {
	valid   bool
	osdName string
	address uint8
}

const cecWaitForScan = time.Duration(2 * time.Second)
const cecIdleCheckDuration = time.Duration(10 * time.Minute)

type cecHadSleep struct{}
type cecIdleTimerReset struct{}
type cecRetry struct {
	command string
}

func cecSetActiveDevice(osdName string, address uint8) {
	cecActiveDevice.valid = true
	cecActiveDevice.osdName = osdName
	cecActiveDevice.address = address
}

func cecCommandToMessage(buf *list.List, command string, isRetry bool) {
	const retryAllowed = true
	cm := cec.CecMsg{}
	from := uint8(1)
	to := uint8(0xff)
	found := false
	switch command {
	case "poll":
		for i := 14; i >= 0; i-- {
			buf.PushFront(&cec.CecMsg{Len: 1, Message: cec.CecMsgMessage{uint8(i<<4 | i)}})
		}
		return
	case "scan":
		buf.PushFront(txPause{duration: cecWaitForScan})
		for i := uint8(0); i <= 14; i++ {
			if i == from {
				continue
			}
			cm := cec.Aliases["give-osd-name"]
			//cm.Reply = 1
			cm.Message[0] = from<<4 | (14 - i)
			buf.PushFront(&cm)
		}
		return
	case "playback-off":
		buf.PushFront("portal-off")
		buf.PushFront("shield-off")
		buf.PushFront("bluray-off")
		cecActiveDevice.valid = false
		return
	case "audio-system-on":
		to, found = cecOsdName["audio-system"]
		if cecDeviceStatus[to].Power.IsUnset() && !isRetry {
			buf.PushFront(cecRetry{command: command})
			cm := cec.Aliases["give-device-power-status"]
			cm.Reply = 1 // tx waits for reply
			cm.Message[0] = uint8(from<<4 | to)
			buf.PushFront(&cm)
			return
		} else {
			if cecDeviceStatus[to].Power.IsOff() {
				buf.PushFront(txPause{duration: receiverOnTime})
			}
			cecDeviceStatus[to].Power.Set(true)
			cm = cec.Aliases["ui-cmd-power-on"]
		}
	case "audio-system-off":
		to, found = cecOsdName["audio-system"]
		cm = cec.Aliases["standby"]
		cecDeviceStatus[to].Power.Set(false)
	case "audio-system-volume-up":
		to, found = cecOsdName["audio-system"]
		cm = cec.Aliases["ui-cmd-volume-up"]
	case "audio-system-volume-down":
		to, found = cecOsdName["audio-system"]
		cm = cec.Aliases["ui-cmd-volume-down"]
	case "audio-system-input-bd-dvd":
		to, found = cecOsdName["audio-system"]
		cm = cec.Aliases["av-input-1"]
	case "audio-system-input-cbl-sat":
		to, found = cecOsdName["audio-system"]
		cm = cec.Aliases["av-input-2"]
	case "audio-system-input-strm-box":
		to, found = cecOsdName["audio-system"]
		cm = cec.Aliases["av-input-3"]
	case "audio-system-input-pc":
		to, found = cecOsdName["audio-system"]
		cm = cec.Aliases["av-input-4"]
	case "audio-system-input-game":
		to, found = cecOsdName["audio-system"]
		cm = cec.Aliases["av-input-5"]
	case "audio-system-input-aux":
		to, found = cecOsdName["audio-system"]
		cm = cec.Aliases["av-input-6"]
	case "audio-system-input-tv-cd":
		to, found = cecOsdName["audio-system"]
		cm = cec.Aliases["av-input-7"]
	case "audio-system-input-phono":
		to, found = cecOsdName["audio-system"]
		cm = cec.Aliases["av-input-8"]
	case "audio-system-input-tv":
		to, found = cecOsdName["audio-system"]
		cm = cec.Aliases["av-input-9"]
	case "audio-system-input-tuner":
		to, found = cecOsdName["audio-system"]
		cm = cec.Aliases["av-input-10"]
	case "audio-system-input-net":
		to, found = cecOsdName["audio-system"]
		cm = cec.Aliases["av-input-11"]
	case "audio-system-input-bluetooth":
		to, found = cecOsdName["audio-system"]
		cm = cec.Aliases["av-input-12"]

	case "shield-on":
		buf.PushFront("portal-off")
		buf.PushFront("bluray-off")
		to, found = cecOsdName[STATE_SHIELD.osdName]
		if found {
			cecSetActiveDevice(STATE_SHIELD.osdName, to)
			buf.PushFront(cecIdleTimerReset{})
		}
		cm = cec.Aliases["ui-cmd-power-on"]
	case "bluray-on":
		buf.PushFront("portal-off")
		buf.PushFront("shield-off")
		to, found = cecOsdName[STATE_BLURAY.osdName]
		if found {
			cecSetActiveDevice(STATE_BLURAY.osdName, to)
			buf.PushFront(cecIdleTimerReset{})
		}
		cm = cec.Aliases["ui-cmd-power-on"]
	case "portal-on":
		buf.PushFront("bluray-off")
		buf.PushFront("shield-off")
		to, found = cecOsdName[STATE_PORTAL.osdName]
		if found {
			cecSetActiveDevice(STATE_PORTAL.osdName, to)
			buf.PushFront(cecIdleTimerReset{})
		}
		cm = cec.Aliases["ui-cmd-power-on"]

	case "shield-off":
		to, found = cecOsdName[STATE_SHIELD.osdName]
		cm = cec.Aliases["standby"]
		if to > 0 {
			cecDeviceStatus[to].Power.Set(false)
		}
	case "bluray-off":
		to, found = cecOsdName[STATE_BLURAY.osdName]
		cm = cec.Aliases["standby"]
		if to > 0 {
			cecDeviceStatus[to].Power.Set(false)
		}
	case "portal-off":
		to, found = cecOsdName[STATE_PORTAL.osdName]
		cm = cec.Aliases["standby"]
		if to > 0 {
			cecDeviceStatus[to].Power.Set(false)
		}
	case "idlecheck":
		if !cecActiveDevice.valid {
			return
		}
		found = true
		to = cecActiveDevice.address
		cm = cec.Aliases["give-device-power-status"]
	case "up":
		if !cecActiveDevice.valid {
			return
		}
		found = true
		to = cecActiveDevice.address
		cm = cec.Aliases["ui-cmd-up"]
	case "down":
		if !cecActiveDevice.valid {
			return
		}
		found = true
		to = cecActiveDevice.address
		cm = cec.Aliases["ui-cmd-down"]
	case "left":
		if !cecActiveDevice.valid {
			return
		}
		found = true
		to = cecActiveDevice.address
		cm = cec.Aliases["ui-cmd-left"]
	case "right":
		if !cecActiveDevice.valid {
			return
		}
		found = true
		to = cecActiveDevice.address
		cm = cec.Aliases["ui-cmd-right"]
	case "select":
		if !cecActiveDevice.valid {
			return
		}
		found = true
		to = cecActiveDevice.address
		cm = cec.Aliases["ui-cmd-select"]
	case "back":
		if !cecActiveDevice.valid {
			return
		}
		found = true
		to = cecActiveDevice.address
		cm = cec.Aliases["ui-cmd-back"]
	case "home":
		if !cecActiveDevice.valid {
			return
		}
		found = true
		to = cecActiveDevice.address
		cm = cec.Aliases["ui-cmd-root-menu"]
	default:
		return
	}
	if !found {
		if retryAllowed && !isRetry {
			msg.Verbose("retrying command after scan...")
			buf.PushFront(cecRetry{command: command})
			buf.PushFront("scan")
		}
		return
	}
	cm.Message[0] = uint8(from<<4 | to)
	buf.PushFront(&cm)
}

func cecReceive(buf *list.List, cm *cec.CecMsg) {
	from := (cm.Message[0] >> 4) & 0xf
	to := (cm.Message[0] >> 0) & 0xf
	if msg.Level >= msg.LvlVerbose2 || cm.Len > 1 {
		// TODO: check if error and choose to retry it? maybe I can wait until
		// response from sequence number assigned during tx?
		msg.Verbose("cec: %s", cm)
	}
	if cm.TxTimestamp != 0 {
		if (cm.TxStatus & cec.CEC_TX_STATUS_OK) != 0 {
			cecDeviceStatus[to].Present.Set(true)
		} else if (cm.TxStatus & cec.CEC_TX_STATUS_NACK) != 0 {
			cecDeviceStatus[to].Present.Set(false)
			cecDeviceStatus[to].Power.Set(false)
		}
	} else if cm.RxTimestamp != 0 {
		cecDeviceStatus[from].Present.Set(true)
		if cm.Len < 2 {
			return
		}
		switch cm.Message[1] {
		case cec.CEC_MSG_SET_OSD_NAME:
			osdName := string(cm.Message[2:cm.Len])
			cecOsdName[osdName] = from
			if cecActiveDevice.valid && cecActiveDevice.osdName == osdName {
				msg.Warning("cec: active device changed address from %v to %v",
					cecActiveDevice.address, from)
				cecActiveDevice.address = from
			}
		case cec.CEC_MSG_REPORT_POWER_STATUS:
			switch cm.Message[2] {
			case cec.CEC_OP_POWER_STATUS_ON:
				fallthrough
			case cec.CEC_OP_POWER_STATUS_TO_ON:
				cecDeviceStatus[from].Power.Set(true)

			case cec.CEC_OP_POWER_STATUS_STANDBY:
				fallthrough
			case cec.CEC_OP_POWER_STATUS_TO_STANDBY:
				cecDeviceStatus[from].Power.Set(false)
				msg.Info("cec: device %v is in standby", from)
				if cecActiveDevice.valid {
					if cecActiveDevice.address == from {
						msg.Info("cec: active device went to sleep, time to die (osdName=%v, addr=%v)",
							cecActiveDevice.osdName, cecActiveDevice.address)
						buf.PushFront(cecHadSleep{})
						cecActiveDevice.valid = false
					}
				}
			}
		}
		//TODO: is there a way to respond to shields request for TV power
		//status? I don't think there is a way to claim logical address 0 from
		//cec-ctl.
	}
}

func cecThread(wakeMain chan any, wake chan any) {

	writer, err := cec.Open(global.cecDevice, cec.CEC_MODE_INITIATOR|cec.CEC_MODE_FOLLOWER)
	if err != nil {
		msg.Fatal("cec: failed to open writer (%v)", err)
	}

	monitor, err := cec.Open(global.cecDevice, cec.CEC_MODE_MONITOR_ALL)
	if err != nil {
		msg.Fatal("cec: failed to open monitor (%v)", err)
	}

	if err := monitor.StartReceiver(); err != nil {
		msg.Fatal("cec: could not start receiver (%v)", err)
	}

	if err := writer.StartTransmitter(); err != nil {
		msg.Fatal("cec: could not start transmitter (%v)", err)
	}

	if msg.Level >= msg.LvlVerbose {
		cec.SetVerbose(true)
	}

	canWrite := false

	pauseTimer := time.NewTimer(time.Duration(0))
	pauseTimer.Stop()
	idleTimer := time.NewTimer(time.Duration(0))
	idleTimer.Stop()
	txPaused := false

	cmds := list.New()
	cmds.PushBack("scan")

	msg.Info("cec: started")
	for {
		keepProcessing := true
		msg.Verbose2("cec: loop canWrite=%v", canWrite)

		for cmds.Front() != nil && keepProcessing {

			elem := cmds.Front()
			switch item := elem.Value.(type) {
			case txPause:
				msg.Verbose("cec: wait for %s", item.duration)
				pauseTimer.Reset(item.duration)
				txPaused = false
			case cecRetry:
				msg.Verbose("cec: retry command \"%s\"", item.command)
				cecCommandToMessage(cmds, item.command, true)
			case string:
				msg.Verbose("cec: command \"%s\"", item)
				cecCommandToMessage(cmds, item, false)
			case *cec.CecMsg:
				// only send one request to the device at a time
				keepProcessing = false
				if !txPaused && canWrite {
					msg.Verbose2("cec: submit %s", item)
					writer.Tx <- item
					cmds.Remove(elem)
					canWrite = false
				}
			case cecHadSleep:
				wakeMain <- STATE_MASTER_OFF
			case cecIdleTimerReset:
				idleTimer.Reset(cecIdleCheckDuration)
			}

			if keepProcessing {
				cmds.Remove(elem)
			}
		}

		select {
		case cm := <-monitor.Rx:
			cecReceive(cmds, cm)
		case ev := <-monitor.Event:
			msg.Verbose2("cec: event %#v", ev)
		case <-writer.TxReady:
			canWrite = true
		case wakeCommand := <-wake:
			switch item := wakeCommand.(type) {
			case *state:
				msg.Verbose("cec: new commands %v", item.cecCommands)
				cmds = list.New()
				for _, cmd := range item.cecCommands {
					cmds.PushBack(cmd)
				}
			case cecCommand:
				cmds.PushBack(string(item))
			default:
				msg.Error("cec: unexpected wake (%v)", item)
			}
		case err := <-monitor.Error:
			msg.Error("cec: monitor error (%v)", err)
		case err := <-writer.Error:
			msg.Error("cec: writer error (%v)", err)
		case <-pauseTimer.C:
			msg.Verbose("cec: done waiting")
			txPaused = false
		case <-idleTimer.C:
			if cecActiveDevice.valid {
				cmds.PushBack("idlecheck")
				idleTimer.Reset(cecIdleCheckDuration)
			}
		}
	}
}

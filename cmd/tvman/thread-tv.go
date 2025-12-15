package main

import (
	"container/list"
	"time"
	"tvman/lgtv"
	"tvman/msg"
)

func tvThread(wakeMain chan any, wake chan any) {
	dev, err := lgtv.Open(global.ttyTV)
	if err != nil {
		msg.Fatal("tv: could not open device (%v)", err)
	}

	if err := dev.StartReceiver(); err != nil {
		msg.Fatal("tv: could not start receiver (%v)", err)
	}

	if err := dev.StartTransmitter(); err != nil {
		msg.Fatal("tv: could not start receiver (%v)", err)
	}

	const waitForRx = time.Duration(2500 * time.Millisecond)

	tvIsOn := false
	canWrite := false

	pauseTimer := time.NewTimer(time.Duration(0))
	pauseTimer.Stop()
	txPaused := false

	rxTimer := time.NewTimer(waitForRx)
	rxDone := false

	cmds := list.New()
	cmds.PushBack("check")

	msg.Info("tv: started")
	for {
		keepProcessing := true
		msg.Verbose2("tv: loop canWrite=%v", canWrite)

		for cmds.Front() != nil && keepProcessing {
			elem := cmds.Front()
			switch item := elem.Value.(type) {
			case txPause:
				msg.Verbose2("tv: wait for %s", item.duration)
				pauseTimer.Reset(item.duration)
				txPaused = true
			case string:
				keepProcessing = false
				if !txPaused && canWrite && rxDone {
					msg.Verbose("tv: submit \"%s\"", item)
					dev.Tx <- item
					cmds.Remove(elem)
					canWrite = false
					rxDone = false
					rxTimer.Reset(waitForRx)
					if item == "on" {
						if !tvIsOn {
							cmds.PushFront(txPause{duration: tvOnTime})
						}
						tvIsOn = true
					} else if item == "off" {
						tvIsOn = false
					}
				}
			default:
				msg.Fatal("tv: unknown cmd %v", item)
			}
			if keepProcessing {
				cmds.Remove(elem)
			}
		}

		select {
		case wakeCommand := <-wake:
			switch item := wakeCommand.(type) {
			case *state:
				msg.Verbose("tv: new commands %v", item.tvCommands)
				cmds = list.New()
				for _, cmd := range item.tvCommands {
					cmds.PushBack(cmd)
				}
			case crashProtect:
				msg.Verbose("tv: enable crash protect")
				dev.CrashProtect()
			case tvCommand:
				msg.Verbose("tv: command %v", item)
				cmds.PushBack(string(item))
			default:
				msg.Fatal("tv: unknown wake %v", item)
			}
		case err := <-dev.Error:
			if err != lgtv.ProtectedError || msg.Level >= msg.LvlVerbose {
				msg.Error("tv: error %v", err)
			}
		case <-dev.TxReady:
			canWrite = true
		case buf := <-dev.Rx:
			if len(buf) == 10 {
				msg.Verbose("tv: rx %s", buf)
				rxDone = true
				rxTimer.Stop()
			} else {
				msg.Verbose2("tv: rx junk %d %s", len(buf), buf)
			}
			if buf == "a 01 OK01x" {
				msg.Info("tv: detected on")
				tvIsOn = true
			} else if buf == "a 01 OK00x" {
				msg.Info("tv: detected off")
				tvIsOn = false
			}
		case <-pauseTimer.C:
			msg.Verbose("tv: done waiting")
			txPaused = false
		case <-rxTimer.C:
			msg.Verbose("tv: rx timeout")
			rxDone = true
		}
	}
}

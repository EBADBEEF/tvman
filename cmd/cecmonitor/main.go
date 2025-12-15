package main

import (
	"flag"
	"fmt"
	"os"
	"tvman/cec"
)

func envOrDefault(e string, def string) string {
	if env := os.Getenv(e); env != "" {
		return env
	}
	return def
}

var cecDevice string

func main() {
	flag.StringVar(&cecDevice, "device", envOrDefault("CECMONITOR_DEVICE", "/dev/cec0"),
		"path to tty device")

	monitor, err := cec.Open(cecDevice, cec.CEC_MODE_MONITOR_ALL)
	if err != nil {
		panic(fmt.Sprintf("cec: failed to open monitor (%v)", err))
	}

	if err := monitor.StartReceiver(); err != nil {
		panic(fmt.Sprintf("cec: could not start receiver (%v)", err))
	}

	cec.SetVerbose(true)

	for {
		select {
		case cm := <-monitor.Rx:
			ts := cm.TxTimestamp | cm.RxTimestamp //only one will be set
			fmt.Printf("%.3fms: %s\n", float64(ts)/1000000000, cm)
		//case ev := <-monitor.Event:
		//	fmt.Printf("device event: %#v\n", ev)
		case err := <-monitor.Error:
			fmt.Printf("error: monitor (%v)", err)
		}
	}
}

// run with: tvman -rc ” -rcprotos 'nec sony' -verbose server
package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"tvman/listen"
	"tvman/msg"
)

var global struct {
	dryRun             bool
	rcUevent           string
	rcProtos           string
	cecDevice          string
	ttyTV              string
	url                string
	menumanUrl         string
	logTimestamp       bool
	enableCrashProtect bool
}

const (
	chanBufferCount = 10
)

type txPause struct {
	duration time.Duration
}

type idleCheck struct{}
type crashProtect struct{}
type tvCommand string
type cecCommand string
type irCommand string
type menuCommand string

func envOrDefault(e string, def string) string {
	if env := os.Getenv(e); env != "" {
		return env
	}
	return def
}

//common_headers = b"Cache-Control: no-cache, no-store, must-revalidate\r\nPragma: no-cache\r\nExpires: 0\r\nConnection: close\r\n"
//response = b"HTTP/1.1 307\r\nLocation: ../\r\n" + common_headers + b"\r\n"

func httpThread(c chan any) {
	var httpCommands = map[string]any{
		"off":          STATE_MASTER_OFF,
		"off-tv":       tvCommand("blank"),
		"on-strmbox":   STATE_SHIELD,
		"on-pc":        STATE_PC,
		"master_off":   STATE_MASTER_OFF,
		"live_tv":      STATE_LIVE_TV,
		"shield":       STATE_SHIELD,
		"bluray":       STATE_BLURAY,
		"portal":       STATE_PORTAL,
		"pc":           STATE_PC,
		"retropie":     STATE_RETROPIE,
		"retropie_off": STATE_RETROPIE_OFF,
		"game":         STATE_GAME,
		"menu":         STATE_MENU,
		"volume": func(v string) {
			num, err := strconv.Atoi(v)
			if err != nil || num < -10 || num > 10 {
				msg.Error("http: bad volume")
				return
			}
			incr := -1
			cmd := irCommand("onkyo-volume-up")
			if num < 0 {
				incr = 1
				cmd = irCommand("onkyo-volume-down")
			}
			for num != 0 {
				num += incr
				c <- cmd
			}
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, param := range strings.Split(r.URL.RawQuery, "&") {
			split := strings.Split(param, "=")
			name := ""
			val := ""
			if len(split) > 0 {
				name = split[0]
			}
			if len(split) > 1 {
				val = split[1]
			}
			switch item := httpCommands[name].(type) {
			case tvCommand:
				c <- item
			case cecCommand:
				c <- item
			case *state:
				c <- item
			case func(string):
				item(val)
			}
		}
	})

	if listener, err := listen.Listen(global.url); err != nil {
		msg.Fatal("listener failed: %v", err)
	} else {
		defer listener.Close()
		msg.Info("server listening on: %s", global.url)
		if err := http.Serve(listener, nil); err != nil {
			msg.Fatal("http server failed: %v ", err)
		}
	}
}

type state struct {
	name         string
	osdName      string
	cecCommands  []string
	irCommands   []string
	tvCommands   []string
	menuCommands []string
	mustBeIn     *state
}

const (
	tvOnTime          = time.Duration(5000 * time.Millisecond)
	receiverOnTime    = time.Duration(5000 * time.Millisecond)
	receiverOffTime   = time.Duration(2000 * time.Millisecond)
	irKeyCooldownTime = time.Duration(250 * time.Millisecond)
)

var (
	STATE_NONE       = &state{name: "none"}
	STATE_MASTER_OFF = &state{
		name:         "master_off",
		cecCommands:  []string{"playback-off"},
		irCommands:   []string{"onkyo-off"},
		tvCommands:   []string{"off"},
		menuCommands: []string{"hide"},
	}
	STATE_LIVE_TV = &state{
		name:         "live_tv",
		cecCommands:  []string{"playback-off"},
		irCommands:   []string{"onkyo-on", "onkyo-input-tv"},
		tvCommands:   []string{"on", "input-tv"},
		menuCommands: []string{"hide"},
	}
	STATE_SHIELD = &state{
		name:         "shield",
		osdName:      "SHIELD",
		cecCommands:  []string{"shield-on"},
		irCommands:   []string{"onkyo-on", "onkyo-input-strm-box"},
		tvCommands:   []string{"on", "input-hdmi1", "unblank"},
		menuCommands: []string{"hide"},
	}
	STATE_BLURAY = &state{
		name:         "bluray",
		osdName:      "BDP-S3500",
		cecCommands:  []string{"bluray-on"},
		irCommands:   []string{"onkyo-on", "onkyo-input-bd-dvd"},
		tvCommands:   []string{"on", "input-hdmi1", "unblank"},
		menuCommands: []string{"hide"},
	}
	STATE_PORTAL = &state{
		name:         "portal",
		osdName:      "PortalTV",
		cecCommands:  []string{"portal-on"},
		irCommands:   []string{"onkyo-on", "onkyo-input-phono"},
		tvCommands:   []string{"on", "input-hdmi1", "unblank"},
		menuCommands: []string{"hide"},
	}
	STATE_PC = &state{
		name:         "pc",
		cecCommands:  []string{"playback-off"},
		irCommands:   []string{"onkyo-on", "onkyo-input-pc"},
		tvCommands:   []string{"on", "input-hdmi1", "unblank"},
		menuCommands: []string{"hide"},
	}
	STATE_RETROPIE = &state{
		name:         "retropie",
		cecCommands:  []string{"playback-off"},
		irCommands:   []string{"onkyo-on", "onkyo-input-pc"},
		tvCommands:   []string{"on", "input-hdmi2", "unblank"},
		menuCommands: []string{"hide"},
	}
	STATE_RETROPIE_OFF = &state{
		name:         "retropie_off",
		cecCommands:  []string{"playback-off"},
		irCommands:   []string{"onkyo-off"},
		tvCommands:   []string{"off"},
		menuCommands: []string{"hide"},
		mustBeIn:     STATE_RETROPIE,
	}
	STATE_GAME = &state{
		name:         "game",
		cecCommands:  []string{"playback-off"},
		irCommands:   []string{"onkyo-on", "onkyo-input-game"},
		tvCommands:   []string{"on", "input-hdmi2", "unblank"},
		menuCommands: []string{"hide"},
	}
	STATE_MENU = &state{
		name:         "menu",
		cecCommands:  []string{"playback-off"},
		irCommands:   []string{"onkyo-on", "onkyo-input-pc"},
		tvCommands:   []string{"on", "input-hdmi2", "unblank"},
		menuCommands: []string{"show"},
	}
)

func server() {

	// used for threads to talk to us
	wakeMain := make(chan any, chanBufferCount)

	// used for us to talk to threads
	var wakeTV chan any
	var wakeCEC chan any
	var wakeIR chan any
	var wakeMenu chan any

	if global.ttyTV != "" {
		wakeTV = make(chan any, chanBufferCount)
		go tvThread(wakeMain, wakeTV)
	}
	if global.cecDevice != "" {
		wakeCEC = make(chan any, chanBufferCount)
		go cecThread(wakeMain, wakeCEC)
	}
	if global.rcUevent != "none" {
		wakeIR = make(chan any, chanBufferCount)
		go irThread(wakeMain, wakeIR)
	}
	if global.url != "" {
		go httpThread(wakeMain)
	}
	if global.menumanUrl != "" {
		wakeMenu = make(chan any, chanBufferCount)
		go menuThread(wakeMain, wakeMenu)
	}

	curstate := STATE_NONE
	lastStateChange := time.Time{}
	const debounceState = time.Duration(2 * time.Second)

	timer := time.NewTimer(0)
	timer.Stop()
	const idleCheckTime = time.Duration(15 * time.Minute)

	for {
		msg.Verbose2("state: main loop")
		var event any

		// wait for event
		select {
		case event = <-wakeMain:
		}

		// process event
		switch item := event.(type) {
		case *state:
			if curstate == item && time.Since(lastStateChange) < debounceState {
				msg.Info("state: already in state %s", curstate.name)
			} else if item.mustBeIn != nil && curstate != item.mustBeIn {
				msg.Info("state: ignore from %s to %s", curstate.name, item.name)
			} else {
				msg.Info("state: change from %s to %s", curstate.name, item.name)
				lastStateChange = time.Now()
				if wakeTV != nil {
					wakeTV <- item
				}
				if wakeCEC != nil {
					wakeCEC <- item
				}
				if wakeIR != nil {
					wakeIR <- item
				}
				if wakeMenu != nil {
					wakeMenu <- item
				}
				curstate = item
				if curstate == STATE_MASTER_OFF {
					timer.Stop()
				} else {
					timer.Reset(idleCheckTime)
				}
			}
		case cecCommand:
			msg.Verbose2("cmd: cec %v", item)
			wakeCEC <- item
		case irCommand:
			msg.Verbose2("cmd: ir %v", item)
			wakeIR <- item
		case tvCommand:
			msg.Verbose2("cmd: tv %v", item)
			wakeTV <- item
		case menuCommand:
			msg.Verbose2("cmd: menu %v", item)
			wakeMenu <- item
		case crashProtect:
			if global.enableCrashProtect {
				wakeTV <- item
			}
		}
	}
}

func client(args []string) {
	msg.Verbose("args = %v", args)
	url := global.url
	transport := &http.Transport{}

	// TODO: refactor this out because it is in three places
	if strings.HasPrefix(url, "unix://") {
		path := url[7:]
		msg.Info("unix path = %v", path)
		transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			return net.Dial("unix", path)
		}
		// fake url because otherwise the client gets confused
		url = "http://localhost/unix/" + path
	}

	client := &http.Client{Transport: transport}

	for _, arg := range args {
		resp, err := client.Get(url + "/?" + arg)
		if err != nil {
			msg.Error("client failed (%v)", err)
		} else {
			resp.Body.Close()
		}
	}
}

func main() {
	var verbose = false
	global.enableCrashProtect = false

	flag.CommandLine.SetOutput(os.Stdout)

	flag.BoolVar(&verbose, "verbose", os.Getenv("TVMAN_VERBOSE") != "",
		"verbose logging")
	flag.BoolVar(&global.dryRun, "dryrun", os.Getenv("TVMAN_DRYRUN") != "",
		"pretend mode")
	flag.StringVar(&global.rcUevent, "rc", envOrDefault("TVMAN_RCUEVENT", ""),
		"string to match /sys/class/rc/rc*/uevent, empty detects the first one with lirc")
	flag.StringVar(&global.ttyTV, "ttytv", envOrDefault("TVMAN_TTYTV", "/dev/ttyTV"),
		"path to tty device")
	flag.StringVar(&global.cecDevice, "cec", envOrDefault("TVMAN_CECDEVICE", "/dev/cec0"),
		"path to cec device")
	flag.StringVar(&global.rcProtos, "rcprotos", envOrDefault("TVMAN_RCPROTOS", "nec sony"),
		"space-separated list of rc decode protocols")
	flag.StringVar(&global.menumanUrl, "menuman", envOrDefault("MENUMAN_URL", ""),
		"url of menuman (optional)")
	flag.StringVar(&global.url, "url", envOrDefault("TVMAN_URL", "http://localhost:8181"),
		"HTTP url to serve on (server) or connect to (client).\n"+
			"The server prefers to use LISTEN_FDS provided by systemd.\n"+
			"The client can use unix:// as a scheme.\n")

	flag.Usage = func() {
		exe := filepath.Base(os.Args[0])
		fmt.Printf("Server:\n  %s [options] server\n", exe)
		fmt.Printf("Client:\n  %s [options] client [commands...]\n", exe)
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if verbose {
		msg.Level = msg.LvlVerbose
	}

	switch {
	case len(flag.Args()) == 0:
		flag.Usage()
	case flag.Args()[0] == "server":
		server()
	case flag.Args()[0] == "client":
		client(flag.Args()[1:])
	}
}

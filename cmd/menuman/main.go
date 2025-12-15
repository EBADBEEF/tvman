package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"tvman/listen"
	"tvman/msg"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var global struct {
	fontPath     string
	baseFontSize int
	renderWidth  int
	renderHeight int
	window       bool
	url          string
	tvman        string
	show         bool
	idleTime     time.Duration
}

type engine struct {
	window *sdl.Window
	font   *ttf.Font
	render *sdl.Renderer
	mode   sdl.DisplayMode
}

type state struct {
	menu         []MenuItem
	top          int
	bottom       int
	selected     int
	redraw       bool
	running      bool
	tvman        chan any
	inactive     bool
	lastActivity time.Time
}

type MenuItem struct {
	Label    string     `json:"label"`
	Children []MenuItem `json:"children,omitempty"`
	Tvman    string
}

var Menu = []MenuItem{
	{Label: "Live TV", Tvman: "live_tv"},
	{Label: "Portal", Tvman: "portal"},
	{Label: "Nintendo Switch", Tvman: "game"},
	{Label: "Off", Tvman: "off"},
}

func initEngine() (*engine, error) {
	e := engine{}
	var err error

	if err = ttf.Init(); err != nil {
		msg.Fatal("ttf init failed: %v", err)
	}

	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		msg.Fatal("sdl init failed: %v", err)
	}

	if e.mode, err = sdl.GetCurrentDisplayMode(0); err != nil {
		msg.Fatal("GetCurrentDisplayMode failed: %v", err)
	}

	var flags uint32

	flags |= sdl.WINDOW_OPENGL
	flags |= sdl.WINDOW_RESIZABLE
	if !global.window {
		flags |= sdl.WINDOW_FULLSCREEN
	}
	if !global.show {
		flags |= sdl.WINDOW_HIDDEN
	}

	if e.window, err = sdl.CreateWindow("Menu", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, e.mode.W, e.mode.H, flags); err != nil {
		msg.Fatal("CreateWindow failed: %v", err)
	}

	if e.render, err = sdl.CreateRenderer(e.window, -1, sdl.RENDERER_ACCELERATED); err != nil {
		msg.Fatal("CreateRenderer failed: %v", err)
	}

	if err = e.render.SetLogicalSize((int32)(global.renderWidth), (int32)(global.renderHeight)); err != nil {
		msg.Fatal("SetLogicalSize failed: %v", err)
	}

	return &e, nil
}

func closeEngine(e *engine) {
	if e.font != nil {
		e.font.Close()
	}
	if e.window != nil {
		e.window.Destroy()
	}
	sdl.Quit()
	ttf.Quit()
}

func draw(e *engine, s *state) (err error) {
	if e.font != nil {
		e.font.Close()
	}

	if e.font, err = ttf.OpenFont(global.fontPath, global.baseFontSize); err != nil {
		msg.Fatal("OpenFont(%v,%v) failed: %v", global.fontPath, global.baseFontSize, err)
	}

	fontHeight := e.font.Height()

	leftRightOffset := 100
	allowedWidth := global.renderWidth - (2 * leftRightOffset)

	topBottomOffset := 50
	betweenHeight := (int32)(1 * fontHeight)
	maxHeight := (int32)(global.renderHeight - (2 * topBottomOffset))

	e.render.SetDrawColor(0, 0, 0, 0xff)
	e.render.Clear()

	// dont draw anything if nobody will see it
	if s.inactive {
		e.render.Present()
		return
	}

	yOffset := (int32)(topBottomOffset)

	for i, mi := range s.menu[s.top:] {
		var fgcolor sdl.Color
		var bgcolor sdl.Color

		if (s.top + i) == s.selected {
			fgcolor = sdl.Color{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
			bgcolor = sdl.Color{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
		} else {
			fgcolor = sdl.Color{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
			bgcolor = sdl.Color{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
		}

		// TODO: cache textures + surfaces

		var surf *sdl.Surface
		if surf, err = e.font.RenderUTF8BlendedWrapped(mi.Label, fgcolor, allowedWidth); err != nil {
			msg.Fatal("RenderUTF8Blended failed: %v", err)
		}

		var texture *sdl.Texture
		if texture, err = e.render.CreateTextureFromSurface(surf); err != nil {
			msg.Fatal("CreateTextureFromSurface failed: %v", err)
		}

		e.render.SetDrawColor(bgcolor.R, bgcolor.G, bgcolor.B, bgcolor.A)
		rectbg := sdl.Rect{X: (int32)(leftRightOffset), Y: yOffset, W: (int32)(allowedWidth), H: surf.H}
		e.render.FillRect(&rectbg)

		rectfg := sdl.Rect{X: (int32)(leftRightOffset), Y: yOffset, W: surf.W, H: surf.H}
		e.render.Copy(texture, nil, &rectfg)

		s.bottom = s.top + i

		yOffset += surf.H + betweenHeight

		texture.Destroy()
		surf.Free()

		if yOffset > maxHeight {
			break
		}
	}

	e.render.Present()
	msg.Verbose2("selected: %v, top: %v, bottom: %v", s.selected, s.top, s.bottom)

	return nil
}

func menuSelect(s *state) {
	s.lastActivity = time.Now()
	if s.inactive {
		s.inactive = false
		s.redraw = true
	} else {
		msg.Verbose("item selected %d (%v)", s.selected, s.menu[s.selected])
		s.tvman <- s.menu[s.selected].Tvman
	}
}

const (
	MENU_HIDE = iota
	MENU_SHOW
	MENU_BLANK
)

func menuHide(s *state, e *engine, mode int) {
	s.lastActivity = time.Now()
	switch mode {
	case MENU_HIDE:
		if global.show {
			msg.Info("not hiding due to -show option")
		} else {
			e.window.Hide()
			s.inactive = true
			s.redraw = true
		}
	case MENU_SHOW:
		e.window.Show()
		s.inactive = false
		s.redraw = true
	case MENU_BLANK:
		s.inactive = true
		s.redraw = true
	}
}

const (
	ITEM = iota
	PAGE
)

func menuMove(s *state, action int, count int) {
	s.redraw = true
	s.lastActivity = time.Now()
	if s.inactive {
		s.inactive = false
		return
	}
	if count > 0 {
		// select down
		if action == PAGE {
			s.selected = s.bottom
		}
		s.selected += count
		if s.selected >= len(s.menu) {
			s.selected = len(s.menu) - 1
		}
		if s.bottom < 0 {
			msg.Fatal("invalid bottom")
		}
		if s.selected > s.bottom {
			// page down
			s.top = s.bottom + 1
			s.selected = s.top
			s.bottom = -1
		}
	} else if count < 0 {
		// select up
		s.selected += count
		if s.selected < 0 {
			s.selected = 0
		}
		if s.selected < s.top {
			// page up
			s.top -= (s.bottom - s.top + 1)
			s.bottom = -1
			if s.top < 0 {
				s.top = 0
			}
			if action == PAGE {
				s.selected = s.top
			}
		}
	}
}

func httpThread(c chan any) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, param := range strings.Split(r.URL.RawQuery, "&") {
			split := strings.Split(param, "=")
			name := ""
			if len(split) > 0 {
				name = split[0]
			}
			if name != "" {
				c <- name
			}
		}
	})

	if listener, err := listen.Listen(global.url); err != nil {
		msg.Fatal("listener failed: %v", err)
	} else {
		defer listener.Close()
		msg.Info("server listening on: %v", global.url)
		if err := http.Serve(listener, nil); err != nil {
			msg.Fatal("http server failed: %v ", err)
		}
	}
}

func envOrDefault(e string, def string) string {
	if env := os.Getenv(e); env != "" {
		return env
	}
	return def
}

func envOrDefaultInt(e string, def int) int {
	if env := os.Getenv(e); env != "" {
		num, err := strconv.Atoi(env)
		if err != nil {
			msg.Fatal("environment variable %v not an integer", e)
		}
		return num
	}
	return def
}

func tvmanThread(c chan any) {
	url := global.tvman
	transport := &http.Transport{}

	// TODO: refactor this out because it is in three places
	if strings.HasPrefix(url, "unix://") {
		path := url[7:]
		msg.Verbose2("unix path = %v", path)
		transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			return net.Dial("unix", path)
		}
		// fake url because otherwise the client gets confused
		url = "http://localhost/unix/" + path
	}

	client := &http.Client{Transport: transport}

	for {
		select {
		case event := <-c:
			switch arg := event.(type) {
			case string:
				msg.Info("sending tvman event: %v", arg)
				if url != "" {
					request := url + "/?" + arg
					msg.Verbose2("http get: %v", request)
					resp, err := client.Get(request)
					if err != nil {
						msg.Error("client failed (%v)", err)
					} else {
						resp.Body.Close()
					}
				}
			}
		}
	}
}

func main() {
	wakeMain := make(chan any, 10)

	flag.CommandLine.SetOutput(os.Stdout)

	var idleTime string
	var verbose bool

	flag.StringVar(&global.fontPath, "font", envOrDefault("MENUMAN_FONTPATH", ""), "path to ttf file")
	flag.IntVar(&global.baseFontSize, "fontsize", envOrDefaultInt("MENUMAN_FONTSIZE", 32), "size of font")
	flag.IntVar(&global.renderWidth, "w", envOrDefaultInt("MENUMAN_WIDTH", 1920), "render width in pixels")
	flag.IntVar(&global.renderHeight, "h", envOrDefaultInt("MENUMAN_HEIGHT", 1080), "render height in pixels")
	flag.BoolVar(&global.window, "window", envOrDefault("MENUMAN_WINDOW", "1") != "", "do not go fullscreen")
	flag.BoolVar(&global.show, "show", envOrDefault("MENUMAN_SHOW", "") != "", "show menu at start")
	flag.BoolVar(&verbose, "verbose", envOrDefault("MENUMAN_VERBOSE", "") != "", "verbose debug")
	flag.StringVar(&global.url, "url", envOrDefault("MENUMAN_URL", "http://localhost:8182"), "HTTP url to serve on")
	flag.StringVar(&global.tvman, "tvman", envOrDefault("TVMAN_URL", "http://localhost:8181"), "how to reach tvman (http:// or unix:// url)")
	flag.StringVar(&idleTime, "idle", envOrDefault("MENUMAN_IDLE", "30s"), "idle time in before blanking the screen, e.g. 5.5s")
	flag.Parse()

	if verbose {
		msg.Level = msg.LvlVerbose
	}

	if global.fontPath == "" {
		msg.Fatal("missing font path")
	}

	s := state{
		menu:         Menu,
		selected:     0,
		top:          0,
		bottom:       0,
		redraw:       true,
		running:      true,
		tvman:        make(chan any, 10),
		lastActivity: time.Now(),
	}

	if idleTime != "" {
		var err error
		global.idleTime, err = time.ParseDuration(idleTime)
		if err != nil {
			msg.Fatal("could not parse idle time")
		}
	}

	if global.url != "" {
		go httpThread(wakeMain)
	}

	go tvmanThread(s.tvman)

	var e *engine
	var err error

	if e, err = initEngine(); err != nil {
		msg.Fatal("initEngine failed: %v", err)
	}

	for s.running {
		if s.redraw {
			draw(e, &s)
			s.redraw = false
		}
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.WindowEvent:
				we := event.(*sdl.WindowEvent)
				switch we.Event {
				case sdl.WINDOWEVENT_SIZE_CHANGED:
					s.redraw = true
				}
			case *sdl.KeyboardEvent:
				ke := event.(*sdl.KeyboardEvent)
				if ke.Type == sdl.KEYDOWN {
					switch ke.Keysym.Sym {
					case sdl.K_j:
						menuMove(&s, ITEM, 1)
					case sdl.K_DOWN:
						menuMove(&s, ITEM, 1)
					case sdl.K_PAGEDOWN:
						menuMove(&s, PAGE, 1)
					case sdl.K_PAGEUP:
						menuMove(&s, PAGE, -1)
					case sdl.K_k:
						menuMove(&s, ITEM, -1)
					case sdl.K_UP:
						menuMove(&s, ITEM, -1)
					case sdl.K_RETURN:
						menuSelect(&s)
					case sdl.K_q:
						s.running = false
					}
				}
			case *sdl.QuitEvent:
				s.running = false
			default:
				msg.Verbose2("unhandled sdl event: %+v", event)
			}
		}
		select {
		case event := <-wakeMain:
			switch event {
			case "up":
				menuMove(&s, ITEM, -1)
			case "down":
				menuMove(&s, ITEM, 1)
			case "pageup":
				menuMove(&s, PAGE, -1)
			case "pagedown":
				menuMove(&s, PAGE, 1)
			case "select":
				menuSelect(&s)
			case "hide":
				menuHide(&s, e, MENU_HIDE)
			case "show":
				menuHide(&s, e, MENU_SHOW)
			case "blank":
				menuHide(&s, e, MENU_BLANK)
			default:
				msg.Error("unknown command: %v", event)
			}
		default:
			sdl.Delay(16)
		}
		if !s.inactive && time.Since(s.lastActivity) > global.idleTime {
			s.inactive = true
			s.redraw = true
		}
	}

	msg.Verbose("byebye")
	os.Exit(0)
}

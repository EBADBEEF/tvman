package listen

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
)

func Listen(arg string) (listener net.Listener, err error) {
	if listenFds := os.Getenv("LISTEN_FDS"); listenFds != "" {
		file := os.NewFile(3, "listener")
		listener, err = net.FileListener(file)
		if err != nil {
			return nil, fmt.Errorf("fd listener failed: %v", err)
		}
	} else {
		listenUrl, err := url.Parse(arg)
		if err != nil {
			return nil, fmt.Errorf("url parse failed: %v", err)
		}
		switch listenUrl.Scheme {
		case "http":
			hostname := listenUrl.Hostname()
			port := listenUrl.Port()
			if port == "" {
				port = "80"
			}
			addr := hostname + ":" + port
			if listener, err = net.Listen("tcp", addr); err != nil {
				return nil, fmt.Errorf("tcp listen failed: %v", err)
			}
		case "unix":
			var mode os.FileMode = 0660
			if modeStr := listenUrl.Query().Get("mode"); modeStr != "" {
				m, err := strconv.ParseUint(modeStr, 8, 32)
				if err != nil {
					return nil, fmt.Errorf("octal mode failed: %v", err)
				}
				mode = os.FileMode(m)
			}
			path := listenUrl.Path
			_ = os.Remove(path)
			if listener, err = net.Listen("unix", path); err != nil {
				return nil, fmt.Errorf("unix listener failed: %v", err)
			}
			if err := os.Chmod(path, mode); err != nil {
				return nil, fmt.Errorf("chmod failed: %v", err)
			}
		default:
			return nil, fmt.Errorf("unsupported URL scheme in %v", arg)
		}
	}
	return listener, nil
}

package main

import (
	"container/list"
	"context"
	"net"
	"net/http"
	"strings"
	"tvman/msg"
)

func sendMenuCommand(client *http.Client, request string, done chan bool) {
	resp, err := client.Get(request)
	if err != nil {
		msg.Error("client failed (%v)", err)
	} else {
		resp.Body.Close()
	}
	done <- true
}

func menuThread(wakeMain chan any, wake chan any) {

	url := global.menumanUrl

	transport := &http.Transport{}

	// TODO: refactor this out because it is in three places
	if strings.HasPrefix(url, "unix://") {
		path := url[7:]
		transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			return net.Dial("unix", path)
		}
		// fake url because otherwise the client gets confused
		url = "http://localhost/unix/" + path
	}

	client := &http.Client{Transport: transport}

	requestDone := make(chan bool, 1)
	requestFree := true

	cmds := list.New()
	msg.Info("menu: started")
	for {
		for cmds.Front() != nil && requestFree {
			elem := cmds.Front()
			item, _ := elem.Value.(string)
			msg.Verbose("menu: command \"%s\"", item)
			cmds.Remove(elem)
			if url != "" {
				request := url + "/?" + string(item)
				msg.Verbose2("http get: %v", request)
				go sendMenuCommand(client, request, requestDone)
				requestFree = false
			}
		}

		select {
		case event := <-wake:
			switch item := event.(type) {
			case *state:
				msg.Verbose2("menu: new commands %v", item.menuCommands)
				for _, cmd := range item.menuCommands {
					cmds.PushBack(cmd)
				}
			case menuCommand:
				cmds.PushBack(string(item))
			}
		case <-requestDone:
			requestFree = true
		}
	}
}

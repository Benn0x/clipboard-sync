package main

import (
	"fmt"

	"github.com/benn0x/clipboard-sync/networking"
	"github.com/benn0x/clipboard-sync/networking/mdns"
	"github.com/benn0x/clipboard-sync/networking/ws"
	"github.com/benn0x/clipboard-sync/security"
	"golang.design/x/clipboard"

	clip "github.com/benn0x/clipboard-sync/clipboard"
)

func main() {
	if err := clipboard.Init(); err != nil {
		fmt.Println(err)
		return
	}

	ip, err := networking.GrabLocalAddress()
	if err != nil {
		fmt.Println(err)
		return
	}

	if running, err := mdns.CheckIsDeviceRunningAlready(ip); running {
		fmt.Println("This device is already running this software.")
		return
	} else if err != nil {
		fmt.Println(err)
		fmt.Println("We couldn't verify whether or not your device is already running this software.")
		return
	}

	host, server, err := mdns.AllowDeviceDiscovery(ip)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer server.Shutdown()

	security.Setup()

	if host {
		ws.StartWebsocketServer()
	} else {
		_, ip, err := mdns.FindHostDevice()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Host ip address: ", ip)
			err := ws.ConnectToWS(ip)
			if err != nil {
				fmt.Println(err)
			} else {
				go clip.StartClipboardMonitor()
				ws.ClipboardOverWS()
			}
		}
	}
}

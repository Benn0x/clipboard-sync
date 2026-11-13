package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/benn0x/clipboard-sync/networking"
	"github.com/benn0x/clipboard-sync/networking/mdns"
	"github.com/benn0x/clipboard-sync/networking/ws"
	"github.com/benn0x/clipboard-sync/security"
	"golang.design/x/clipboard"

	clip "github.com/benn0x/clipboard-sync/clipboard"
)

const VERSION string = "1.0.0"
const VERSION_LISTING string = "https://raw.githubusercontent.com/Benn0x/clipboard-sync/refs/heads/development/version"

func main() {
	resp, err := http.Get(VERSION_LISTING)
	if err != nil {
		fmt.Printf("Couldn't check if an update was available. %v\n", err)
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Couldn't check if an update was available. %v\n", err)
		} else {
			latest_version := string(body)
			if latest_version != VERSION {
				fmt.Printf("You don't have the most up-to-date version of the software. %s -> %s\n", VERSION, latest_version)
				fmt.Printf("You can download the newest version @ https://github.com/Benn0x/clipboard-sync/releases/latest\n")
			} else {
				fmt.Printf("Your current version (%s) is up to date.\n", VERSION)
			}
		}
	}

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
		fmt.Printf("We couldn't verify whether or not your device is already running this software. Err: %v\n", err)
		return
	}

	host, server, err := mdns.AllowDeviceDiscovery(ip)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer server.Shutdown()

	if host {
		fmt.Println("Your device is the host. The passphrase used MUST be given to client devices for connection.")
	}

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

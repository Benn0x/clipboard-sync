package main

import (
	"fmt"

	"github.com/benn0x/clipboard-sync/networking/mdns"
	"github.com/benn0x/clipboard-sync/networking/websocket"
)

func main() {
	hasHost, _, err := mdns.GrabHostDevice()
	if err != nil {
		fmt.Println(err)
	}

	mdnsServer, err := mdns.AllowNetworkDiscovery(!hasHost)
	if err != nil {
		fmt.Println(err)
	}

	defer mdnsServer.Shutdown()
	if !hasHost {
		if err := websocket.HostWebsocket(); err != nil {
			fmt.Println("ListenAndServe: ", err)
		}
	}
}

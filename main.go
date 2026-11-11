package main

import (
	"fmt"

	"github.com/benn0x/clipboard-sync/networking/mdns"
	"github.com/benn0x/clipboard-sync/networking/websocket"
)

func main() {
	// Check whether there are any known hosts already
	hasHost, hostIp, err := mdns.GrabHostDevice()
	if err != nil {
		fmt.Println(err)
	}

	// Allows the device to be discovered locally
	mdnsServer, err := mdns.AllowNetworkDiscovery(!hasHost)
	if err != nil {
		fmt.Println(err)
	}

	defer mdnsServer.Shutdown() // Stops the device from cancelling the broadcast

	// Checks if the device is a host or not.
	if !hasHost {
		// Starts the websocket connection
		if err := websocket.HostWebsocket(); err != nil {
			fmt.Println("ListenAndServe: ", err)
		}
	} else {
		// Tries to connect to the websocket being ran
		fmt.Println("This device is not currently a host, attempting to connect to", hostIp)
	}
}

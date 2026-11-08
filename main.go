package main

import (
	"fmt"

	"github.com/benn0x/clipboard-sync/networking/mdns"
)

func main() {
	mdnsServer, err := mdns.AllowNetworkDiscovery()
	if err != nil {
		fmt.Println(err)
	}

	defer mdnsServer.Shutdown()
	entries, err := mdns.GrabNetworkDevices()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(entries)
}

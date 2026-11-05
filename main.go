package main

import (
	"fmt"

	"github.com/benn0x/clipboard-sync/mdns"
)

func main() {

	if err := mdns.IdentifyIsHost(); err != nil {
		fmt.Println(err)
	}

	server, err := mdns.RegisterDevice()
	if err != nil {
		fmt.Println(err)
	}
	defer server.Shutdown()

	entries, err := mdns.SearchDevices()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(entries)
}

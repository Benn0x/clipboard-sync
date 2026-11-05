package mdns

import (
	"fmt"

	"github.com/hashicorp/mdns"
)

func SearchDevices() ([]*mdns.ServiceEntry, error) {
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	entryListing := []*mdns.ServiceEntry{}
	go func() {
		for entry := range entriesCh {
			fmt.Println(entry)
			entryListing = append(entryListing, entry)
		}
	}()

	err := mdns.Lookup(ServiceName, entriesCh)
	if err != nil {
		return nil, err
	}

	close(entriesCh)
	return entryListing, nil
}

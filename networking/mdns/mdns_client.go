package mdns

import "github.com/hashicorp/mdns"

func GrabNetworkDevices() ([]*mdns.ServiceEntry, error) {
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	entries := []*mdns.ServiceEntry{}
	go func() {
		for entry := range entriesCh {
			entries = append(entries, entry)
		}
	}()

	if err := mdns.Lookup(SERVICE_NAME, entriesCh); err != nil {
		return nil, err
	}

	close(entriesCh)
	return entries, nil
}

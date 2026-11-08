package mdns

import (
	"strings"

	"github.com/hashicorp/mdns"
)

func GrabHostDevice() (bool, string, error) {
	entries, err := GrabNetworkDevices()
	if err != nil {
		return false, "", err
	}

	hostFound := false
	hostIp := ""
	for _, entry := range entries {
		splitEntry := strings.Split(entry.Info, ";")
		for i := 0; i < len(splitEntry); i++ {
			infoEntry := splitEntry[i]
			isHost := strings.TrimPrefix("host:", infoEntry)
			ipAddress := strings.TrimPrefix("ip:", infoEntry)
			if isHost == "true" {
				hostFound = true
			}

			if ipAddress != "" {
				hostIp = ipAddress
			}
		}

		if hostFound == true {
			break
		}
	}

	return hostFound, hostIp, nil
}

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

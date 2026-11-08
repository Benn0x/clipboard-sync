package mdns

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/mdns"
)

// Grabs a device that is known to be a host, or returns False in first argument if no device is found.
func GrabHostDevice() (bool, string, error) {
	entries, err := GrabNetworkDevices()
	if err != nil {
		return false, "", err
	}

	hostFound := false
	hostIp := ""
	for _, entry := range entries {
		splitEntry := strings.Split(entry.Info, ";")
		for i := range splitEntry {
			infoEntry := splitEntry[i]
			isHost := strings.TrimPrefix(infoEntry, "host:")
			ipAddress := strings.TrimPrefix(infoEntry, "ip:")
			if isHost == "true" {
				hostFound = true
			}

			if ipAddress != "" && (validator.New().Var(ipAddress, "required,ip4_addr") == nil) {
				hostIp = ipAddress
			}
		}

		if hostFound == true {
			break
		}
	}

	return hostFound, hostIp, nil
}

// Grabs and returns a list of devices that are running the software on your network.
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

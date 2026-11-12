package mdns

import (
	"fmt"
	"strings"

	"github.com/hashicorp/mdns"
)

// Checks whether or not the device is already running the software
func CheckIsDeviceRunningAlready(ip string) (bool, error) {
	entries, err := ScanLocalDevices()
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	for _, entry := range entries {
		text := entry.Info
		textSplit := strings.Split(text, ";")
		for _, s := range textSplit {
			ipValue, isIpFound := strings.CutPrefix(s, "ip:")
			if isIpFound {
				if ipValue == ip {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// Finds the host device information if it exists
func FindHostDevice() (*mdns.ServiceEntry, string, error) {
	entries, err := ScanLocalDevices()
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}
	var host *mdns.ServiceEntry
	var ip string

	for _, entry := range entries {
		text := entry.Info
		textSplit := strings.Split(text, ";")
		for _, s := range textSplit {
			hostValue, isHostFound := strings.CutPrefix(s, "host:")
			ipValue, isIpFound := strings.CutPrefix(s, "ip:")
			if isHostFound {
				if hostValue == "true" {
					host = entry
				}
			}

			if isIpFound {
				ip = ipValue
			}
		}

		if host != nil {
			break
		}
	}

	return host, ip, nil
}

// Scans the local network for devices running the software
func ScanLocalDevices() ([]*mdns.ServiceEntry, error) {
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	listOfDevices := []*mdns.ServiceEntry{}
	go func() {
		for entry := range entriesCh {
			listOfDevices = append(listOfDevices, entry)
		}
	}()

	if err := mdns.Lookup(SERVICE_NAME, entriesCh); err != nil {
		return nil, fmt.Errorf("MDNs scan error %v", err)
	}
	close(entriesCh)

	return listOfDevices, nil
}

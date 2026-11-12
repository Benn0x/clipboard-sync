package mdns

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/mdns"
)

func hostOnNetwork() (bool, error) {
	devices, err := ScanLocalDevices()
	if err != nil {
		return false, err
	}
	for _, entry := range devices {
		infoText := entry.Info
		splitText := strings.Split(infoText, ";")
		for _, text := range splitText {
			host := strings.TrimPrefix(text, "host:")
			if host == "true" {
				return true, nil
			}
		}
	}

	return false, err
}

func AllowDeviceDiscovery(ip string) (bool, *mdns.Server, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return false, nil, fmt.Errorf("Hostname error %v", err)
	}
	host, err := hostOnNetwork()
	if err != nil {
		return !host, nil, err
	}
	info := []string{fmt.Sprintf("host:%t;ip:%s", !host, ip)}
	service, err := mdns.NewMDNSService(hostname, SERVICE_NAME, "", "", SERVICE_PORT, nil, info)
	if err != nil {
		return !host, nil, fmt.Errorf("Service error %v", err)
	}

	server, err := mdns.NewServer(&mdns.Config{
		Zone: service,
	})
	if err != nil {
		return !host, nil, fmt.Errorf("Server error %v", err)
	}

	return !host, server, nil
}

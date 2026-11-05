package mdns

import (
	"fmt"
	"os"

	"github.com/benn0x/clipboard-sync/networking"
	"github.com/hashicorp/mdns"
)

func RegisterDevice() (*mdns.Server, error) {
	text := TextRecord{
		IPAddress: networking.GetLocalAddress(),
		Host:      DeviceIsHost,
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	service, err := mdns.NewMDNSService(hostname, ServiceName, "", "", ServicePort, nil, []string{fmt.Sprintf("Ip:%s;Host:%t", text.IPAddress, text.Host)})
	if err != nil {
		return nil, err
	}

	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return nil, err
	}

	return server, nil
}

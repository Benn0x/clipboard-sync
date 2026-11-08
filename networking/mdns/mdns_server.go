package mdns

import (
	"fmt"
	"os"

	"github.com/benn0x/clipboard-sync/networking"
	"github.com/hashicorp/mdns"
)

// This broadcasts the MDNs device across the local network
func AllowNetworkDiscovery(isHost bool) (*mdns.Server, error) {
	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	localAddress, err := networking.GrabLocalIp()
	if err != nil {
		return nil, err
	}

	info := []string{fmt.Sprintf("ip:%s;host:%t", localAddress.IP.String(), isHost)}
	service, err := mdns.NewMDNSService(host, SERVICE_NAME, "", "", SERVICE_PORT, nil, info)
	if err != nil {
		return nil, err
	}

	server, _ := mdns.NewServer(&mdns.Config{
		Zone: service,
	})

	return server, nil
}

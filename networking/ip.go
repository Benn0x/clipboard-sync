package networking

import (
	"fmt"
	"net"
)

// Grabs the local ip address of the device
func GrabLocalAddress() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", fmt.Errorf("IP grabbing error %v", err)
	}

	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP.String(), err
}

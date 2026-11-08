package networking

import "net"

func GrabLocalIp() (*net.UDPAddr, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr), nil
}

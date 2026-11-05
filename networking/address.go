package networking

import (
	"net"
	"strings"
)

func GetLocalAddress() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}

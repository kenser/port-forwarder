package utils

import (
	"net"
)

func GetLocalHostAddress() (ip string) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ip
	}

	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

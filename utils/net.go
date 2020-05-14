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


func CheckPortIsAvailable(network string, address string) bool {
	ln, err := net.Listen(network, address)

	if err != nil {
		return false
	}

	_ = ln.Close()
	return true
}
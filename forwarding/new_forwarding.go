package forwarding

import (
	"fmt"
	"github.com/cloverzrg/go-portforward/logger"
	"io"
	"net"
	"time"
)

// network: tcp, tcp4, tcp6, udp, udp4, udp6, ip, ip4, ip6, unix, unixgram, unixpacket
// listenAddress: :8080, 127.0.0.1:8080
func New(network, listenAddress string, listenPort int, targetAddress string, targetPort int) (err error) {
	var quit = make(chan struct{})
	listen := fmt.Sprintf("%s:%d", listenAddress, listenPort)
	ln, err := net.Listen(network, listen)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer ln.Close()
	logger.Infof("new forwarding:%s:%d -> %s:%d", listenAddress, listenPort, targetAddress, targetPort)

	var connChan = make(chan net.Conn)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				panic(err)
			}
			connChan <- conn
		}
	}()

	for {
		select {
		case <-quit:
			return
		case conn := <-connChan:
			go handleRequest(network, targetAddress, targetPort, conn)
		}

	}
}

func handleRequest(network string, targetAddress string, targetPort int, conn net.Conn) {
	target := fmt.Sprintf("%s:%d", targetAddress, targetPort)
	proxy, err := net.Dial(network, target)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Infof("forward:%v-->%v-->%v", conn.RemoteAddr(), conn.LocalAddr(), target)
	go copyIO(conn, proxy, 1)
	go copyIO(proxy, conn, 2)
}

func copyIO(src, dest net.Conn, connType int) {
	var n int64
	start := time.Now()
	defer src.Close()
	defer dest.Close()
	n, _ = io.Copy(src, dest)
	end := time.Now()
	if connType == 1 {
		logger.Info("%s --> %s conn close, 发送 %d bytes, dur %+v\n", src.LocalAddr(), dest.RemoteAddr(), n, end.Sub(start))
	} else {
		logger.Info("%s --> %s conn close, 接收 %d bytes, dur %+v\n", src.RemoteAddr(), dest.LocalAddr(), n, end.Sub(start))
	}

}

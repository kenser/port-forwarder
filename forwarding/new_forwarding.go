package forwarding

import (
	"fmt"
	"github.com/cloverzrg/go-portforward/logger"
	"io"
	"net"
	"time"
)

// network: tcp, tcp4, tcp6, udp, udp4, udp6, ip, ip4, ip6, unix, unixgram, unixpacket
// listenAddress + listenPort: :8080, 127.0.0.1:8080
func New2(network, listenAddress string, listenPort int, targetAddress string, targetPort int, stopChan chan struct{}) (err error) {
	listen := fmt.Sprintf("%s:%d", listenAddress, listenPort)
	ln, err := net.Listen(network, listen)
	if err != nil {
		logger.Error(err)
		return err
	}

	connStopChanArr := make([]*chan struct{}, 0)
	var connChan = make(chan net.Conn)

	defer func() {
		var err error
		err = ln.Close()
		if err != nil {
			logger.Error(err)
		}
		for _, v := range connStopChanArr {
			*v <- struct{}{}
		}
		close(connChan)
		close(stopChan)

	}()
	logger.Infof("new forwarding:%s:%d -> %s:%d", listenAddress, listenPort, targetAddress, targetPort)

	go func() {
		defer func() {
			_ = ln.Close()
		}()
		for {
			conn, err := ln.Accept()
			if err != nil {
				select {
				case <-stopChan:
					logger.Info("stopChan had closed")
				case <-time.After(1 * time.Second):
					logger.Error(err)
					logger.Info("now send stop signal to stopChan")
					stopChan <- struct{}{}
				}
				return
			}
			connChan <- conn
		}
	}()

	for {
		select {
		case <-stopChan:
			return
		case conn := <-connChan:
			c := make(chan struct{})
			connStopChanArr = append(connStopChanArr, &c)
			go handleRequest(network, targetAddress, targetPort, conn, c)
		}
	}
}

func handleRequest(network string, targetAddress string, targetPort int, conn net.Conn, stopChan chan struct{}) {
	target := fmt.Sprintf("%s:%d", targetAddress, targetPort)
	proxy, err := net.Dial(network, target)
	if err != nil {
		logger.Error(err)
		return
	}
	defer func() {
		var err error
		err = proxy.Close()
		if err != nil {
			logger.Error(err)
		}
		err = conn.Close()
		if err != nil {
			logger.Error(err)
		}
	}()

	logger.Infof("forward:%v-->%v-->%v", conn.RemoteAddr(), conn.LocalAddr(), target)
	go copyIO(conn, proxy, 1)
	go copyIO(proxy, conn, 2)
	<-stopChan
}

func copyIO(src, dest net.Conn, connType int) {
	var n int64
	start := time.Now()
	n, _ = io.Copy(src, dest)
	end := time.Now()
	if connType == 1 {
		logger.Infof("%s --> %s conn close, 发送 %d bytes, dur %+v\n", src.LocalAddr(), dest.RemoteAddr(), n, end.Sub(start))
	} else {
		logger.Infof("%s --> %s conn close, 接收 %d bytes, dur %+v\n", src.RemoteAddr(), dest.LocalAddr(), n, end.Sub(start))
	}

}

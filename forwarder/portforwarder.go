package forwarder

import (
	"fmt"
	"github.com/cloverzrg/go-portforward/logger"
	"io"
	"math"
	"net"
	"sync"
	"time"
)

type PortForwarder struct {
	Network          string // listen network:tcp, tcp4, tcp6, udp, udp4, udp6, ip, ip4, ip6, unix, unixgram, unixpacket
	ListenAddress    string // '0.0.0.0','127.0.0.1', ''
	ListenPort       int
	Listener         net.Listener
	TargetAddress    string        // forward target ip
	TargetPort       int           // forward target port
	StopChan         chan struct{} // stop chan, when received will stop listen and close all conn
	ConnChan         chan net.Conn // conn chan for new accepted conn
	ConnMap          map[uint]net.Conn
	ConnMapPointer   uint
	Mutex            sync.Mutex
	ConnCount        uint
	CurrentConnCount uint
	IsClosed         bool
}

func (pf *PortForwarder) nextConnMapPointer() uint {
	pf.Mutex.Lock()
	defer pf.Mutex.Unlock()
	pf.ConnMapPointer = (pf.ConnMapPointer + 1) % math.MaxUint32
	return pf.ConnMapPointer
}

func (pf *PortForwarder) checkListenPortAvailable() error {
	var err error
	listen := fmt.Sprintf("[%s]:%d", pf.ListenAddress, pf.ListenPort)
	li, err := net.Listen(pf.Network, listen)
	if err != nil {
		return err
	}
	li.Close()
	return err
}

func (pf *PortForwarder) checkTargetPortAvailable() error {
	var err error
	listen := fmt.Sprintf("[%s]:%d", pf.TargetAddress, pf.TargetPort)
	li, err := net.Dial(pf.Network, listen)
	if err != nil {
		return err
	}
	li.Close()
	return err
}

func New(network, listenAddress string, listenPort int, targetAddress string, targetPort int) (pf *PortForwarder, err error) {
	if listenAddress != "" {
		listenIP := net.ParseIP(listenAddress)
		if listenIP == nil {
			return pf, fmt.Errorf("listenAddress %s is not a valid IP", listenAddress)
		}
		listenAddress = listenIP.String()
	}
	targetIP := net.ParseIP(targetAddress)
	if targetIP == nil {
		return pf, fmt.Errorf("targetAddress %s is not a valid IP", targetAddress)
	}
	if !(listenPort >= 0 && listenPort <= 65535) {
		return pf, fmt.Errorf("listenPort %d is invalid", listenPort)
	}
	if !(targetPort >= 0 && targetPort <= 65535) {
		return pf, fmt.Errorf("targetPort %d is invalid", targetPort)
	}

	pf = &PortForwarder{
		Network:       network,
		ListenAddress: listenAddress,
		ListenPort:    listenPort,
		TargetAddress: targetAddress,
		TargetPort:    targetPort,
		ConnChan:      make(chan net.Conn, 100), // 100 buffer
		ConnMap:       make(map[uint]net.Conn),
		StopChan:      make(chan struct{}),
	}

	err = pf.checkListenPortAvailable()
	if err != nil {
		return nil, fmt.Errorf("the listen port is not available:%+v", err)
	}

	err = pf.checkTargetPortAvailable()
	if err != nil {
		logger.Warn("target port is not available:", err)
	}
	return pf, err
}

func (pf *PortForwarder) Start() (err error) {
	if pf.IsClosed {
		return fmt.Errorf("the portforwarder id closed. please use New() to new one")
	}
	listen := fmt.Sprintf("[%s]:%d", pf.ListenAddress, pf.ListenPort)
	pf.Listener, err = net.Listen(pf.Network, listen)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Infof("start forwarding:%s:%d -> %s:%d", pf.ListenAddress, pf.ListenPort, pf.TargetAddress, pf.TargetPort)

	go func() {
		for {
			conn, err := pf.Listener.Accept()
			if err != nil {
				//select {
				//case <-pf.StopChan:
				//	logger.Info("stopChan had closed")
				//case <-time.After(1 * time.Second):
				//	logger.Error(err)
				//	logger.Info("now send stop signal to stopChan")
				//	pf.StopChan <- struct{}{}
				//}
				return
			}
			pf.ConnChan <- conn
		}
	}()

	go func() {
		for {
			select {
			case <-pf.StopChan:
				return
			case conn := <-pf.ConnChan:
				pf.ConnCount++
				pf.CurrentConnCount++
				go pf.handleRequest(conn, pf.nextConnMapPointer())
			}
		}
	}()
	return
}

func (pf *PortForwarder) AddConn(id uint, conn net.Conn) {
	pf.Mutex.Lock()
	defer pf.Mutex.Unlock()
	pf.ConnMap[id] = conn
}

func (pf *PortForwarder) DecCurrentConnCount() {
	pf.Mutex.Lock()
	defer pf.Mutex.Unlock()
	pf.CurrentConnCount--
}

func (pf *PortForwarder) DelConn(id uint) {
	pf.Mutex.Lock()
	defer pf.Mutex.Unlock()
	delete(pf.ConnMap, id)
}

func (pf *PortForwarder) handleRequest(conn net.Conn, id uint) {
	pf.AddConn(id, conn)
	defer func() {
		var err error
		if !pf.IsClosed {
			err = conn.Close()
			if err != nil {
				logger.Errorf("close conn err:", id, err)
			}
			pf.DelConn(id)
		}
		pf.DecCurrentConnCount()
	}()
	target := fmt.Sprintf("[%s]:%d", pf.TargetAddress, pf.TargetPort)
	proxy, err := net.Dial(pf.Network, target)
	if err != nil {
		logger.Error(err)
		return
	}
	//proxyId := pf.nextConnMapPointer()
	//pf.ConnMap[proxyId] = proxy
	defer func() {
		var err error
		err = proxy.Close()
		//logger.Infof("conn closed(%d)", proxyId)
		if err != nil {
			logger.Errorf("close proxy conn(%d) err:", id, err)
		}
	}()

	logger.Infof("new connection(%d):%v-->%v-->%v", id, conn.RemoteAddr(), conn.LocalAddr(), target)
	c1 := make(chan struct{})
	c2 := make(chan struct{})
	go pf.copyIO(conn, proxy, 1, c1)
	go pf.copyIO(proxy, conn, 2, c2)

	select {
	case <-c1:
	case <-c2:
	}
}

func (pf *PortForwarder) copyIO(src, dest net.Conn, connType int, c chan struct{}) {
	defer func() {
		c <- struct{}{}
	}()
	//var n int64
	//start := time.Now()
	_, _ = io.Copy(src, dest)
	//end := time.Now()
	//if connType == 1 {
	//	logger.Infof("%s --> %s conn close, 发送 %d bytes, dur %+v\n", src.LocalAddr(), dest.RemoteAddr(), n, end.Sub(start))
	//} else {
	//	logger.Infof("%s --> %s conn close, 接收 %d bytes, dur %+v\n", src.RemoteAddr(), dest.LocalAddr(), n, end.Sub(start))
	//}

}

func (pf *PortForwarder) Close() (err error) {
	// stop listen
	err = pf.Listener.Close()
	if err != nil {
		logger.Error(err)
	}
	pf.StopChan <- struct{}{}
	// close channel
	close(pf.ConnChan)
	close(pf.StopChan)
	for {
		_, ok := <-pf.ConnChan
		if !ok {
			break
		}
	}

	pf.IsClosed = true
	// stop all connection
	logger.Infof("conn map len:%d", len(pf.ConnMap))
	for k, v := range pf.ConnMap {
		err = v.Close()
		if err != nil {
			logger.Errorf("close conn(%d) by Close() err:", k, err)
		}
	}

	pf.ConnMap = nil

	time.Sleep(100 * time.Millisecond)
	logger.Info("close port forwarder done")
	return err
}

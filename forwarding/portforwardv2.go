package forwarding

import (
	"fmt"
	"github.com/cloverzrg/go-portforward/logger"
	"io"
	"math"
	"net"
	"sync"
	"time"
)

type PortForward struct {
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
	ClosedCount      uint
	CurrentConnCount uint
}

func (pf *PortForward) nextConnMapPointer() uint {
	pf.Mutex.Lock()
	defer pf.Mutex.Unlock()
	pf.ConnMapPointer = (pf.ConnMapPointer + 1) % math.MaxUint32
	return pf.ConnMapPointer
}

func New3(network, listenAddress string, listenPort int, targetAddress string, targetPort int) (pf *PortForward) {
	pf = &PortForward{
		Network:       network,
		ListenAddress: listenAddress,
		ListenPort:    listenPort,
		TargetAddress: targetAddress,
		TargetPort:    targetPort,
		ConnChan:      make(chan net.Conn, 100), // 100 buffer
		ConnMap:       make(map[uint]net.Conn),
		StopChan:      make(chan struct{}),
	}
	return pf
}

func (pf *PortForward) Start() (err error) {
	listen := fmt.Sprintf("%s:%d", pf.ListenAddress, pf.ListenPort)
	pf.Listener, err = net.Listen(pf.Network, listen)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Infof("start forwarding:%s:%d -> %s:%d", pf.ListenAddress, pf.ListenPort, pf.TargetAddress, pf.TargetPort)

	go func() {
		//defer pf.Stop()
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
				go pf.handleRequest(conn, pf.nextConnMapPointer())
			}
		}
	}()
	return
}

func (pf *PortForward) AddConn(id uint, conn net.Conn) {
	pf.Mutex.Lock()
	defer pf.Mutex.Unlock()
	pf.ConnMap[id] = conn
}

func (pf *PortForward) DelConn(id uint) {
	pf.Mutex.Lock()
	defer pf.Mutex.Unlock()
	delete(pf.ConnMap, id)
}

func (pf *PortForward) handleRequest(conn net.Conn, id uint) {
	pf.AddConn(id, conn)
	defer func() {
		var err error
		select {
		case <-pf.StopChan:
			// 通过 Stop() 关闭
		case <-time.After(5 * time.Millisecond):
			// 正常关闭
			err = conn.Close()
			if err != nil {
				logger.Errorf("close conn err:", id, err)
			} else {
				logger.Infof("conn closed(%d)", id)
			}
			pf.DelConn(id)
		}

	}()
	target := fmt.Sprintf("%s:%d", pf.TargetAddress, pf.TargetPort)
	proxy, err := net.Dial(pf.Network, target)
	if err != nil {
		logger.Error(err)
		return
	}
	//proxyId := pf.nextConnMapPointer()
	//pf.ConnMap[proxyId] = proxy
	defer func() {
		var err error
		if proxy != nil {
			err = proxy.Close()
			//logger.Infof("conn closed(%d)", proxyId)
			if err != nil {
				logger.Errorf("close proxy conn(%d) err:", id, err)
			} else {
				logger.Infof("proxy conn closed(%d)", id)
			}
		} else {
			logger.Infof("proxy conn had closed(%d)", id)
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

func (pf *PortForward) copyIO(src, dest net.Conn, connType int, c chan struct{}) {
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

func (pf *PortForward) Stop() (err error) {
	logger.Info("forward stop")
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

	// stop listen
	err = pf.Listener.Close()
	if err != nil {
		logger.Error(err)
	}
	// stop all connection
	logger.Infof("map len:%d", len(pf.ConnMap))
	for k, v := range pf.ConnMap {
		if v != nil {
			err = v.Close()
			if err != nil {
				logger.Errorf("close conn(%d) by Stop() err:", k, err)
			} else {
				logger.Infof("conn closed(%d) by Stop()", k)
			}

		} else {
			logger.Infof("conn(%d) is closed(%d)", k)
		}
	}

	time.Sleep(100 * time.Millisecond)

	return err
}
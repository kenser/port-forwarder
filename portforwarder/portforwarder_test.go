package portforwarder

import (
	"fmt"
	"github.com/cloverzrg/go-portforward/logger"
	"net"
	"testing"
	"time"
)

func TestParseIP(t *testing.T) {
	ip := net.ParseIP("2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b")
	t.Log(ip.To4())
	t.Log(ip.To16())
	t.Log(ip.String())
}

func TestDialIPv6(t *testing.T) {
	_, err := net.Dial("tcp", "[127.0.0.1]:80")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestPortForward(t *testing.T) {
	var err error
	pf, err := New("tcp", "127.0.0.1", 8080, "47.52.114.182", 80)
	if err != nil {
		t.Error(err)
	}
	err = pf.Start()
	if err != nil {
		logger.Error(err)
	}
	for i := 0; i < 1000; i++ {
		time.Sleep(1 * time.Second)
		fmt.Printf("conn count:%d, cur conn:%d\n", pf.ConnCount, pf.CurrentConnCount)
	}

	err = pf.Close()
	if err != nil {
		logger.Error(err)
	}
}

func TestDial(t *testing.T) {
	_, err := net.Dial("tcp", "47.52.114.182:800")
	if err != nil {
		logger.Error(err)
		return
	}
}

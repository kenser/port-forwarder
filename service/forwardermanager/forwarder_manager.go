package forwardermanager

import (
	"context"
	"fmt"
	"github.com/cloverzrg/go-portforward/dns"
	"github.com/cloverzrg/go-portforward/logger"
	"github.com/cloverzrg/go-portforward/model/forwarddao"
	"github.com/cloverzrg/go-portforward/portforwarder"
)

var ForwardingMap map[int]*portforwarder.PortForwarder

func init() {
	ForwardingMap = make(map[int]*portforwarder.PortForwarder)
}

func CloseById(ctx context.Context, id int) (err error) {
	defer func() {
		if err != nil {
			logger.Error(err)
		}
	}()
	if ForwardingMap[id] != nil {
		err = ForwardingMap[id].Close()
		if err != nil {
			return
		}
		delete(ForwardingMap, id)
		return
	} else {
		return fmt.Errorf("the forward is not running")
	}
	//return err
}

func StartById(ctx context.Context, id int) (err error) {
	defer func() {
		if err != nil {
			logger.Error(err)
		}
	}()
	data, err := forwarddao.GetById(id)
	if err != nil {
		return err
	}
	if ForwardingMap[data.Id] != nil {
		if ForwardingMap[data.Id].IsClosed != true {
			return fmt.Errorf("the forward is already running")
		}
		err = ForwardingMap[data.Id].Close()
		if err != nil {
			return err
		}
	}
	targetIp, err := dns.LookupIP(data.TargetAddress)
	if err != nil {
		return err
	}
	newForwarder, err := portforwarder.New(data.Network, data.ListenAddress, data.ListenPort, targetIp, data.TargetPort)
	if err != nil {
		return err
	}
	err = newForwarder.Start()
	if err != nil {
		return err
	}
	ForwardingMap[data.Id] = newForwarder
	return err
}

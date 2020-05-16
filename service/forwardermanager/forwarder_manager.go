package forwardermanager

import (
	"context"
	"github.com/cloverzrg/go-portforward/dns"
	"github.com/cloverzrg/go-portforward/logger"
	"github.com/cloverzrg/go-portforward/model/forwarddao"
	"github.com/cloverzrg/go-portforward/portforwarder"
	"net"
)

var ForwardingMap map[int]*portforwarder.PortForwarder

func init() {
	ForwardingMap = make(map[int]*portforwarder.PortForwarder)
}

func CloseById(ctx context.Context, id int) (err error) {
	return err
}

func StartById(ctx context.Context, id int) (err error) {
	defer func() {
		if err != nil {
			logger.Error(err)
		}
	}()
	data, err := forwarddao.GetByID(id)
	if err != nil {
		return err
	}
	if ForwardingMap[data.Id] != nil {
		err = ForwardingMap[data.Id].Close()
		if err != nil {
			return
		}
	}
	targetIp := data.TargetAddress
	if net.ParseIP(data.TargetAddress) == nil {
		// 当识别ip失败,尝试使用dns解析
		targetIp, err = dns.LookupIP(data.TargetAddress)
		if err != nil {
			return err
		}
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
	data.Status = 1
	err = forwarddao.UpdateByID(data.Id, data)
	if err != nil {
		return err
	}
	return err
}

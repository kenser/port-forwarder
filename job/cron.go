package job

import (
	"context"
	"github.com/cloverzrg/go-portforwarder/dns"
	"github.com/cloverzrg/go-portforwarder/logger"
	"github.com/cloverzrg/go-portforwarder/memtool"
	"github.com/cloverzrg/go-portforwarder/model/forwarddao"
	"github.com/cloverzrg/go-portforwarder/service/forwardermanager"
	"github.com/robfig/cron/v3"
)

func CheckDomainIPChange() {
	list, err := forwarddao.FindAllRunning()
	if err != nil {
		logger.Error(err)
		return
	}
	for _, v := range list {
		ip, err := dns.LookupIP(v.TargetAddress)
		if err != nil {
			logger.Error("cron CheckDomainIPChange err :", err)
			continue
		}
		if forwardermanager.ForwardingMap[v.Id] != nil && ip != forwardermanager.ForwardingMap[v.Id].TargetAddress {
			err = forwardermanager.CloseById(context.Background(), v.Id)
			if err != nil {
				logger.Error(err)
				continue
			}
			err = forwardermanager.StartById(context.Background(), v.Id)
			if err != nil {
				logger.Error(err)
				continue
			}
		}
	}
}

func Start() {
	// Seconds field, required
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 */1 * * * ?", memtool.PrintMemUsage)
	c.AddFunc("0 */1 * * * ?", CheckDomainIPChange)
	c.Start()
}

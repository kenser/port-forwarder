package job

import (
	"github.com/cloverzrg/go-portforward/memtool"
	"github.com/robfig/cron/v3"
)

func CheckDomainIPChange() {

}

func Start() {
	// Seconds field, required
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 */1 * * * ?", memtool.PrintMemUsage)
	c.Start()
}

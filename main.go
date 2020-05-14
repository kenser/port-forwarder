package main

import (
	"fmt"
	"github.com/cloverzrg/go-portforward/config"
	"github.com/cloverzrg/go-portforward/forwarding"
	"github.com/cloverzrg/go-portforward/logger"
	"github.com/cloverzrg/go-portforward/web"
	"time"
)

// @title go-portforward
// @version 1.0

// @contact.name API Support
// @contact.url https://github.com/cloverzrg/go-portforward
// @contact.email cloverzrg@gmail.com

// @license.name go-portforward

var (
	BuildTime string
	GoVersion string
	GitHead   string
)

func main() {
	var err error
	for i:= 0; i< 100; i++ {
		logger.Info(i)
	}
	pf, _ := forwarding.New3("tcp", "", 8080, "47.52.114.182", 80)
	err = pf.Start()
	if err != nil {
		logger.Error(err)
	}
	time.Sleep(15 * time.Second)
	err = pf.Stop()
	if err != nil {
		logger.Error(err)
	}
	return

	err = web.Start()
	if err != nil {
		logger.Panic(err)
	}
}

func init() {
	fmt.Printf("BuildTime: %s\nGoVersion: %s\nGitHead: %s\n", BuildTime, GoVersion, GitHead)
	config.Parse("./config.json")
}

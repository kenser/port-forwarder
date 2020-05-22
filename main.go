package main

import (
	"fmt"
	"github.com/cloverzrg/go-portforwarder/config"
	"github.com/cloverzrg/go-portforwarder/db"
	"github.com/cloverzrg/go-portforwarder/logger"
	"github.com/cloverzrg/go-portforwarder/model"
	"github.com/cloverzrg/go-portforwarder/service/forward"
	"github.com/cloverzrg/go-portforwarder/web"
)

// @title go-portforwarder
// @version 1.0

// @contact.name API Support
// @contact.url https://github.com/cloverzrg/go-portforwarder
// @contact.email cloverzrg@gmail.com

// @license.name go-portforwarder

var (
	BuildTime string
	GoVersion string
	GitHead   string
)

func main() {
	var err error
	err = web.Start()
	if err != nil {
		logger.Panic(err)
	}
}

func init() {
	var err error
	fmt.Printf("BuildTime: %s\nGoVersion: %s\nGitHead: %s\n", BuildTime, GoVersion, GitHead)
	err = config.Parse("./data/config.json")
	if err != nil {
		logger.Panic(err)
	}
	err = db.Connect()
	if err != nil {
		logger.Panic(err)
	}
	err = model.CreateAllTable()
	if err != nil {
		logger.Panic(err)
	}
	err = forward.StartUp()
	if err != nil {
		logger.Error(err)
	}
}

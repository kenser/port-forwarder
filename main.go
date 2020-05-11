package main

import (
	"fmt"
	"github.com/cloverzrg/go-portforward/config"
	"github.com/cloverzrg/go-portforward/logger"
	"github.com/cloverzrg/go-portforward/web"
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
	err = web.Start()
	if err != nil {
		logger.Panic(err)
	}

	fmt.Println("forward:", conn.RemoteAddr().(*net.TCPAddr).IP, "-->", conn.LocalAddr(), "-->", "47.52.114.182:80")
	go copyIO(conn, proxy, 1)
	go copyIO(proxy, conn, 2)
}

func init() {
	fmt.Printf("BuildTime: %s\nGoVersion: %s\nGitHead: %s\n", BuildTime, GoVersion, GitHead)
	config.Parse("./config.json")
}

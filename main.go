package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	var quit = make(chan struct{})
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	defer ln.Close()
	if err != nil {
		panic(err)
	}

	var connChan = make(chan net.Conn)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				panic(err)
			}
			connChan <- conn
		}
	}()

	for {
		select {
		case <-quit:
			return
		case conn := <-connChan:
			go handleRequest(conn)
		}

	}
}

func handleRequest(conn net.Conn) {
	proxy, err := net.Dial("tcp", "47.52.114.182:80")
	if err != nil {
		panic(err)
	}

	fmt.Println("forward:", conn.RemoteAddr(), "-->", conn.LocalAddr(), "-->", "47.52.114.182:80")
	go copyIO(conn, proxy, 1)
	go copyIO(proxy, conn, 2)
}

func copyIO(src, dest net.Conn, connType int) {
	var n int64
	start := time.Now()
	defer src.Close()
	defer dest.Close()
	n, _ = io.Copy(src, dest)
	end := time.Now()
	if connType == 1 {
		fmt.Printf("%s --> %s conn close, 发送 %d bytes, dur %+v\n", src.LocalAddr(), dest.RemoteAddr(), n, end.Sub(start))
	} else {
		fmt.Printf("%s --> %s conn close, 接收 %d bytes, dur %+v\n", src.RemoteAddr(), dest.LocalAddr(), n, end.Sub(start))
	}

}

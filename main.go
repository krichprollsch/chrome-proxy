package main

import (
	"io"
	"log"
	"net"
	"os"
	"sync"
)

var logger *log.Logger

func main() {
	logger = log.New(os.Stderr, "chrome-proxy: ", log.Lshortfile)

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Fatalf("impossible to start the server: %v", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Printf("impossible to start the server: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(cliConn net.Conn) {
	defer cliConn.Close()
	logger.Printf("handling client connection: %v", cliConn.RemoteAddr())

	// dial a tcp conn to the google chrome
	logger.Printf("starting chrome connection")
	chromeConn, err := net.Dial("tcp", "127.0.0.1:9222")
	if err != nil {
		logger.Printf("impossible to connect to the chrome: %v", err)
		return
	}
	defer chromeConn.Close()

	// create a wait group to wait write and read goroutines
	var wg sync.WaitGroup

	logger.Printf("start copy from cli to chrome")
	wg.Add(1)
	go copy(chromeConn, cliConn, wg)

	logger.Printf("start copy from chrome to cli")
	wg.Add(1)
	go copy(cliConn, chromeConn, wg)

	wg.Wait()
}

func copy(dst io.Writer, src io.Reader, wg sync.WaitGroup) {
	defer wg.Done()
	if _, err := io.Copy(dst, src); err != nil {
		logger.Printf("impossible to copy from chrome to cli: %v", err)
	}
}

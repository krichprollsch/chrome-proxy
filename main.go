package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
)

var (
	flagListen string
	flagBind   string
	flagKey    string
	logger     *log.Logger
)

func main() {
	logger = log.New(os.Stderr, "chrome-proxy: ", log.Lshortfile)

	flag.StringVar(&flagListen, "listen", "127.0.0.1:8080", "proxy's address to listen")
	flag.StringVar(&flagBind, "bind", "127.0.0.1:9222", "chrome's server address to bind")
	flag.StringVar(&flagKey, "key", "secret", "http Api-Key header secret to check")
	flag.Parse()

	logger.Printf("start the proxy listening %s, binding %s", flagListen, flagBind)
	ln, err := net.Listen("tcp", flagListen)
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
	chromeConn, err := net.Dial("tcp", flagBind)
	if err != nil {
		logger.Printf("impossible to connect to the chrome: %v", err)
		return
	}
	defer chromeConn.Close()

	logger.Printf("start copy from chrome to cli")
	go copy(cliConn, chromeConn)

	logger.Printf("start copy from cli to chrome")
	copy(chromeConn, cliConn)

	logger.Printf("end of handler")
}

// copy sends bytes read from src to dest
func copy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		logger.Printf("impossible to copy: %v", err)
	}
	logger.Printf("end of copy")
}

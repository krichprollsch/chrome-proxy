package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
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

	// check the Api-Key header
	var cliBuf bytes.Buffer
	r, err := http.ReadRequest(bufio.NewReader(io.TeeReader(cliConn, &cliBuf)))
	if err != nil {
		logger.Printf("impossible to read to the cli request: %v", err)
		return
	}
	if "foo" != r.Header.Get("Api-Key") {
		logger.Printf("invalid Api-Key receveived")
		resp := http.Response{
			StatusCode: 401,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Unauthorized\n")),
		}
		resp.Write(cliConn)
		return
	}

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
	go copy(chromeConn, &cliBuf, wg)

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

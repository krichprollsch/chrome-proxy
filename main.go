package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
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

	// check the Api-Key header
	var cliBuf bytes.Buffer
	r, err := fw(io.TeeReader(cliConn, &cliBuf))
	if err != nil {
		logger.Printf("authentication error: %v", err)
		resp := http.Response{
			StatusCode: 401,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Unauthorized\n")),
		}
		resp.Write(cliConn)
		return
	}
	r.Body.Close()

	// dial a tcp conn to the google chrome
	logger.Printf("starting chrome connection")
	chromeConn, err := net.Dial("tcp", flagBind)
	if err != nil {
		logger.Printf("impossible to connect to the chrome: %v", err)
		return
	}
	defer chromeConn.Close()

	var wg sync.WaitGroup

	logger.Printf("start copy from cli to chrome")
	wg.Add(1)
	go copy(chromeConn, ioutil.NopCloser(&cliBuf), &wg)

	logger.Printf("start copy from chrome to cli")
	wg.Add(1)
	go copyTimeout(cliConn, chromeConn, &wg)

	wg.Wait()
	logger.Printf("end of handler")

}

// copy sends bytes read from src to dest
func copy(dst io.Writer, src io.ReadCloser, wg *sync.WaitGroup) {
	defer wg.Done()
	if _, err := io.Copy(dst, src); err != nil {
		logger.Printf("impossible to copy: %v", err)
	}
}

func copyTimeout(dst io.Writer, src net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	buffer := make([]byte, 512)
	for {
		// define a deadline
		err := src.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		n, err := src.Read(buffer)
		if err != nil {
			logger.Printf("error during read: %v", err)
			break
		}
		n, err = dst.Write(buffer[:n])
		if err != nil {
			logger.Printf("error during write: %v", err)
			break
		}
	}
}

// fw parses a http request from the reader and checks the Api-Key token
func fw(reader io.Reader) (*http.Request, error) {
	r, err := http.ReadRequest(bufio.NewReader(reader))
	if err != nil {
		return nil, fmt.Errorf("impossible to read to the cli request: %v", err)
	}

	if flagKey != r.Header.Get("Api-Key") {
		return r, fmt.Errorf("invalid Api-Key receveived")
	}

	return r, nil
}

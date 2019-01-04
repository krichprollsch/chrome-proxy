package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	flagListen string
	flagBind   string
	flagKey    string
	logger     *log.Logger
)

const HTTPHeaderSeparator string = "\r\n\r\n"

func main() {
	logger = log.New(os.Stderr, "", log.LstdFlags)

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

	// does the headers contains the valid secret?
	headers, err := fw(cliConn)
	if err != nil {
		logger.Printf("authentication error: %v", err)
		resp := http.Response{
			StatusCode: 401,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Unauthorized\n")),
		}
		resp.Write(cliConn)
		return
	}

	// dial a tcp conn to the google chrome
	logger.Printf("starting chrome connection")
	chromeConn, err := net.Dial("tcp", flagBind)
	if err != nil {
		logger.Printf("impossible to connect to the chrome: %v", err)
		return
	}
	defer chromeConn.Close()

	// first write the read header to chromeConn
	chromeConn.Write(headers)

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

func fw(r io.Reader) ([]byte, error) {
	headers, err := readHttpHeaders(r)
	if err != nil {
		return nil, fmt.Errorf("impossible to read to the cli request: %v", err)
	}

	if !bytes.Contains(headers, []byte(fmt.Sprintf("Api-Key: %s", flagKey))) {
		return headers, fmt.Errorf("invalid Api-Key receveived")
	}

	return headers, nil
}

// read the reader until http header separator or error
// returns the read bytes slice
// the function doesn't return EOF error if found
func readHttpHeaders(r io.Reader) ([]byte, error) {
	read := make([]byte, 64)
	buf := make([]byte, 0, 1024)
	// we read until we found "\n\n" separator
	for {
		n, err := r.Read(read)
		buf = append(buf, read[:n]...)
		if err == io.EOF || n == 0 {
			// end of file w/o finding end of headers
			break
		}
		if err != nil {
			return buf, err
		}

		// do we read the http header separator?
		if bytes.Contains(buf, []byte(HTTPHeaderSeparator)) {
			break
		}

		// TODO introduce a size limit
	}

	return buf, nil
}

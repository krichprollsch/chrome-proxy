// +build integration

package main

import (
	"bytes"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestGetJsonVersion(t *testing.T) {
	flagBind = "127.0.0.1:9222"

	var buf bytes.Buffer
	logger = log.New(&buf, "chrome-proxy: ", log.Lshortfile)

	req, err := http.NewRequest("GET", "http://127.0.0.1:9222/json/version/json/version", nil)
	req.Close = true
	if err != nil {
		t.Fatalf("%v", err)
	}

	c := &mockConn{}
	req.Write(&c.Body)
	req.Write(os.Stdout)

	handleConnection(c)
}

type mockConn struct {
	Body     bytes.Buffer
	Response bytes.Buffer
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return m.Body.Read(b)
}
func (m *mockConn) Write(b []byte) (n int, err error) {
	return m.Response.Write(b)
}
func (m *mockConn) Close() error {
	return nil
}
func (m *mockConn) LocalAddr() net.Addr {
	return &mockAddr{}
}
func (m *mockConn) RemoteAddr() net.Addr {
	return &mockAddr{}
}
func (m *mockConn) SetDeadline(t time.Time) error {
	return nil
}
func (m *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}
func (m *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

type mockAddr struct{}

func (a *mockAddr) Network() string {
	return "mockNet"
}
func (a *mockAddr) String() string {
	return "mockAddr"
}

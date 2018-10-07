package main

import (
	"bytes"
	"testing"
)

func TestFwRefuseNonHttp(t *testing.T) {
	_, e := fw(bytes.NewBufferString("hello world"))
	if e == nil {
		t.Fatalf("fw must return an error if the input is not http")
	}
}

func TestFwRefuseRequestWithoutApiKeyHeader(t *testing.T) {
	flagKey = "secret"
	_, e := fw(bytes.NewBufferString(`GET / HTTP/1.0

`))
	if e == nil {
		t.Fatalf("fw must return an error if the Api-Key header is missing")
	}
}

func TestFwRefuseRequestWithBadApiKeyHeader(t *testing.T) {
	flagKey = "secret"
	_, e := fw(bytes.NewBufferString(`GET / HTTP/1.0
Api-Key: bad

`))
	if e == nil {
		t.Fatalf("fw must return an error if the Api-Key header is invalid")
	}
}

func TestFwQcceptRequestValidApiKeyHeader(t *testing.T) {
	flagKey = "secret"
	_, e := fw(bytes.NewBufferString(`GET / HTTP/1.0
Api-Key: secret

`))
	if e != nil {
		t.Fatalf("fw must not return an error if the Api-Key header is valid: %v", e)
	}
}

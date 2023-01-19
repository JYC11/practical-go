package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func startBadTestHTTPServer(shutdownServer chan struct{}) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-shutdownServer
		fmt.Fprint(w, "Hello World")
	}))
	return ts
}

func TestFetchBadRemoteResource(t *testing.T) {
	shutdownServer := make(chan struct{})
	ts := startBadTestHTTPServer(shutdownServer)
	defer ts.Close()
	defer func() {
		shutdownServer <- struct{}{}
	}()

	client := createHttpClientWithTimeout(200 * time.Millisecond)
	_, err := fetchRemoteResource(client, ts.URL)
	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	if !strings.Contains(err.Error(), "Client.Timeout exceeded while awaiting headers") {
		t.Fatalf("Expected error to contain: Client.Timeout exceeded while awaiting headers, Got: %v", err.Error())
	}
}

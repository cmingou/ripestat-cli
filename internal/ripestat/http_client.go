package ripestat

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/http2"
)

const (
	defaultHTTPTimeout           = 15 * time.Second
	defaultMaxIdleConns          = 64
	defaultMaxIdleConnsPerHost   = 8
	defaultMaxConnsPerHost       = 8
	defaultIdleConnTimeout       = 90 * time.Second
	defaultDialTimeout           = 5 * time.Second
	defaultDialKeepAlive         = 30 * time.Second
	defaultTLSHandshakeTimeout   = 10 * time.Second
	defaultExpectContinueTimeout = 1 * time.Second
	defaultHTTP2ReadIdleTimeout  = 30 * time.Second
	defaultHTTP2PingTimeout      = 15 * time.Second
)

var (
	httpClientMu     sync.RWMutex
	sharedHTTPClient = newHTTPClient()
)

func newHTTPClient() *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   defaultDialTimeout,
			KeepAlive: defaultDialKeepAlive,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          defaultMaxIdleConns,
		MaxIdleConnsPerHost:   defaultMaxIdleConnsPerHost,
		MaxConnsPerHost:       defaultMaxConnsPerHost,
		IdleConnTimeout:       defaultIdleConnTimeout,
		TLSHandshakeTimeout:   defaultTLSHandshakeTimeout,
		ExpectContinueTimeout: defaultExpectContinueTimeout,
	}

	if http2Transport, err := http2.ConfigureTransports(transport); err != nil {
		debugf("failed to configure HTTP/2 transport: %v", err)
	} else {
		http2Transport.ReadIdleTimeout = defaultHTTP2ReadIdleTimeout
		http2Transport.PingTimeout = defaultHTTP2PingTimeout
	}

	return &http.Client{
		Timeout:   defaultHTTPTimeout,
		Transport: transport,
	}
}

func getHTTPClient() *http.Client {
	httpClientMu.RLock()
	client := sharedHTTPClient
	httpClientMu.RUnlock()
	return client
}

func setHTTPClient(client *http.Client) {
	httpClientMu.Lock()
	sharedHTTPClient = client
	httpClientMu.Unlock()
}

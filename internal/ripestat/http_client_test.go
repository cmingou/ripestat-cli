package ripestat

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/net/http2"
)

type testPayload struct {
	Status string `json:"status"`
}

func TestGetHttpGetResponseUsesHTTP2(t *testing.T) {
	restoreClient := getHTTPClient()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !r.ProtoAtLeast(2, 0) {
			t.Errorf("expected HTTP/2 request, got %s", r.Proto)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(testPayload{Status: "ok"}); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})

	ts, err := newHTTP2TLSServer(handler)
	if err != nil {
		t.Skipf("skipping HTTP/2 test: %v", err)
	}
	t.Cleanup(ts.Close)

	client := newHTTPClient()
	transport := client.Transport.(*http.Transport)
	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{}
	}
	transport.TLSClientConfig.InsecureSkipVerify = true
	transport.TLSClientConfig.NextProtos = []string{"h2", "http/1.1"}
	setHTTPClient(client)
	t.Cleanup(func() {
		setHTTPClient(restoreClient)
	})

	var logBuf bytes.Buffer
	setDebugOutput(&logBuf)
	setHTTPDebug(true)
	t.Cleanup(func() {
		setHTTPDebug(false)
		setDebugOutput(io.Discard)
	})

	body, err := getHttpGetResponse(ts.URL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(body) == 0 {
		t.Fatalf("expected response body, got empty result")
	}

	if !strings.Contains(logBuf.String(), "proto=HTTP/2.0") {
		t.Fatalf("expected log to contain HTTP/2 protocol, got %q", logBuf.String())
	}
}

func newHTTP2TLSServer(handler http.Handler) (*httptest.Server, error) {
	var (
		ts  *httptest.Server
		err error
	)

	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("failed to create httptest server: %v", r)
			}
		}()
		ts = httptest.NewUnstartedServer(handler)
	}()
	if err != nil {
		return nil, err
	}

	// Configure HTTP/2 on the test server before starting
	if cfgErr := http2.ConfigureServer(ts.Config, &http2.Server{}); cfgErr != nil {
		ts.Close()
		return nil, cfgErr
	}

	// Set NextProtos to support HTTP/2
	ts.TLS = &tls.Config{
		NextProtos: []string{"h2"},
	}

	ts.StartTLS()
	return ts, nil
}

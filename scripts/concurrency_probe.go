package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/http2"
)

const defaultEndpoint = "https://stat.ripe.net/data/network-info/data.json?resource=1.1.1.0/24"

var allowedProtocols = map[string]struct{}{
	"http2": {},
	"http1": {},
}

type connTracker struct {
	mu    sync.Mutex
	seen  map[string]int
	total atomic.Int64
}

func newConnTracker() *connTracker {
	return &connTracker{seen: make(map[string]int)}
}

func (c *connTracker) record(conn net.Conn) net.Conn {
	key := fmt.Sprintf("%s->%s", conn.LocalAddr(), conn.RemoteAddr())
	c.mu.Lock()
	c.seen[key]++
	c.mu.Unlock()
	c.total.Add(1)
	return conn
}

func (c *connTracker) uniqueConnections() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.seen)
}

func (c *connTracker) visits() int {
	return int(c.total.Load())
}

type result struct {
	status int
	proto  string
	dur    time.Duration
	err    error
}

func main() {
	protocol := flag.String("protocol", "http2", "Which protocol to use: http2 or http1")
	concurrency := flag.Int("concurrency", 16, "Maximum number of in-flight requests")
	requests := flag.Int("requests", 32, "Total number of requests to issue")
	endpoint := flag.String("endpoint", defaultEndpoint, "RIPEstat endpoint to hit")
	timeout := flag.Duration("timeout", 10*time.Second, "Per-request timeout")
	verbose := flag.Bool("verbose", false, "If true, print each response status as it arrives")
	flag.Parse()

	if _, ok := allowedProtocols[*protocol]; !ok {
		log.Fatalf("unsupported protocol %q; allowed values: http2, http1", *protocol)
	}

	if *concurrency < 1 {
		log.Fatal("concurrency must be at least 1")
	}

	if *requests < 1 {
		log.Fatal("requests must be at least 1")
	}

	if _, err := url.ParseRequestURI(*endpoint); err != nil {
		log.Fatalf("invalid endpoint URL: %v", err)
	}

	// Ensure we only create as many goroutines as requests.
	if *concurrency > *requests {
		*concurrency = *requests
	}

	tracker := newConnTracker()
	client, err := buildClient(*protocol, tracker, *timeout)
	if err != nil {
		log.Fatalf("failed to construct http client: %v", err)
	}

	log.Printf("Running %d requests with max %d in-flight using %s against %s", *requests, *concurrency, *protocol, *endpoint)

	ctx := context.Background()
	sem := make(chan struct{}, *concurrency)
	results := make(chan result, *requests)
	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < *requests; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, *endpoint, nil)
			if err != nil {
				results <- result{err: fmt.Errorf("build request %d: %w", idx, err)}
				return
			}

			issued := time.Now()
			resp, err := client.Do(req)
			if err != nil {
				results <- result{err: fmt.Errorf("request %d failed: %w", idx, err)}
				return
			}
			defer resp.Body.Close()
			if _, err = io.Copy(io.Discard, resp.Body); err != nil {
				results <- result{err: fmt.Errorf("drain response %d: %w", idx, err)}
				return
			}

			r := result{status: resp.StatusCode, proto: resp.Proto, dur: time.Since(issued)}
			if *verbose {
				log.Printf("[%03d] %s %d in %s", idx, resp.Proto, resp.StatusCode, r.dur)
			}
			results <- r
		}(i)
	}

	wg.Wait()
	close(results)

	total := time.Since(start)

	statusCounts := make(map[int]int)
	protoCounts := make(map[string]int)
	var durations []time.Duration
	var errorsEncountered []error

	for r := range results {
		if r.err != nil {
			errorsEncountered = append(errorsEncountered, r.err)
			continue
		}
		statusCounts[r.status]++
		protoCounts[r.proto]++
		durations = append(durations, r.dur)
	}

	if len(durations) > 0 {
		sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	}

	fmt.Println("--- Summary ---")
	fmt.Printf("Wall clock elapsed: %s\n", total)
	fmt.Printf("Total successes: %d / %d\n", len(durations), *requests)
	fmt.Printf("Unique TCP connections: %d (tracked %d dials)\n", tracker.uniqueConnections(), tracker.visits())

	if len(protoCounts) > 0 {
		fmt.Println("Protocols observed:")
		for proto, count := range protoCounts {
			fmt.Printf("  %s: %d\n", proto, count)
		}
	}

	if len(statusCounts) > 0 {
		fmt.Println("HTTP status distribution:")
		statuses := make([]int, 0, len(statusCounts))
		for status := range statusCounts {
			statuses = append(statuses, status)
		}
		sort.Ints(statuses)
		for _, status := range statuses {
			fmt.Printf("  %d: %d\n", status, statusCounts[status])
		}
	}

	if len(durations) > 0 {
		fmt.Println("Latency percentiles:")
		fmt.Printf("  p50: %s\n", percentile(durations, 0.50))
		fmt.Printf("  p90: %s\n", percentile(durations, 0.90))
		fmt.Printf("  p99: %s\n", percentile(durations, 0.99))
	}

	if len(errorsEncountered) > 0 {
		fmt.Printf("Errors (%d):\n", len(errorsEncountered))
		for _, err := range errorsEncountered {
			fmt.Printf("  %v\n", err)
		}
		os.Exit(1)
	}
}

func buildClient(protocol string, tracker *connTracker, perRequestTimeout time.Duration) (*http.Client, error) {
	switch protocol {
	case "http2":
		dialTLS := func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
			baseCfg := cfg
			if baseCfg == nil {
				baseCfg = &tls.Config{}
			}
			dialer := tls.Dialer{NetDialer: &net.Dialer{Timeout: 10 * time.Second}, Config: baseCfg}
			conn, err := dialer.DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}
			return tracker.record(conn), nil
		}

		transport := &http2.Transport{
			DialTLSContext:     dialTLS,
			AllowHTTP:          false,
			DisableCompression: false,
		}

		return &http.Client{
			Transport: transport,
			Timeout:   perRequestTimeout,
		}, nil

	case "http1":
		dialer := &net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}
		transport := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				conn, err := dialer.DialContext(ctx, network, addr)
				if err != nil {
					return nil, err
				}
				return tracker.record(conn), nil
			},
			ForceAttemptHTTP2:     false,
			MaxIdleConns:          0,
			MaxIdleConnsPerHost:   0,
			MaxConnsPerHost:       0,
			DisableKeepAlives:     true,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}

		return &http.Client{
			Transport: transport,
			Timeout:   perRequestTimeout,
		}, nil
	default:
		return nil, errors.New("unexpected protocol value")
	}
}

func percentile(vals []time.Duration, p float64) time.Duration {
	if len(vals) == 0 {
		return 0
	}
	if p <= 0 {
		return vals[0]
	}
	if p >= 1 {
		return vals[len(vals)-1]
	}
	pos := p * float64(len(vals)-1)
	idx := int(pos)
	frac := pos - float64(idx)
	if idx+1 >= len(vals) {
		return vals[idx]
	}
	lower := vals[idx]
	upper := vals[idx+1]
	return lower + time.Duration(frac*float64(upper-lower))
}

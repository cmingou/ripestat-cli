package ripestat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	httpSemMu sync.RWMutex
	httpSem   chan struct{} // nil = no limit
)

// SetMaxConcurrentHTTPRequests sets the global HTTP request concurrency limit.
func SetMaxConcurrentHTTPRequests(n int) {
	httpSemMu.Lock()
	defer httpSemMu.Unlock()
	if n <= 0 {
		httpSem = nil
		return
	}
	httpSem = make(chan struct{}, n)
}

func acquireHTTPSlot() {
	httpSemMu.RLock()
	sem := httpSem
	httpSemMu.RUnlock()
	if sem != nil {
		sem <- struct{}{}
	}
}

func releaseHTTPSlot() {
	httpSemMu.RLock()
	sem := httpSem
	httpSemMu.RUnlock()
	if sem != nil {
		<-sem
	}
}

func GetAsOverview(as int) (*AsOverview, error) {
	url := fmt.Sprintf("https://stat.ripe.net/data/as-overview/data.json?resource=AS%v", as)

	body, err := getHttpGetResponse(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to get HTTP GET response: %v", err)
	}

	var asOverview AsOverview
	err = json.Unmarshal(body, &asOverview)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal response body: %v", err)
	}

	return &asOverview, nil
}

func GetRIR(resource string) (*RIR, error) {
	url := fmt.Sprintf("https://stat.ripe.net/data/rir/data.json?resource=%v&lod=2", resource)

	body, err := getHttpGetResponse(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to get HTTP GET response: %v", err)
	}

	var rir RIR
	err = json.Unmarshal(body, &rir)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal response body: %v", err)
	}

	return &rir, nil
}

func GetPrefixRoutingConsistency(resource string) (*PrefixRoutingConsistency, error) {
	url := fmt.Sprintf("https://stat.ripe.net/data/prefix-routing-consistency/data.json?resource=%v", resource)

	body, err := getHttpGetResponse(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to get HTTP GET response: %v", err)
	}

	var prefixRoutingConsistency PrefixRoutingConsistency
	err = json.Unmarshal(body, &prefixRoutingConsistency)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal response body: %v", err)
	}

	return &prefixRoutingConsistency, nil
}

func GetMaxmindGeoLite(resource string) (*MaxmindGeoLite, error) {
	url := fmt.Sprintf("https://stat.ripe.net/data/maxmind-geo-lite/data.json?resource=%v", resource)

	body, err := getHttpGetResponse(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to get HTTP GET response: %v", err)
	}

	var maxmindGeoLite MaxmindGeoLite
	err = json.Unmarshal(body, &maxmindGeoLite)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal response body: %v", err)
	}

	return &maxmindGeoLite, nil
}

func GetIpGeoLocation(resource string) (string, error) {
	ipLocation := ""
	locationRsp, err := GetMaxmindGeoLite(resource)
	if err != nil {
		return "", fmt.Errorf("Failed to get Maxmind GeoLite: %v\n", err)
	}

	if locationRsp.Data.LocatedResources != nil && len(locationRsp.Data.LocatedResources) > 0 {
		if locationRsp.Data.LocatedResources[0].Locations != nil && len(locationRsp.Data.LocatedResources[0].Locations) > 0 {
			city := locationRsp.Data.LocatedResources[0].Locations[0].City
			country := locationRsp.Data.LocatedResources[0].Locations[0].Country
			if len(city) != 0 && len(country) != 0 {
				ipLocation = fmt.Sprintf("%v, %v", city, country)
			} else if len(city) == 0 && len(country) != 0 {
				ipLocation = fmt.Sprintf("%v", country)
			}
		}
	}
	return ipLocation, nil
}

func getHttpGetResponse(url string) ([]byte, error) {
	acquireHTTPSlot()
	defer releaseHTTPSlot()

	const maxRetries = 3
	client := getHTTPClient()

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(1<<(attempt-1)) * time.Second
			debugf("Retrying GET %s (attempt %d/%d) after %v", url, attempt+1, maxRetries+1, backoff)
			time.Sleep(backoff)
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("Failed to create HTTP GET request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("Failed to send request: %v", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("Failed to read response body: %v", err)
			continue
		}

		debugf("GET %s => %s (proto=%s)", url, resp.Status, resp.Proto)

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
			continue
		}

		if len(body) > 0 && bytes.HasPrefix(bytes.TrimSpace(body), []byte("<")) {
			lastErr = fmt.Errorf("received HTML response instead of JSON from %s", url)
			continue
		}

		return body, nil
	}

	return nil, lastErr
}

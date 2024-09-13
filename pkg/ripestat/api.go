package ripestat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

	if locationRsp.Data.LocatedResources != nil {
		if locationRsp.Data.LocatedResources[0].Locations != nil {
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
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create HTTP GET request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}

	return body, nil
}

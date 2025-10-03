package ripestat

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAsOverview(t *testing.T) {
	// Test with real API call
	result, err := GetAsOverview(13335)
	if err != nil {
		t.Fatalf("GetAsOverview failed: %v", err)
	}

	// Verify basic structure
	if result.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", result.Status)
	}
	if result.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", result.StatusCode)
	}
	if result.Data.Resource != "13335" {
		t.Errorf("Expected resource '13335', got '%s'", result.Data.Resource)
	}
	if result.Data.Holder != "CLOUDFLARENET" {
		t.Errorf("Expected holder 'CLOUDFLARENET', got '%s'", result.Data.Holder)
	}
	if result.Data.Type != "as" {
		t.Errorf("Expected type 'as', got '%s'", result.Data.Type)
	}
}

func TestGetAsOverview_MockServer(t *testing.T) {
	// Mock server with sample response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{
			"messages": [],
			"see_also": [],
			"version": "1.3",
			"data_call_name": "as-overview",
			"data_call_status": "supported - based on 2.1",
			"cached": false,
			"data": {
				"type": "as",
				"resource": "13335",
				"block": {
					"resource": "13312-15359",
					"desc": "Assigned by ARIN",
					"name": "IANA 16-bit Autonomous System (AS) Numbers Registry"
				},
				"holder": "CLOUDFLARENET",
				"announced": true,
				"query_starttime": "2025-06-13T16:00:00",
				"query_endtime": "2025-06-13T16:00:00"
			},
			"query_id": "test-query-id",
			"process_time": 129,
			"server_id": "test-server",
			"build_version": "test-version",
			"status": "ok",
			"status_code": 200,
			"time": "2025-06-14T01:07:42.224257"
		}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	// This test demonstrates the structure but uses real API since we can't easily mock the internal function
	// In a real-world scenario, you'd refactor the code to accept a client interface
	result, err := GetAsOverview(13335)
	if err != nil {
		t.Fatalf("GetAsOverview failed: %v", err)
	}
	if result.Data.Holder != "CLOUDFLARENET" {
		t.Errorf("Expected holder 'CLOUDFLARENET', got '%s'", result.Data.Holder)
	}
}

func TestHTTPClientReuse(t *testing.T) {
	client1 := getHTTPClient()
	if client1 == nil {
		t.Fatal("expected non-nil HTTP client")
	}

	client2 := getHTTPClient()
	if client1 != client2 {
		t.Error("expected getHTTPClient to return shared instance")
	}
}

func TestSetHTTPClient(t *testing.T) {
	original := getHTTPClient()
	if original == nil {
		t.Fatal("expected non-nil original HTTP client")
	}

	custom := &http.Client{}
	setHTTPClient(custom)
	t.Cleanup(func() {
		setHTTPClient(original)
	})

	if got := getHTTPClient(); got != custom {
		t.Errorf("expected custom client instance, got %p (want %p)", got, custom)
	}
}

func TestGetRIR(t *testing.T) {
	// Test with real API call using Google DNS IP
	result, err := GetRIR("8.8.8.8")
	if err != nil {
		t.Fatalf("GetRIR failed: %v", err)
	}

	// Verify basic structure
	if result.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", result.Status)
	}
	if result.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", result.StatusCode)
	}
	if len(result.Data.Rirs) == 0 {
		t.Error("Expected at least one RIR entry")
	} else {
		// Google DNS should be allocated by ARIN in US
		if result.Data.Rirs[0].Rir != "ARIN" {
			t.Errorf("Expected RIR 'ARIN', got '%s'", result.Data.Rirs[0].Rir)
		}
		if result.Data.Rirs[0].Country != "US" {
			t.Errorf("Expected country 'US', got '%s'", result.Data.Rirs[0].Country)
		}
	}
}

func TestGetRIR_ASN(t *testing.T) {
	// Test with ASN
	result, err := GetRIR("AS13335")
	if err != nil {
		t.Fatalf("GetRIR with ASN failed: %v", err)
	}

	if result.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", result.Status)
	}
	if len(result.Data.Rirs) == 0 {
		t.Error("Expected at least one RIR entry for ASN")
	}
}

func TestGetPrefixRoutingConsistency(t *testing.T) {
	// Test with real API call
	result, err := GetPrefixRoutingConsistency("8.8.8.8")
	if err != nil {
		t.Fatalf("GetPrefixRoutingConsistency failed: %v", err)
	}

	// Verify basic structure
	if result.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", result.Status)
	}
	if result.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", result.StatusCode)
	}
	if result.Data.Resource != "8.8.8.8" {
		t.Errorf("Expected resource '8.8.8.8', got '%s'", result.Data.Resource)
	}
	if len(result.Data.Routes) == 0 {
		t.Error("Expected at least one route entry")
	} else {
		// Verify first route contains Google's ASN
		found := false
		for _, route := range result.Data.Routes {
			if route.Origin == 15169 { // Google's ASN
				found = true
				if !route.InBgp {
					t.Error("Expected Google's route to be in BGP")
				}
				if route.Prefix != "8.8.8.0/24" {
					t.Errorf("Expected prefix '8.8.8.0/24', got '%s'", route.Prefix)
				}
				break
			}
		}
		if !found {
			t.Error("Expected to find Google's ASN 15169 in routes")
		}
	}
}

func TestGetMaxmindGeoLite(t *testing.T) {
	// Test with real API call
	result, err := GetMaxmindGeoLite("8.8.8.8")
	if err != nil {
		t.Fatalf("GetMaxmindGeoLite failed: %v", err)
	}

	// Verify basic structure
	if result.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", result.Status)
	}
	if result.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", result.StatusCode)
	}
	if len(result.Data.LocatedResources) == 0 {
		t.Error("Expected at least one located resource")
	} else {
		if len(result.Data.LocatedResources[0].Locations) == 0 {
			t.Error("Expected at least one location")
		} else {
			// Google DNS should be located in US
			location := result.Data.LocatedResources[0].Locations[0]
			if location.Country != "US" {
				t.Errorf("Expected country 'US', got '%s'", location.Country)
			}
		}
	}
}

func TestGetMaxmindGeoLite_IPv6(t *testing.T) {
	// Test with IPv6 address
	result, err := GetMaxmindGeoLite("2001:4860:4860::8888")
	if err != nil {
		t.Fatalf("GetMaxmindGeoLite with IPv6 failed: %v", err)
	}

	if result.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", result.Status)
	}
}

func TestGetIpGeoLocation(t *testing.T) {
	// Test with real API call
	location, err := GetIpGeoLocation("8.8.8.8")
	if err != nil {
		t.Fatalf("GetIpGeoLocation failed: %v", err)
	}

	// Google DNS should return US location
	if location != "US" && location == "" {
		t.Errorf("Expected non-empty location for Google DNS, got '%s'", location)
	}

	// The location should contain "US" somewhere
	if len(location) > 0 && location[len(location)-2:] != "US" {
		t.Errorf("Expected location to end with 'US', got '%s'", location)
	}
}

func TestGetIpGeoLocation_WithCity(t *testing.T) {
	// Test with an IP that typically has city information
	// Using a different IP that might have city data
	location, err := GetIpGeoLocation("1.1.1.1")
	if err != nil {
		t.Fatalf("GetIpGeoLocation failed: %v", err)
	}

	// Should return some location information
	if location == "" {
		t.Error("Expected non-empty location for 1.1.1.1")
	}
}

func TestGetIpGeoLocation_IPv6(t *testing.T) {
	// Test with IPv6
	location, err := GetIpGeoLocation("2001:4860:4860::8888")
	if err != nil {
		t.Fatalf("GetIpGeoLocation with IPv6 failed: %v", err)
	}

	// Should return some location information
	if location == "" {
		t.Error("Expected non-empty location for Google IPv6 DNS")
	}
}

func TestGetHttpGetResponse(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"test": "response"}`))
	}))
	defer server.Close()

	// Test the function
	body, err := getHttpGetResponse(server.URL)
	if err != nil {
		t.Fatalf("getHttpGetResponse failed: %v", err)
	}

	expected := `{"test": "response"}`
	if string(body) != expected {
		t.Errorf("Expected response '%s', got '%s'", expected, string(body))
	}
}

func TestGetHttpGetResponse_InvalidURL(t *testing.T) {
	// Test with invalid URL
	_, err := getHttpGetResponse("invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestGetHttpGetResponse_404(t *testing.T) {
	// Create a test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}))
	defer server.Close()

	// Test the function - it should still work as it doesn't check status codes
	body, err := getHttpGetResponse(server.URL)
	if err != nil {
		t.Fatalf("getHttpGetResponse failed: %v", err)
	}

	if string(body) != "Not Found" {
		t.Errorf("Expected 'Not Found', got '%s'", string(body))
	}
}

// Integration tests that verify the complete flow
func TestAPIIntegration_Cloudflare(t *testing.T) {
	// Test Cloudflare ASN
	asn := 13335

	// Test AS Overview
	asOverview, err := GetAsOverview(asn)
	if err != nil {
		t.Fatalf("GetAsOverview failed: %v", err)
	}

	// Test RIR for same ASN
	rir, err := GetRIR("AS13335")
	if err != nil {
		t.Fatalf("GetRIR failed: %v", err)
	}

	// Verify consistency
	if asOverview.Data.Holder != "CLOUDFLARENET" {
		t.Errorf("Expected Cloudflare holder name, got '%s'", asOverview.Data.Holder)
	}

	if len(rir.Data.Rirs) == 0 {
		t.Error("Expected RIR data for Cloudflare ASN")
	}
}

func TestAPIIntegration_GoogleDNS(t *testing.T) {
	// Test Google DNS IP
	ip := "8.8.8.8"

	// Test all IP-related APIs
	rir, err := GetRIR(ip)
	if err != nil {
		t.Fatalf("GetRIR failed: %v", err)
	}

	prc, err := GetPrefixRoutingConsistency(ip)
	if err != nil {
		t.Fatalf("GetPrefixRoutingConsistency failed: %v", err)
	}

	geo, err := GetMaxmindGeoLite(ip)
	if err != nil {
		t.Fatalf("GetMaxmindGeoLite failed: %v", err)
	}

	location, err := GetIpGeoLocation(ip)
	if err != nil {
		t.Fatalf("GetIpGeoLocation failed: %v", err)
	}

	// Verify consistency across APIs
	if len(rir.Data.Rirs) == 0 {
		t.Error("Expected RIR data for Google DNS")
	}

	if len(prc.Data.Routes) == 0 {
		t.Error("Expected routing data for Google DNS")
	}

	if len(geo.Data.LocatedResources) == 0 {
		t.Error("Expected geo data for Google DNS")
	}

	if location == "" {
		t.Error("Expected location string for Google DNS")
	}

	// All should indicate US location
	if rir.Data.Rirs[0].Country != "US" {
		t.Errorf("RIR should show US, got '%s'", rir.Data.Rirs[0].Country)
	}

	if geo.Data.LocatedResources[0].Locations[0].Country != "US" {
		t.Errorf("Geo should show US, got '%s'", geo.Data.LocatedResources[0].Locations[0].Country)
	}
}

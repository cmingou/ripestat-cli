package utils

import (
	"bytes"
	"net/netip"
	"os"
	"strings"
	"testing"
)

// TestSearchAsnInfo tests ASN information search functionality
func TestSearchAsnInfo(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test with Cloudflare ASN
	asns := []int{13335}
	
	// This will call the real API and print to stdout
	SearchAsnInfo(asns)

	// Restore stdout and capture output
	w.Close()
	os.Stdout = oldStdout
	
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify expected content
	if !strings.Contains(output, "13335") {
		t.Errorf("Expected ASN 13335 in output, got: %s", output)
	}
	if !strings.Contains(output, "CLOUDFLARENET") {
		t.Errorf("Expected CLOUDFLARENET in output, got: %s", output)
	}
	if !strings.Contains(output, "ARIN") {
		t.Errorf("Expected ARIN in output, got: %s", output)
	}

	// Verify table format
	if !strings.Contains(output, "|") {
		t.Errorf("Expected table format with '|' separators, got: %s", output)
	}
	if !strings.Contains(output, "AS NAME") {
		t.Errorf("Expected table header 'AS NAME', got: %s", output)
	}
}

// TestSearchIpv4Info tests IPv4 information search functionality  
func TestSearchIpv4Info(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test with Google DNS
	ip, _ := netip.ParseAddr("8.8.8.8")
	ips := []netip.Addr{ip}
	
	SearchIpv4Info(ips)

	// Restore stdout and capture output
	w.Close()
	os.Stdout = oldStdout
	
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify expected content
	if !strings.Contains(output, "8.8.8.8") {
		t.Errorf("Expected IP 8.8.8.8 in output, got: %s", output)
	}
	if !strings.Contains(output, "15169") {
		t.Errorf("Expected Google ASN 15169 in output, got: %s", output)
	}
	if !strings.Contains(output, "GOOGLE") {
		t.Errorf("Expected GOOGLE in output, got: %s", output)
	}

	// Verify table format and headers
	expectedHeaders := []string{"IP", "LOCATION", "PREFIX", "IN BGP", "AS NUMBER", "AS NAME"}
	for _, header := range expectedHeaders {
		if !strings.Contains(output, header) {
			t.Errorf("Expected table header '%s', got: %s", header, output)
		}
	}
}

// TestSearchIpv6Info tests IPv6 information search functionality
func TestSearchIpv6Info(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test with Google IPv6 DNS
	ip, _ := netip.ParseAddr("2001:4860:4860::8888")
	ips := []netip.Addr{ip}
	
	SearchIpv6Info(ips)

	// Restore stdout and capture output
	w.Close()
	os.Stdout = oldStdout
	
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify expected content
	if !strings.Contains(output, "2001:4860:4860::8888") {
		t.Errorf("Expected IPv6 2001:4860:4860::8888 in output, got: %s", output)
	}
	if !strings.Contains(output, "15169") {
		t.Errorf("Expected Google ASN 15169 in output, got: %s", output)
	}

	// Verify table format
	if !strings.Contains(output, "|") {
		t.Errorf("Expected table format with '|' separators, got: %s", output)
	}
}

// TestMultipleAsns tests handling multiple ASNs
func TestMultipleAsns(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test with multiple ASNs: Cloudflare and Google
	asns := []int{13335, 15169}
	
	SearchAsnInfo(asns)

	// Restore stdout and capture output
	w.Close()
	os.Stdout = oldStdout
	
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify both ASNs are present
	if !strings.Contains(output, "13335") {
		t.Errorf("Expected Cloudflare ASN 13335 in output")
	}
	if !strings.Contains(output, "15169") {
		t.Errorf("Expected Google ASN 15169 in output")
	}
	if !strings.Contains(output, "CLOUDFLARENET") {
		t.Errorf("Expected CLOUDFLARENET in output")
	}
	if !strings.Contains(output, "GOOGLE") {
		t.Errorf("Expected GOOGLE in output")
	}

	// Count number of data rows (should be 2)
	lines := strings.Split(output, "\n")
	dataRows := 0
	for _, line := range lines {
		// Count lines that contain ASN data (contain both number and text)
		if strings.Contains(line, "|") && (strings.Contains(line, "13335") || strings.Contains(line, "15169")) {
			dataRows++
		}
	}
	if dataRows < 2 {
		t.Errorf("Expected at least 2 data rows for 2 ASNs, got %d", dataRows)
	}
}

// TestMultipleIPs tests handling multiple IP addresses
func TestMultipleIPs(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test with multiple IPs: Google DNS and Cloudflare DNS
	ip1, _ := netip.ParseAddr("8.8.8.8")
	ip2, _ := netip.ParseAddr("1.1.1.1")
	ips := []netip.Addr{ip1, ip2}
	
	SearchIpv4Info(ips)

	// Restore stdout and capture output
	w.Close()
	os.Stdout = oldStdout
	
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify both IPs are present
	if !strings.Contains(output, "8.8.8.8") {
		t.Errorf("Expected Google DNS 8.8.8.8 in output")
	}
	if !strings.Contains(output, "1.1.1.1") {
		t.Errorf("Expected Cloudflare DNS 1.1.1.1 in output")
	}

	// Should contain information about both providers
	if !strings.Contains(output, "GOOGLE") && !strings.Contains(output, "15169") {
		t.Errorf("Expected Google information in output")
	}
	if !strings.Contains(output, "CLOUDFLARE") && !strings.Contains(output, "13335") {
		t.Errorf("Expected Cloudflare information in output")
	}
}

// TestPrintInvalidArgs tests invalid argument handling
func TestPrintInvalidArgs(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test with invalid arguments
	invalidArgs := []string{"invalid", "not-an-ip", "999.999.999.999"}
	
	PrintInvalidArgs(invalidArgs)

	// Restore stdout and capture output
	w.Close()
	os.Stdout = oldStdout
	
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify invalid arguments are shown
	for _, invalid := range invalidArgs {
		if !strings.Contains(output, invalid) {
			t.Errorf("Expected invalid argument '%s' in output, got: %s", invalid, output)
		}
	}

	// Should contain "Invalid" message
	if !strings.Contains(output, "Invalid") {
		t.Errorf("Expected 'Invalid' message in output, got: %s", output)
	}
}

// TestTableFormatConsistency tests that all tables have consistent formatting
func TestTableFormatConsistency(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func() string
	}{
		{
			name: "ASN table format",
			testFunc: func() string {
				oldStdout := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w
				
				SearchAsnInfo([]int{13335})
				
				w.Close()
				os.Stdout = oldStdout
				var buf bytes.Buffer
				buf.ReadFrom(r)
				return buf.String()
			},
		},
		{
			name: "IPv4 table format",
			testFunc: func() string {
				oldStdout := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w
				
				ip, _ := netip.ParseAddr("8.8.8.8")
				SearchIpv4Info([]netip.Addr{ip})
				
				w.Close()
				os.Stdout = oldStdout
				var buf bytes.Buffer
				buf.ReadFrom(r)
				return buf.String()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.testFunc()
			
			// All tables should use | as separator
			if !strings.Contains(output, "|") {
				t.Errorf("Expected table to use '|' separator")
			}
			
			// Should have proper table structure (multiple lines)
			lines := strings.Split(output, "\n")
			nonEmptyLines := 0
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					nonEmptyLines++
				}
			}
			if nonEmptyLines < 2 {
				t.Errorf("Expected at least 2 non-empty lines for table, got %d", nonEmptyLines)
			}
		})
	}
}

// TestDataAccuracy verifies the accuracy of returned data
func TestDataAccuracy(t *testing.T) {
	// Test known ASN data
	t.Run("Cloudflare ASN accuracy", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		
		SearchAsnInfo([]int{13335})
		
		w.Close()
		os.Stdout = oldStdout
		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		// Known facts about Cloudflare's ASN
		expectedData := map[string]bool{
			"13335":         false,
			"CLOUDFLARENET": false,
			"ARIN":          false,
		}

		for expected := range expectedData {
			if strings.Contains(output, expected) {
				expectedData[expected] = true
			}
		}

		for expected, found := range expectedData {
			if !found {
				t.Errorf("Expected accurate data '%s' not found in ASN output", expected)
			}
		}
	})

	// Test known IP data
	t.Run("Google DNS accuracy", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		
		ip, _ := netip.ParseAddr("8.8.8.8")
		SearchIpv4Info([]netip.Addr{ip})
		
		w.Close()
		os.Stdout = oldStdout
		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		// Known facts about Google DNS
		expectedData := map[string]bool{
			"8.8.8.8": false,
			"15169":   false,
			"GOOGLE":  false,
		}

		for expected := range expectedData {
			if strings.Contains(output, expected) {
				expectedData[expected] = true
			}
		}

		for expected, found := range expectedData {
			if !found {
				t.Errorf("Expected accurate data '%s' not found in IP output", expected)
			}
		}
	})
}
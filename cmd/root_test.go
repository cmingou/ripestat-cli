package cmd

import (
	"testing"
)

func TestIsASN(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid ASNs
		{"Valid ASN - 16509 (Amazon)", "16509", true},
		{"Valid ASN - 13335 (Cloudflare)", "13335", true},
		{"Valid ASN - 15169 (Google)", "15169", true},
		{"Valid ASN - 8075 (Microsoft)", "8075", true},
		{"Valid ASN - Min value", "0", true},
		{"Valid ASN - Max value", "4294967295", true},
		{"Valid ASN - 32768", "32768", true},
		
		// Invalid ASNs
		{"Invalid ASN - Negative", "-1", false},
		{"Invalid ASN - Above max", "4294967296", false},
		{"Invalid ASN - Non-numeric", "abc", false},
		{"Invalid ASN - Empty string", "", false},
		{"Invalid ASN - Mixed chars", "123abc", false},
		{"Invalid ASN - Float", "123.45", false},
		{"Invalid ASN - With spaces", " 16509 ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isASN(tt.input)
			if result != tt.expected {
				t.Errorf("isASN(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsIPv4(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid public IPv4 addresses
		{"Valid IPv4 - Google DNS", "8.8.8.8", true},
		{"Valid IPv4 - Cloudflare DNS", "1.1.1.1", true},
		{"Valid IPv4 - OpenDNS", "208.67.222.222", true},
		{"Valid IPv4 - Quad9", "9.9.9.9", true},
		
		// Valid private IPv4 addresses (RFC 1918)
		{"Valid IPv4 - Private Class A", "10.0.0.1", true},
		{"Valid IPv4 - Private Class A Max", "10.255.255.254", true},
		{"Valid IPv4 - Private Class B", "172.16.0.1", true},
		{"Valid IPv4 - Private Class B Mid", "172.20.1.100", true},
		{"Valid IPv4 - Private Class B Max", "172.31.255.254", true},
		{"Valid IPv4 - Private Class C", "192.168.1.1", true},
		{"Valid IPv4 - Private Class C Alt", "192.168.0.254", true},
		
		// Loopback and special addresses
		{"Valid IPv4 - Localhost", "127.0.0.1", true},
		{"Valid IPv4 - Broadcast", "255.255.255.255", true},
		{"Valid IPv4 - Network zero", "0.0.0.0", true},
		
		// Link-local addresses (RFC 3927)
		{"Valid IPv4 - Link Local", "169.254.1.1", true},
		
		// Invalid IPv4 addresses
		{"Invalid IPv4 - Out of range octet", "256.1.1.1", false},
		{"Invalid IPv4 - Negative octet", "-1.1.1.1", false},
		{"Invalid IPv4 - Too few octets", "192.168.1", false},
		{"Invalid IPv4 - Too many octets", "192.168.1.1.1", false},
		{"Invalid IPv4 - Non-numeric", "192.168.a.1", false},
		{"Invalid IPv4 - Empty string", "", false},
		{"Invalid IPv4 - Only dots", "...", false},
		{"Invalid IPv4 - IPv6 address", "2001:db8::1", false},
		{"Invalid IPv4 - With port", "192.168.1.1:80", false},
		{"Invalid IPv4 - Leading zeros", "192.168.001.001", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isIPv4(tt.input)
			if result != tt.expected {
				t.Errorf("isIPv4(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsIPv6(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid IPv6 addresses
		{"Valid IPv6 - Full address", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", true},
		{"Valid IPv6 - Compressed", "2001:db8:85a3::8a2e:370:7334", true},
		{"Valid IPv6 - Loopback", "::1", true},
		{"Valid IPv6 - All zeros", "::", true},
		{"Valid IPv6 - Google DNS", "2001:4860:4860::8888", true},
		{"Valid IPv6 - Cloudflare DNS", "2606:4700:4700::1111", true},
		{"Valid IPv6 - Link local", "fe80::1", true},
		{"Valid IPv6 - Unique local", "fc00::1", true},
		{"Valid IPv6 - Documentation", "2001:db8::1", true},
		
		// IPv4-mapped IPv6 addresses
		{"Valid IPv6 - IPv4 mapped", "::ffff:192.168.1.1", true},
		{"Valid IPv6 - IPv4 compatible", "::192.168.1.1", true},
		
		// Invalid IPv6 addresses
		{"Invalid IPv6 - Too many groups", "2001:0db8:85a3:0000:0000:8a2e:0370:7334:extra", false},
		{"Invalid IPv6 - Invalid characters", "2001:0db8:85a3:0000:0000:8a2e:0370:733g", false},
		{"Invalid IPv6 - Empty string", "", false},
		{"Invalid IPv6 - IPv4 address", "192.168.1.1", false},
		{"Invalid IPv6 - Multiple double colons", "2001::db8::1", false},
		{"Invalid IPv6 - Leading colon only", ":2001:db8::1", false},
		{"Invalid IPv6 - Trailing colon only", "2001:db8::1:", false},
		{"Invalid IPv6 - Too long group", "2001:0db8:85a30:0000:0000:8a2e:0370:7334", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isIPv6(tt.input)
			if result != tt.expected {
				t.Errorf("isIPv6(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestInputClassification(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedASN bool
		expectedV4  bool
		expectedV6  bool
	}{
		// Mixed classification tests
		{"ASN 16509", "16509", true, false, false},
		{"IPv4 Google DNS", "8.8.8.8", false, true, false},
		{"IPv6 Google DNS", "2001:4860:4860::8888", false, false, true},
		{"Private IPv4", "192.168.1.1", false, true, false},
		{"Private IPv4 Class A", "10.0.0.1", false, true, false},
		{"IPv6 loopback", "::1", false, false, true},
		{"Invalid input", "invalid", false, false, false},
		{"Empty input", "", false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asnResult := isASN(tt.input)
			ipv4Result := isIPv4(tt.input)
			ipv6Result := isIPv6(tt.input)

			if asnResult != tt.expectedASN {
				t.Errorf("isASN(%q) = %v, expected %v", tt.input, asnResult, tt.expectedASN)
			}
			if ipv4Result != tt.expectedV4 {
				t.Errorf("isIPv4(%q) = %v, expected %v", tt.input, ipv4Result, tt.expectedV4)
			}
			if ipv6Result != tt.expectedV6 {
				t.Errorf("isIPv6(%q) = %v, expected %v", tt.input, ipv6Result, tt.expectedV6)
			}

			// Ensure mutual exclusivity (an input can only be one type)
			trueCount := 0
			if asnResult {
				trueCount++
			}
			if ipv4Result {
				trueCount++
			}
			if ipv6Result {
				trueCount++
			}
			
			if trueCount > 1 {
				t.Errorf("Input %q classified as multiple types: ASN=%v, IPv4=%v, IPv6=%v", 
					tt.input, asnResult, ipv4Result, ipv6Result)
			}
		})
	}
}
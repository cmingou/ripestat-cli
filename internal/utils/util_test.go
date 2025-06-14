package utils

import (
	"net/netip"
	"testing"
)

func TestCheckArgsNonExist(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{"Empty args", []string{}, true},
		{"Nil args", nil, true},
		{"Single arg", []string{"8.8.8.8"}, false},
		{"Multiple args", []string{"8.8.8.8", "16509"}, false},
		{"Empty string arg", []string{""}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckArgsNonExist(tt.args)
			if result != tt.expected {
				t.Errorf("CheckArgsNonExist(%v) = %v, expected %v", tt.args, result, tt.expected)
			}
		})
	}
}

func TestCnovertStringToAsn(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedAsn int
		expectError bool
	}{
		// Valid ASNs
		{"Valid ASN - 16509", "16509", 16509, false},
		{"Valid ASN - 13335", "13335", 13335, false},
		{"Valid ASN - 15169", "15169", 15169, false},
		{"Valid ASN - 1", "1", 1, false},
		{"Valid ASN - 65535", "65535", 65535, false},
		{"Valid ASN - 4294967295", "4294967295", 4294967295, false},
		
		// Invalid ASNs
		{"Invalid ASN - Zero", "0", 0, true},
		{"Invalid ASN - Negative", "-1", 0, true},
		{"Invalid ASN - Non-numeric", "abc", 0, true},
		{"Invalid ASN - Empty", "", 0, true},
		{"Invalid ASN - Float", "123.45", 0, true},
		{"Invalid ASN - Mixed", "123abc", 0, true},
		{"Invalid ASN - Spaces", " 16509 ", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CnovertStringToAsn(tt.input)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("CnovertStringToAsn(%q) expected error but got none", tt.input)
				}
				if result != 0 {
					t.Errorf("CnovertStringToAsn(%q) expected 0 on error but got %d", tt.input, result)
				}
			} else {
				if err != nil {
					t.Errorf("CnovertStringToAsn(%q) unexpected error: %v", tt.input, err)
				}
				if result != tt.expectedAsn {
					t.Errorf("CnovertStringToAsn(%q) = %d, expected %d", tt.input, result, tt.expectedAsn)
				}
			}
		})
	}
}

func TestCnovertStringToIp(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		// Valid IPv4 addresses
		{"Valid IPv4 - Google DNS", "8.8.8.8", false},
		{"Valid IPv4 - Cloudflare DNS", "1.1.1.1", false},
		{"Valid IPv4 - Private Class A", "10.0.0.1", false},
		{"Valid IPv4 - Private Class B", "172.16.0.1", false},
		{"Valid IPv4 - Private Class C", "192.168.1.1", false},
		{"Valid IPv4 - Localhost", "127.0.0.1", false},
		{"Valid IPv4 - Link Local", "169.254.1.1", false},
		{"Valid IPv4 - Broadcast", "255.255.255.255", false},
		{"Valid IPv4 - Zero", "0.0.0.0", false},
		
		// Valid IPv6 addresses
		{"Valid IPv6 - Full", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", false},
		{"Valid IPv6 - Compressed", "2001:db8:85a3::8a2e:370:7334", false},
		{"Valid IPv6 - Loopback", "::1", false},
		{"Valid IPv6 - All zeros", "::", false},
		{"Valid IPv6 - Google DNS", "2001:4860:4860::8888", false},
		{"Valid IPv6 - Cloudflare DNS", "2606:4700:4700::1111", false},
		{"Valid IPv6 - Link local", "fe80::1", false},
		{"Valid IPv6 - IPv4 mapped", "::ffff:192.168.1.1", false},
		
		// Invalid IP addresses
		{"Invalid IP - Out of range IPv4", "256.1.1.1", true},
		{"Invalid IP - Negative IPv4", "-1.1.1.1", true},
		{"Invalid IP - Too few octets", "192.168.1", true},
		{"Invalid IP - Too many octets", "192.168.1.1.1", true},
		{"Invalid IP - Non-numeric IPv4", "192.168.a.1", true},
		{"Invalid IP - Invalid IPv6", "2001::db8::1", true},
		{"Invalid IP - Empty string", "", true},
		{"Invalid IP - Non-IP text", "not-an-ip", true},
		{"Invalid IP - ASN", "16509", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CnovertStringToIp(tt.input)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("CnovertStringToIp(%q) expected error but got none", tt.input)
				}
				if result.IsValid() {
					t.Errorf("CnovertStringToIp(%q) expected invalid IP on error but got valid: %v", tt.input, result)
				}
			} else {
				if err != nil {
					t.Errorf("CnovertStringToIp(%q) unexpected error: %v", tt.input, err)
				}
				if !result.IsValid() {
					t.Errorf("CnovertStringToIp(%q) expected valid IP but got invalid", tt.input)
				}
				
				// Verify the result matches the input
				expectedAddr, _ := netip.ParseAddr(tt.input)
				if result != expectedAddr {
					t.Errorf("CnovertStringToIp(%q) = %v, expected %v", tt.input, result, expectedAddr)
				}
			}
		})
	}
}

func TestIPAddressTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		isV4     bool
		isV6     bool
		isPrivate bool
		isLoopback bool
	}{
		// IPv4 tests
		{"IPv4 Public - Google DNS", "8.8.8.8", true, false, false, false},
		{"IPv4 Public - Cloudflare", "1.1.1.1", true, false, false, false},
		{"IPv4 Private - Class A", "10.0.0.1", true, false, true, false},
		{"IPv4 Private - Class B", "172.16.0.1", true, false, true, false},
		{"IPv4 Private - Class C", "192.168.1.1", true, false, true, false},
		{"IPv4 Loopback", "127.0.0.1", true, false, false, true},
		{"IPv4 Link Local", "169.254.1.1", true, false, false, false},
		
		// IPv6 tests
		{"IPv6 Public - Google DNS", "2001:4860:4860::8888", false, true, false, false},
		{"IPv6 Loopback", "::1", false, true, false, true},
		{"IPv6 Link Local", "fe80::1", false, true, false, false},
		{"IPv6 Unique Local", "fc00::1", false, true, true, false},
		{"IPv6 Documentation", "2001:db8::1", false, true, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := CnovertStringToIp(tt.input)
			if err != nil {
				t.Fatalf("CnovertStringToIp(%q) failed: %v", tt.input, err)
			}

			if addr.Is4() != tt.isV4 {
				t.Errorf("Address %q Is4() = %v, expected %v", tt.input, addr.Is4(), tt.isV4)
			}
			
			if addr.Is6() != tt.isV6 {
				t.Errorf("Address %q Is6() = %v, expected %v", tt.input, addr.Is6(), tt.isV6)
			}
			
			if addr.IsPrivate() != tt.isPrivate {
				t.Errorf("Address %q IsPrivate() = %v, expected %v", tt.input, addr.IsPrivate(), tt.isPrivate)
			}
			
			if addr.IsLoopback() != tt.isLoopback {
				t.Errorf("Address %q IsLoopback() = %v, expected %v", tt.input, addr.IsLoopback(), tt.isLoopback)
			}
		})
	}
}
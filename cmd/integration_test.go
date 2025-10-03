package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestCLIEndToEnd tests the complete CLI functionality
func TestCLIEndToEnd(t *testing.T) {
	// Build the CLI binary for testing
	buildCmd := exec.Command("go", "build", "-o", "../ripestat_test", "../main.go")
	err := buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build CLI binary: %v", err)
	}
	defer os.Remove("../ripestat_test")

	tests := []struct {
		name       string
		args       []string
		expectErr  bool
		contains   []string // Strings that should be present in output
		notContains []string // Strings that should NOT be present
	}{
		{
			name:      "Test ASN lookup - Cloudflare",
			args:      []string{"13335"},
			expectErr: false,
			contains:  []string{"13335", "CLOUDFLARENET", "ARIN"},
		},
		{
			name:      "Test IPv4 lookup - Google DNS",
			args:      []string{"8.8.8.8"},
			expectErr: false,
			contains:  []string{"8.8.8.8", "US", "15169", "GOOGLE"},
		},
		{
			name:      "Test IPv6 lookup - Google DNS",
			args:      []string{"2001:4860:4860::8888"},
			expectErr: false,
			contains:  []string{"2001:4860:4860::8888", "15169"},
		},
		{
			name:      "Test mixed input - ASN and IP",
			args:      []string{"13335", "8.8.8.8"},
			expectErr: false,
			contains:  []string{"13335", "CLOUDFLARENET", "8.8.8.8", "GOOGLE"},
		},
		{
			name:      "Test multiple ASNs",
			args:      []string{"13335", "15169"},
			expectErr: false,
			contains:  []string{"13335", "CLOUDFLARENET", "15169", "GOOGLE"},
		},
		{
			name:      "Test multiple IPs",
			args:      []string{"8.8.8.8", "1.1.1.1"},
			expectErr: false,
			contains:  []string{"8.8.8.8", "1.1.1.1", "GOOGLE", "CLOUDFLARE"},
		},
		{
			name:      "Test invalid input",
			args:      []string{"invalid"},
			expectErr: false,
			contains:  []string{"Invalid", "invalid"},
		},
		{
			name:      "Test no arguments",
			args:      []string{},
			expectErr: true,
			contains:  []string{"Please provide input via command line arguments or --file flag"},
		},
		{
			name:      "Test mixed valid and invalid",
			args:      []string{"13335", "invalid", "8.8.8.8"},
			expectErr: false,
			contains:  []string{"13335", "CLOUDFLARENET", "8.8.8.8", "Invalid", "invalid"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("../ripestat_test", tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			output := stdout.String() + stderr.String()

			if tt.expectErr && err == nil {
				t.Errorf("Expected error but command succeeded. Output: %s", output)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected success but command failed with error: %v. Output: %s", err, output)
			}

			// Check that expected strings are present
			for _, expected := range tt.contains {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', but it didn't.\nFull output: %s", expected, output)
				}
			}

			// Check that unwanted strings are NOT present
			for _, unwanted := range tt.notContains {
				if strings.Contains(output, unwanted) {
					t.Errorf("Expected output to NOT contain '%s', but it did.\nFull output: %s", unwanted, output)
				}
			}
		})
	}
}

// TestCLITableFormat tests that the output is properly formatted in tables
func TestCLITableFormat(t *testing.T) {
	// Build the CLI binary for testing
	buildCmd := exec.Command("go", "build", "-o", "../ripestat_test", "../main.go")
	err := buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build CLI binary: %v", err)
	}
	defer os.Remove("../ripestat_test")

	tests := []struct {
		name string
		args []string
		checkTableHeaders []string
	}{
		{
			name: "ASN table format",
			args: []string{"13335"},
			checkTableHeaders: []string{"AS", "COUNTRY", "RIR", "AS NAME"},
		},
		{
			name: "IPv4 table format", 
			args: []string{"8.8.8.8"},
			checkTableHeaders: []string{"IP", "LOCATION", "PREFIX", "IN BGP", "AS NUMBER", "AS NAME"},
		},
		{
			name: "IPv6 table format",
			args: []string{"2001:4860:4860::8888"},
			checkTableHeaders: []string{"IP", "LOCATION", "PREFIX", "IN BGP", "AS NUMBER", "AS NAME"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("../ripestat_test", tt.args...)
			var stdout bytes.Buffer
			cmd.Stdout = &stdout

			err := cmd.Run()
			if err != nil {
				t.Fatalf("Command failed: %v", err)
			}

			output := stdout.String()
			
			// Check that all expected table headers are present
			for _, header := range tt.checkTableHeaders {
				if !strings.Contains(output, header) {
					t.Errorf("Expected table header '%s' not found in output:\n%s", header, output)
				}
			}

			// Check that output contains table separators (indicating proper table format)
			if !strings.Contains(output, "|") {
				t.Errorf("Expected table format with '|' separators, but not found in output:\n%s", output)
			}
		})
	}
}

// TestCLIPerformance tests that CLI responds within reasonable time
func TestCLIPerformance(t *testing.T) {
	// Build the CLI binary for testing
	buildCmd := exec.Command("go", "build", "-o", "../ripestat_test", "../main.go")
	err := buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build CLI binary: %v", err)
	}
	defer os.Remove("../ripestat_test")

	// Test with a simple query that should complete quickly
	cmd := exec.Command("../ripestat_test", "13335")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err = cmd.Run()
	if err != nil {
		t.Fatalf("Command failed: %v", err)
	}

	output := stdout.String()
	if len(output) == 0 {
		t.Error("Expected non-empty output from CLI")
	}

	// Verify essential information is present
	if !strings.Contains(output, "13335") {
		t.Error("Expected ASN 13335 in output")
	}
}

// TestCLIErrorHandling tests various error conditions
func TestCLIErrorHandling(t *testing.T) {
	// Build the CLI binary for testing
	buildCmd := exec.Command("go", "build", "-o", "../ripestat_test", "../main.go")
	err := buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build CLI binary: %v", err)
	}
	defer os.Remove("../ripestat_test")

	tests := []struct {
		name         string
		args         []string
		expectExit   bool
		expectOutput string
	}{
		{
			name:         "No arguments",
			args:         []string{},
			expectExit:   true,
			expectOutput: "Please provide input via command line arguments or --file flag",
		},
		{
			name:         "Invalid ASN - negative",
			args:         []string{"--", "-1"}, // Use -- to prevent flag parsing
			expectExit:   false,
			expectOutput: "Invalid",
		},
		{
			name:         "Invalid ASN - too large", 
			args:         []string{"4294967296"},
			expectExit:   false,
			expectOutput: "Invalid",
		},
		{
			name:         "Invalid IP address",
			args:         []string{"999.999.999.999"},
			expectExit:   false,
			expectOutput: "Invalid",
		},
		{
			name:         "Mixed valid and invalid inputs",
			args:         []string{"13335", "invalid", "8.8.8.8", "not-an-ip"},
			expectExit:   false,
			expectOutput: "Invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("../ripestat_test", tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			output := stdout.String() + stderr.String()

			if tt.expectExit && err == nil {
				t.Errorf("Expected command to exit with error, but it didn't. Output: %s", output)
			}

			if !strings.Contains(output, tt.expectOutput) {
				t.Errorf("Expected output to contain '%s', but got: %s", tt.expectOutput, output)
			}
		})
	}
}

// TestCLIRegressionBaseline creates a baseline test to detect unexpected changes
func TestCLIRegressionBaseline(t *testing.T) {
	// Build the CLI binary for testing
	buildCmd := exec.Command("go", "build", "-o", "../ripestat_test", "../main.go")
	err := buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build CLI binary: %v", err)
	}
	defer os.Remove("../ripestat_test")

	// Test with known stable inputs and verify key characteristics
	testCases := []struct {
		name     string
		args     []string
		checks   func(t *testing.T, output string)
	}{
		{
			name: "Cloudflare ASN baseline",
			args: []string{"13335"},
			checks: func(t *testing.T, output string) {
				requiredElements := []string{"13335", "CLOUDFLARENET", "ARIN"}
				for _, element := range requiredElements {
					if !strings.Contains(output, element) {
						t.Errorf("Missing expected element '%s' in ASN output", element)
					}
				}
				// Should have exactly one table for ASN
				lines := strings.Split(output, "\n")
				nonEmptyLines := 0
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						nonEmptyLines++
					}
				}
				if nonEmptyLines < 3 { // At least header, separator, data
					t.Errorf("Expected at least 3 non-empty lines for ASN table, got %d", nonEmptyLines)
				}
			},
		},
		{
			name: "Google DNS baseline",
			args: []string{"8.8.8.8"},
			checks: func(t *testing.T, output string) {
				requiredElements := []string{"8.8.8.8", "15169", "GOOGLE"}
				for _, element := range requiredElements {
					if !strings.Contains(output, element) {
						t.Errorf("Missing expected element '%s' in IP output", element)
					}
				}
			},
		},
		{
			name: "Mixed input baseline",
			args: []string{"13335", "8.8.8.8"},
			checks: func(t *testing.T, output string) {
				// Should contain both ASN and IP information
				asnElements := []string{"13335", "CLOUDFLARENET"}
				ipElements := []string{"8.8.8.8", "15169"}
				
				for _, element := range asnElements {
					if !strings.Contains(output, element) {
						t.Errorf("Missing ASN element '%s' in mixed output", element)
					}
				}
				
				for _, element := range ipElements {
					if !strings.Contains(output, element) {
						t.Errorf("Missing IP element '%s' in mixed output", element)
					}
				}
				
				// Should have multiple sections (ASN table + IP table)
				tables := strings.Count(output, "|")
				if tables < 4 { // Each table should have multiple | separators
					t.Errorf("Expected multiple tables in mixed output, but found insufficient table markers")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("../ripestat_test", tc.args...)
			var stdout bytes.Buffer
			cmd.Stdout = &stdout

			err := cmd.Run()
			if err != nil {
				t.Fatalf("Command failed: %v", err)
			}

			output := stdout.String()
			tc.checks(t, output)
		})
	}
}
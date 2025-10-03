package utils

import (
	"bufio"
	"fmt"
	"net/netip"
	"os"
	"strconv"
	"strings"

	"github.com/cmingou/ripestat-cli/internal/ripestat"
)

func CheckArgsNonExist(args []string) bool {
	return len(args) < 1
}

func CnovertStringToAsn(str string) (int, error) {
	asn, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("Invalid ASN: %v", str)
	}

	if asn < 1 {
		return 0, fmt.Errorf("Invalid ASN: %v", asn)
	}

	return asn, nil
}

func CnovertStringToIp(str string) (netip.Addr, error) {
	ipAddress, err := netip.ParseAddr(str)
	if err != nil {
		return netip.Addr{}, fmt.Errorf("Invalid IP address: %v", str)
	}

	return ipAddress, nil
}

func PrintInvalidArgs(invalidSlice []string) {
	fmt.Printf("Invalid argument:\n")
	for _, invalid := range invalidSlice {
		fmt.Printf("%v\n", invalid)
	}
}

// ReadInputFile reads a file and returns a slice of strings
// Each line in the file should contain one IP address or ASN
// Empty lines and lines starting with # are ignored
func ReadInputFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var inputs []string
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Remove inline comments (everything after #)
		if idx := strings.Index(line, "#"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
			// If nothing left after removing comment, skip this line
			if line == "" {
				continue
			}
		}

		inputs = append(inputs, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file at line %d: %w", lineNum, err)
	}

	return inputs, nil
}

// GetPrefixLength extracts the prefix length from a CIDR notation string
// e.g., "8.8.8.0/24" returns 24, "2001:db8::/32" returns 32
func GetPrefixLength(prefix string) int {
	parts := strings.Split(prefix, "/")
	if len(parts) != 2 {
		return -1
	}

	length, err := strconv.Atoi(parts[1])
	if err != nil {
		return -1
	}

	return length
}

// FindLongestPrefixIndex returns the index of the route with the longest prefix
// from PrefixRoutingConsistency routes. Returns -1 if routes is empty.
func FindLongestPrefixIndex(rsp *ripestat.PrefixRoutingConsistency) int {
	if rsp == nil || len(rsp.Data.Routes) == 0 {
		return -1
	}

	longestIdx := 0
	longestLength := GetPrefixLength(rsp.Data.Routes[0].Prefix)

	for i := 1; i < len(rsp.Data.Routes); i++ {
		currentLength := GetPrefixLength(rsp.Data.Routes[i].Prefix)
		if currentLength > longestLength {
			longestLength = currentLength
			longestIdx = i
		}
	}

	return longestIdx
}

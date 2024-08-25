package utils

import (
	"fmt"
	"net/netip"
	"strconv"
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

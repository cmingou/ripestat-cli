/*
Copyright © 2024 Jimmy Ou

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"net/netip"
	"os"
	"strconv"

	"github.com/cmingou/ripestat-cli/internal/utils"
	"github.com/spf13/cobra"
)

var (
	asnSlice     []int
	ipv4Slice    []netip.Addr
	ipv6Slice    []netip.Addr
	invalidSlice []string
)

var rootCmd = &cobra.Command{
	Use:  "ripestat",
	Long: `This command will help to check the information about ASN, IPv4 and IPv6 from RIPEstat.`,
	Run: func(cmd *cobra.Command, args []string) {
		if utils.CheckArgsNonExist(args) {
			fmt.Printf("Please check parameter\n")
			os.Exit(1)
		}

		for _, arg := range args {
			if isASN(arg) {
				asn, _ := strconv.Atoi(arg)
				asnSlice = append(asnSlice, asn)
			} else if isIPv4(arg) {
				ipv4, _ := netip.ParseAddr(arg)
				ipv4Slice = append(ipv4Slice, ipv4)
			} else if isIPv6(arg) {
				ipv6, _ := netip.ParseAddr(arg)
				ipv6Slice = append(ipv6Slice, ipv6)
			} else {
				invalidSlice = append(invalidSlice, arg)
			}
		}

		if len(asnSlice) > 0 {
			fmt.Printf("## ASN\n")
			utils.SearchAsnInfo(asnSlice)
			fmt.Printf("\n")
		}

		if len(ipv4Slice) > 0 {
			fmt.Printf("## IPv4\n")
			utils.SearchIpv4Info(ipv4Slice)
			fmt.Printf("\n")
		}

		if len(ipv6Slice) > 0 {
			fmt.Printf("## IPv6\n")
			utils.SearchIpv6Info(ipv6Slice)
			fmt.Printf("\n")
		}

		if len(invalidSlice) > 0 {
			utils.PrintInvalidArgs(invalidSlice)
		}
	},
}

// Check if the input is an ASN
func isASN(arg string) bool {
	// Try to parse the string as a number
	asn, err := strconv.Atoi(arg)
	if err != nil {
		return false
	}
	// Check if it falls within the valid ASN range (0 to 4294967295)
	return asn >= 0 && asn <= 4294967295
}

// Check if the input is an IPv4 address
func isIPv4(arg string) bool {
	ip, err := netip.ParseAddr(arg)
	if err != nil {
		return false
	}

	return ip.Is4()
}

// Check if the input is an IPv6 address
func isIPv6(arg string) bool {
	ip, err := netip.ParseAddr(arg)
	if err != nil {
		return false
	}

	return ip.Is6()
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

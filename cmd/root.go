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
	maxConcurrency int
	inputFile      string
	longestPrefix  bool
	groupByPrefix  bool
	groupByAsn     bool
)

var rootCmd = &cobra.Command{
	Use:  "ripestat",
	Long: `This command will help to check the information about ASN, IPv4 and IPv6 from RIPEstat.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			asnSlice     []int
			ipv4Slice    []netip.Addr
			ipv6Slice    []netip.Addr
			invalidSlice []string
			allArgs      []string
		)

		utils.SetMaxConcurrentRequests(maxConcurrency)

		// Read from file if --file flag is provided
		if inputFile != "" {
			fileArgs, err := utils.ReadInputFile(inputFile)
			if err != nil {
				fmt.Printf("Error reading file: %v\n", err)
				os.Exit(1)
			}
			allArgs = append(allArgs, fileArgs...)
		}

		// Append command line arguments
		allArgs = append(allArgs, args...)

		// Check if we have any input
		if len(allArgs) == 0 {
			fmt.Printf("Please provide input via command line arguments or --file flag\n")
			os.Exit(1)
		}

		// Remove duplicates while preserving order
		seen := make(map[string]bool)
		var uniqueArgs []string
		for _, arg := range allArgs {
			if !seen[arg] {
				seen[arg] = true
				uniqueArgs = append(uniqueArgs, arg)
			}
		}

		for _, arg := range uniqueArgs {
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
			utils.SearchAsnInfo(asnSlice)
			fmt.Printf("\n")
		}

		if len(ipv4Slice) > 0 {
			if groupByPrefix {
				utils.SearchIpv4InfoGroupByPrefix(ipv4Slice)
			} else if groupByAsn {
				utils.SearchIpv4InfoGroupByAsn(ipv4Slice)
			} else {
				utils.SearchIpv4Info(ipv4Slice, longestPrefix)
			}
			fmt.Printf("\n")
		}

		if len(ipv6Slice) > 0 {
			if groupByPrefix {
				utils.SearchIpv6InfoGroupByPrefix(ipv6Slice)
			} else if groupByAsn {
				utils.SearchIpv6InfoGroupByAsn(ipv6Slice)
			} else {
				utils.SearchIpv6Info(ipv6Slice, longestPrefix)
			}
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
	maxConcurrency = utils.GetMaxConcurrentRequests()
	rootCmd.PersistentFlags().IntVar(&maxConcurrency, "max-concurrency", maxConcurrency, "Maximum concurrent RIPEstat requests (1-8)")
	rootCmd.PersistentFlags().StringVarP(&inputFile, "file", "f", "", "Read input from file (one IP/ASN per line)")
	rootCmd.PersistentFlags().BoolVarP(&longestPrefix, "longest-prefix", "l", false, "Show only the longest prefix route")
	rootCmd.PersistentFlags().BoolVar(&groupByPrefix, "group-by-prefix", false, "Group results by prefix (show each unique prefix once)")
	rootCmd.PersistentFlags().BoolVar(&groupByAsn, "group-by-asn", false, "Group results by ASN (show each unique ASN once)")

	// Make these flags mutually exclusive
	rootCmd.MarkFlagsMutuallyExclusive("longest-prefix", "group-by-prefix")
	rootCmd.MarkFlagsMutuallyExclusive("longest-prefix", "group-by-asn")
	rootCmd.MarkFlagsMutuallyExclusive("group-by-prefix", "group-by-asn")
}

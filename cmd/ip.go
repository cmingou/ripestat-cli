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
	"github.com/cmingou/ripestat-cli/pkg/ripestat"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	ipAddresses []netip.Addr
)

// ipCmd represents the ip command
var ipCmd = &cobra.Command{
	Use:   "ip [IP]",
	Short: "Search for IP address",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if utils.CheckArgsNonExist(args) {
			fmt.Printf("Please check parameter\n")
			os.Exit(1)
		}

		for _, ip := range args {
			ipAddress, err := utils.CnovertStringToIp(ip)
			if err != nil {
				fmt.Printf("Please check parameter, err: %v\n", err)
				os.Exit(1)
			}
			ipAddresses = append(ipAddresses, ipAddress)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"IP", "Prefix", "AS Number", "AS Name"})
		table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
		table.SetCenterSeparator("|")
		table.SetAutoWrapText(false)

		for _, ip := range ipAddresses {
			rsp, err := ripestat.GetPrefixRoutingConsistency(ip.String())
			if err != nil {
				fmt.Printf("Failed to get Prefix Routing Consistency: %v\n", err)
				os.Exit(1)

			}

			for idx, route := range rsp.Data.Routes {
				if route.InBgp {
					if idx == 0 {
						table.Append([]string{rsp.Data.Resource, route.Prefix, strconv.Itoa(route.Origin), route.AsnName})
					} else {
						table.Append([]string{"", route.Prefix, strconv.Itoa(route.Origin), route.AsnName})
					}
				}
			}
		}
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(ipCmd)
}

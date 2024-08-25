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
	"os"
	"strconv"

	"github.com/cmingou/ripestat-cli/internal/utils"
	"github.com/cmingou/ripestat-cli/pkg/ripestat"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	asNumbers []int
)

var asnCmd = &cobra.Command{
	Use:   "asn [AS number]",
	Short: "Search for AS number",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if utils.CheckArgsNonExist(args) {
			fmt.Printf("Please check parameter\n")
			os.Exit(1)
		}

		for _, as := range args {
			asNumber, err := utils.CnovertStringToAsn(as)
			if err != nil {
				fmt.Printf("Please check parameter, err: %v\n", err)
				os.Exit(1)
			}
			asNumbers = append(asNumbers, asNumber)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"AS", "Country", "RIR", "AS Name"})
		table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
		table.SetCenterSeparator("|")
		table.SetAutoWrapText(false)

		for _, as := range asNumbers {
			asOverview, err := ripestat.GetAsOverview(as)
			if err != nil {
				fmt.Printf("Failed to get AS overview: %v\n", err)
				os.Exit(1)
			}

			rir, err := ripestat.GetRIR(strconv.Itoa(as))
			if err != nil {
				fmt.Printf("Failed to get RIR: %v\n", err)
				os.Exit(1)
			}

			if len(rir.Data.Rirs) == 0 {
				table.Append([]string{asOverview.Data.Resource, "", "", asOverview.Data.Holder})
			} else {
				table.Append([]string{asOverview.Data.Resource, rir.Data.Rirs[0].Country, rir.Data.Rirs[0].Rir, asOverview.Data.Holder})
			}
		}

		table.Render()

	},
}

func init() {
	rootCmd.AddCommand(asnCmd)
}

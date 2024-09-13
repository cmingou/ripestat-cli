package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cmingou/ripestat-cli/pkg/ripestat"
	"github.com/olekukonko/tablewriter"
)

func SearchAsnInfo(asns []int) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"AS", "Country", "RIR", "AS Name"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)

	for _, as := range asns {
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
	fmt.Printf("## ASN\n")
	table.Render()
}

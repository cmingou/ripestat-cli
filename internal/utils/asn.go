package utils

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/cmingou/ripestat-cli/internal/ripestat"
	"github.com/olekukonko/tablewriter"
)

func SearchAsnInfo(asns []int) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"AS", "Country", "RIR", "AS Name"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)

	results, err := runWithConcurrency(context.Background(), asns, func(ctx context.Context, idx int, asn int) ([][]string, error) {
		asOverview, err := ripestat.GetAsOverview(asn)
		if err != nil {
			return nil, fmt.Errorf("AS%d: failed to get overview: %w", asn, err)
		}

		rir, err := ripestat.GetRIR(strconv.Itoa(asn))
		if err != nil {
			return nil, fmt.Errorf("AS%d: failed to get RIR data: %w", asn, err)
		}

		if len(rir.Data.Rirs) == 0 {
			return [][]string{{asOverview.Data.Resource, "", "", asOverview.Data.Holder}}, nil
		}

		entry := rir.Data.Rirs[0]
		return [][]string{{asOverview.Data.Resource, entry.Country, entry.Rir, asOverview.Data.Holder}}, nil
	})

	if err != nil {
		fmt.Printf("Failed to get ASN info: %v\n", err)
		os.Exit(1)
	}

	for _, rows := range results {
		for _, row := range rows {
			table.Append(row)
		}
	}
	fmt.Printf("## ASN\n")
	table.Render()
}

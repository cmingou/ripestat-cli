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
		var (
			overview    *ripestat.AsOverview
			overviewErr error
			rir         *ripestat.RIR
			rirErr      error
		)

		done := make(chan struct{}, 2)
		go func() {
			overview, overviewErr = ripestat.GetAsOverview(asn)
			done <- struct{}{}
		}()
		go func() {
			rir, rirErr = ripestat.GetRIR(strconv.Itoa(asn))
			done <- struct{}{}
		}()

		for i := 0; i < 2; i++ {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-done:
			}
		}

		if overviewErr != nil {
			return nil, fmt.Errorf("AS%d: failed to get overview: %w", asn, overviewErr)
		}
		if rirErr != nil {
			return nil, fmt.Errorf("AS%d: failed to get RIR data: %w", asn, rirErr)
		}

		if len(rir.Data.Rirs) == 0 {
			return [][]string{{overview.Data.Resource, "", "", overview.Data.Holder}}, nil
		}

		entry := rir.Data.Rirs[0]
		return [][]string{{overview.Data.Resource, entry.Country, entry.Rir, overview.Data.Holder}}, nil
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

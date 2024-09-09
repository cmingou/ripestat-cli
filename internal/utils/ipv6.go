package utils

import (
	"fmt"
	"net/netip"
	"os"
	"strconv"

	"github.com/cmingou/ripestat-cli/pkg/ripestat"
	"github.com/olekukonko/tablewriter"
)

func SearchIpv6Info(ipv6s []netip.Addr) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"IP", "Prefix", "AS Number", "AS Name"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)

	for _, ip := range ipv6s {
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
}

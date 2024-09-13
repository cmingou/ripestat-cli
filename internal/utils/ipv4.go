package utils

import (
	"fmt"
	"net/netip"
	"os"
	"strconv"

	"github.com/cmingou/ripestat-cli/pkg/ripestat"
	"github.com/olekukonko/tablewriter"
)

func SearchIpv4Info(ipv4s []netip.Addr) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"IP", "Location", "Prefix", "In BGP", "AS Number", "AS Name"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)

	for _, ip := range ipv4s {
		ipLocation, err := ripestat.GetIpGeoLocation(ip.String())
		if err != nil {
			fmt.Printf("Failed to get IP Geo Location: %v\n", err)
			os.Exit(1)
		}

		rsp, err := ripestat.GetPrefixRoutingConsistency(ip.String())
		if err != nil {
			fmt.Printf("Failed to get Prefix Routing Consistency: %v\n", err)
			os.Exit(1)
		}

		for idx, route := range rsp.Data.Routes {
			if idx == 0 {
				table.Append([]string{rsp.Data.Resource, ipLocation, route.Prefix, strconv.FormatBool(route.InBgp), strconv.Itoa(route.Origin), route.AsnName})
			} else {
				table.Append([]string{"", "", route.Prefix, strconv.FormatBool(route.InBgp), strconv.Itoa(route.Origin), route.AsnName})
			}

		}
	}
	fmt.Printf("## IPv4\n")
	table.Render()
}

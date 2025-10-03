package utils

import (
	"context"
	"fmt"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/cmingou/ripestat-cli/internal/ripestat"
	"github.com/olekukonko/tablewriter"
)

type PrefixInfo struct {
	Prefix string
	InBGP  bool
	ASN    int
	ASName string
}

type ASNInfo struct {
	ASN    int
	ASName string
}

func SearchIpv4Info(ipv4s []netip.Addr, longestPrefix bool) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"IP", "Location", "Prefix", "In BGP", "AS Number", "AS Name"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)

	var asnNameCache sync.Map

	results, err := runWithConcurrency(context.Background(), ipv4s, func(ctx context.Context, idx int, ip netip.Addr) ([][]string, error) {
		resource := ip.String()

		var (
			ipLocation         string
			ipLocationErr      error
			routingConsistency *ripestat.PrefixRoutingConsistency
			routingErr         error
		)

		done := make(chan struct{}, 2)
		go func() {
			ipLocation, ipLocationErr = ripestat.GetIpGeoLocation(resource)
			done <- struct{}{}
		}()
		go func() {
			routingConsistency, routingErr = ripestat.GetPrefixRoutingConsistency(resource)
			done <- struct{}{}
		}()

		for i := 0; i < 2; i++ {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-done:
			}
		}

		if ipLocationErr != nil {
			return nil, fmt.Errorf("%s: failed to get Geo location: %w", resource, ipLocationErr)
		}
		if routingErr != nil {
			return nil, fmt.Errorf("%s: failed to get routing consistency: %w", resource, routingErr)
		}

		rsp := routingConsistency

		// Filter to longest prefix if flag is set
		var routesToProcess []int
		if longestPrefix {
			longestIdx := FindLongestPrefixIndex(rsp)
			if longestIdx != -1 {
				routesToProcess = []int{longestIdx}
			}
		} else {
			routesToProcess = make([]int, len(rsp.Data.Routes))
			for i := range rsp.Data.Routes {
				routesToProcess[i] = i
			}
		}

		rows := make([][]string, 0, len(routesToProcess))
		for i, routeIdx := range routesToProcess {
			route := rsp.Data.Routes[routeIdx]
			asnName := strings.TrimSpace(route.AsnName)
			if asnName == "" || asnName == "-" {
				if cached, ok := asnNameCache.Load(route.Origin); ok {
					asnName = cached.(string)
				} else {
					overview, err := ripestat.GetAsOverview(route.Origin)
					if err != nil {
						return nil, fmt.Errorf("%s: failed to resolve AS name for %d: %w", resource, route.Origin, err)
					}
					asnName = overview.Data.Holder
					asnNameCache.Store(route.Origin, asnName)
				}
			}

			row := []string{"", "", route.Prefix, strconv.FormatBool(route.InBgp), strconv.Itoa(route.Origin), asnName}
			if i == 0 {
				row[0] = rsp.Data.Resource
				row[1] = ipLocation
			}
			rows = append(rows, row)
		}

		return rows, nil
	})

	if err != nil {
		fmt.Printf("Failed to get IPv4 info: %v\n", err)
		os.Exit(1)
	}

	for _, rows := range results {
		for _, row := range rows {
			table.Append(row)
		}
	}
	fmt.Printf("## IPv4\n")
	table.Render()
}

func SearchIpv4InfoGroupByPrefix(ipv4s []netip.Addr) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Prefix", "In BGP", "AS Number", "AS Name"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)

	var (
		prefixMap    sync.Map
		asnNameCache sync.Map
	)

	results, err := runWithConcurrency(context.Background(), ipv4s, func(ctx context.Context, idx int, ip netip.Addr) (*PrefixInfo, error) {
		resource := ip.String()

		rsp, err := ripestat.GetPrefixRoutingConsistency(resource)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get routing consistency: %w", resource, err)
		}

		// Always use the longest prefix
		longestIdx := FindLongestPrefixIndex(rsp)
		if longestIdx == -1 {
			return nil, nil
		}

		route := rsp.Data.Routes[longestIdx]

		// Check if this prefix already exists
		if _, exists := prefixMap.LoadOrStore(route.Prefix, true); exists {
			return nil, nil
		}

		// Resolve ASN name if needed
		asnName := strings.TrimSpace(route.AsnName)
		if asnName == "" || asnName == "-" {
			if cached, ok := asnNameCache.Load(route.Origin); ok {
				asnName = cached.(string)
			} else {
				overview, err := ripestat.GetAsOverview(route.Origin)
				if err != nil {
					return nil, fmt.Errorf("%s: failed to resolve AS name for %d: %w", resource, route.Origin, err)
				}
				asnName = overview.Data.Holder
				asnNameCache.Store(route.Origin, asnName)
			}
		}

		return &PrefixInfo{
			Prefix: route.Prefix,
			InBGP:  route.InBgp,
			ASN:    route.Origin,
			ASName: asnName,
		}, nil
	})

	if err != nil {
		fmt.Printf("Failed to get IPv4 info: %v\n", err)
		os.Exit(1)
	}

	for _, result := range results {
		if result != nil {
			table.Append([]string{
				result.Prefix,
				strconv.FormatBool(result.InBGP),
				strconv.Itoa(result.ASN),
				result.ASName,
			})
		}
	}
	fmt.Printf("## IPv4\n")
	table.Render()
}

func SearchIpv4InfoGroupByAsn(ipv4s []netip.Addr) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"AS Number", "AS Name"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)

	var (
		asnMap       sync.Map
		asnNameCache sync.Map
	)

	results, err := runWithConcurrency(context.Background(), ipv4s, func(ctx context.Context, idx int, ip netip.Addr) (*ASNInfo, error) {
		resource := ip.String()

		rsp, err := ripestat.GetPrefixRoutingConsistency(resource)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get routing consistency: %w", resource, err)
		}

		// Always use the longest prefix
		longestIdx := FindLongestPrefixIndex(rsp)
		if longestIdx == -1 {
			return nil, nil
		}

		route := rsp.Data.Routes[longestIdx]

		// Check if this ASN already exists
		if _, exists := asnMap.LoadOrStore(route.Origin, true); exists {
			return nil, nil
		}

		// Resolve ASN name if needed
		asnName := strings.TrimSpace(route.AsnName)
		if asnName == "" || asnName == "-" {
			if cached, ok := asnNameCache.Load(route.Origin); ok {
				asnName = cached.(string)
			} else {
				overview, err := ripestat.GetAsOverview(route.Origin)
				if err != nil {
					return nil, fmt.Errorf("%s: failed to resolve AS name for %d: %w", resource, route.Origin, err)
				}
				asnName = overview.Data.Holder
				asnNameCache.Store(route.Origin, asnName)
			}
		}

		return &ASNInfo{
			ASN:    route.Origin,
			ASName: asnName,
		}, nil
	})

	if err != nil {
		fmt.Printf("Failed to get IPv4 info: %v\n", err)
		os.Exit(1)
	}

	for _, result := range results {
		if result != nil {
			table.Append([]string{
				strconv.Itoa(result.ASN),
				result.ASName,
			})
		}
	}
	fmt.Printf("## IPv4\n")
	table.Render()
}

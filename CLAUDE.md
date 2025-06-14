# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ripestat-cli is a Go-based command-line tool that provides a wrapper for the RIPEstat API. It allows users to query information about ASNs (Autonomous System Numbers), IPv4, and IPv6 addresses directly from the command line.

## Architecture

- **main.go**: Entry point that calls cmd.Execute()
- **cmd/root.go**: Main command logic using Cobra CLI framework
  - Parses input arguments to identify ASNs, IPv4, and IPv6 addresses
  - Routes queries to appropriate utility functions
- **pkg/ripestat/**: API client package
  - `api.go`: HTTP client functions for RIPEstat API endpoints
  - `struts.go`: Go structs for JSON response unmarshaling
- **internal/utils/**: Utility functions for different resource types
  - `asn.go`: ASN information queries
  - `ipv4.go`: IPv4 address queries  
  - `ipv6.go`: IPv6 address queries
  - `util.go`: Common utility functions

## Common Commands

### Development
```bash
# Build for current platform
go build -o ripestat main.go

# Install locally
make install

# Run tests
make test

# Lint code
make lint
```

### Building
```bash
# Clean previous builds
make clean

# Build for all platforms (Mac and Linux)
make all

# Build specific platforms
make darwin  # Mac
make linux   # Linux
```

### Usage
The tool accepts ASNs, IPv4, and IPv6 addresses as arguments and automatically categorizes them:
```bash
./ripestat 8.8.8.8 13335 2001:db8::1
```

## Key Dependencies

- **github.com/spf13/cobra**: CLI framework
- **github.com/olekukonko/tablewriter**: Table output formatting
- **net/netip**: Standard library for IP address parsing

## API Integration

The tool integrates with several RIPEstat API endpoints:
- AS Overview: Basic ASN information
- RIR: Regional Internet Registry data
- Prefix Routing Consistency: BGP routing information
- MaxMind GeoLite: Geographic location data
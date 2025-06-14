# ripestat-cli

A Go-based command-line tool that provides a simple wrapper for the RIPEstat API. Query information about ASNs (Autonomous System Numbers), IPv4, and IPv6 addresses directly from your terminal.

[![Go](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## Features

- **Multi-input support**: Query ASNs, IPv4, and IPv6 addresses in a single command
- **Automatic detection**: Automatically categorizes input as ASN, IPv4, or IPv6
- **Rich information**: Provides comprehensive data including:
  - ASN overview and details
  - Regional Internet Registry (RIR) information
  - BGP routing consistency data
  - Geographic location information
- **Table output**: Clean, formatted table display for easy reading

## Installation

### From Source

```bash
git clone https://github.com/yourusername/ripestat-cli.git
cd ripestat-cli
make install
```

### Build from Source

```bash
# Build for current platform
go build -o ripestat main.go

# Build for all platforms
make all

# Build for specific platforms
make darwin  # macOS
make linux   # Linux
```

## Usage

The tool accepts multiple arguments and automatically detects their types:

```bash
# Query multiple resources at once
./ripestat 8.8.8.8 13335 2001:db8::1

# Query individual resources
./ripestat 8.8.8.8                # IPv4 address
./ripestat 13335                   # ASN
./ripestat 2001:db8::1            # IPv6 address
```

### Example Output

```
╭─────────────────────────────────────────────────────────────╮
│                        AS Information                        │
├─────────────────────────────────────────────────────────────┤
│ ASN: 13335                                                  │
│ Name: CLOUDFLARENET                                         │
│ Country: US                                                 │
│ Registry: ARIN                                              │
╰─────────────────────────────────────────────────────────────╯
```

## API Integration

This tool integrates with several RIPEstat API endpoints:

- **AS Overview**: Basic ASN information and metadata
- **RIR**: Regional Internet Registry data
- **Prefix Routing Consistency**: BGP routing information
- **MaxMind GeoLite**: Geographic location data for IP addresses

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional, for using Makefile commands)

### Building

```bash
# Clean previous builds
make clean

# Run tests
make test

# Lint code
make lint

# Build for development
go build -o ripestat main.go
```

### Project Structure

```
├── main.go              # Entry point
├── cmd/
│   ├── root.go         # Main command logic (Cobra CLI)
│   └── root_test.go    # Tests for input validation functions
├── pkg/
│   └── ripestat/       # API client package
│       ├── api.go      # HTTP client functions
│       └── struts.go   # Response structs
└── internal/
    └── utils/          # Utility functions
        ├── asn.go      # ASN queries
        ├── ipv4.go     # IPv4 queries
        ├── ipv6.go     # IPv6 queries
        ├── util.go     # Common utilities
        └── util_test.go # Tests for utility functions
```

## Dependencies

- [github.com/spf13/cobra](https://github.com/spf13/cobra) - CLI framework
- [github.com/olekukonko/tablewriter](https://github.com/olekukonko/tablewriter) - Table formatting
- Standard Go libraries for networking and HTTP

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [RIPE NCC](https://www.ripe.net/) for providing the RIPEstat API
- The Go community for excellent libraries and tools

---

For the Chinese version of this README, see [README.zh_tw.md](README.zh_tw.md).
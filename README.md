# ripestat-cli

A comprehensive Go-based command-line tool that provides a powerful wrapper for the RIPEstat API. Query detailed information about ASNs (Autonomous System Numbers), IPv4, and IPv6 addresses directly from your terminal with beautifully formatted output.

[![Go](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/Tests-Passing-brightgreen.svg)](#testing)

## Features

- **🔍 Multi-input support**: Query multiple ASNs, IPv4, and IPv6 addresses in a single command
- **🤖 Automatic detection**: Intelligent categorization of input as ASN, IPv4, or IPv6
- **📊 Rich information**: Comprehensive data including:
  - ASN overview, holder information, and registry details
  - Regional Internet Registry (RIR) allocation data
  - BGP routing consistency and prefix information  
  - Geographic location data with city/country details
- **📋 Clean table output**: Professional formatted tables for easy reading
- **⚡ Fast performance**: Concurrent API calls for optimal speed
- **🧪 Comprehensive testing**: Full test suite ensuring reliability

## Quick Start

```bash
# Clone and build
git clone https://github.com/cmingou/ripestat-cli.git
cd ripestat-cli
go build -o ripestat main.go

# Query multiple resources at once
./ripestat 8.8.8.8 13335 2001:4860:4860::8888
```

## Installation

### From GitHub Releases (Recommended)

Download the latest pre-built binary for your platform from the [Releases page](https://github.com/cmingou/ripestat-cli/releases):

**macOS:**
```bash
# For Intel Macs (x86_64)
curl -L https://github.com/cmingou/ripestat-cli/releases/latest/download/ripestat_Darwin_x86_64.tar.gz | tar xz
sudo mv ripestat /usr/local/bin/

# For Apple Silicon (M1/M2/M3 - arm64)
curl -L https://github.com/cmingou/ripestat-cli/releases/latest/download/ripestat_Darwin_arm64.tar.gz | tar xz
sudo mv ripestat /usr/local/bin/
```

**Linux:**
```bash
# For x86_64 (Debian, Ubuntu, CentOS, etc.)
curl -L https://github.com/cmingou/ripestat-cli/releases/latest/download/ripestat_Linux_x86_64.tar.gz | tar xz
sudo mv ripestat /usr/local/bin/

# For ARM64 (Raspberry Pi, ARM servers, etc.)
curl -L https://github.com/cmingou/ripestat-cli/releases/latest/download/ripestat_Linux_arm64.tar.gz | tar xz
sudo mv ripestat /usr/local/bin/
```

### From Source

```bash
git clone https://github.com/cmingou/ripestat-cli.git
cd ripestat-cli
make install
```

### Development Build

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
# Query multiple resources simultaneously
./ripestat 8.8.8.8 13335 2001:db8::1

# Query individual resources
./ripestat 8.8.8.8                    # IPv4 address (Google DNS)
./ripestat 13335                       # ASN (Cloudflare)
./ripestat 2001:4860:4860::8888       # IPv6 address (Google DNS)
./ripestat 1.1.1.1 15169 8.8.8.8      # Mixed input types
```

### Options

| Flag | Short | Description |
|------|-------|-------------|
| `--all-prefixes` | `-a` | Show all prefix routes (default shows only the longest BGP prefix) |
| `--group-by-prefix` | | Group results by prefix (show each unique prefix once) |
| `--group-by-asn` | | Group results by ASN (show each unique ASN once) |
| `--file` | `-f` | Read input from file (one IP/ASN per line) |
| `--max-concurrency` | | Maximum concurrent RIPEstat requests (1-8) |

```bash
# Default: show only the longest BGP prefix per IP
./ripestat 8.8.8.8

# Show all prefix routes
./ripestat -a 8.8.8.8
./ripestat --all-prefixes 140.113.1.1

# Group by prefix or ASN
./ripestat --group-by-prefix 8.8.8.8 1.1.1.1
./ripestat --group-by-asn 8.8.8.8 1.1.1.1

# Read from file
./ripestat -f ips.txt
```

### Example Output

#### ASN Information
```
## ASN
   AS   | COUNTRY | RIR  |    AS NAME     
--------|---------|------|----------------
  13335 | US      | ARIN | CLOUDFLARENET  
```

#### IPv4 Information
```
## IPv4
    IP    | LOCATION |   PREFIX   | IN BGP | AS NUMBER |           AS NAME             
----------|----------|------------|--------|-----------|-------------------------------
  8.8.8.8 | US       | 8.8.8.0/24 | true   |     15169 | GOOGLE - Google LLC           
```

#### IPv6 Information
```
## IPv6
                   IP          | LOCATION |     PREFIX     | IN BGP | AS NUMBER |       AS NAME        
-----------------------|----------|----------------|--------|-----------|----------------------
  2001:4860:4860::8888 | US       | 2001:4860::/32 | true   |     15169 | GOOGLE - Google LLC   
```

## Architecture

<div align="center">
  <img src="docs/images/architecture.svg" alt="ripestat-cli Architecture Diagram" width="800">
</div>

The tool follows a clean, modular architecture:

1. **Input Processing**: CLI arguments are parsed and categorized by type
2. **Type Detection**: Automatic identification of ASN, IPv4, IPv6, or invalid inputs
3. **API Integration**: Unified client interface to multiple RIPEstat endpoints
4. **Data Processing**: Structured handling of API responses
5. **Output Formatting**: Professional table formatting for console display

## API Integration

This tool integrates with multiple RIPEstat API endpoints:

| API Endpoint | Purpose | Data Provided |
|--------------|---------|---------------|
| **AS Overview** | Basic ASN information | ASN details, holder name, allocation status |
| **RIR** | Registry data | Regional registry allocation and country info |
| **Prefix Routing Consistency** | BGP information | Route announcements, prefixes, origin ASNs |
| **MaxMind GeoLite** | Geographic location | City, country, coordinates for IP addresses |

## Development

### Prerequisites

- **Go 1.21+** - Required for building
- **Make** - Optional, for using Makefile commands
- **Internet connection** - Required for API calls

### Building and Testing

```bash
# Development workflow
go build -o ripestat main.go           # Quick build

# Comprehensive testing
go test ./... -v                       # Run all tests
go test ./cmd/ -v                      # CLI integration tests
go test ./internal/utils/ -v           # Business logic tests
go test ./internal/ripestat/ -v        # API client tests

# Build management
make clean                             # Clean previous builds
make test                              # Run tests via Makefile
make lint                              # Code linting
```

### Concurrency controls

RIPEstat limits each IP address to eight concurrent requests. The CLI batches lookups with that cap by default. Adjust the concurrency ceiling with the `--max-concurrency` flag or the `RIPESTAT_MAX_CONCURRENCY` environment variable; both values are clamped between 1 and 8.

### Network optimizations

- **HTTP/2 by default** – the shared client negotiates HTTP/2 (via `golang.org/x/net/http2`) and keeps a single connection multiplexed up to the RIPEstat per-IP limit.
- **Fine-tuned pooling** – idle and active connection caps mirror the 8-stream concurrency guard so the transport does not overrun provider limits.
- **Debug visibility** – export `RIPESTAT_DEBUG_HTTP=1` to log the protocol negotiated per request (e.g. `proto=HTTP/2.0`) to stderr without changing CLI output.
- **Sanity check support** – verify the upstream endpoint with `curl --http2 -I https://stat.ripe.net` (network access required) before attributing slowdowns to the client.
### HTTP/2 diagnostics

```bash
# Enable protocol logging while exercising the CLI
RIPESTAT_DEBUG_HTTP=1 ./ripestat 8.8.8.8 13335

# Isolated unit probe (skips automatically if the sandbox blocks loopback listeners)
GOCACHE=$(pwd)/.gocache go test ./internal/ripestat -run TestGetHttpGetResponseUsesHTTP2 -v
```

If the test reports a skip, the local sandbox denied creating the loopback HTTP/2 server; re-run on a developer machine to observe the full negotiation log.
### Project Structure

```
├── main.go                    # Application entry point
├── cmd/
│   ├── root.go              # Main CLI logic (Cobra framework)
│   ├── root_test.go         # Input validation tests
│   └── integration_test.go  # End-to-end CLI tests
├── internal/
│   ├── ripestat/            # API client package
│   │   ├── api.go          # HTTP client and API functions
│   │   ├── api_test.go     # API integration tests
│   │   ├── http_client.go  # Shared HTTP client with HTTP/2 support
│   │   ├── http_client_test.go
│   │   ├── logging.go      # Debug logging utilities
│   │   └── structs.go      # JSON response structures
│   └── utils/               # Business logic utilities
│       ├── asn.go          # ASN information processing
│       ├── concurrency.go  # Concurrent API request management
│       ├── ipv4.go         # IPv4 address processing
│       ├── ipv6.go         # IPv6 address processing
│       ├── util.go         # Common utility functions
│       ├── util_test.go    # Unit tests for utilities
│       └── utils_integration_test.go # Integration tests
└── CLAUDE.md               # Development guidance for AI assistants
```

## Testing

The project includes a comprehensive test suite ensuring reliability:

### Test Categories

- **🔧 Unit Tests**: Core functionality and input validation
- **🔗 Integration Tests**: API client behavior and data accuracy
- **🎯 End-to-End Tests**: Complete CLI workflow testing
- **⚡ Performance Tests**: Response time and efficiency
- **🛡️ Error Handling Tests**: Invalid input and edge cases

### Running Tests

```bash
# All tests with verbose output
go test ./... -v

# Specific test suites
go test ./cmd/ -v                    # CLI functionality
go test ./internal/utils/ -v         # Business logic  
go test ./internal/ripestat/ -v      # API client

# Test with coverage
go test ./... -cover
```

### Test Coverage

The test suite covers:
- ✅ All API endpoints with real data
- ✅ Input validation and categorization
- ✅ Table formatting and output
- ✅ Error handling and edge cases
- ✅ Mixed input scenarios
- ✅ Performance benchmarks

## Dependencies

| Package | Purpose | License |
|---------|---------|---------|
| [github.com/spf13/cobra](https://github.com/spf13/cobra) | CLI framework and command parsing | Apache 2.0 |
| [github.com/olekukonko/tablewriter](https://github.com/olekukonko/tablewriter) | ASCII table formatting | MIT |
| [golang.org/x/net](https://golang.org/x/net) | HTTP/2 transport support and network utilities | BSD-3-Clause |
| **Standard Go Libraries** | HTTP client, JSON parsing, networking | BSD-3-Clause |

## Performance

- **Concurrent API calls** for multiple inputs
- **HTTP/2 connection multiplexing** for reduced latency
- **Fine-tuned connection pooling** respecting RIPEstat quotas
- **Response caching** where appropriate
- **Optimized table rendering** for large datasets
- **Minimal memory footprint** (~10MB runtime)

## Error Handling

The tool gracefully handles various error conditions:

- **Invalid input formats** → Clear error messages
- **Network timeouts** → Retry logic with backoff
- **API rate limits** → Automatic throttling
- **Missing data** → Partial results with warnings

## Contributing

We welcome contributions! Please follow these steps:

1. **Fork** the repository
2. **Create** your feature branch (`git checkout -b feature/amazing-feature`)
3. **Add tests** for new functionality
4. **Ensure** all tests pass (`go test ./... -v`)
5. **Commit** your changes (`git commit -m 'Add amazing feature'`)
6. **Push** to the branch (`git push origin feature/amazing-feature`)
7. **Open** a Pull Request

### Development Guidelines

- Follow Go best practices and conventions
- Add tests for all new functionality
- Update documentation for user-facing changes
- Ensure code passes `go vet` and `go fmt`

## License

This project is licensed under the **Apache License 2.0** - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- **[RIPE NCC](https://www.ripe.net/)** - For providing the comprehensive RIPEstat API
- **[Go Community](https://golang.org/community/)** - For excellent libraries and development tools
- **[MaxMind](https://www.maxmind.com/)** - For geographic location data

## Related Projects

- [RIPEstat API Documentation](https://stat.ripe.net/docs/data_api)
- [RIPE Database](https://www.ripe.net/manage-ips-and-asns/db/)
- [BGP Tools](https://bgp.tools/)

---

**Language Versions:** [English](README.md) | [繁體中文](README.zh_tw.md)

For questions, issues, or feature requests, please [open an issue](https://github.com/cmingou/ripestat-cli/issues).

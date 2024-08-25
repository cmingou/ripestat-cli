# Variables for the project
BINARY_NAME=ripestat
VERSION=1.0.0
BUILD_DIR=build

# Default task
all: clean darwin linux

# Clean up previous builds
clean:
	rm -rf $(BUILD_DIR)

# Build for Mac (darwin)
darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64-$(VERSION) main.go

# Build for Linux (amd64)
linux:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64-$(VERSION) main.go

# Install locally
install:
	go install

# Run tests
test:
	go test ./...

# Lint the code
lint:
	golangci-lint run

# Cross-compile for all targets
build-all: clean darwin linux

# Help command
help:
	@echo "Makefile commands:"
	@echo "  all        - Clean and build for all platforms (Mac, Linux)"
	@echo "  clean      - Clean up previous builds"
	@echo "  darwin     - Build for MacOS (darwin)"
	@echo "  linux      - Build for Linux (amd64)"
	@echo "  install    - Install the binary locally"
	@echo "  test       - Run unit tests"
	@echo "  lint       - Lint the code using golangci-lint"
	@echo "  build-all  - Clean and build for all platforms (Mac, Linux)"

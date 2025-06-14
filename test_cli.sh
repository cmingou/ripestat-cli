#!/bin/bash

# Test script for ripestat CLI functionality
# This script helps you verify that your CLI is working correctly after modifications

set -e

echo "🧪 Testing ripestat CLI functionality..."

# Build the CLI
echo "📦 Building CLI..."
go build -o ripestat_test main.go

echo "✅ CLI built successfully"

echo ""
echo "🔍 Running functionality tests..."

# Test ASN lookup
echo "1️⃣  Testing ASN lookup (Cloudflare 13335)..."
./ripestat_test 13335 | head -10
echo ""

# Test IPv4 lookup  
echo "2️⃣  Testing IPv4 lookup (Google DNS 8.8.8.8)..."
./ripestat_test 8.8.8.8 | head -10
echo ""

# Test IPv6 lookup
echo "3️⃣  Testing IPv6 lookup (Google IPv6 DNS)..."
./ripestat_test 2001:4860:4860::8888 | head -10
echo ""

# Test mixed input
echo "4️⃣  Testing mixed input (ASN + IPv4)..."
./ripestat_test 13335 8.8.8.8 | head -15
echo ""

# Test invalid input
echo "5️⃣  Testing invalid input handling..."
./ripestat_test invalid-input 2>/dev/null || echo "✅ Invalid input handled correctly"
echo ""

# Test no arguments
echo "6️⃣  Testing no arguments..."
./ripestat_test 2>/dev/null || echo "✅ No arguments handled correctly"
echo ""

echo "🧹 Cleaning up..."
rm -f ripestat_test

echo ""
echo "✅ All basic functionality tests completed!"
echo ""
echo "To run comprehensive tests, use:"
echo "  go test ./... -v"
echo ""
echo "To run specific test suites:"
echo "  go test ./cmd/ -v                    # CLI integration tests"
echo "  go test ./internal/utils/ -v         # Business logic tests"
echo "  go test ./internal/ripestat/ -v      # API client tests"
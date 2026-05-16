#!/usr/bin/env bash
set -e
echo "Fetching dependencies..."
go mod tidy
echo
echo "Building CLI tool..."
go build -o make_colors ./cmd/make_colors
echo
echo "Done! Run: ./make_colors -t"

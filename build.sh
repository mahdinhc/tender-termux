#!/bin/bash
set -e

mkdir -p build/android_arm64

export GOARCH=arm64
export GOOS=android

go build -o build/tender -ldflags "-s -w" cli/tender/main.go

echo "Build complete. Press enter to exit."
read
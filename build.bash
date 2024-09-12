#!/bin/bash

# Set the output binary name
OUTPUT_BINARY="whoknows"

# Detect the operating system
OS=$(uname -s)

# Check if the operating system is Windows
if [[ "$OS" == "MINGW"* || "$OS" == "CYGWIN"* || "$OS" == "MSYS"* || "$OS" == "Windows_NT" ]]; then
    OUTPUT_BINARY="${OUTPUT_BINARY}.exe"
fi

# Cd to backend dir
cd ./backend/

# Run the Go build command
echo "Building for OS: ${GOOS:-$(go env GOOS)}, Arch: ${GOARCH:-$(go env GOARCH)}"
go build -o "$OUTPUT_BINARY"

# Output the name of the generated binary
echo "Build complete: $OUTPUT_BINARY"

# Run binary
./"${OUTPUT_BINARY}"

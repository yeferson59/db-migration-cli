#!/bin/bash

# Limpia builds anteriores
rm -rf dist/

# Crea directorio
mkdir -p dist/

# Install dependencies
go mod download
go mod tidy

# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o dist/db-migration-cli_linux_amd64 main.go

# macOS AMD64
GOOS=darwin GOARCH=amd64 go build -o dist/db-migration-cli_darwin_amd64 main.go

# macOS ARM64 (para Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o dist/db-migration-cli_darwin_arm64 main.go

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o dist/db-migration-cli_windows_amd64.exe main.go

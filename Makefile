.PHONY: build build-all linux windows darwin compress clean install sizes help

# Build for current platform
build:
	go build -ldflags "-s -w" -o soundplay

# Build for all platforms
build-all: linux windows darwin

# Linux builds
linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/soundplay-linux-amd64
	GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o dist/soundplay-linux-arm64

# Windows build
windows:
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o dist/soundplay-windows-amd64.exe

# macOS builds
darwin:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o dist/soundplay-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o dist/soundplay-darwin-arm64

# Compress binaries with UPX (optional, requires upx installed)
compress:
	@command -v upx >/dev/null 2>&1 || { echo "UPX not found. Install it first."; exit 1; }
	upx --best --lzma dist/*

# Clean build artifacts
clean:
	rm -rf soundplay dist/

# Install to /usr/local/bin (Unix-like systems)
install: build
	install -m 755 soundplay /usr/local/bin/

# Show binary sizes
sizes:
	@ls -lh dist/* 2>/dev/null || echo "No binaries found in dist/. Run 'make build-all' first."

# Show help
help:
	@echo "Soundplay Build System"
	@echo ""
	@echo "Targets:"
	@echo "  build       - Build for current platform"
	@echo "  build-all   - Cross-compile for all platforms"
	@echo "  linux       - Build for Linux (amd64, arm64)"
	@echo "  windows     - Build for Windows (amd64)"
	@echo "  darwin      - Build for macOS (amd64, arm64)"
	@echo "  compress    - Compress binaries with UPX"
	@echo "  clean       - Remove build artifacts"
	@echo "  install     - Install to /usr/local/bin"
	@echo "  sizes       - Show binary sizes"
	@echo "  help        - Show this help message"

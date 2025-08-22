.PHONY: build clean install test

# Build the gosecret binary
build:
	go build -o gosecret

# Build with optimizations for release
build-release:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o gosecret

# Clean build artifacts
clean:
	rm -f gosecret

# Install to system PATH (requires sudo)
install: build-release
	install -m 755 gosecret /usr/local/bin/

# Remove from system PATH (requires sudo)
uninstall:
	rm -f /usr/local/bin/gosecret

# Run basic functionality test
test: build
	@echo "Testing gosecret functionality..."
	@echo "Testing set and get..."
	@echo "testpassword123" | ./gosecret set gosecret.test
	@./gosecret get gosecret.test
	@echo ""
	@echo "Testing list..."
	@./gosecret list gosecret
	@echo ""
	@echo "Cleaning up test secret..."
	@./gosecret delete gosecret.test
	@echo "Test completed successfully!"

# Download dependencies
deps:
	go mod tidy
	go mod download

# Format Go code
fmt:
	go fmt ./...

# Run Go vet
vet:
	go vet ./...

# Run all checks
check: fmt vet

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the gosecret binary"
	@echo "  build-release - Build optimized release binary"
	@echo "  clean         - Remove build artifacts"
	@echo "  install       - Install to /usr/local/bin (requires sudo)"
	@echo "  uninstall     - Remove from /usr/local/bin (requires sudo)"
	@echo "  test          - Run basic functionality test"
	@echo "  deps          - Download dependencies"
	@echo "  fmt           - Format Go code"
	@echo "  vet           - Run Go vet"
	@echo "  check         - Run fmt and vet"
	@echo "  help          - Show this help message"
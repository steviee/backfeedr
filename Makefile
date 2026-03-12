# Modern Go Makefile (2026 style)
# Uses go tool for everything possible

.PHONY: all build test lint fmt vet clean run dev docker

BINARY_NAME := backfeedr
CLIENT_NAME := backfeedr-client

# Default target
all: build

# Build using go build (modern idiom)
build:
	go build -o $(BINARY_NAME) ./cmd/backfeedr
	go build -o $(CLIENT_NAME) ./cmd/backfeedr-client

# Test with gotestsum or standard go test
test:
	go test -v -race ./...

# Run integration tests locally
test-integration: build
	./test/integration.sh

# Run all tests (unit + integration)
test-all: test test-integration

# Lint with golangci-lint (standard 2026)
lint:
	golangci-lint run ./...

# Format with gofmt
fmt:
	gofmt -s -w .

# Vet with go vet
vet:
	go vet ./...

# Tidy modules
tidy:
	go mod tidy

# Download dependencies
deps:
	go mod download

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME) $(CLIENT_NAME)
	go clean -cache

# Run server (dev mode with reload)
dev:
	go run ./cmd/backfeedr

# Run client (dev mode)
dev-client:
	go run ./cmd/backfeedr-client

# Docker build
docker:
	docker build -t $(BINARY_NAME):latest .

# Run the built binary
run: build
	./$(BINARY_NAME)

# Full CI pipeline
ci: fmt vet test build

.PHONY: build test lint coverage install uninstall clean run docs version

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -X main.version=$(VERSION)

# Build both binaries
build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o build/saiad ./cmd/saiad
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o build/saia ./cmd/saia

# Run all tests with race detector
test:
	go test ./... -race -count=1 -v

# Run tests with coverage report
test-coverage:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out | grep total
	go tool cover -html=coverage.out -o coverage.html

# Lint
lint:
	golangci-lint run ./...

# Format
fmt:
	go fmt ./...

# Tidy dependencies
tidy:
	go mod tidy

# Run saiad from source (dev mode)
run:
	go run ./cmd/saiad

# Run saia TUI from source (dev mode)
run-tui:
	go run ./cmd/saia

# Generate coverage report
coverage: test-coverage

# Install binaries to /usr/local/bin
install:
	install -d /usr/local/bin
	install -m 755 build/saiad /usr/local/bin/saiad
	install -m 755 build/saia /usr/local/bin/saia

# Uninstall binaries from /usr/local/bin
uninstall:
	rm -f /usr/local/bin/saiad
	rm -f /usr/local/bin/saia

# Remove build artifacts
clean:
	rm -rf build/
	rm -f coverage.out coverage.html

# Show version
version:
	@echo "saia version $(VERSION)"

# Dev setup - first time
dev: tidy fmt lint build

# Verify all pre-commit checks pass
verify: fmt lint tidy build test

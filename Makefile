.PHONY: build run test clean install lint help

# Build the binary
build:
	@echo "Building anaphase..."
	@go build -o bin/anaphase ./cmd/anaphase
	@echo "Build complete: bin/anaphase"

# Run the CLI
run: build
	@./bin/anaphase

# Run tests
test:
	@echo "Running tests..."
	@go test ./... -v -race -coverprofile=coverage.out

# Quick test (fast verification)
test-quick:
	@echo "Quick test..."
	@go test ./... -short

# Verify Phase 0
verify:
	@echo "Verifying Phase 0 setup..."
	@go mod verify
	@go build ./...
	@go test ./...
	@make build > /dev/null
	@./bin/anaphase --version
	@echo "✅ Phase 0 verified successfully!"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out
	@echo "Clean complete"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Install the CLI globally
install:
	@echo "Installing anaphase globally..."
	@go install ./cmd/anaphase

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Show help
help:
	@echo "Available targets:"
	@echo "  build    - Build the binary to bin/anaphase"
	@echo "  run      - Build and run the CLI"
	@echo "  test     - Run tests with coverage"
	@echo "  clean    - Remove build artifacts"
	@echo "  deps     - Install dependencies"
	@echo "  install  - Install CLI globally"
	@echo "  lint     - Run linter"
	@echo "  fmt      - Format code"
	@echo "  help     - Show this help message"

# Default target
.DEFAULT_GOAL := help

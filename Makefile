# Makefile for kubeprobes

# Default target
.PHONY: all
all: build

# Variables
BINARY_NAME=kubeprobes
BUILD_DIR=./bin
CMD_DIR=./cmd/kubeprobes
INSTALL_DIR?=/usr/local/bin
ARGS?=

# Build the binary for current platform
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

# Build for multiple platforms  
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)

# Install the binary to system
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)

# Run the application with optional arguments
# Usage: make run ARGS="scan --help"
.PHONY: run
run:
	@go run $(CMD_DIR) $(ARGS)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	@go vet ./...

# Download dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all        - Default target, builds the binary"
	@echo "  build      - Build the binary for current platform"
	@echo "  build-all  - Build binaries for multiple platforms"
	@echo "  install    - Install binary to $(INSTALL_DIR) (configurable with INSTALL_DIR)"
	@echo "  run        - Run the application (use ARGS to pass parameters)"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"
	@echo "  fmt        - Format code"
	@echo "  lint       - Lint code"
	@echo "  deps       - Download and tidy dependencies"
	@echo "  help       - Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make run ARGS=\"scan --help\""
	@echo "  make install INSTALL_DIR=/opt/bin"

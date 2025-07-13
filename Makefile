# Makefile for kubeprobes

# Default target
.PHONY: all
all: build

# Variables
BINARY_NAME=kubeprobes
SRC_DIR=src
INSTALL_DIR?=/usr/local/bin
ARGS?=

# Build the binary
.PHONY: build
build:
	cd $(SRC_DIR) && go build -o $(BINARY_NAME) probes.go

# Install the binary
.PHONY: install
install: build
	sudo cp $(SRC_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)

# Run the application with optional arguments
# Usage: make run ARGS="scan --help"
.PHONY: run
run: build
	cd $(SRC_DIR) && ./$(BINARY_NAME) $(ARGS)

# Clean build artifacts
.PHONY: clean
clean:
	cd $(SRC_DIR) && rm -f $(BINARY_NAME)

# Run tests
.PHONY: test
test:
	cd $(SRC_DIR) && go test -v ./...

# Format code
.PHONY: fmt
fmt:
	cd $(SRC_DIR) && go fmt ./...

# Lint code
.PHONY: lint
lint:
	cd $(SRC_DIR) && go vet ./...

# Download dependencies
.PHONY: deps
deps:
	cd $(SRC_DIR) && go mod download && go mod tidy

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all        - Default target, builds the binary"
	@echo "  build      - Build the binary"
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
# Build stage
FROM cgr.dev/chainguard/go:1.24 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./
COPY go.sum ./

# Download dependencies and generate go.sum
RUN go mod download && go mod tidy

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/kubeprobes ./cmd/kubeprobes

# Final stage
FROM cgr.dev/chainguard/static:latest

# Create a non-root user
# Note: chainguard/static already includes a non-root user 'nonroot' with UID 65532
USER 65532:65532

# Copy the binary from builder
COPY --from=builder /app/kubeprobes /usr/local/bin/kubeprobes

# Set the entrypoint
ENTRYPOINT ["kubeprobes"]

# Simple healthcheck using the binary itself
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD ["kubeprobes", "--help"]

# Default command
CMD ["--help"]
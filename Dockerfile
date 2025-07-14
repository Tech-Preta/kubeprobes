# Build stage
FROM golang:1.24-alpine AS builder

# Install ca-certificates for HTTPS requests during build
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy go mod files to leverage layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/kubeprobes ./cmd/kubeprobes

# Final stage - using alpine for better compatibility
FROM alpine:3.22.0

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S kubeprobes && \
    adduser -u 1001 -S kubeprobes -G kubeprobes

# Copy the binary from builder
COPY --from=builder /app/kubeprobes /usr/local/bin/kubeprobes

# Set user
USER kubeprobes:kubeprobes

# Set the entrypoint
ENTRYPOINT ["kubeprobes"]

# Simple healthcheck using the binary itself
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD ["kubeprobes", "--help"]

# Default command
CMD ["--help"]
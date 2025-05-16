# Build stage
FROM cgr.dev/chainguard/go:latest AS builder

WORKDIR /app/src

# Copy go mod files
COPY src/go.mod ./
COPY src/go.sum ./

# Download dependencies and generate go.sum
RUN go mod download && go mod tidy

# Copy source code
COPY src/. .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/kubeprobes probes.go

# Final stage
FROM cgr.dev/chainguard/static:latest

# Copy the binary from builder
COPY --from=builder /app/kubeprobes /usr/local/bin/kubeprobes

# Set the entrypoint
ENTRYPOINT ["kubeprobes"]

# Default command
CMD ["--help"]
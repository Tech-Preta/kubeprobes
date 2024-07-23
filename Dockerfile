# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.23rc2 as builder

# Copy local code to the container image.
WORKDIR /go/src/projeto
COPY . .

# Build the command inside the container.
RUN go get -d -v ./...
RUN go install -v ./...

# Use a Docker multi-stage build to create a lean production image.
FROM golang:1.23rc2
COPY --from=builder /go/bin/projeto /projeto

# Run the web service on container startup.
CMD ["/projeto"]
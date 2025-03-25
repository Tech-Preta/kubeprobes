# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
<<<<<<< HEAD
FROM golang:1.24.1-alpine3.21 as builder
=======
FROM FROM golang:1.23.1 as builder

>>>>>>> main

# Copy local code to the container image.
WORKDIR /go/src/projeto
COPY . .

# Build the command inside the container.
RUN go get -d -v ./...
RUN go install -v ./...

# Use a Docker multi-stage build to create a lean production image.
FROM golang:1.24.1-alpine3.21
COPY --from=builder /go/bin/projeto /projeto

# Run the web service on container startup.
CMD ["/projeto"]
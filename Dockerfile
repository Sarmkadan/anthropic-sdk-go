# syntax=docker/dockerfile:1.4

# Development environment for building and testing the Go SDK
FROM golang:1.23.0-alpine3.20


# Install development tools
RUN apk add --no-cache \
    ca-certificates \
    git \
    make \
    bash \
    curl \
    jq \
    gcc \
    musl-dev

WORKDIR /app

# Copy go mod and sum files first for better layer caching
COPY go.mod go.sum ./


# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the SDK
RUN go build -v ./...

# Install the SDK
RUN go install ./...

# Set environment variables
ENV GOPATH=/go
ENV PATH="$GOPATH/bin:$PATH"

# Default command - show SDK version
CMD ["go", "version"]
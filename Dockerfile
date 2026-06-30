# syntax=docker/dockerfile:1.4

# Build stage for the Go SDK
FROM golang:1.23.0-alpine3.20 AS builder

WORKDIR /app

# Copy go mod and sum files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the SDK binary
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /anthropic-sdk-go

# Final stage - minimal runtime image
FROM alpine:3.20

# Install CA certificates for HTTPS support
RUN apk add --no-cache ca-certificates

# Copy the compiled binary from the builder
COPY --from=builder /anthropic-sdk-go /usr/local/bin/anthropic-sdk-go

# Set environment variables
ENV ANTHROPIC_API_KEY=""
ENV GOPATH=/go
ENV PATH="$GOPATH/bin:$PATH"

# Default command
CMD ["anthropic-sdk-go"]

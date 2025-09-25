# Stage 1: Build

FROM golang:1.25.1-alpine AS builder

# Install git for dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build a static Linux binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o subway-service ./cmd/server/main.go

# Stage 2: Minimal runtime

FROM alpine:latest

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/subway-service .

# Set ownership and switch to non-root
USER appuser

# Expose port
EXPOSE 8080

# Environment variable
ENV MBTA_API_KEY=

# Run the service
CMD ["./subway-service"]


# Use Go 1.25.1 image
FROM golang:1.25.1

# Set working directory
WORKDIR /app

# Copy all source code
COPY . .

# Download dependencies
RUN go mod download

# Build the binary
RUN go build -o subway-service ./cmd/server/main.go

# Expose the port
EXPOSE 8080

# Set environment variable (replace with your key)
ENV MBTA_API_KEY="81058e8b660047bbb09238cc53958f27"

# Run the binary
CMD ["./subway-service"]
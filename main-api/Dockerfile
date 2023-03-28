# Build stage
FROM golang:1.20 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code
COPY . .

# Build the application as a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main cmd/api/main.go

# Production stage
FROM alpine:latest

# Install ca-certificates for HTTPS support
RUN apk --no-cache add ca-certificates

# Copy the main executable from the builder stage
COPY --from=builder /app/main /app/

# Copy the configuration file
COPY config-dev.yaml /app/

# Make the main executable executable
RUN chmod +x /app/main

# Expose the application port
EXPOSE 8085

# Set the entrypoint to run the main executable
ENTRYPOINT ["/app/main"]

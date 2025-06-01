# Start from the official Go image for building
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Install 'file' utility
RUN apt-get update && apt-get install -y file

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o app main.go

# Use a minimal Debian image for running
FROM debian:bullseye-slim

# Install minimal packages needed (chrome dependencies etc.)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/app .

# Ensure executable permissions
RUN chmod +x ./app

# Set the entrypoint
ENTRYPOINT ["./app"]
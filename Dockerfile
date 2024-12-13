FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o ininal-importer

FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/ininal-importer .

# Run the application
ENTRYPOINT ["./ininal-importer"]

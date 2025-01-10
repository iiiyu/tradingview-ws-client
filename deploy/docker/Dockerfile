# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/app

# Final stage
FROM alpine:3.19

# Add non root user
RUN adduser -D -g '' appuser

# Install certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Use non root user
USER appuser

# Expose port
EXPOSE 3333

# Command to run
CMD ["./main"]

# Build stage
FROM golang:1.23.3-alpine AS builder

WORKDIR /app

# Install required packages
RUN apk add --no-cache git gcc musl-dev file

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/server/main.go

# Verify the build
RUN ls -la /app/main && file /app/main && pwd

# Final stage
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache file

# Copy binary and required files
COPY --from=builder /app/main /app/main
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/views ./views
COPY --from=builder /app/static ./static
COPY --from=builder /app/internal/config/app.conf ./conf/app.conf

# Set proper permissions
RUN chmod +x /app/main

# Debug: Show contents and file info
RUN pwd && \
    ls -la /app && \
    file /app/main && \
    echo "Binary exists at:" && \
    which main || echo "Binary not in PATH" && \
    echo "Current directory contains:" && \
    ls -la

EXPOSE 8080

# Change to use absolute path
CMD ["/app/main"]
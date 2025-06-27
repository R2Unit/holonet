# Dockerfile for Holonet Core
# Multi-stage build:
# 1. Builder stage: Compile the Go application
# 2. Final stage: Create a minimal runtime image using debian:bookworm-slim

FROM golang:1.23.4 AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the main application
RUN CGO_ENABLED=0 GOOS=linux go build -o holonet-core ./cmd/core

# RUN ls -la /app

# Second stage: debian:bookworm-slim
FROM debian:bookworm-slim

WORKDIR /app

# Install necessary packages
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    netcat-openbsd \
    bash \
    && rm -rf /var/lib/apt/lists/*

# Copy the binaries from the builder stage
COPY --from=builder /app/holonet-core /app/
RUN chmod +x /app/holonet-core

# Set environment variables
ENV PATH="/usr/local/bin:${PATH}"
ENV LOG_LEVEL=info
ENV HOLONET_DEVELOP=true
ENV NETBOX_HOST=http://netbox:8000
ENV NETBOX_API_TOKEN=your_token_here
# Database connection (customize these values as needed)
ENV DB_HOST=postgres
ENV DB_PORT=5432
ENV DB_USER=holonet
ENV DB_PASSWORD=insecure
ENV DB_NAME=holonet
# Cache connection
ENV VALKEY_HOST=valkey
ENV VALKEY_PORT=6379
# Admin user configuration (defaults will be used if not set)
ENV ADMIN_USERNAME=admin
ENV ADMIN_EMAIL=admin@example.com
ENV ADMIN_PASSWORD=insecure

# Expose the port the app runs on
EXPOSE 3000

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD nc -z localhost 3000 || exit 1

# Command to run the application
CMD ["./holonet-core"]

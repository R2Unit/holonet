# syntax=docker/dockerfile:1

##########################
# Build Stage (builder)
##########################
FROM golang:1.23.6 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum if your project uses modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the application (output binary is named "holonet")
RUN CGO_ENABLED=0 GOOS=linux go build -o holonet .

##########################
# Final Stage
##########################
FROM debian:bookworm-slim

# Set working directory for the runtime image
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/holonet .

# Copy the entire template directory from the builder stage
COPY --from=builder /app/template ./template

# Optionally expose a port (adjust if necessary)
EXPOSE 8080

# Command to run the application
CMD ["./holonet"]

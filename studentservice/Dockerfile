# Stage 1: Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum to leverage Docker cache during builds
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code (after dependencies are downloaded)
COPY . .

# Tidy up the dependencies and build the binary
RUN go mod tidy

# Build the Go binary with specific flags to strip debugging symbols and optimize size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Stage 2: Lightweight runtime stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates to allow HTTPS requests (for accessing external APIs)
RUN apk --no-cache add ca-certificates

# Define build-time arguments
ARG VAULT_ADDR
ARG VAULT_TOKEN
ARG SERVICE_NAME

# Set the environment variables that will be available at runtime
ENV VAULT_ADDR=${VAULT_ADDR}
ENV VAULT_TOKEN=${VAULT_TOKEN}
ENV SERVICE_NAME=${SERVICE_NAME}

# Copy the Go binary from the builder stage
COPY --from=builder /app/main .

# Create a non-root user for security
RUN adduser -D -g '' appuser

# Switch to the non-root user for running the application
USER appuser

# Set the entrypoint for the container to run the compiled binary
ENTRYPOINT ["./main"]

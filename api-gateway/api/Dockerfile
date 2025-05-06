# Use Golang Alpine base image for building
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install git (needed for go mod downloads)
RUN apk add --no-cache git

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./

# Download Go modules
RUN go mod tidy && go mod vendor

# Copy the rest of the application source code
COPY . .

# Build the application binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server/main.go

# Use a minimal lightweight image for running the app
FROM alpine:latest

# Create a non-root user for security
#RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Ensure the config directory exists inside the container
RUN mkdir -p /app/config

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main /app/main

# Set permissions
RUN chmod +x /app/main

# Set runtime environment variables (to be overridden in docker-compose)
COPY .env.development /app/.env.development
COPY config/config.development.yaml /app/config/config.development.yaml


# Expose application port
EXPOSE 50051
COPY ./scripts/wait-for-kafka.sh /usr/bin/wait-for-kafka.sh
RUN chmod +x /usr/bin/wait-for-kafka.sh

ENTRYPOINT ["/usr/bin/wait-for-kafka.sh"]
# Set the command to run the application
CMD ["/app/main"]

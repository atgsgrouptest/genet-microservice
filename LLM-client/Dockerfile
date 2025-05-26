# ---------- Stage 1: Build ----------
FROM golang:1.24.2 AS builder

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Create working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go Fiber app
RUN go build -o fiber-app .

# ---------- Stage 2: Run ----------
FROM alpine:latest

# Install certificate dependencies (required for HTTPS requests)
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/fiber-app .

# Set permissions
RUN chown -R appuser /app
USER appuser

# Expose application port (adjust based on your app)
EXPOSE 8002

# Start the app
ENTRYPOINT ["./fiber-app"]

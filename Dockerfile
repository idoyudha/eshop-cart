# Step 1: Modules caching
FROM golang:1.23.4 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.23.4 as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
# Build the application with optimization flags
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/main ./cmd/app

# Step 3: Final for production
FROM debian:bullseye-slim as production
# Add CA certificates and timezone data
RUN apt-get update && \
    apt-get install -y ca-certificates tzdata && \
    rm -rf /var/lib/apt/lists/*

# Create a non-root user
RUN useradd -r -u 1001 -g root appuser

# Create app directory and set permissions
RUN mkdir /app && chown appuser:root /app

# Copy the binary from builder
COPY --from=builder /go/bin/main /app/

# Use the non-root user
USER appuser

# Set the working directory
WORKDIR /app

# Command to run the application
CMD ["./main"]
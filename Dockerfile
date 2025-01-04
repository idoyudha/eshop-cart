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
# Build with CGO enabled for Kafka
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o main .

# Step 3: Final for production
FROM redhat/ubi8-minimal

# Create a non-root user
RUN microdnf install shadow-utils -y && \
    useradd -r -u 1001 -g root appuser && \
    microdnf remove shadow-utils -y && \
    microdnf clean all

# Copy the binary from builder
COPY --from=builder /app/main /app/main

# Use the non-root user
USER appuser

# Set the working directory
WORKDIR /app

# Command to run the application
ENTRYPOINT ["/app/main"]
FROM golang:1.25.1 AS builder
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server/main.go

# Use a smaller image for the final container
FROM alpine:latest
WORKDIR /app

# Copy the binary and config file from builder
COPY --from=builder /app/server .
COPY config.yaml .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./server"]
# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -o temperature-api

# Runtime stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/temperature-api .

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Expose the port the app runs on
EXPOSE 8081

# Command to run the application
CMD ["./temperature-api"]

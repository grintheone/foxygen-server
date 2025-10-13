# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy SSL certificates
COPY --from=builder /app/ssl ./ssl
# Copy built binary
COPY --from=builder /app/main .
# Copy setup.sql (if needed by the app)
COPY --from=builder /app/setup.sql .
# Copy .env file
COPY --from=builder /app/.env .

# Expose ports
# EXPOSE 8080 443
EXPOSE 8080 8443

# Command to run the executable
CMD ["./main"]

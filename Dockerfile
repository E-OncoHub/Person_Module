#Build stage
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# Start a new stage from scratch
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy the .env file
COPY .env .

# Copy the wallet directory
COPY wallet_oncodb ./wallet_oncodb

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
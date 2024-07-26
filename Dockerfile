# Start from the official Go image
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
# If you have a go.sum file, uncomment the next line
# COPY go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Start a new stage from scratch
FROM alpine:latest

# Add ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

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
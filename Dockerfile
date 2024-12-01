# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod ./

# Copy the rest of the application's source code
COPY . .

# Build the Go app
RUN go build -o mp-app ./cmd/main

# Stage 2: Create a smaller image for running the Go app
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/mp-app .

# Install any necessary certificates (if your app needs to make HTTPS requests)
RUN apk add --no-cache ca-certificates

# Expose the port the app will run on
EXPOSE 8080

# Command to run the app
CMD ["./mp-app"]
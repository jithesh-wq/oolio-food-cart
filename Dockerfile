# Use official Go base image
FROM golang:1.23-alpine

# Set working directory inside the container
WORKDIR /app

# Install required packages
RUN apk add --no-cache git tzdata

# Copy go.mod and go.sum separately to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy all other files
COPY . .

# Build the binary (assuming main.go is in /cmd)
RUN go build -o server ./cmd

# Expose port (optional; match your app port)
EXPOSE 8080

# Command to run the app
CMD ["./server"]

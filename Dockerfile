# Use official Go image
FROM golang:1.22-alpine


# Set working directory
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Set to /app/cmd where main.go lives
WORKDIR /app/cmd

# Build the binary
RUN go build -o server

# Expose port 
EXPOSE 8080

# Start the server
CMD ["./server"]

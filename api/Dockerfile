# Use official Go image as builder
FROM golang:1.21-alpine

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o upload-service

# Expose port
EXPOSE 8080

# Create directory for images
RUN mkdir -p content/images/posts

# Run the service
CMD ["./upload-service"]
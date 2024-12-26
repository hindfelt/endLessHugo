# Use official Go image as builder
FROM golang:1.21-alpine

# Set working directory
WORKDIR /app

# Install git for private repos (if needed)
RUN apk add --no-cache git

# Copy go mod file
COPY go.mod ./

# Download dependencies and show any errors
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Show directory contents and try to build with verbose output
RUN ls -la && \
    CGO_ENABLED=0 GOOS=linux go build -v -o upload-service

# Expose port
EXPOSE 8080

# Create directory for images
RUN mkdir -p content/images/posts

# Run the service
CMD ["./upload-service"]
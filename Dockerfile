# Create build stage based on buster image
FROM golang:1.22-alpine AS builder
# Create working directory under /app
WORKDIR /app
# Copy over all go config (go.mod, go.sum etc.)
COPY . ./
# Install any required modules
RUN go mod download
# Run the Go build and output binary under hello_go_http
RUN go build -o /server
# Make sure to expose the port the HTTP server is using
EXPOSE 8080
# Run the app binary when we run the container
ENTRYPOINT ["/server"]

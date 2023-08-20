# Use the official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the entire current directory into the container's working directory
COPY . .

# Build the Go application
RUN go build -o main .

# Expose port 8080 for the HTTP server
EXPOSE 8080

# Start the application
CMD ["./main"]

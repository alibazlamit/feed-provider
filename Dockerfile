# Use the official Golang image as the base image
FROM golang:1.16-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the necessary files to the working directory
COPY go.mod go.sum ./
COPY main.go ./
COPY . .

# Build the Go application
RUN go build -o feed-provider ./main.go

# Expose the port that the application listens on
EXPOSE 8080

# Set the entrypoint command to run the application
CMD ["./feed-provider"]
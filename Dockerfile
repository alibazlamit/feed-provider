# Use an official Golang runtime as the base image
FROM golang:1.19.0-alpine3.16 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum . /
RUN go mod download

# Copy the rest of the project files
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o ./feed-reader/main.go
# RUN go build -o feed-reader.
EXPOSE 8080
# Set the entrypoint command to run the Go application
CMD ["./feed-reader"]
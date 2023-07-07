# Use an official Golang runtime as a parent image
FROM golang:1.16-alpine

# Set the working directory to /go/src/app
WORKDIR /go/src/app

# Copy the current directory contents into the container at /go/src/app
COPY . .

# Install required packages for building sqlite driver
RUN apk add git gcc musl-dev

# Download dependencies
RUN go mod download

# Build the Go app
RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64 GOARM=7 go build -ldflags="-s -w" -o /go/bin/app cmd/server/main.go

# Use an official Alpine Linux image as a parent image
FROM alpine:latest

# Install SQLite
RUN apk add ca-certificates sqlite

# Set the working directory to /app
WORKDIR /app

# Copy the binary from the first stage
COPY --from=0 /go/bin/app .

# Expose port 8000 for the Go app to listen on
EXPOSE 8000

# Start the Go app
CMD ["./app", "server"]

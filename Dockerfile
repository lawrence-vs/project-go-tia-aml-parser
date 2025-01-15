# Use the official Go image as the base image
FROM golang:1.20

# Copy go.mod and go.sum files to the container
COPY /app ./

# Set the working directory in the container
WORKDIR /app/src

# Download dependencies
RUN go mod download

# Build the Go application
# RUN go build -o app .

# Command to run the application
CMD ["bash"]

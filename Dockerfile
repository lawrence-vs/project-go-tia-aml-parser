# Use the official Go image as the base image
FROM golang:1.20

# Copy go.mod and go.sum files to the container
COPY /app ./run.sh ./

# Set the working directory in the container
WORKDIR /app/src

# Download dependencies
RUN mod init tia-aml-parser \
&& go mod download \
&& go install github.com/xuri/excelize/v2@latest \
&& go mod tidy \
&& go build -o tia-aml-parser

# Command to run the application
CMD ["bash"]

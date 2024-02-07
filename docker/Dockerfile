# Start from the official Golang image.
FROM golang:1.20

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod files first to leverage Docker cache.
COPY go.mod ./

# Download necessary Go modules.
RUN go mod download

# Copy the rest of the application source code.
COPY . .

# Build the application.
RUN go build -o main .

# Expose port 8080 to the outside world.
EXPOSE 8080

# Command to run the executable.
CMD ["./main"]

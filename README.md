# Go HTTP Request Counter Server

## Introduction

This Go application is an HTTP server that counts the number of requests received in the last 60 seconds (moving window). The count is persisted to a file, ensuring that the data is not lost when the server restarts. This project demonstrates the use of Go's standard library to create a web server, manage concurrent requests, and handle file operations.

## Prerequisites

Before running this application, ensure you have the following installed:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

These tools are required to build and run the application inside a Docker container, ensuring a consistent environment and easy deployment.

## Running the Application

To run the application, follow these steps:

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/mtaha1996/request-counter
   cd request-counter
   ```

2. **Build and Run with Docker Compose:**

   Use Docker Compose to build and run the application:

   ```bash
   docker-compose up --build
   ```

   This command builds the Docker image and starts the server. The server will be accessible at `http://localhost:8080`.

3. **Accessing the Application:**

   Send requests to the server using a web browser or tools like `curl`:

   ```bash
   curl http://localhost:8080
   ```

   The server will respond with the number of requests received in the past 60 seconds.

## Persistent Data

The application persists request count data in a file named `storage.json`. This file is stored outside the Docker container in a mounted volume to ensure data persistence across container restarts.

## Stopping the Application

To stop the application, use the following Docker Compose command:

```bash
docker-compose down
```

## Additional Notes

- The application is configured to run on port 8080. If this port is occupied on your machine, you may need to modify the port configuration in the `docker-compose.yml` file.
- The base image used in the Dockerfile is `golang:1.20`. You can change this to a different version if required, but ensure it's compatible with the application.

## Configuration and Flags

This application supports various configurations through command-line flags. These flags allow you to customize certain behaviors and settings without modifying the code.

### Available Flags

- **window-size**: Sets the size of the moving window in seconds. This is the time frame for counting requests. Default is 60 seconds.
  
  Usage: `-window-size=<value>`

- **persist-interval**: Defines the interval in seconds at which the data is persisted to the file. Default is 1 second.
  
  Usage: `-persist-interval=<value>`

- **data-ttl**: Time-to-live for data in seconds. This is the duration after which the data is considered outdated and is removed. Default is 60 seconds.
  
  Usage: `-data-ttl=<value>`

- **storage-path**: Specifies the path to the storage file where request counts are persisted. Default is `storage/storage.json`.
  
  Usage: `-storage-path=<path>`

#### Running the Application with Flags

You can start the application with custom configurations by passing these flags via the command line. For example:

```sh
go run main.go -window-size=120 -persist-interval=2 -data-ttl=120 -storage-path="path/to/storage.json"
```

In this example, the application is started with a 120-second window size, a persistence interval of 2 seconds, a data TTL of 120 seconds, and a custom storage path. You can omit any flag to use its default value.

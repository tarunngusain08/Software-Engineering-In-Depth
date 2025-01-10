# Redis Streams with Go - Producer and Consumer

This project demonstrates how to use Redis Streams with a producer-consumer pattern in a Go application. It uses Docker Compose to set up Redis and the Go app as services.

## Project Overview

The project includes:

1. A **Redis service** running in Docker.
2. A **Go application** that acts as both a producer and a consumer of messages in Redis Streams.
   - The **Producer** sends messages to a Redis stream.
   - The **Consumer** reads from the stream and processes messages.

### Technologies Used
- **Go 1.20** for the backend application
- **Redis** for stream-based messaging
- **Docker Compose** to manage services and their dependencies

## Getting Started

### Prerequisites
Ensure you have the following tools installed on your machine:
- [Docker](https://www.docker.com/products/docker-desktop)
- [Docker Compose](https://docs.docker.com/compose/install/)
- Go 1.20 (if you wish to run the app locally outside Docker)

### Setup and Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/tarunngusain08/Software-Engineering-In-Depth
   cd AsyncQueueing/redis
   ```

2. **Build and run the application** using Docker Compose:

   ```bash
   docker-compose up --build
   ```

   This will:
   - Build the Go application using the `Dockerfile`
   - Start a Redis container
   - Start the Go app container that connects to Redis and performs the producer-consumer tasks.

   The Go app will listen for messages on a Redis stream (`dataStream`), produce messages, and consume them concurrently.

3. **Stopping the services**:

   To stop the containers and remove them:

   ```bash
   docker-compose down
   ```

### Configuration

The Docker Compose file (`docker-compose.yml`) defines two services:
1. **Redis service**:
   - The Redis service runs the latest Redis image (`redis:latest`).
   - It exposes port `6379` for communication.
   - It has a health check to ensure the Redis container is up and running.

2. **Go app service**:
   - The Go app service is built from the `Dockerfile` and exposes port `8081`.
   - It depends on the Redis service, meaning it will only start after Redis is healthy.

### Application Logic

#### Producer
The producer generates 10 integer messages (0-9) and adds them to a Redis stream (`dataStream`). It uses the Redis `XAdd` command to add messages to the stream.

#### Consumer
The consumer listens to the stream (`dataStream`) using a Redis consumer group (`consumerGroup`). It reads messages from the stream, acknowledges them using `XAck`, and processes them. Once the consumer receives and processes the value `9`, it stops.

### Dockerfile
This Dockerfile sets up a Go environment and builds the application:

```dockerfile
FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o app .

EXPOSE 8081

CMD ["./app"]
```

### Go Code Overview
The Go application performs the following steps:
1. **Setting up Redis connection**:
   The Go application connects to the Redis service running in Docker using the hostname `redis_service` on port `6379`.

2. **Setting up the stream**:
   The `XGroupCreateMkStream` command is used to create a stream with a consumer group.

3. **Producer logic**:
   The producer sends 10 messages (integer values 0 to 9) to the stream.

4. **Consumer logic**:
   The consumer continuously listens to the stream, processes the messages, and acknowledges them once processed.

### Go Modules
The Go modules required for this project:

```go
module redis-streams

go 1.20

require github.com/redis/go-redis/v9 v9.7.0

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)
```

### Docker Compose Configuration

```yaml
version: '3.9'
services:
  redis:
    image: redis:latest
    container_name: redis_service
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      retries: 3
      timeout: 5s
      start_period: 5s

  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: go_app
    ports:
      - "8081:8081"
    environment:
      REDIS_HOST: redis
    depends_on:
      redis:
        condition: service_healthy
```

### Troubleshooting

- **Connection refused errors**:
  Ensure the Redis service is running and accessible at `localhost:6379` or `redis_service:6379` (inside the Docker network).
  
- **Health check failures**:
  Check the Redis logs for any startup issues, and ensure Redis is starting correctly.

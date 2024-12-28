# Kafka Write-Through Cache with Consumer Pool and Throughput Testing

This project implements a write-through cache system using Kafka, Redis, and MySQL. The architecture was designed to allow for high throughput while ensuring the data is immediately visible in MySQL, Redis, and metrics on the Confluent Dashboard.

## Project Overview

The project consists of the following key components:
1. **Write-Through Cache**: The system uses Kafka as a messaging layer to buffer requests before writing them to MySQL and Redis.
2. **Producer and Consumer**: A producer sends data to Kafka, and a consumer processes this data and writes it to MySQL and Redis.
3. **Automation Script**: An automation script was created to generate a load of 10 requests per second to simulate real-world traffic.
4. **Consumer Pool**: The consumer pool was implemented to increase throughput by scaling up the number of consumers processing messages from Kafka.
5. **Testing Different Throughput Scenarios**: Various scenarios were tested by changing the producer throughput and scaling the number of consumers.

## Key Steps and Performance Testing

### 1. Write-Through Cache with Producer and Consumer
- The system was designed using Kafka, Redis, and MySQL to implement a write-through cache.
- **Kafka** was used to buffer data between the producer and consumer.
- **Redis** was used for fast in-memory caching, and **MySQL** served as the persistent data store.
- Upon sending a request via HTTP, the data was pushed to Kafka by the producer, and immediately processed and written to MySQL and Redis by the consumer.
- The **Confluent Dashboard** showed real-time metrics, and the data was visible in MySQL and Redis immediately.

### 2. Automation Script to Send 10 Requests/Second
- An automation script was written to simulate traffic by sending 10 requests per second to the `write-through` endpoint.
- The requests were sent using randomly generated data, and each request was directed to the Kafka topic.
- This scenario was run for 4-5 minutes with only **1 consumer**.
- The topic had **6 partitions**, and the consumer was able to process the incoming data at a reasonable rate.

### 3. Consumer Pool with 6 Consumers
- After establishing the baseline performance with 1 consumer, the next step was to scale up the consumer pool.
- The number of consumers was increased from **1 to 6** to process data from all **6 partitions**.
- However, despite increasing the number of consumers, the **consumer lag** started growing significantly.
- The lag increased to **4-5k messages** with a **800 message lag per partition**. This test was run for **10 minutes**.
  
### 4. Increased Producer Throughput (1000 Requests/Second)
- To test how the system would handle increased traffic, the producer throughput was raised to **1000 requests per second** while keeping the number of consumers at **6**.
- This change caused a significant increase in lag, with the total lag growing to **150k messages**, and approximately **25k message lag per partition**.
- The consumer pool was still struggling to keep up with the high message rate, and the lag continued to grow.

### 5. Consumer Processing After Stopping Producer
- After running the producer at 1000 requests/sec for **4-5 minutes**, the producer was stopped to test how well the consumers could catch up on the backlog.
- The **consumer** was allowed to continue consuming messages asynchronously.
- As expected, the **consumer lag** began to decrease gradually as the consumer slowly processed the pending messages.
- The lag continued to reduce as more messages were consumed, and the system eventually reached a steady state.

## Key Observations

- **Throughput Scalability**: The system handled **10 requests/sec** effectively with a single consumer, but the performance started to degrade as the throughput increased to **1000 requests/sec**.
- **Consumer Lag**: The consumer lag grew significantly when the producer throughput was increased while keeping the number of consumers the same. This highlighted the need to scale the consumer pool or optimize consumer processing to keep up with the increased traffic.
- **Consumer Pooling**: Increasing the number of consumers helped alleviate the lag, but the system still faced challenges with lag as the producer throughput was raised.
- **Asynchronous Consumer Processing**: After stopping the producer, the consumer processed the backlog asynchronously, and the lag gradually decreased, showing that the system can catch up when the load is reduced.

## Conclusion

This testing allowed us to understand the system's behavior under different load conditions:
1. The Kafka consumer pool can handle a moderate rate of requests, but as the producer throughput increases, the system needs further optimizations.
2. Adding more consumers can help alleviate lag, but there is still a limit to how quickly the consumers can process messages when the throughput is too high.
3. Stopping the producer and allowing the consumer to process the backlog demonstrates that the system can eventually catch up with the pending messages.

## Future Improvements

- **Scaling Consumers Dynamically**: Implement automatic scaling of consumers based on the message lag or queue length.
- **Optimizing Consumer Processing**: Improve the consumer's ability to process messages faster or in parallel to reduce lag in high throughput scenarios.
- **Load Balancing**: Implement load balancing strategies for Kafka consumers to ensure even distribution of messages across consumers and partitions.
- **Monitoring and Alerts**: Set up automated monitoring and alerts to proactively handle lag and consumer issues before they become critical.

## Requirements

- Go (for building the producer, consumer, and automation script)
- Kafka (with a Confluent Cloud account for monitoring metrics)
- Redis and MySQL for persistent storage and caching
- Docker (optional, for running services locally)

## How to Run

### 1. Set up Kafka, MySQL, and Redis
Start Redis and MySQL containers using Docker:

If you have Docker installed, run the following commands to start Redis and MySQL containers.
```bash
docker run --name redis-container -p 6379:6379 -d redis
docker run --name mysql-container -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -d mysql:latest
```

For Kafka 
Create confluent account -> <https://confluent.cloud/home>
```bash
brew install confluentinc/tap/cli
confluent login
```

Authenticate the confluent login
```
export CONFLUENT_REST_URL=http://localhost:8082
confluent kafka cluster list
confluent kafka cluster use <cluster-id>
confluent kafka topic create users
```

#### Create a cluster in confluent dashboard if not done already!

#### Copy the api-key, secret key and bootstrap-servers and use them to intialize your kafka consumers and producers!

### 2. Running the Producer and Consumer
Run the following Go commands to start the producer and consumer:

```bash
go run producer.go  # Start the producer to send messages to Kafka
go run consumer.go  # Start the consumer to process messages from Kafka
```

### 3. Running the Automation Script
Run the automation script to send 10 requests/sec:

```bash
go run send_requests.go  # This sends random data to the 'write-through' endpoint
```

### 4. Monitor Metrics on Confluent Dashboard
You can monitor the system's performance and Kafka metrics on the **Confluent Dashboard**.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
```

This `README.md` file provides a clear, detailed overview of your project, including the setup, testing scenarios, observations, and areas for future improvement. Let me know if you'd like to make any further adjustments!

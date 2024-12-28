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

### 2. Running the Code
Run the following Go commands to start the code!:

```bash
git clone https://github.com/tarunngusain08/Software-Engineering-In-Depth
cd Software-Engineering-In-Depth/Caching/WriteThrough/Kafka
go run main.go
```

#### Use postman with sample curl - 
```
curl --location 'http://localhost:8080/write-through' \
--header 'Content-Type: application/json' \
--data '{
    "name": "John Doe",
    "age": 34,
    "occupation": "Engineer"
}'
```

### 3. Running the Automation Script
Run the automation script to send 10 requests/sec:
Open a new terminal and create the file in different location!
```
vi main.go
```

```go
package main

import (
        "bytes"
        "encoding/json"
        "fmt"
        "math/rand"
        "net/http"
        "time"
)

func generateRandomData() map[string]interface{} {
        names := []string{"John Doe", "Jane Smith", "Alice Brown", "Bob White", "Charlie Green"}
        occupations := []string{"Engineer", "Doctor", "Teacher", "Artist", "Scientist"}

        return map[string]interface{}{
                "name":       names[rand.Intn(len(names))],
                "age":        rand.Intn(60) + 18, // Random age between 18 and 77
                "occupation": occupations[rand.Intn(len(occupations))],
        }
}

func sendRequest() {
        url := "http://localhost:8080/write-through"
        data := generateRandomData()

        // Marshal data to JSON
        jsonData, err := json.Marshal(data)
        if err != nil {
                fmt.Printf("Error marshaling JSON: %v\n", err)
                return
        }

        // Send the request
        resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
        if err != nil {
                fmt.Printf("Error sending request: %v\n", err)
                return
        }
        defer resp.Body.Close()

        // Log the response status
        fmt.Printf("Response status: %s\n", resp.Status)
}

func main() {
        // Seed the random number generator
        rand.Seed(time.Now().UnixNano())

        // Create a ticker to send a request every 100ms (10 requests/sec)
        ticker := time.NewTicker(100 * time.Millisecond)
        defer ticker.Stop()

        // Start sending requests
        for {
                select {
                case <-ticker.C:
                        sendRequest()
                }
        }
}
```

```bash
go run main.go
```

#### Customize the ticker time as per your requirement!

### 4. Monitor Metrics on Confluent Dashboard
You can monitor the system's performance and Kafka metrics on the **Confluent Dashboard**.


### 5. Screenshots 

### With 1 producer and 1 consumer

#### Producer and Consumer
<img width="500" alt="Screenshot 2024-12-28 at 3 03 42 PM" src="https://github.com/user-attachments/assets/2c11e4dd-4be6-455e-b53d-a4b1579765c1" /> <img width="500" alt="Screenshot 2024-12-28 at 3 03 50 PM" src="https://github.com/user-attachments/assets/4f93d354-9d12-4e3b-8665-c37cb35b5ec0" />

#### DB and Redis
<img width="374" alt="Screenshot 2024-12-28 at 3 04 05 PM" src="https://github.com/user-attachments/assets/04edb037-f20f-42cb-821d-52ef55114124" /> <img width="590" alt="Screenshot 2024-12-28 at 3 04 22 PM" src="https://github.com/user-attachments/assets/48222803-8224-4ecd-8ff6-1f26b57b2bbd" />

#### Producer, Consumer and Consumer Lag Metrics
<img width="330" alt="Screenshot 2024-12-28 at 3 15 42 PM" src="https://github.com/user-attachments/assets/c83698e3-3aeb-493d-8157-3808d12b9dd1" />
<img width="330" alt="Screenshot 2024-12-28 at 3 15 49 PM" src="https://github.com/user-attachments/assets/134112fb-0158-4dee-ac80-d9e8c7a9666b" />
<img width="330" alt="Screenshot 2024-12-28 at 3 16 13 PM" src="https://github.com/user-attachments/assets/32f9bae6-afdb-49f5-9656-ad040ccab812" />


#### With 1 producer and consumer pool of 6 consumers

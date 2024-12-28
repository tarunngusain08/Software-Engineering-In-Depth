package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"sync"

	"github.com/redis/go-redis/v9"
	_ "github.com/go-sql-driver/mysql"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Struct to hold request data
type requestData struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Occupation string `json:"occupation"`
}

var (
	cache          *redis.Client
	db             *sql.DB
	kafkaProducer  *kafka.Producer
	kafkaConsumer  *kafka.Consumer
	consumerGroup  = "user-consumer-group"
	topic          = "users"
	workerPoolSize = 6
)

func init() {
	// Initialize Redis client
	cache = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       0,  // Default DB
	})

	// Initialize MySQL connection
	var err error
	dsn := "root:1234@tcp(localhost:3306)/users"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// Verify MySQL connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("MySQL connection failed: %v", err)
	}

	// Initialize Kafka producer
	kafkaProducer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  "pkc-619z3.us-east1.gcp.confluent.cloud:9092", // Replace with your bootstrap servers
		"sasl.username":      "OAJPGNTLBH6KR2HF",  // Your Confluent Cloud API Key
		"sasl.password":      "UpsT5OHHkk0GdFak/EEAuBYFLEkpgHvqBm6YKxwF3my2DU06OFHoWLuTFUOOa24S",  // Your Confluent Cloud API Secret
		"security.protocol":  "SASL_SSL",
		"sasl.mechanism":     "PLAIN",
		"acks":               "all", // Ensure reliability
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}

	// Initialize Kafka consumer
	kafkaConsumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "pkc-619z3.us-east1.gcp.confluent.cloud:9092", // Replace with your bootstrap servers
		"sasl.username":      "OAJPGNTLBH6KR2HF",  // Your Confluent Cloud API Key
		"sasl.password":      "UpsT5OHHkk0GdFak/EEAuBYFLEkpgHvqBm6YKxwF3my2DU06OFHoWLuTFUOOa24S",  // Your Confluent Cloud API Secret
		"group.id":           consumerGroup,
		"auto.offset.reset":  "earliest", // Start reading from the earliest offset
		"security.protocol":  "SASL_SSL",
		"sasl.mechanism":     "PLAIN",
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	// Subscribe to Kafka topic
	err = kafkaConsumer.Subscribe(topic, nil)
	if err != nil {
		log.Printf("Consumer failed to subscribe to topic: %v", err)
		return
	}

	log.Println("Initialized Redis, MySQL, and Kafka")
}

func main() {
	// Setup shutdown signal handling
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the worker pool for consuming Kafka messages
	go startConsumerWorkers()

	http.HandleFunc("/write-through", writeThroughHandler)
	// Start the HTTP server in a goroutine
	go func() {
		log.Println("Server started at :8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-stopChan

	// Graceful shutdown
	shutdown()
}

// Handler for the /write-through endpoint
func writeThroughHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userData requestData
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(&userData)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	// Write to Redis cache
	cacheErr := cache.Set(ctx, userData.Name, fmt.Sprintf("%v", userData), 0).Err()
	if cacheErr != nil {
		http.Error(w, "Redis error: "+cacheErr.Error(), http.StatusInternalServerError)
		return
	}

	// Write to Kafka (asynchronously)
	err = produceData(userData)
	if err != nil {
		log.Printf("Produce error: %v", err)
		http.Error(w, "Kafka produce error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Data write successful!")
}

// Produce message to Kafka
func produceData(userData requestData) error {
	// Marshal userData to JSON
	value, err := json.Marshal(userData)
	if err != nil {
		log.Printf("Failed to marshal userData to JSON: %v", err)
		return err
	}

	// Prepare Kafka message
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}

	// Produce message to Kafka (acks=all for reliability)
	return kafkaProducer.Produce(message, nil)
}

// Consumer worker function
func startConsumerWorkers() {
	// Create multiple workers to consume messages from Kafka topic
	var wg sync.WaitGroup

	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			consumeData()
		}(i)
	}

	// Wait for all consumers to finish (if needed for graceful shutdown)
	wg.Wait()
}

// Consumer function to read messages from Kafka and write to MySQL
func consumeData() {
	log.Printf("Consumer started")

	for {
		// Consume messages from Kafka
		msg, err := kafkaConsumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer failed to read message: %v", err)
			continue
		}

		log.Printf("Consumer received message: %s", string(msg.Value))

		// Process the message (write to MySQL)
		var userData requestData
		err = json.Unmarshal(msg.Value, &userData)
		if err != nil {
			log.Printf("Consumer failed to unmarshal message: %v", err)
			continue
		}

		// Write to MySQL
		err = writeToDatabase(userData)
		if err != nil {
			log.Printf("Consumer failed to write to database: %v", err)
			continue
		}

		// Manually commit the offset after successful processing
		_, err = kafkaConsumer.CommitOffsets([]kafka.TopicPartition{
			{Topic: &topic, Partition: msg.TopicPartition.Partition, Offset: msg.TopicPartition.Offset + 1},
		})
		if err != nil {
			log.Printf("Consumer failed to commit offset: %v", err)
		}
	}
}

// Write to MySQL database
func writeToDatabase(userData requestData) error {
	query := "INSERT INTO users (name, age, occupation) VALUES (?, ?, ?)"
	_, err := db.Exec(query, userData.Name, userData.Age, userData.Occupation)
	return err
}

// Graceful shutdown
func shutdown() {
	log.Println("Shutting down gracefully...")
	if kafkaConsumer != nil {
		kafkaConsumer.Close()
	}
	if kafkaProducer != nil {
		kafkaProducer.Close()
	}
	if db != nil {
		db.Close()
	}
	if cache != nil {
		cache.Close()
	}
	log.Println("Shutdown complete")
}

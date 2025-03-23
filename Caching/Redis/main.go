package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		PoolSize: 1000,
	})

	// Check Redis connectivity
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v\n", err)
	}
	log.Println("Connected to Redis successfully")
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Failed to decode request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for key, value := range data {
		log.Printf("Setting key: %s, value: %s\n", key, value)
		err := rdb.Set(ctx, key, value, 0).Err()
		if err != nil {
			log.Printf("Failed to set key %s in Redis: %v\n", key, err)
			http.Error(w, fmt.Sprintf("Failed to set value in Redis: %v", err), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Values set successfully")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	keys := r.URL.Query()["key"]
	if len(keys) == 0 {
		log.Println("No key parameter provided")
		http.Error(w, "Key parameter is required", http.StatusBadRequest)
		return
	}

	log.Printf("Received /get request for keys: %v\n", keys)

	result := make(map[string]string)
	for _, key := range keys {
		value, err := rdb.Get(ctx, key).Result()
		if err == redis.Nil {
			log.Printf("Key not found: %s\n", key)
			result[key] = "Key not found"
		} else if err != nil {
			log.Printf("Failed to get key %s from Redis: %v\n", key, err)
			http.Error(w, "Failed to get value from Redis", http.StatusInternalServerError)
			return
		} else {
			log.Printf("Retrieved key: %s, value: %s\n", key, value)
			result[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	initRedis()

	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get", getHandler)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

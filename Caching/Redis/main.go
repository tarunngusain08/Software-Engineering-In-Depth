package main

import (
	"encoding/json"
	"fmt"
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
		Addr: "localhost:6379", // Redis server address
	})
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for key, value := range data {
		err := rdb.Set(ctx, key, value, 0).Err()
		if err != nil {
			http.Error(w, "Failed to set value in Redis", http.StatusInternalServerError)
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
		http.Error(w, "Key parameter is required", http.StatusBadRequest)
		return
	}

	result := make(map[string]string)
	for _, key := range keys {
		value, err := rdb.Get(ctx, key).Result()
		if err == redis.Nil {
			result[key] = "Key not found"
		} else if err != nil {
			http.Error(w, "Failed to get value from Redis", http.StatusInternalServerError)
			return
		} else {
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
	http.ListenAndServe(":8080", nil)
}

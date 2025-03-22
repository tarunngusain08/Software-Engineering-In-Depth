package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	_ "github.com/go-sql-driver/mysql"
)

// Struct to hold user data
type requestData struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Occupation string `json:"occupation"`
}

var cache *redis.Client
var db *sql.DB

func init() {
	// Initialize Redis client
	cache = redis.NewClient(&redis.Options{
		// Addr: "localhost:6379", // Update with your Redis address
		Addr: "redis_service:6379",
	})

	// Initialize MySQL connection
	var err error
	// db, err = sql.Open("mysql", "root:1234@tcp(localhost:3306)/users")
	db, err = sql.Open("mysql", "root:1234@tcp(mysql_service:3306)/users")
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	log.Println("Mysql and Redis client init!")
}

func main() {
	http.HandleFunc("/read-cache-aside", cacheAsideReadHandler)
	http.HandleFunc("/write-cache-aside", cacheAsideWriteHandler)
	log.Println("Server started at :8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// Handler for reading data (cache-aside pattern)
func cacheAsideReadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data requestData
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	// Check Redis cache
	val, err := cache.Get(ctx, data.Name).Result()
	if err == nil {
		// Data found in cache
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(val))
		fmt.Println("Data retrieved from cache")
		return
	} else if err != redis.Nil {
		// Redis error
		http.Error(w, "Redis error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Data not found in cache, retrieve from database
	userData, err := readFromDatabase(data.Name)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if userData == nil {
		http.Error(w, "Data not found", http.StatusNotFound)
		fmt.Println("Data not present in Database")
		return
	}

	// Store data in Redis
	userJSON, _ := json.Marshal(userData)
	cache.Set(ctx, data.Name, string(userJSON), 5*time.Minute) // Cache with a 5-minute TTL

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.Write(userJSON)
	fmt.Println("Data retrieved from database and cached")
}

// Handler for writing data (cache-aside pattern)
func cacheAsideWriteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data requestData
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Write directly to the database
	err = writeToDatabase(data)
	if err != nil {
		http.Error(w, "Database write error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Invalidate the cache for this key
	ctx := context.Background()
	cache.Del(ctx, data.Name)

	fmt.Fprintf(w, "Data written successfully!")
	fmt.Println("Cache invalidated for key:", data.Name)
}

// Read from MySQL database
func readFromDatabase(name string) (*requestData, error) {
	query := "SELECT name, age, occupation FROM users WHERE name = ?"
	row := db.QueryRow(query, name)

	var user requestData
	err := row.Scan(&user.Name, &user.Age, &user.Occupation)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Write to MySQL database
func writeToDatabase(data requestData) error {
	query := "INSERT INTO users (name, age, occupation) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE age = ?, occupation = ?"
	_, err := db.Exec(query, data.Name, data.Age, data.Occupation, data.Age, data.Occupation)
	return err
}


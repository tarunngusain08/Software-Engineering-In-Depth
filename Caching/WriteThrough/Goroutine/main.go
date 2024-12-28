package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
	_ "github.com/go-sql-driver/mysql"
)

// Struct to hold request data
type requestData struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Occupation string `json:"occupation"`
}

var cache *redis.Client
var db *sql.DB

func main() {
	http.HandleFunc("/write-through", writeThroughHandler)

	log.Println("Server started at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// Handler for the /write-through endpoint
func writeThroughHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		fmt.Println("Method not allowed")
		return
	}

	var data requestData
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		fmt.Println("Invalid request body")
		return
	}

	ctx := context.Background()

	// Write to Redis cache
	cacheErr := cache.Set(ctx, data.Name, fmt.Sprintf("%v", data), 0).Err()
	if cacheErr != nil {
		http.Error(w, "Redis error: "+cacheErr.Error(), http.StatusInternalServerError)
		fmt.Println("Redis error")
		return
	}

	// Write asynchronously to MySQL
	go func() {
		err := writeToDatabase(data)
		if err != nil {
			log.Printf("MySQL write error: %v", err)
			fmt.Println("MySQL write error")
		}
	}()

	fmt.Fprintf(w, "Data write successful!")
}

// Write to MySQL database
func writeToDatabase(data requestData) error {
	query := "INSERT INTO users (name, age, occupation) VALUES (?, ?, ?)"
	_, err := db.Exec(query, data.Name, data.Age, data.Occupation)
	return err
}

// Initialize Redis and MySQL
func init() {
	// Initialize Redis client
	cache = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", 
		DB:       0,
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

	log.Println("Initialized Redis and MySQL")
}
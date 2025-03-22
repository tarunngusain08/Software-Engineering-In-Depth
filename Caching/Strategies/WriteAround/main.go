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

func init() {
	// Initialize Redis client
	cache = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Update with your Redis address
	})

	// Initialize MySQL connection
	var err error
	db, err = sql.Open("mysql", "root:1234@tcp(localhost:3306)/users")
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	log.Println("Mysql and Redis client init!")
}

func main() {
	http.HandleFunc("/write-around", writeAroundHandler)
	http.HandleFunc("/read", readHandler)
	log.Println("Server started at :8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// Handler for the /write-around endpoint
func writeAroundHandler(w http.ResponseWriter, r *http.Request) {
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

	// Write only to the database
	err = writeToDatabase(data)
	if err != nil {
		http.Error(w, "Database write error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Do not update the cache here (Write-Around logic)
	fmt.Fprintf(w, "Data written to database successfully!")
}

// Read handler to demonstrate lazy cache update
func readHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing 'name' query parameter", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	// Check the cache
	val, err := cache.Get(ctx, name).Result()
	if err == nil {
		// Data found in cache
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(val))
		return
	} else if err != redis.Nil {
		// Redis error
		http.Error(w, "Redis error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Data not found in cache, fetch from database
	userData, err := readFromDatabase(name)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if userData == nil {
		http.Error(w, "Name not found", http.StatusNotFound)
		return
	}

	// Update cache lazily
	userJSON, _ := json.Marshal(userData)
	cache.Set(ctx, name, userJSON, 0)

	// Return the data
	w.Header().Set("Content-Type", "application/json")
	w.Write(userJSON)
}

// Write to MySQL database
func writeToDatabase(data requestData) error {
	query := "INSERT INTO users (name, age, occupation) VALUES (?, ?, ?)"
	_, err := db.Exec(query, data.Name, data.Age, data.Occupation)
	return err
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


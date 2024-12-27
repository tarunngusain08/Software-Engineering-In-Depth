# Write-Through Cache

This project demonstrates the implementation of a **Write-Through Cache** pattern using **Redis** as the cache and **MySQL** as the persistent storage. The application receives data via an HTTP endpoint and writes the data to both Redis and MySQL, ensuring that the cache is updated synchronously and the database is updated asynchronously.

## Features

- **Write-Through Cache**: Data is written to the cache immediately upon receiving a request.
- **Asynchronous Database Write**: After updating the cache, the data is asynchronously written to the MySQL database.
- **Redis Cache**: Utilizes Redis to store and retrieve data quickly.
- **MySQL Database**: Stores data persistently and retrieves it asynchronously when needed.

## Technologies Used

- **Go (Golang)**: Programming language used for backend development.
- **Redis**: Used as the in-memory cache.
- **MySQL**: Used as the relational database for persistent storage.
- **HTTP Server (Gin)**: Used to expose the API endpoint for data input.

## Installation

### Prerequisites

- Docker
- Go (Golang)
- Redis running locally or in Docker
- MySQL running locally or in Docker

### Steps to Run

1. **Clone the repository:**

   ```bash
   git clone https://github.com/tarunngusain08/WriteThrough.git
   cd WriteThrough
   ```

2. **Start Redis and MySQL containers using Docker:**

   If you have Docker installed, run the following commands to start Redis and MySQL containers.

   ```bash
   docker run --name redis-container -p 6379:6379 -d redis
   docker run --name mysql-container -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -d mysql:latest
   ```

3. **Configure MySQL:**

   - Enter the MySQL container:

     ```bash
     docker exec -it mysql-container mysql -uroot -p
     ```

   - Inside MySQL, create a database and a table:

     ```sql
     CREATE DATABASE users;
     USE users;
     CREATE TABLE users (
         id INT AUTO_INCREMENT PRIMARY KEY,
         name VARCHAR(255) NOT NULL,
         age INT NOT NULL,
         occupation VARCHAR(255) NOT NULL
     );
     ```
<img width="1512" alt="Screenshot 2024-12-27 at 7 32 15 PM" src="https://github.com/user-attachments/assets/5d36a6ea-97b0-4622-99fa-3792de4d9729" />


4. **Run the Go application:**

   ```bash
   go run main.go
   ```

   The application will start a web server on port `8080`.

6. **Testing the API:**

   You can now send a POST request to the `/write-through` endpoint to store data.

   Example request:

   ```bash
   curl -X POST http://localhost:8080/write-through \
     -H "Content-Type: application/json" \
     -d '{"name": "John Doe", "age": 30, "occupation": "Engineer"}'
   ```

   The data will be stored in both Redis and MySQL.

## API Endpoint

### `POST /write-through`

- **Request Body:**

  ```json
  {
    "name": "string",
    "age": "integer",
    "occupation": "string"
  }
  ```

- **Response:**

  ```json
  {
    "message": "Data write successful!"
  }
  ```

## Screenshots

###

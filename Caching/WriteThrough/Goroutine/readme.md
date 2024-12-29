# Read/Write-Through Cache using Goroutine

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

- **Request Body:**
  
  ```
  curl --location 'http://localhost:8080/read-through' --header 'Content-Type: application/json' --data '{"name": "John Doe"}'
  ```

- **Response:**

  ```
  {Tarunn Gusain 25 SWE}
  ```

## Screenshots

### Fail cases -
<img width="1512" alt="fail2" src="https://github.com/user-attachments/assets/a027bc9b-ac7a-4d22-9b88-4e9c7c7481f8" />
<img width="1512" alt="fail1" src="https://github.com/user-attachments/assets/c5718585-add3-4e0d-aa9e-b50b2a435898" />

### Success cases -
<img width="1512" alt="success1" src="https://github.com/user-attachments/assets/ec7b7e73-a832-42db-a2a9-d6041f85ded0" />
<img width="1512" alt="success2" src="https://github.com/user-attachments/assets/dec6f722-f8c0-4b24-ba7d-a6e20e96b0de" />

### Redis state - 
<img width="690" alt="redis success" src="https://github.com/user-attachments/assets/cf7d707d-cc92-431d-a9fa-f36088158467" />
<img width="1512" alt="Screenshot 2024-12-27 at 8 25 28 PM" src="https://github.com/user-attachments/assets/d5886f1e-3f0b-47db-8cb9-c56a7be06af9" />

### Mysql state - 
<img width="576" alt="mysql success" src="https://github.com/user-attachments/assets/5a3a7bb2-7570-4f69-b4fb-e8436756abeb" />
<img width="1512" alt="mysql" src="https://github.com/user-attachments/assets/88847666-effc-4c78-85c8-29524cd11bbd" />

### Data fetch from redis - 
<img width="668" alt="Screenshot 2024-12-29 at 6 38 06 PM" src="https://github.com/user-attachments/assets/e719f7fe-b6b8-4417-b090-9bd38b87fe25" />

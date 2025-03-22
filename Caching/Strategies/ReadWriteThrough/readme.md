# Write-Through and Read-Through Caching

This project demonstrates a simple implementation of **Write-Through** and **Read-Through** caching strategies using Golang, Redis, and MySQL.

---

## Features

1. **Write-Through Caching**:
   - Data is written synchronously to both the cache (Redis) and the database (MySQL).
   - Ensures cache consistency with the database.

2. **Read-Through Caching**:
   - Data is read from the cache first (Redis).
   - On a cache miss, the data is fetched from the database and updated in the cache.

---

## Requirements

- Go (Golang) 1.19 or higher
- MySQL database
- Redis server

---

## Setup

### Step 1: Clone the Repository
```bash
git clone https://github.com/tarunngusain08/Software-Engineering-In-Depth
cd Software-Engineering-In-Depth/Caching/ReadWriteThrough
go run main.go
```

### Step 2: **Start Redis and MySQL containers using Docker:**

1. If you have Docker installed, run the following commands to start Redis and MySQL containers.

   ```bash
   docker run --name redis-container -p 6379:6379 -d redis
   docker run --name mysql-container -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -d mysql:latest
   ```

2. **Configure MySQL:**

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

3. Update the MySQL DSN in the code:
   - Replace `"root:1234@tcp(localhost:3306)/users"` with your MySQL credentials.

### Step 3: Run the Application

1. Install dependencies:
```bash
go mod tidy
```

2. Run the server:
```bash
go run main.go
```

3. The server will start at `http://localhost:8080`.

---

## API Endpoints

### 1. Write-Through Caching

**Endpoint:** `/write-through`  
**Method:** `POST`

**Request Body:**
```json
{
  "name": "John Doe",
  "age": 30,
  "occupation": "Engineer"
}
```

**Response:**
```
Data write successful!
```

### 2. Read-Through Caching

**Endpoint:** `/read-through`  
**Method:** `POST`

**Request Body:**
```json
{
  "name": "John Doe"
}
```

**Response:** (From Cache or Database)
```json
{
  "name": "John Doe",
  "age": 30,
  "occupation": "Engineer"
}
```

---

## How It Works

1. **Write-Through**:
   - Data is synchronously written to both Redis and MySQL.
   - Ensures the cache remains up-to-date.

2. **Read-Through**:
   - The system attempts to read data from Redis first.
   - If data is not in the cache (cache miss), it fetches the data from MySQL and updates the cache.

---

## Example Usage

1. Write data using the `/write-through` endpoint.
2. Read the same data using the `/read-through` endpoint. On the first read, the system fetches the data from the database and caches it. Subsequent reads will fetch data from the cache.

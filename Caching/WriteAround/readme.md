# Write-Around Caching Strategy

This project demonstrates the **Write-Around** caching strategy using Golang, Redis, and MySQL.

---

## Features

1. **Write-Around Pattern**:
   - Writes are directly handled by the database (MySQL), without updating the cache.
   - Reads check the cache (Redis) first. On a cache miss, data is fetched from the database and lazily added to the cache.

---

## Requirements

- Go (Golang) 1.19 or higher
- MySQL database
- Redis server
- Docker (optional, for containerized Redis/MySQL setup)

---

## Setup

### Step 1: Clone the Repository
```bash
git clone https://github.com/tarunngusain08/Software-Engineering-In-Depth
cd Software-Engineering-In-Depth/Caching/WriteAround
```

### Step 2: Start Redis and MySQL Containers (Optional)

If Docker is available, run:
```bash
docker run --name redis-container -p 6379:6379 -d redis
docker run --name mysql-container -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -d mysql:latest
```

### Step 3: Configure MySQL

1. Enter the MySQL container (if using Docker):
   ```bash
   docker exec -it mysql-container mysql -uroot -p
   ```

2. Create the database and table:
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
   Replace `"root:1234@tcp(localhost:3306)/users"` with your actual MySQL credentials.

### Step 4: Install Dependencies
```bash
go mod tidy
```

### Step 5: Run the Application
```bash
go run main.go
```

The server will start at `http://localhost:8081`.

---

## API Endpoints

### 1. Write (Write-Around Pattern)

**Endpoint:** `/write-around`  
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
Data written to database successfully!
```

- Data is written directly to MySQL.
- Cache is not updated as per the Write-Around pattern.

---

### 2. Read (Lazy Cache Update)

**Endpoint:** `/read`  
**Method:** `GET`

**Query Parameter:**
```text
name=John Doe
```

**Response:** (From Cache or Database)
```json
{
  "name": "John Doe",
  "age": 30,
  "occupation": "Engineer"
}
```

- If the data exists in Redis, it is returned directly.
- On a cache miss, data is retrieved from MySQL, cached in Redis, and then returned.

---

## How It Works

1. **Write Operations:**
   - Writes bypass Redis and update MySQL directly.

2. **Read Operations:**
   - The system checks Redis for the requested data.
   - On a cache miss, it fetches data from MySQL, updates Redis lazily, and returns the data.

---

## Example Usage

1. Write data using the `/write-around` endpoint.
2. Read the same data using the `/read` endpoint:
   - On the first read, data is fetched from MySQL and added to Redis.
   - Subsequent reads fetch the data directly from Redis.

---

## Notes

- This implementation demonstrates the lazy update approach in caching.
- Suitable for scenarios where write operations are infrequent, and cache consistency is not critical.
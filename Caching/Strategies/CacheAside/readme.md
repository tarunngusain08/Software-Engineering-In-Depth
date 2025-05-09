
# Cache-Aside Caching Strategy

This project demonstrates the **Cache-Aside** caching strategy implementation using Golang, Redis, and MySQL.

---

## Features

1. **Cache-Aside Pattern**:
   - Data is read from the cache (Redis) first.
   - On a cache miss, data is retrieved from the database (MySQL) and added to the cache with a TTL.
   - For writes, data is directly updated in the database, and the corresponding cache entry is invalidated to maintain consistency.

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
cd Software-Engineering-In-Depth/Caching/CacheAside
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

### 1. Read (Cache-Aside Pattern)

**Endpoint:** `/read-cache-aside`  
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

- If the data is found in Redis, it is returned directly.
- If a cache miss occurs, the data is retrieved from MySQL, stored in Redis, and then returned.

---

### 2. Write (Cache-Aside Pattern)

**Endpoint:** `/write-cache-aside`  
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
Data written successfully!
```

- Data is written directly to MySQL.
- Corresponding cache entry is invalidated to maintain consistency.

---

## How It Works

1. **Cache Reads:**
   - The system checks Redis for the requested data.
   - On a cache miss, it fetches data from MySQL, caches it in Redis with a 5-minute TTL, and returns it.

2. **Cache Writes:**
   - Data is updated in MySQL.
   - The corresponding cache key in Redis is invalidated to ensure consistency between the cache and database.

---

## Example Usage

1. Write data using the `/write-cache-aside` endpoint.
2. Read the same data using the `/read-cache-aside` endpoint:
   - On the first read, data is fetched from MySQL and cached in Redis.
   - Subsequent reads fetch the data from Redis.


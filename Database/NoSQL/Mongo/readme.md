# MongoDB Metrics Collection System

## Overview

This document outlines the implementation of a high-performance metrics collection system using MongoDB, focusing on:
- Connection pooling optimization
- Bulk upsert operations
- Index strategies
- Performance considerations

## Features

### Connection Pool Configuration
- Maximum pool size: 150 connections
- Minimum pool size: 25 connections
- Connection idle timeout: 1 hour

### Data Structure
- Compound indexes on:
  - source_ip
  - dest_ip
  - source_port
  - dest_port
  - protocol
  - timestamp

### Performance Optimizations
1. **Bulk Write Operations**
   - Unordered bulk writes for better performance
   - Upsert capability for atomic operations

2. **Index Strategy**
   - Compound index on all key fields
   - Supports efficient upsert operations
   - Ensures data uniqueness

3. **Connection Management**
   - Efficient connection pooling
   - Automatic connection cleanup
   - Connection health monitoring

## Usage Example

```go
agent, err := NewMetricsAgent(
    "mongodb://localhost:27017",
    "metrics_db",
    "metrics",
)
if err != nil {
    log.Fatal(err)
}
defer agent.Close(context.Background())

// Create indexes
if err := agent.CreateMetricsIndexes(); err != nil {
    log.Fatal(err)
}

// Batch upsert metrics
metrics := []Metric{
    {
        SourceIP:   "192.168.1.1",
        DestIP:     "192.168.1.2",
        SourcePort: 8080,
        DestPort:   80,
        Protocol:   "TCP",
        Timestamp:  time.Now(),
        Value:      123.45,
    },
}

if err := agent.BatchUpsertMetrics(metrics); err != nil {
    log.Fatal(err)
}
```

## Best Practices

1. **Error Handling**
   - Always check for errors on database operations
   - Use context with timeouts for operations
   - Implement proper cleanup on errors

2. **Resource Management**
   - Close connections when done
   - Monitor connection pool metrics
   - Implement proper error recovery

3. **Performance Tips**
   - Use bulk operations for multiple documents
   - Implement proper indexes
   - Monitor query performance

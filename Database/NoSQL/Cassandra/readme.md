# Cassandra Metrics Collection System

## Overview

This document outlines the implementation of a high-performance metrics collection system using Apache Cassandra, focusing on:
- Efficient data modeling
- Batch operations
- Query optimization
- Connection management

## Features

### Connection Configuration
- Token-aware routing with round-robin fallback
- Exponential backoff retry policy
- Configurable consistency levels
- Multiple connections per host

### Data Model
Optimized for Cassandra's distributed nature:
```cql
CREATE TABLE metrics (
    source_ip text,
    dest_ip text,
    source_port int,
    dest_port int,
    protocol text,
    timestamp timestamp,
    value double,
    PRIMARY KEY ((source_ip, dest_ip, protocol), timestamp, source_port, dest_port)
) WITH CLUSTERING ORDER BY (timestamp DESC, source_port ASC, dest_port ASC)
```

#### Key Design Decisions:
1. **Partition Key**: `(source_ip, dest_ip, protocol)`
   - Ensures even distribution
   - Enables efficient queries by IP pairs

2. **Clustering Columns**: `timestamp, source_port, dest_port`
   - Timestamp-based ordering
   - Efficient range queries
   - Port numbers for uniqueness

### Performance Optimizations

1. **Batch Operations**
   - Unlogged batches for better performance
   - Atomic operations within partition key

2. **Query Patterns**
   - Efficient time-range queries
   - Partition key optimization
   - Minimal cross-partition queries

3. **Connection Management**
   - Connection pooling per host
   - Load balancing
   - Automatic failover

## Usage Example

```go
hosts := []string{"localhost:9042"}
agent, err := NewMetricsAgent(hosts, "metrics_keyspace")
if err != nil {
    log.Fatal(err)
}
defer agent.Close()

// Create table
if err := agent.CreateMetricsTable(); err != nil {
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

1. **Data Modeling**
   - Design for query patterns
   - Avoid large partitions
   - Use appropriate compound keys

2. **Query Optimization**
   - Always include partition key
   - Limit cross-partition queries
   - Use prepared statements

3. **Resource Management**
   - Configure appropriate timeouts
   - Implement retry policies
   - Monitor cluster health

4. **Performance Tips**
   - Use unlogged batches
   - Configure consistency appropriately
   - Monitor partition sizes

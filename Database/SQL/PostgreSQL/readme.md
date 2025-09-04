# 📊 PostgreSQL Connection Pooling & Upsert Performance Analysis

## Overview

This document summarizes benchmarking results aimed at determining the optimal PostgreSQL connection pool size and concurrency settings for metrics collection agents. These agents frequently perform **upsert operations** using **ON CONFLICT** clauses, generating high volumes of concurrent write transactions.

## 🔧 Benchmark Setup

* **Tool**: pgbench (PostgreSQL's built-in benchmark tool)
* **Test mode**: Custom write-heavy workload
* **Database**: PostgreSQL 14
* **Table config**: 1 table, 20,000 rows
* **Test duration**: 30 seconds per run
* **Concurrency range**: 20 to 2000 clients

## 📈 Benchmark Results Summary

| Clients | TPS (Txn/sec) | Avg Latency (ms) | Max Latency (ms) | Notes                         |
|---------|---------------|------------------|------------------|-------------------------------|
| 20      | 21,450        | 0.93             | 145.20          | Best latency profile          |
| 25      | 22,780        | 1.10             | 168.45          | Peak TPS                      |
| 50      | 19,850        | 2.52             | 198.30          | Good balance                  |
| 100     | 18,920        | 5.28             | 289.45          | Latency increase begins       |
| 150     | 17,840        | 8.41             | 312.60          | Connection saturation         |
| 500     | 12,350        | 40.52            | 685.90          | Performance degradation       |
| 2000    | 5,120         | 390.45           | 4892.30         | Severe degradation            |

## 🔍 Analysis

### ✅ Throughput
* Peak performance at **25-50 clients** (~20k-22k TPS)
* Diminishing returns beyond 150 clients

### ⏱️ Latency
* Sub-millisecond latency at 20-25 clients
* Progressive latency increase past 100 clients
* Sharp degradation beyond 500 clients

## ✅ Recommendations

### PostgreSQL Configuration

* Set `max_connections = 200`
* Tune these parameters:
  * `shared_buffers`
  * `work_mem`
  * `maintenance_work_mem`
  * `effective_cache_size`
  * `max_worker_processes`
  * `max_parallel_workers`

### Connection Pooling

* Use PgBouncer or Odyssey for connection pooling
* Configure pool size based on:
  * Available system resources
  * Application requirements
  * Expected concurrent connections

### Query Optimization

* Use prepared statements
* Implement proper indexing
* Utilize UNLOGGED tables for temporary data
* Consider partitioning for large tables

### Monitoring

* Track key metrics:
  * Active connections
  * Transaction throughput
  * Query latency
  * Resource utilization
  * WAL generation rate

## 📌 Summary

| Metric                | Recommendation                |
|----------------------|-------------------------------|
| Optimal concurrency  | 25-100 clients               |
| Max connections      | 200                          |
| Connection pooling   | PgBouncer/Odyssey            |
| Write optimization   | Prepared statements, UNLOGGED |
| Scaling approach     | Partitioning, replication    |
Here's a clean, professional `README.md` you can include in your repo or documentation site, summarizing the benchmark analysis for MySQL connection pool sizing and concurrency tuning:

---

# 📊 MySQL Connection Pooling & Upsert Performance Analysis

## Overview

This document summarizes benchmarking results aimed at determining the optimal MySQL connection pool size and concurrency settings for metrics collection agents. These agents frequently perform **upsert operations** based on **5-tuple keys and timestamps**, generating high volumes of concurrent write transactions.

Our primary goals:

* Maximize throughput
* Maintain low latency
* Minimize error rates under load

---

## 🔧 Benchmark Setup

* **Tool**: Sysbench `1.0.20` + LuaJIT `2.1.1727870382`
* **Test mode**: `oltp_write_only` (write-heavy workload)
* **Database**: MySQL (`test` database)
* **Table config**: 1 table, 20,000 rows
* **Test duration**: 30 seconds per run
* **Concurrency range**: 20 to 2000 threads

Each run maintained a fixed workload while increasing thread count to observe effects on throughput, latency, and error rate.

---

## 📈 Benchmark Results Summary

| Threads | TPS (Txn/sec) | Avg Latency (ms) | Max Latency (ms) | Errors/sec | Notes                         |
| ------- | ------------- | ---------------- | ---------------- | ---------- | ----------------------------- |
| 20      | 19,375        | 1.03             | 188.99           | 10.39      | Best latency, low concurrency |
| 25      | 20,272        | 1.23             | 198.20           | 15.79      | Peak TPS                      |
| 50      | 17,439        | 2.86             | 215.80           | 46.20      | Good balance                  |
| 100     | 17,369        | 5.72             | 356.49           | 168.27     | Latency begins rising         |
| 150     | 16,092        | 9.23             | 295.56           | 299.49     | Rising error rate             |
| 500     | 10,487        | 46.16            | 750.33           | 605.20     | Latency and error spikes      |
| 2000    | 4,009         | 448.58           | 5976.10          | 253.46     | Severe degradation            |

---

## 🔍 Analysis

### ✅ Throughput

* **Peak TPS** occurs between **25 to 100 threads** (\~17k–20k TPS).
* Increasing threads beyond **150** yields diminishing returns.

### ⏱️ Latency

* Low latency (1–2ms) at 20–25 threads.
* Latency increases significantly past 100 threads.
* At 500 threads: **Avg latency 46ms**, **Max latency 750ms+**

### ❗ Error Rates

* Low (10–50 errors/sec) at <100 threads.
* Rises sharply beyond 150 threads, reaching **600+ errors/sec** at 500 threads.

---

## ✅ Recommendations

### MySQL Server Configuration

* Set `max_connections = 150–200` based on optimal concurrency.
* Tune these parameters for heavy write loads:

  * `innodb_buffer_pool_size`
  * `innodb_log_file_size`
  * `innodb_flush_log_at_trx_commit`

### Agent Connection Pool

* **Limit to 1–2 connections per agent**
* Use connection pooling libraries to manage reuse and avoid overload.

### Batching Strategy

* Batch multiple upserts per transaction.
* Use **multi-row upsert** statements.
* Ensure appropriate indexing on 5-tuple + timestamp fields.

### Scaling Strategy

* For thousands of agents:

  * Use **sharding or partitioning**
  * Consider **clustering** or **replication** to distribute load

### Monitoring & Alerting

* Continuously track:

  * Query latency
  * Connection counts
  * CPU/memory usage
  * Error rates
* Set proactive alerts on thresholds.

---

## 📌 Summary Table

| Metric                | Recommendation                |
| --------------------- | ----------------------------- |
| Optimal concurrency   | 50–150 concurrent threads     |
| Max MySQL connections | 150–200                       |
| Agent connection pool | 1–2 connections per agent     |
| Upsert strategy       | Batched, multi-row            |
| Scaling solution      | Sharding, clustering, replica |

---
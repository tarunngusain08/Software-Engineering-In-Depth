package postgresql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type MetricsAgent struct {
	pool *sql.DB
}

type Metric struct {
	SourceIP   string
	DestIP     string
	SourcePort int
	DestPort   int
	Protocol   string
	Timestamp  time.Time
	Value      float64
}

func NewMetricsAgent(connStr string) (*MetricsAgent, error) {
	// Initialize connection pool
	pool, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %v", err)
	}

	// Configure pool settings based on benchmarks
	pool.SetMaxOpenConns(200) // Optimal from benchmarks
	pool.SetMaxIdleConns(50)  // Good balance for idle connections
	pool.SetConnMaxLifetime(time.Hour)

	return &MetricsAgent{
		pool: pool,
	}, nil
}

func (ma *MetricsAgent) CreateMetricsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS metrics (
			source_ip VARCHAR(45),
			dest_ip VARCHAR(45),
			source_port INTEGER,
			dest_port INTEGER,
			protocol VARCHAR(10),
			timestamp TIMESTAMP,
			value DOUBLE PRECISION,
			PRIMARY KEY (source_ip, dest_ip, source_port, dest_port, protocol, timestamp)
		)`

	_, err := ma.pool.Exec(query)
	return err
}

func (ma *MetricsAgent) BatchUpsertMetrics(metrics []Metric) error {
	// Using UNLOGGED table for better performance as mentioned in readme
	query := `
		INSERT INTO metrics (
			source_ip, dest_ip, source_port, dest_port, protocol, timestamp, value
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (source_ip, dest_ip, source_port, dest_port, protocol, timestamp)
		DO UPDATE SET value = EXCLUDED.value`

	tx, err := ma.pool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute batch upserts
	for _, m := range metrics {
		_, err = stmt.Exec(
			m.SourceIP, m.DestIP, m.SourcePort, m.DestPort,
			m.Protocol, m.Timestamp, m.Value,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (ma *MetricsAgent) Close() error {
	return ma.pool.Close()
}

func Example() {
	// Initialize metrics agent with PgBouncer connection string
	agent, err := NewMetricsAgent("postgres://user:password@localhost:6432/metrics_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer agent.Close()

	// Create metrics table
	if err := agent.CreateMetricsTable(); err != nil {
		log.Fatal(err)
	}

	// Sample metrics data
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
		// Add more metrics...
	}

	// Perform batch upsert
	if err := agent.BatchUpsertMetrics(metrics); err != nil {
		log.Fatal(err)
	}
}

package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MetricsAgent struct {
	db   *sql.DB
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

func NewMetricsAgent(dsn string) (*MetricsAgent, error) {
	// Initialize connection pool
	pool, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %v", err)
	}

	// Configure pool settings based on benchmarks
	pool.SetMaxOpenConns(150) // Optimal concurrency range
	pool.SetMaxIdleConns(25)  // Good balance for idle connections
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
			source_port INT,
			dest_port INT,
			protocol VARCHAR(10),
			timestamp DATETIME,
			value DOUBLE,
			PRIMARY KEY (source_ip, dest_ip, source_port, dest_port, protocol, timestamp)
		)`

	_, err := ma.pool.Exec(query)
	return err
}

func (ma *MetricsAgent) BatchUpsertMetrics(metrics []Metric) error {
	// Prepare multi-row upsert statement
	query := `
		INSERT INTO metrics (
			source_ip, dest_ip, source_port, dest_port, protocol, timestamp, value
		) VALUES (?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE value = VALUES(value)`

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
	// Initialize metrics agent
	agent, err := NewMetricsAgent("user:password@tcp(localhost:3306)/metrics_db")
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

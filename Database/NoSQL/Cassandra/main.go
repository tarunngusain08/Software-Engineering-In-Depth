package cassandra

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

type MetricsAgent struct {
	session *gocql.Session
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

func NewMetricsAgent(hosts []string, keyspace string) (*MetricsAgent, error) {
	// Initialize cluster config
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum // Strong consistency for metrics
	cluster.NumConns = 2               // Connections per host
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())

	// Configure timeouts and retry policy
	cluster.Timeout = 5 * time.Second
	cluster.ConnectTimeout = 10 * time.Second
	cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{
		NumRetries: 3,
		Min:        100 * time.Millisecond,
		Max:        2 * time.Second,
	}

	// Create session
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create Cassandra session: %v", err)
	}

	return &MetricsAgent{
		session: session,
	}, nil
}

func (ma *MetricsAgent) CreateMetricsTable() error {
	// Create table with appropriate clustering order and partition key
	query := `
		CREATE TABLE IF NOT EXISTS metrics (
			source_ip text,
			dest_ip text,
			source_port int,
			dest_port int,
			protocol text,
			timestamp timestamp,
			value double,
			PRIMARY KEY ((source_ip, dest_ip, protocol), timestamp, source_port, dest_port)
		) WITH CLUSTERING ORDER BY (timestamp DESC, source_port ASC, dest_port ASC)
	`

	if err := ma.session.Query(query).Exec(); err != nil {
		return fmt.Errorf("failed to create metrics table: %v", err)
	}

	return nil
}

func (ma *MetricsAgent) BatchUpsertMetrics(metrics []Metric) error {
	// Create a batch operation
	batch := ma.session.NewBatch(gocql.UnloggedBatch)

	// Prepare batch query
	query := `
		INSERT INTO metrics (
			source_ip, dest_ip, source_port, dest_port,
			protocol, timestamp, value
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	// Add all metrics to the batch
	for _, m := range metrics {
		batch.Query(query,
			m.SourceIP, m.DestIP, m.SourcePort, m.DestPort,
			m.Protocol, m.Timestamp, m.Value,
		)
	}

	// Execute the batch
	if err := ma.session.ExecuteBatch(batch); err != nil {
		return fmt.Errorf("failed to execute batch upsert: %v", err)
	}

	return nil
}

func (ma *MetricsAgent) QueryMetricsByTimeRange(
	sourceIP, destIP, protocol string,
	startTime, endTime time.Time,
) ([]Metric, error) {
	var metrics []Metric

	// Query metrics within time range for specific partition
	query := `
		SELECT source_ip, dest_ip, source_port, dest_port,
			   protocol, timestamp, value
		FROM metrics
		WHERE source_ip = ?
		  AND dest_ip = ?
		  AND protocol = ?
		  AND timestamp >= ?
		  AND timestamp <= ?
	`

	scanner := ma.session.Query(query,
		sourceIP, destIP, protocol, startTime, endTime,
	).Iter().Scanner()

	for scanner.Next() {
		var m Metric
		err := scanner.Scan(
			&m.SourceIP, &m.DestIP, &m.SourcePort, &m.DestPort,
			&m.Protocol, &m.Timestamp, &m.Value,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan metric: %v", err)
		}
		metrics = append(metrics, m)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error during metrics scan: %v", err)
	}

	return metrics, nil
}

func (ma *MetricsAgent) Close() {
	ma.session.Close()
}

func Example() {
	// Initialize metrics agent
	hosts := []string{"localhost:9042"}
	agent, err := NewMetricsAgent(hosts, "metrics_keyspace")
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

	// Query metrics for the last hour
	endTime := time.Now()
	startTime := endTime.Add(-1 * time.Hour)

	results, err := agent.QueryMetricsByTimeRange(
		"192.168.1.1",
		"192.168.1.2",
		"TCP",
		startTime,
		endTime,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Process results
	for _, metric := range results {
		fmt.Printf("Metric at %v: %v\n", metric.Timestamp, metric.Value)
	}
}

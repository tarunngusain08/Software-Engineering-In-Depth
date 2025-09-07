package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MetricsAgent struct {
	client     *mongo.Client
	collection *mongo.Collection
}

type Metric struct {
	SourceIP   string    `bson:"source_ip"`
	DestIP     string    `bson:"dest_ip"`
	SourcePort int       `bson:"source_port"`
	DestPort   int       `bson:"dest_port"`
	Protocol   string    `bson:"protocol"`
	Timestamp  time.Time `bson:"timestamp"`
	Value      float64   `bson:"value"`
}

func NewMetricsAgent(uri string, database, collection string) (*MetricsAgent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize connection with MongoDB
	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(150). // Optimal connection pool size
		SetMinPoolSize(25).  // Minimum connections to maintain
		SetMaxConnIdleTime(time.Hour)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %v", err)
	}

	// Ping the database to verify connection
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Get collection reference
	coll := client.Database(database).Collection(collection)

	return &MetricsAgent{
		client:     client,
		collection: coll,
	}, nil
}

func (ma *MetricsAgent) CreateMetricsIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a compound index on the fields that form our unique key
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "source_ip", Value: 1},
			{Key: "dest_ip", Value: 1},
			{Key: "source_port", Value: 1},
			{Key: "dest_port", Value: 1},
			{Key: "protocol", Value: 1},
			{Key: "timestamp", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := ma.collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %v", err)
	}

	return nil
}

func (ma *MetricsAgent) BatchUpsertMetrics(metrics []Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Prepare bulk write operations
	operations := make([]mongo.WriteModel, len(metrics))
	for i, metric := range metrics {
		filter := bson.D{
			{Key: "source_ip", Value: metric.SourceIP},
			{Key: "dest_ip", Value: metric.DestIP},
			{Key: "source_port", Value: metric.SourcePort},
			{Key: "dest_port", Value: metric.DestPort},
			{Key: "protocol", Value: metric.Protocol},
			{Key: "timestamp", Value: metric.Timestamp},
		}
		update := bson.D{
			{Key: "$set", Value: metric},
		}
		operations[i] = mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)
	}

	// Execute bulk write with ordered=false for better performance
	opts := options.BulkWrite().SetOrdered(false)
	_, err := ma.collection.BulkWrite(ctx, operations, opts)
	if err != nil {
		return fmt.Errorf("failed to perform bulk upsert: %v", err)
	}

	return nil
}

func (ma *MetricsAgent) Close(ctx context.Context) error {
	return ma.client.Disconnect(ctx)
}

func Example() {
	// Initialize metrics agent
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

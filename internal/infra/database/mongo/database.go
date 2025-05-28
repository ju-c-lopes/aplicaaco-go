package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database interface defines the methods that a MongoDB database should implement
type Database interface {
	Collection(name string) *mongo.Collection
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Client() *mongo.Client
}

// MongoDatabase implements the Database interface
type MongoDatabase struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewMongoDatabase creates a new MongoDatabase instance
func NewMongoDatabase(uri, dbName string) (*MongoDatabase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &MongoDatabase{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

// Collection returns a handle to a MongoDB collection
func (m *MongoDatabase) Collection(name string) *mongo.Collection {
	return m.database.Collection(name)
}

// Connect establishes a connection to the MongoDB server
func (m *MongoDatabase) Connect(ctx context.Context) error {
	return m.client.Connect(ctx)
}

// Disconnect closes the connection to the MongoDB server
func (m *MongoDatabase) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

// Client returns the MongoDB client
func (m *MongoDatabase) Client() *mongo.Client {
	return m.client
} 
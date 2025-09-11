package storage

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocumentStore struct {
	client *mongo.Client
	db     *mongo.Database
	mu     sync.RWMutex
}

func newMongoClient(mongoURI string) *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: mongoURI = %s\n", mongoURI)
		panic(err)
	}
	return client
}
func NewDocumentStore(dbName string) *DocumentStore {
	client := newMongoClient(cfg.MongoURI)
	db := client.Database(dbName)
	return &DocumentStore{
		client: client,
		db:     db,
	}
}

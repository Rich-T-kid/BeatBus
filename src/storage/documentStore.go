package storage

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
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

func (ds *DocumentStore) InsertNewUser(username, hashedPassword string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	collection := ds.db.Collection("users")
	ctx := context.Background()

	// Check if username already exists
	count, err := collection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("username already taken")
	}

	// Insert new user
	_, err = collection.InsertOne(ctx, bson.M{
		"username": username,
		"password": hashedPassword,
	})
	return err
}
func (ds *DocumentStore) ValidateUser(username, hashedPassword string) error {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	collection := ds.db.Collection("users")
	ctx := context.Background()

	var result bson.M
	err := collection.FindOne(ctx, bson.M{"username": username, "password": hashedPassword}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return fmt.Errorf("user not found with provided username and password")
	} else if err != nil {
		return err // Some other error
	}

	return nil // Valid credentials
}

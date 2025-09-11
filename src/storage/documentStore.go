package storage

import (
	"BeatBus/internal"
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrUserNameTaken                      = fmt.Errorf("username already taken")
	CannotCreateRoomAlreadyInSessionError = fmt.Errorf("cannot create room, user already hosting a session")
)

const (
	MongoDBName        = "BeatBus"
	UsersCollection    = "users"
	UserInfoCollection = "usersInfo"
	RoomsCollection    = "rooms"
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
func NewDocumentStore() *DocumentStore {
	client := newMongoClient(cfg.MongoURI)
	db := client.Database(MongoDBName)
	return &DocumentStore{
		client: client,
		db:     db,
	}
}

func (ds *DocumentStore) InsertNewUser(username, hashedPassword string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	collection := ds.db.Collection(UsersCollection)
	ctx := context.Background()

	// Check if username already exists
	count, err := collection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrUserNameTaken
	}

	// Insert new user
	res, err := collection.InsertOne(ctx, bson.M{
		"username": username,
		"password": hashedPassword,
	})
	if err != nil {
		fmt.Println("Error inserting new user:", err)
		return err
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	err = ds.insertNewUserInfo(id) //insert user info in a separate go routine
	return err
}
func (ds *DocumentStore) insertNewUserInfo(id string) error {
	collection := ds.db.Collection(UserInfoCollection)
	ctx := context.Background()

	//dont need to check if it exists because of foreign key relationship with users collection
	_, err := collection.InsertOne(ctx, bson.M{
		"user_id":           id,
		"inSession":         false,
		"join_date":         time.Now(),
		"previous_sessions": []string{},
		"listened_to":       []string{},
		"liked_songs":       []string{},
		"disliked_songs":    []string{},
	})
	return err
}

func (ds *DocumentStore) ValidateUser(username, hashedPassword string) error {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	collection := ds.db.Collection(UsersCollection)
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
func (ds *DocumentStore) inSession(username string) bool {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	collection := ds.db.Collection(UsersCollection)
	ctx := context.Background()

	var user bson.M
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return false // User not found or some error occurred
	}

	userID := user["_id"].(primitive.ObjectID).Hex()
	fmt.Println("userID:", userID)
	var userInfo bson.M
	infoCollection := ds.db.Collection("usersInfo")
	err = infoCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&userInfo)
	if err != nil {
		return false // User info not found or some error occurred
	}

	inSession, ok := userInfo["inSession"].(bool)
	if !ok {
		return false // inSession field missing or of wrong type
	}

	return inSession
}

// need to lock b4 calling this but this wont do it by default
func (ds *DocumentStore) setInSession(username string, inSession bool) error {
	ctx := context.Background()
	coll := ds.db.Collection(UsersCollection)

	// Quick visibility checks
	var user bson.M
	err := coll.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return err
	}
	fmt.Println("No errors finding user -> ", user)
	coll = ds.db.Collection(UserInfoCollection)
	err = coll.FindOneAndUpdate(ctx, bson.M{"user_id": user["_id"].(primitive.ObjectID).Hex()}, bson.M{"$set": bson.M{"inSession": inSession}}).Err()
	return err
}
func (ds *DocumentStore) CreateRoom(hostUsername, roomName string, lifetime, maxUsers uint, public bool) (map[string]interface{}, error) {
	fmt.Printf(
		"CreateRoom called with hostUsername=%s, roomName=%s, lifetime=%d, maxUsers=%d, public=%t\n",
		hostUsername, roomName, lifetime, maxUsers, public,
	)
	if ds.inSession(hostUsername) {
		fmt.Printf("user %s is already in a session, cannot create room\n", hostUsername)
		return nil, CannotCreateRoomAlreadyInSessionError
	}

	fmt.Printf("user %s is not in a session, proceeding to create room\n", hostUsername)
	err := ds.setInSession(hostUsername, true)
	if err != nil {
		fmt.Printf("failed to set user %s inSession to true: %v\n", hostUsername, err)
		return nil, err
	}
	roomsCollection := ds.db.Collection(RoomsCollection)
	ctx := context.Background()
	roomID := randomHash()
	roomPassword := randomHash()
	token := internal.NewJWTHandler().CreateToken(hostUsername, roomID, time.Duration(lifetime)*time.Minute)
	_, err = roomsCollection.InsertOne(ctx, bson.M{
		"roomID":       roomID,
		"hostID":       hostUsername,
		"accessToken":  token,
		"CurrentQueue": []string{},
		"playedSongs":  []string{},
		"usersJoined":  []string{hostUsername},
		"RoomStats": bson.M{
			"name":         roomName,
			"lifetime":     lifetime,
			"maxUsers":     maxUsers,
			"public":       public,
			"created":      time.Now(),
			"roomPassword": roomPassword,
		},
	})
	if err != nil {
		fmt.Printf("Error creating room: %v\n", err)
		return nil, err
	}

	return map[string]interface{}{
		"roomProps": map[string]interface{}{
			"roomID":       roomID,
			"roomPassword": roomPassword,
			"hostID":       hostUsername,
			"roomName":     roomName,
			"maxUsers":     maxUsers,
			"isPublic":     public,
		},
		"accessToken": map[string]interface{}{
			"token":     token,
			"expiresIn": time.Now().Add(time.Duration(lifetime) * time.Minute).Unix(),
		},
		"timeStamp": time.Now().Unix(),
	}, nil
}

type CreateRoomResult struct{}

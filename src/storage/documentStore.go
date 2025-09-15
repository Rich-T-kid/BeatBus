package storage

import (
	"BeatBus/internal"
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrUserNameTaken                    = fmt.Errorf("username already taken")
	ErrCannotCreateRoomAlreadyInSession = fmt.Errorf("cannot create room, user already hosting a session")
	ErrRoomDoesntExist                  = fmt.Errorf("room does not exist")
	ErrInvalidRoomPassword              = fmt.Errorf("invalid room password")
	ErrUserAlreadyInRoom                = fmt.Errorf("user already in room")
	ErrRoomFull                         = fmt.Errorf("room is full")
	ErrInvalidSongOperation             = func(operation string) error {
		return fmt.Errorf("[%s] is not a valid song action | Valid actions are [like, unlike, dislike, undislike]", operation)
	}
	ErrQueueIsEmpty  = fmt.Errorf("the queue is empty")
	ErrNoSongsPlayed = fmt.Errorf("no songs have been played in this room yet")
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
	logger *log.Logger
	mu     sync.RWMutex
}

var (
	mongoClient *mongo.Client
)

func newMongoClient(mongoURI string) *mongo.Client {
	if mongoClient != nil {
		return mongoClient
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: mongoURI = %s\n", mongoURI)
		panic(err)
	}
	mongoClient = client
	return client
}
func NewDocumentStore(l *log.Logger) *DocumentStore {
	client := newMongoClient(cfg.MongoURI)
	db := client.Database(MongoDBName)
	return &DocumentStore{
		client: client,
		db:     db,
		logger: l,
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
		ds.logger.Println("Error inserting new user:", err)
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
	ds.logger.Println("userID:", userID)
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

func (ds *DocumentStore) setInSession(username string, inSession bool) error {
	ctx := context.Background()
	coll := ds.db.Collection(UsersCollection)

	// Quick visibility checks
	var user bson.M
	err := coll.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return err
	}
	ds.logger.Println("No errors finding user -> ", user)
	coll = ds.db.Collection(UserInfoCollection)
	err = coll.FindOneAndUpdate(ctx, bson.M{"user_id": user["_id"].(primitive.ObjectID).Hex()}, bson.M{"$set": bson.M{"inSession": inSession}}).Err()
	return err
}
func (ds *DocumentStore) CreateRoom(hostUsername, roomName string, lifetime, maxUsers uint, public bool) (map[string]interface{}, error) {
	ds.logger.Printf(
		"CreateRoom called with hostUsername=%s, roomName=%s, lifetime=%d, maxUsers=%d, public=%t\n",
		hostUsername, roomName, lifetime, maxUsers, public,
	)
	if ds.inSession(hostUsername) {
		ds.logger.Printf("user %s is already in a session, cannot create room\n", hostUsername)
		return nil, ErrCannotCreateRoomAlreadyInSession
	}

	ds.logger.Printf("user %s is not in a session, proceeding to create room\n", hostUsername)
	err := ds.setInSession(hostUsername, true)
	if err != nil {
		ds.logger.Printf("failed to set user %s inSession to true: %v\n", hostUsername, err)
		return nil, err
	}
	roomsCollection := ds.db.Collection(RoomsCollection)
	ctx := context.Background()
	roomID := internal.RandomHash()
	roomPassword := internal.RandomHash()
	token := internal.NewJWTHandler().CreateToken(hostUsername, roomID, time.Duration(lifetime)*time.Minute)
	_, err = roomsCollection.InsertOne(ctx, bson.M{
		"roomID":       roomID,
		"hostID":       hostUsername,
		"accessToken":  token,
		"CurrentQueue": []interface{}{},
		"playedSongs":  []interface{}{},
		"songCount":    0,
		"usersJoined":  []string{hostUsername},
		"RoomStats": bson.M{
			"name":         roomName,
			"lifetime":     int64(lifetime),
			"maxUsers":     maxUsers,
			"public":       public,
			"createdAt":    time.Now(),
			"roomPassword": roomPassword,
		},
	})
	if err != nil {
		ds.logger.Printf("Error creating room: %v\n", err)
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

func (ds *DocumentStore) UpdateRoomSettings(hostUsername, roomName string, maxUsers uint, public bool) (map[string]interface{}, error) {
	userColl := ds.db.Collection(UsersCollection)
	ctx := context.Background()

	var user bson.M
	err := userColl.FindOne(ctx, bson.M{"username": hostUsername}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	roomColl := ds.db.Collection(RoomsCollection)
	var room bson.M
	err = roomColl.FindOne(ctx, bson.M{"hostID": hostUsername}).Decode(&room)
	if err != nil {
		return nil, fmt.Errorf("room not found")
	}
	_, err = roomColl.UpdateOne(ctx, bson.M{"hostID": hostUsername}, bson.M{
		"$set": bson.M{
			"RoomStats.name":     roomName,
			"RoomStats.maxUsers": maxUsers,
			"RoomStats.public":   public,
		},
	})
	if err != nil {
		return nil, err
	}
	createdAt := room["RoomStats"].(bson.M)["createdAt"].(primitive.DateTime)
	// Use createdAt here
	duration := time.Since(createdAt.Time())
	totalMinutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) % 60
	fmt.Printf("time since createdAt: %d:%02d\n", totalMinutes, seconds)
	fmt.Printf("%v\n type: %T\n", room["RoomStats"].(bson.M)["lifetime"], room["RoomStats"].(bson.M)["lifetime"])
	originalMinutes := room["RoomStats"].(bson.M)["lifetime"].(int64)
	difference := originalMinutes - int64(totalMinutes)
	timeLeft := difference
	return map[string]interface{}{
		"roomProps": map[string]interface{}{
			"roomID":       room["roomID"],
			"roomPassword": room["RoomStats"].(bson.M)["roomPassword"],
			"hostID":       hostUsername,
			"roomName":     roomName,
			"maxUsers":     maxUsers,
			"isPublic":     public,
			"timeLeft":     timeLeft,
		},
		"timeStamp": time.Now().Unix(),
	}, nil
}

func (ds *DocumentStore) DeleteRoom(accessToken, hostUsername, roomID string) (map[string]interface{}, error) {
	RoomsCollection := ds.db.Collection(RoomsCollection)
	ctx := context.Background()

	// Verify room exists and hostUsername matches
	var room bson.M
	ds.logger.Printf("Attempting to delete room with roomID=%s by hostUsername=%s\n", roomID, hostUsername)
	err := RoomsCollection.FindOne(ctx, bson.M{"roomID": roomID, "hostID": hostUsername}).Decode(&room)
	if err == mongo.ErrNoDocuments {
		return nil, ErrRoomDoesntExist
	} else if err != nil {
		return nil, err
	}

	// Verify accessToken matches
	if room["accessToken"] != accessToken {
		return nil, fmt.Errorf("invalid access token")
	}
	// (1) get the user with the most likes across all the songs,
	// (2) and the user with the most likes on any given song,
	// (3) user with most dislikes,
	// (4) user with most disliked song
	playedSongs := room["playedSongs"].(primitive.A)
	userLikeCount := make(map[string]int32)    // user:like count
	userDislikeCount := make(map[string]int32) // user:dislike count
	MostLikedSong := make(map[string]int32)    // songID:like count
	MostDislikedSong := make(map[string]int32) // songID:dislike count
	SongTable := make(map[string]interface{})  // songID:song object
	type resultStruct struct {
		MostLikedUser    string
		MostDislikedUser string
		MostLikedSong    interface{}
		MostDislikedSong interface{}
	}
	var result resultStruct
	if len(playedSongs) > 0 {
		SongEntry := playedSongs[0]
		songMap := SongEntry.(bson.M)
		song := songMap["song"].(bson.M)
		metadata := song["metadata"].(bson.M)
		//stats := song["stats"].(bson.M)

		addedBy := metadata["addedBy"].(string)
		result.MostLikedSong = addedBy
		result.MostDislikedSong = addedBy
		result.MostLikedSong = SongEntry
		result.MostDislikedSong = SongEntry
	} else {
		return nil, ErrNoSongsPlayed
	}
	for _, SongEntry := range playedSongs {
		ds.logger.Println("SongEntry -> ", SongEntry)

		songMap := SongEntry.(bson.M)
		song := songMap["song"].(bson.M)
		metadata := song["metadata"].(bson.M)
		stats := song["stats"].(bson.M)

		// Extract the values
		songID := song["songId"].(string)
		likes := metadata["likes"].(int32)
		dislikes := metadata["dislikes"].(int32)
		addedBy := metadata["addedBy"].(string)

		SongTable[songID] = SongEntry
		userLikeCount[addedBy] += likes
		userDislikeCount[addedBy] += dislikes
		MostLikedSong[songID] = likes
		MostDislikedSong[songID] = dislikes

		infoString := fmt.Sprintf("%s - %s - %s", stats["artist"], stats["title"], stats["album"])
		ds.logger.Printf("Artist Info string : %s\n", infoString)
		ds.logger.Printf("\n SongID: %s, Likes: %d, Dislikes: %d, AddedBy: %s\n", songID, likes, dislikes, addedBy)
	}
	ds.logger.Printf(" %v, %v,%v,%v,%v\n", userLikeCount, userDislikeCount, MostLikedSong, MostDislikedSong, SongTable)
	sortedUserLikes := toSlice(userLikeCount)
	sortSlice(sortedUserLikes)
	sortedUserDislikes := toSlice(userDislikeCount)
	sortSlice(sortedUserDislikes)
	mostLikedSorted := toSlice(MostLikedSong)
	sortSlice(mostLikedSorted)
	mostDislikedSorted := toSlice(MostDislikedSong)
	sortSlice(mostDislikedSorted)
	//
	result.MostLikedUser = sortedUserLikes[0].txt
	result.MostDislikedUser = sortedUserDislikes[0].txt
	result.MostLikedSong = SongTable[mostLikedSorted[0].txt]
	result.MostDislikedSong = SongTable[mostDislikedSorted[0].txt]

	err = RoomsCollection.FindOneAndDelete(ctx, bson.M{"roomID": roomID}).Err()
	if err != nil {
		return nil, err
	}
	// Set host user's inSession to false
	err = ds.setInSession(hostUsername, false)
	if err != nil {
		ds.logger.Printf("failed to set user %s inSession to false: %v\n", hostUsername, err)
		// Not returning error here because room deletion was successful
	}
	userInfoColl := ds.db.Collection(UserInfoCollection)
	err = userInfoColl.FindOneAndUpdate(ctx, bson.M{"user_id": room["hostID"]}, bson.M{
		"$push": bson.M{"previous_sessions": roomID},
	}).Err()
	if err != nil {
		ds.logger.Printf("failed to update previousSessions for user %s: %v\n", hostUsername, err)
		// Not returning error here because room deletion was successful
	}
	return map[string]interface{}{
		"mostLikedUser": map[string]interface{}{
			"username":   result.MostLikedUser,
			"like Count": sortedUserLikes[0].sortValue,
		},
		"mostDislikedUser": map[string]interface{}{
			"username":      result.MostDislikedUser,
			"dislike Count": sortedUserDislikes[0].sortValue,
		},
		"mostLikedSong":    result.MostLikedSong,
		"mostDislikedSong": result.MostDislikedSong,
	}, nil
}
func (ds *DocumentStore) RoomExist(roomID string) bool {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	collection := ds.db.Collection(RoomsCollection)
	ctx := context.Background()

	count, err := collection.CountDocuments(ctx, bson.M{"roomID": roomID})
	if err != nil {
		return false
	}
	return count > 0
}
func (ds *DocumentStore) AddUserToRoom(roomID, roomPassword, username string) error {
	roomsColl := ds.db.Collection(RoomsCollection)
	ctx := context.Background()

	var room bson.M
	err := roomsColl.FindOne(ctx, bson.M{"roomID": roomID}).Decode(&room)
	if err == mongo.ErrNoDocuments {
		return ErrRoomDoesntExist
	} else if err != nil {
		return err
	}
	if room["RoomStats"].(bson.M)["roomPassword"] != roomPassword {
		return ErrInvalidRoomPassword
	}
	// Check if user is already in the room
	usersJoined := room["usersJoined"].(primitive.A)
	for _, user := range usersJoined {
		if user == username {
			return ErrUserAlreadyInRoom
		}
	}
	// Check if room is full
	ds.logger.Printf("size of room is %d and maxUsers is %d\n", len(usersJoined), room["RoomStats"].(bson.M)["maxUsers"].(int64))
	if int64(len(usersJoined)) >= room["RoomStats"].(bson.M)["maxUsers"].(int64) {
		return ErrRoomFull
	}

	// Add user to the room
	_, err = roomsColl.UpdateOne(ctx, bson.M{"roomID": roomID}, bson.M{
		"$push": bson.M{"usersJoined": username},
	})
	if err != nil {
		return err
	}
	return nil
}
func (ds *DocumentStore) AddSongToQueue(roomID string, song map[string]interface{}) error {
	roomCol := ds.db.Collection(RoomsCollection)
	ctx := context.Background()
	var room bson.M
	err := roomCol.FindOne(ctx, bson.M{"roomID": roomID}).Decode(&room)
	if err != nil {
		return err
	}
	position := room["songCount"].(int32)

	// Build the nested song structure
	stats, _ := song["stats"].(map[string]interface{})
	metadata := map[string]interface{}{
		"addedBy":  stats["addedBy"],
		"likes":    0,
		"dislikes": 0,
	}

	songDoc := map[string]interface{}{
		"song": map[string]interface{}{
			"songId": song["songID"],
			"stats": map[string]interface{}{
				"title":  stats["songName"],
				"artist": stats["artistName"],
				"album":  stats["albumName"],
				// "duration": // add if available
			},
			"metadata": metadata,
		},
		"alreadyPlayed": false,
		"position":      position,
	}

	ds.logger.Printf("Adding song to room %s queue: %+v\n", roomID, songDoc)
	_, err = roomCol.UpdateOne(ctx, bson.M{"roomID": roomID}, bson.M{
		"$push": bson.M{"CurrentQueue": songDoc},
		"$inc":  bson.M{"songCount": 1},
	})
	return err
}
func (ds *DocumentStore) GetCurrentQueue(roomID string) (primitive.A, error) {
	roomCol := ds.db.Collection(RoomsCollection)
	ctx := context.Background()

	var room bson.M
	err := roomCol.FindOne(ctx, bson.M{"roomID": roomID}).Decode(&room)
	if err != nil {
		return nil, err
	}

	currentQueue := room["CurrentQueue"].(primitive.A)
	return currentQueue, nil
}

func (ds *DocumentStore) UpdateQueue(roomID string, newQueue []string) ([]interface{}, error) {
	roomColl := ds.db.Collection(RoomsCollection)
	ctx := context.Background()

	var room bson.M
	err := roomColl.FindOne(ctx, bson.M{"roomID": roomID}).Decode(&room)
	if err != nil {
		return nil, err
	}
	var temporder = make(map[string]int)
	for index, songId := range newQueue {
		temporder[songId] = index
	}

	currentQueue := room["CurrentQueue"].(primitive.A)
	var newOrderQueue = make([]interface{}, len(newQueue))
	for _, songMap := range currentQueue {
		songId := songMap.(primitive.M)["song"].(primitive.M)["songId"].(string)
		pos := temporder[songId]
		newOrderQueue[pos] = songMap
	}

	// Create a map from songId to song object for fast lookup
	// Update the database with the new ordered queue
	_, err = roomColl.UpdateOne(ctx, bson.M{"roomID": roomID}, bson.M{
		"$set": bson.M{"CurrentQueue": newOrderQueue},
	})
	if err != nil {
		return nil, err
	}

	return newOrderQueue, nil
}

// removed from function just so it doesnt get called each time for non changing values
var validOperations = map[string]bool{
	"like":       true,
	"dislike":    true,
	"un-like":    true,
	"un-dislike": true,
}

func (ds *DocumentStore) SongOperation(roomID, songID, userID, operation string) error {
	roomCol := ds.db.Collection(RoomsCollection)
	ctx := context.Background()

	var room bson.M
	err := roomCol.FindOne(ctx, bson.M{"roomID": roomID}).Decode(&room)
	if err != nil {
		return err
	}
	if !validOperations[operation] {
		return ErrInvalidSongOperation(operation)
	}
	ds.logger.Printf("%s is Performing operation '%s' on songs in room '%s'\n", userID, operation, roomID)
	// dont want the this holding up the main operation
	if userID == room["hostID"] {
		go func() {
			userInfoColl := ds.db.Collection(UserInfoCollection)
			if operation == "like" {
				_, err := userInfoColl.UpdateOne(ctx, bson.M{"userID": userID}, bson.M{"$addToSet": bson.M{"liked_songs": songID}})
				if err != nil {
					ds.logger.Printf("Failed to update user info for %s: %v\n", userID, err)
				}
			}
			if operation == "dislike" {
				_, err := userInfoColl.UpdateOne(ctx, bson.M{"userID": userID}, bson.M{"$addToSet": bson.M{"disliked_songs": songID}})
				if err != nil {
					ds.logger.Printf("Failed to update user info for %s: %v\n", userID, err)
				}
			}
		}()

	}
	switch operation {
	case "like":
		_, err = roomCol.UpdateOne(ctx,
			bson.M{
				"roomID":                   roomID,
				"CurrentQueue.song.songId": songID,
			},
			bson.M{"$inc": bson.M{"CurrentQueue.$.song.metadata.likes": 1}},
		)
		return err
	case "dislike":
		_, err = roomCol.UpdateOne(ctx,
			bson.M{
				"roomID":                   roomID,
				"CurrentQueue.song.songId": songID,
			},
			bson.M{"$inc": bson.M{"CurrentQueue.$.song.metadata.dislikes": 1}},
		)
		return err
	case "un-like":
		_, err = roomCol.UpdateOne(ctx,
			bson.M{
				"roomID":                   roomID,
				"CurrentQueue.song.songId": songID,
			},
			bson.M{"$inc": bson.M{"CurrentQueue.$.song.metadata.likes": -1}},
		)
		return err
	case "un-dislike":
		_, err = roomCol.UpdateOne(ctx,
			bson.M{
				"roomID":                   roomID,
				"CurrentQueue.song.songId": songID,
			},
			bson.M{"$inc": bson.M{"CurrentQueue.$.song.metadata.dislikes": -1}},
		)
		return err
	default:
		return ErrInvalidSongOperation(operation)
	}

}

// Most Liked songs, Most disliked songs, User with most likes/dislikes, room size , queue legth
// if no one has any likes there will be no userwith most likes/dislikes. only for likes/dislikes > 0
func (ds *DocumentStore) RoomMetrics(roomID string) (bson.M, error) {
	roomCol := ds.db.Collection(RoomsCollection)
	ctx := context.Background()

	var room bson.M
	err := roomCol.FindOne(ctx, bson.M{"roomID": roomID}).Decode(&room)
	if err != nil {
		return nil, err
	}

	currentQueue := room["CurrentQueue"].(primitive.A)
	// Track songs and user statistics
	type SongMetric struct {
		Song     interface{}
		Likes    int32
		Dislikes int32
	}

	var songs []SongMetric
	userLikes := make(map[string]int32)
	userDislikes := make(map[string]int32)

	// Process each song in the queue
	for _, songItem := range currentQueue {
		songMap := songItem.(primitive.M)
		song := songMap["song"].(primitive.M)
		metadata := song["metadata"].(primitive.M)
		// Extract likes, dislikes, and addedBy
		likes := metadata["likes"].(int32)
		dislikes := metadata["dislikes"].(int32)
		addedBy := metadata["addedBy"].(string)

		songs = append(songs, SongMetric{
			Song:     songItem,
			Likes:    likes,
			Dislikes: dislikes,
		})

		// Accumulate user statistics
		if addedBy != "" {
			userLikes[addedBy] += likes
			userDislikes[addedBy] += dislikes
		}
	}

	// Sort songs by likes (descending)
	sort.Slice(songs, func(i, j int) bool {
		return songs[i].Likes > songs[j].Likes
	})

	mostLikedSongs := make([]interface{}, 0)
	for i, song := range songs {
		if i >= 5 { // top 5 or stop if we exceed queue length
			break
		}
		mostLikedSongs = append(mostLikedSongs, song.Song)
	}

	// Sort songs by dislikes (descending)
	sort.Slice(songs, func(i, j int) bool {
		return songs[i].Dislikes > songs[j].Dislikes
	})

	mostDislikedSongs := make([]interface{}, 0)
	for i, song := range songs {
		if i >= 5 { // Top 5 or stop if we exceed queue length
			break
		}
		mostDislikedSongs = append(mostDislikedSongs, song.Song)
	}

	// Find user with most likes
	userWithMostLikes := ""
	maxLikes := int32(0)
	var tiedUsersLikes []string
	for user, likes := range userLikes {
		if likes > maxLikes {
			maxLikes = likes
			tiedUsersLikes = []string{user}
		} else if likes == maxLikes && likes > 0 {
			tiedUsersLikes = append(tiedUsersLikes, user)
		}
	}
	userWithMostLikes = strings.Join(tiedUsersLikes, ", ")

	// Find user with most dislikes
	userWithMostDislikes := ""
	maxDislikes := int32(0)
	var tiedUsersDislikes []string
	for user, dislikes := range userDislikes {
		if dislikes > maxDislikes {
			maxDislikes = dislikes
			tiedUsersDislikes = []string{user}
		} else if dislikes == maxDislikes && dislikes > 0 {
			tiedUsersDislikes = append(tiedUsersDislikes, user)
		}
	}
	userWithMostDislikes = strings.Join(tiedUsersDislikes, ", ")

	var result = map[string]interface{}{
		"mostLikedSongs":       mostLikedSongs,
		"mostDislikedSongs":    mostDislikedSongs,
		"userWithMostLikes":    userWithMostLikes,
		"userWithMostDislikes": userWithMostDislikes,
		"roomSize":             len(room["usersJoined"].(primitive.A)),
		"queueLength":          len(room["CurrentQueue"].(primitive.A)),
	}
	return result, nil
}

func (ds *DocumentStore) QueueHistory(roomID string) (primitive.A, error) {
	roomCol := ds.db.Collection(RoomsCollection)
	ctx := context.Background()

	var room bson.M
	err := roomCol.FindOne(ctx, bson.M{"roomID": roomID}).Decode(&room)
	if err != nil {
		return nil, err
	}

	history := room["playedSongs"].(primitive.A)
	return history, nil
}

func (ds *DocumentStore) GetRoomsPlaylist(roomID string) ([]interface{}, []interface{}, error) {
	roomCol := ds.db.Collection(RoomsCollection)
	ctx := context.Background()
	type SongWithLikes struct {
		Song  interface{}
		Likes int
	}

	var room bson.M
	err := roomCol.FindOne(ctx, bson.M{"roomID": roomID}).Decode(&room)
	if err != nil {
		return nil, nil, err
	}

	currentQueue := room["playedSongs"].(primitive.A)
	var songsWithLikes []SongWithLikes
	for _, songMap := range currentQueue {
		if songMap == nil {
			continue
		}

		songPrimitive, ok := songMap.(primitive.M)
		if !ok {
			continue
		}

		// Safely extract likes count
		likes := 0
		if song, ok := songPrimitive["song"].(primitive.M); ok {
			if metadata, ok := song["metadata"].(primitive.M); ok {
				if likesField, ok := metadata["likes"]; ok {
					if likesInt, ok := likesField.(int32); ok {
						likes = int(likesInt)
					} else if likesInt, ok := likesField.(int); ok {
						likes = likesInt
					}
				}
			}
		}

		songsWithLikes = append(songsWithLikes, SongWithLikes{
			Song:  songMap,
			Likes: likes,
		})
	}

	// Sort by likes (descending order)
	sort.Slice(songsWithLikes, func(i, j int) bool {
		return songsWithLikes[i].Likes > songsWithLikes[j].Likes
	})
	// Sort by likes (descending order)

	// Extract sorted songs
	sortedQueue := make([]interface{}, len(songsWithLikes))
	for i, item := range songsWithLikes {
		sortedQueue[i] = item.Song
	}

	return sortedQueue, currentQueue, nil
}
func (ds *DocumentStore) RoomState(roomID string) (map[string]interface{}, error) {
	roomColl := ds.db.Collection(RoomsCollection)
	ctx := context.Background()

	var room bson.M
	err := roomColl.FindOne(ctx, bson.M{"roomID": roomID}).Decode(&room)
	if err != nil {
		return nil, err
	}
	// get the song current playing
	currentQ := room["CurrentQueue"].(primitive.A)
	var currentSong interface{}
	if len(currentQ) > 0 {
		currentSong = currentQ[0]
	} else {
		currentSong = nil
	}
	if len(currentQ) <= 1 {
		currentQ = []interface{}{}
	}
	return map[string]interface{}{
		"roomID":        roomID,
		"currentSong":   currentSong,
		"queue":         currentQ, // Exclude the currently playing song
		"numberOfUsers": len(room["usersJoined"].(primitive.A)),
		"RoomSettings":  room["RoomStats"],
	}, nil
}

func (ds *DocumentStore) NextSong(roomID string) error {
	roomCol := ds.db.Collection(RoomsCollection)
	ctx := context.Background()

	// Find the room
	var room bson.M
	err := roomCol.FindOne(ctx, bson.M{"roomID": roomID}).Decode(&room)
	if err != nil {
		return err
	}

	// Get the current queue
	currentQ := room["CurrentQueue"].(primitive.A)
	if len(currentQ) == 0 {
		return ErrQueueIsEmpty
	}

	// Move the first song to the played songs
	playedSong := currentQ[0]
	if len(currentQ) >= 2 {
		currentQ = currentQ[1:]
	} else {
		currentQ = []interface{}{}
	}

	// Mark the song as already played
	playedSong.(primitive.M)["alreadyPlayed"] = true

	// Update the room
	_, err = roomCol.UpdateOne(ctx, bson.M{"roomID": roomID}, bson.M{"$set": bson.M{"CurrentQueue": currentQ}})
	if err != nil {
		return err
	}

	// Add the played song to the played songs
	_, err = roomCol.UpdateOne(ctx, bson.M{"roomID": roomID}, bson.M{"$push": bson.M{"playedSongs": playedSong}})
	if err != nil {
		return err
	}

	return nil
}

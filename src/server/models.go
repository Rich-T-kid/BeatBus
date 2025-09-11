package server

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type CreateRoomRequest struct {
	HostUserName string `json:"hostUsername"`
	RoomName     string `json:"roomName"`
	LifeTime     int    `json:"lifetime"` // in minutes
	MaxUsers     int    `json:"maxUsers"`
	IsPublic     bool   `json:"isPublic"`
}
type CreateRoomResponse struct {
	Properties  RoomProperties  `json:"roomProperties"`
	TimeStamp   int64           `json:"timeStamp"`
	AccessToken JWT_AccessToken `json:"accessToken"`
}

type RoomProperties struct {
	RoomID       string `json:"roomID"`
	RoomPassword string `json:"roomPassword"`
	HostID       string `json:"hostID"`
	RoomName     string `json:"roomName"`
	MaxUsers     int    `json:"maxUsers"`
	IsPublic     bool   `json:"isPublic"`
}
type JWT_AccessToken struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expiresIn"`
}

//MOngo db

type UserInfo struct {
	ID               primitive.ObjectID   `bson:"_id,omitempty" json:"id"`                    // MongoDB ObjectID
	UserID           primitive.ObjectID   `bson:"user_id" json:"user_id"`                     // Reference to user in "users" collection
	InSession        bool                 `bson:"inSession" json:"inSession"`                 // Whether user is currently in a session
	JoinDate         time.Time            `bson:"join_date" json:"join_date"`                 // Date user joined
	PreviousSessions []primitive.ObjectID `bson:"previous_sessions" json:"previous_sessions"` // List of room IDs
	ListenedTo       []string             `bson:"listened_to" json:"listened_to"`             // Songs listened to
	LikedSongs       []string             `bson:"liked_songs" json:"liked_songs"`             // Songs user liked
	DislikedSongs    []string             `bson:"disliked_songs" json:"disliked_songs"`       // Songs user disliked
}

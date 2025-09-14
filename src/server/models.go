package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	textBeltAPIURL = "https://textbelt.com/text"
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
type DeleteRoomRequest struct {
	HostUsername string `json:"hostUsername"`
	RoomID       string `json:"roomID"`
	AccessToken  string `json:"accessToken"`
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
type AddSongRequest struct {
	SongName   string `json:"songName"`
	ArtistName string `json:"artistName"`
	AlbumName  string `json:"albumName"`
	AddedBy    string `json:"addedBy"`
}
type SongMetricRequest struct {
	UserID string `json:"userID"`
	SongID string `json:"songID"`
	Action string `json:"action"` // [like, unlike,dislike, undislike, skip, play]
}

type NewOrderRequest struct {
	NewOrder []string `json:"newOrder"`
}

type txtBeltResponse struct {
	Success        bool   `json:"success"`
	TextId         string `json:"textId"`
	QuotaRemaining int    `json:"quotaRemaining"`
}

type NotifyUserRequest struct {
	UserIds []UserNotify `json:"userIds"`
}

func (nwr *NotifyUserRequest) songDataToString(songData []interface{}) string {
	// convert songData to a string representation
	fmt.Println("inside songDataToString function", 3)
	var sb strings.Builder
	for i, songItem := range songData {
		if i > 0 {
			sb.WriteString(", ") // separator between songs
		}

		// Extract song info from the nested structure
		songMap := songItem.(primitive.M)
		song := songMap["song"].(primitive.M)
		stats := song["stats"].(primitive.M)

		title := stats["title"].(string)
		artist := stats["artist"].(string)
		album := stats["album"].(string)

		// Format as "title-artist-album"
		sb.WriteString(fmt.Sprintf("\n[song %d] %s-%s-%s", i+1, title, artist, album))
	}
	return sb.String()
}

// return whether the email was sent successfully or not
func (nwr *NotifyUserRequest) sendEmail(userID, email string, songData []interface{}) (bool, string) {
	fmt.Println("inside sendEmail function", 2)
	msg := nwr.songDataToString(songData)
	fmt.Println("out of songDataToString function, msg:", msg)
	fmt.Println("Sending email to:", email)
	fmt.Println("Email content:", msg)
	// Implement actual email sending logic here
	// For now, we assume it's always successful
	return true, ""
}

// return whether the sms was sent successfully or not
func (nwr *NotifyUserRequest) sendSMS(userID, phoneNumber string, songData []interface{}) (bool, string) {
	msg := nwr.songDataToString(songData)
	fmt.Println("out of songDataToString function, msg:", msg)
	fmt.Println("Sending SMS to:", phoneNumber)
	// Implement actual SMS sending logic here
	// For now, we assume it's always successful
	values := url.Values{
		"phone":   {phoneNumber},
		"message": {msg},
		"key":     {cfg.TxtBeltAPIKey},
	}

	resp, err := http.PostForm(textBeltAPIURL, values)
	if err != nil {
		fmt.Println("Error sending SMS:", err)
		return false, ""
	}
	defer resp.Body.Close()
	var result txtBeltResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding SMS response:", err)
		return false, ""
	}
	fmt.Printf("SMS response: %+v\n", result)

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to send SMS, status code:", resp.StatusCode)
		return false, ""
	}
	if !result.Success {
		fmt.Println("SMS sending failed, response:", result)
		return false, ""
	}
	fmt.Println("SMS sent successfully, response:", result)
	return true, ""
}

type FailNotification struct {
	UserID string `json:"userID"`
	Reason string `json:"reason"`
}

func (nwr *NotifyUserRequest) Sendmessages(allSongs, mostLiked []interface{}) map[string]interface{} {
	var successful []string
	var failed []FailNotification
	fmt.Println("inside Sendmessages function", 1)
	// Implementation for sending messages
	for _, user := range nwr.UserIds {
		// validate that the users mean and method is correct
		if user.Means == "" || user.UserID == "" || user.Method == "" {
			failed = append(failed, FailNotification{UserID: user.UserID, Reason: "missing means, userID or method"})
			continue
		}
		switch user.Method {
		case "email":
			if user.IncludeMostLikedOnly {
				if suc, failureReason := nwr.sendEmail(user.UserID, user.Means, mostLiked); suc {
					successful = append(successful, user.UserID)
				} else {
					failed = append(failed, FailNotification{UserID: user.UserID, Reason: failureReason})
				}
			} else {
				if suc, failureReason := nwr.sendEmail(user.UserID, user.Means, allSongs); suc {
					successful = append(successful, user.UserID)
				} else {
					failed = append(failed, FailNotification{UserID: user.UserID, Reason: failureReason})
				}
			}
			// validate email format
		case "sms":
			if user.IncludeMostLikedOnly {
				if suc, failureReason := nwr.sendSMS(user.UserID, user.Means, mostLiked); suc {
					successful = append(successful, user.UserID)
				} else {
					failed = append(failed, FailNotification{UserID: user.UserID, Reason: failureReason})
				}
			} else {
				if suc, failureReason := nwr.sendSMS(user.UserID, user.Means, allSongs); suc {
					successful = append(successful, user.UserID)
				} else {
					failed = append(failed, FailNotification{UserID: user.UserID, Reason: failureReason})
				}
			}
			// validate phone number format
		default:
			failed = append(failed, FailNotification{UserID: user.UserID, Reason: "invalid method"})
			continue
		}
	}
	return map[string]interface{}{
		"success": successful,
		"failed":  failed,
	}
}

type UserNotify struct {
	UserID               string `json:"userID"`
	Method               string `json:"method"` // email or sms
	Means                string `json:"means"`  // email address or phone number
	IncludeMostLikedOnly bool   `json:"includeMostLikedOnly"`
}

//Mongo db

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

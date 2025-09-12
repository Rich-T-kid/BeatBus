package server

import (
	"BeatBus/internal"
	"BeatBus/storage"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

/* in the future we could make all of these methods of the server struct so that all the logs go to the same place but the doesnt matter too much*/

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// Authentication
func SignUp(w http.ResponseWriter, r *http.Request) {
	var reqBody AuthRequest
	// Parse the JSON request body
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = storage.NewDocumentStore().InsertNewUser(reqBody.Username, hashPassword(reqBody.Password))
	if err != nil {
		if err == storage.ErrUserNameTaken {
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))

}
func LogIn(w http.ResponseWriter, r *http.Request) {
	var reqBody AuthRequest
	// Parse the JSON request body
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = storage.NewDocumentStore().ValidateUser(reqBody.Username, hashPassword(reqBody.Password))
	if err != nil {
		http.Error(w, "[Invalid Creds] "+err.Error(), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

// Rooms
func JoinRoom(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomId"]
	if roomID == "" {
		http.Error(w, "Missing roomID parameter", http.StatusBadRequest)
		return
	}
	roomPassword := r.URL.Query().Get("roomPassword")
	if roomPassword == "" {
		http.Error(w, "Missing roomPassword parameter", http.StatusBadRequest)
		return
	}
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Missing username parameter", http.StatusBadRequest)
		return
	}
	err := storage.NewDocumentStore().AddUserToRoom(roomID, roomPassword, username)
	if err != nil {
		switch err {
		case storage.ErrRoomDoesntExist:
			http.Error(w, fmt.Sprintf("[The Room you are attempting to join doesn't exist] -> %s \n check that you have permission to join this room and that the provided information is correct", roomID), http.StatusNotFound)
			return
		case storage.ErrInvalidRoomPassword:
			http.Error(w, "[The Room Password you provided is incorrect] -> please try again or contact the room host", http.StatusUnauthorized)
			return
		case storage.ErrRoomFull:
			http.Error(w, "[The Room you are attempting to join is full] -> please try again later or contact the room host", http.StatusForbidden)
			return
		case storage.ErrUserAlreadyInRoom:
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("User already in room"))
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}
	resp := map[string]string{
		"username": username,
		"roomID":   roomID,
		"message":  "Successfully joined room",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
func Rooms(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var reqBody CreateRoomRequest
		// Parse the JSON request body
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		storage := storage.NewDocumentStore()
		res, err := storage.CreateRoom(reqBody.HostUserName, reqBody.RoomName, uint(reqBody.LifeTime), uint(reqBody.MaxUsers), reqBody.IsPublic)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(res)
		log.Printf("Received CreateRoom request: %+v\n", reqBody)

	case "PUT":
		// Update room settings
		err := jwtValidation(*r)
		if err != nil {
			http.Error(w, "[Invalid Token] "+err.Error(), http.StatusUnauthorized)
			return
		}
		var reqBody CreateRoomRequest
		err = json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		response, err := storage.NewDocumentStore().UpdateRoomSettings(reqBody.HostUserName, reqBody.RoomName, uint(reqBody.LifeTime), uint(reqBody.MaxUsers), reqBody.IsPublic)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(response)
	case "DELETE":
		// Delete a room
		err := jwtValidation(*r)
		if err != nil {
			http.Error(w, "[Invalid Token] "+err.Error(), http.StatusUnauthorized)
			return
		}
		var reqBody map[string]string
		err = json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		fmt.Printf("received DELETE request for room: %+v\n", reqBody)
		Winner, err := storage.NewDocumentStore().DeleteRoom(reqBody["accessToken"], reqBody["hostUsername"], reqBody["roomID"])
		if err != nil {
			if err == storage.ErrRoomDoesntExist {
				http.Error(w, fmt.Sprintf("[The Room you are attempting to delete doesn't exist] -> %s \n check that you have permission to delete this room and that the provided information is correct. \n You may have already deleted this", reqBody["roomID"]), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"MrPutOn": Winner})
	}
}
func RoomState(w http.ResponseWriter, r *http.Request) {}

// Queue
func QueuesPlaylist(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]
	if roomID == "" {
		http.Error(w, "Missing roomID parameter", http.StatusBadRequest)
		return
	}
	if !storage.NewDocumentStore().RoomExist(roomID) {
		http.Error(w, fmt.Sprintf("Room with ID [ %s ] does not exist", roomID), http.StatusNotFound)
		return
	}
	fmt.Printf("Handling playlist for roomID: %s\n", roomID)
	switch r.Method {
	case "POST":
		var reqBody AddSongRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		err = storage.NewDocumentStore().AddSongToQueue(roomID, map[string]interface{}{
			"songID": internal.RandomHash(),
			"stats": map[string]interface{}{
				"songName":   reqBody.SongName,
				"artistName": reqBody.ArtistName,
				"albumName":  reqBody.AlbumName,
				"addedBy":    reqBody.AddedBy,
			},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		NewDownloadQueue().RetrieveSong(reqBody)
		w.WriteHeader(http.StatusCreated)
		// Add song to queue
	case "GET":
		// Get current queue
		resp, err := storage.NewDocumentStore().GetCurrentQueue(roomID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(resp)
	case "PUT":
		//TODO: do this endpoint skipping it for now because its annoying to implement
		// Update queue (e.g., reorder songs)
		err := jwtValidation(*r)
		if err != nil {
			http.Error(w, "[Invalid Token] "+err.Error(), http.StatusUnauthorized)
			return
		}
		var reqBody []string
		err = json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		updatedQueue, err := storage.NewDocumentStore().UpdateQueue(roomID, reqBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(updatedQueue)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Metrics
func Metrics(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]
	if roomID == "" {
		http.Error(w, "Missing roomID parameter", http.StatusBadRequest)
		return
	}
	if !storage.NewDocumentStore().RoomExist(roomID) {
		http.Error(w, fmt.Sprintf("Room with ID [ %s ] does not exist", roomID), http.StatusNotFound)
		return
	}
	switch r.Method {
	case "GET":
		resp, err := storage.NewDocumentStore().RoomMetrics(roomID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(resp)
	case "POST":
		// Create new metric entry
	}
}

// TODO: This needs to be idempotent
func MetricsPlaylistSend(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]
	if roomID == "" {
		http.Error(w, "Missing roomID parameter", http.StatusBadRequest)
		return
	}
	if !storage.NewDocumentStore().RoomExist(roomID) {
		http.Error(w, fmt.Sprintf("Room with ID [ %s ] does not exist", roomID), http.StatusNotFound)
		return
	}
	err := jwtValidation(*r)
	if err != nil {
		http.Error(w, "[Invalid Token] "+err.Error(), http.StatusUnauthorized)
		return
	}
	var reqBody NotifyUserRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received notify request for roomID: %s with body: %+v\n", roomID, reqBody)
}
func MetricsHistory(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]
	if roomID == "" {
		http.Error(w, "Missing roomID parameter", http.StatusBadRequest)
		return
	}
	if !storage.NewDocumentStore().RoomExist(roomID) {
		http.Error(w, fmt.Sprintf("Room with ID [ %s ] does not exist", roomID), http.StatusNotFound)
		return
	}
	resp, err := storage.NewDocumentStore().QueueHistory(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

// Handlers

// Middleware
func jwtValidation(r http.Request) error {
	token := r.Header.Get("Authorization")
	if !strings.HasPrefix(token, "Bearer ") {
		return fmt.Errorf("invalid token format")
	}
	token = strings.TrimPrefix(token, "Bearer ")
	err := internal.NewJWTHandler().VerifyToken(token)
	return err
}

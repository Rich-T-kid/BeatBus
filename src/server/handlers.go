package server

import (
	"BeatBus/internal"
	"BeatBus/storage"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
func (s *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	var reqBody AuthRequest
	// Parse the JSON request body
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if reqBody.Username == "" || reqBody.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}
	err = storage.NewDocumentStore(s.documentLogger).InsertNewUser(reqBody.Username, hashPassword(reqBody.Password))
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
func (s *Server) LogIn(w http.ResponseWriter, r *http.Request) {
	var reqBody AuthRequest
	// Parse the JSON request body
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if reqBody.Username == "" || reqBody.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}
	err = storage.NewDocumentStore(s.documentLogger).ValidateUser(reqBody.Username, hashPassword(reqBody.Password))
	if err != nil {
		http.Error(w, "[Invalid Creds] "+err.Error(), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

// Rooms
func (s *Server) JoinRoom(w http.ResponseWriter, r *http.Request) {
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
	err := storage.NewDocumentStore(s.documentLogger).AddUserToRoom(roomID, roomPassword, username)
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
func (s *Server) Rooms(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var reqBody CreateRoomRequest
		// Parse the JSON request body
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if reqBody.HostUserName == "" || reqBody.RoomName == "" || reqBody.LifeTime <= 0 || reqBody.LifeTime > 300 || reqBody.MaxUsers <= 0 {
			http.Error(w, "HostUserName, RoomName, LifeTime and MaxUsers are required and must be greater than 0. Lifetime must be between 1 and 300 (minutes)", http.StatusBadRequest)
			return
		}
		storage := storage.NewDocumentStore(s.documentLogger)
		res, err := storage.CreateRoom(reqBody.HostUserName, reqBody.RoomName, uint(reqBody.LifeTime), uint(reqBody.MaxUsers), reqBody.IsPublic)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(res)
		s.logger.Printf("Received CreateRoom request: %+v\n", reqBody)

	case "PUT":
		// Update room settings -> user cannot increase/update the room lifetime once its made
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
		if reqBody.HostUserName == "" || reqBody.RoomName == "" || reqBody.MaxUsers <= 0 {
			http.Error(w, "HostUserName, RoomName, LifeTime and MaxUsers are required and must be greater than 0. Lifetime must be between 1 and 300 (minutes)", http.StatusBadRequest)
			return
		}
		response, err := storage.NewDocumentStore(s.documentLogger).UpdateRoomSettings(reqBody.HostUserName, reqBody.RoomName, uint(reqBody.MaxUsers), reqBody.IsPublic)
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
		var reqBody DeleteRoomRequest
		err = json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if reqBody.HostUsername == "" || reqBody.RoomID == "" || reqBody.AccessToken == "" {
			http.Error(w, "hostUsername, roomID and accessToken are required", http.StatusBadRequest)
			return
		}
		s.logger.Printf("received DELETE request for room: %+v\n", reqBody)
		Winner, err := storage.NewDocumentStore(s.documentLogger).DeleteRoom(reqBody.AccessToken, reqBody.HostUsername, reqBody.RoomID)
		if err != nil {
			if err == storage.ErrRoomDoesntExist {
				http.Error(w, fmt.Sprintf("[The Room you are attempting to delete doesn't exist] -> %s \n check that you have permission to delete this room and that the provided information is correct. \n You may have already deleted this", reqBody.RoomID), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"MrPutOn": Winner})
	}
}

// TODO:  This needs to be implemented using websock so that the client can get real time updates
func (s *Server) RoomState(w http.ResponseWriter, r *http.Request) {}

// Queue
func (s *Server) QueuesPlaylist(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]
	if roomID == "" {
		http.Error(w, "Missing roomID parameter", http.StatusBadRequest)
		return
	}
	if !storage.NewDocumentStore(s.documentLogger).RoomExist(roomID) {
		http.Error(w, fmt.Sprintf("Room with ID [ %s ] does not exist", roomID), http.StatusNotFound)
		return
	}
	s.logger.Printf("Handling playlist for roomID: %s\n", roomID)
	switch r.Method {
	//TODO: add more input validation here. if the fields are empty return an error all fields must be present
	case "POST":
		var reqBody AddSongRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if reqBody.SongName == "" || reqBody.ArtistName == "" || reqBody.AlbumName == "" || reqBody.AddedBy == "" {
			http.Error(w, "SongName, ArtistName, AlbumName and AddedBy are required", http.StatusBadRequest)
			return
		}
		err = storage.NewDocumentStore(s.documentLogger).AddSongToQueue(roomID, map[string]interface{}{
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
		go NewDownloadQueue().RetrieveSong(reqBody)
		w.WriteHeader(http.StatusCreated)
	case "GET":
		// Get current queue
		resp, err := storage.NewDocumentStore(s.documentLogger).GetCurrentQueue(roomID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(resp)
	case "PUT":
		err := jwtValidation(*r)
		if err != nil {
			http.Error(w, "[Invalid Token] "+err.Error(), http.StatusUnauthorized)
			return
		}
		var reqBody NewOrderRequest
		err = json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		updatedQueue, err := storage.NewDocumentStore(s.documentLogger).UpdateQueue(roomID, reqBody.NewOrder)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(updatedQueue)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) Metrics(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]
	if roomID == "" {
		http.Error(w, "Missing roomID parameter", http.StatusBadRequest)
		return
	}
	if !storage.NewDocumentStore(s.documentLogger).RoomExist(roomID) {
		http.Error(w, fmt.Sprintf("Room with ID [ %s ] does not exist", roomID), http.StatusNotFound)
		return
	}
	switch r.Method {
	case "GET":
		// TODO: Mabey could make websocket connection but not the prio right not at all
		resp, err := storage.NewDocumentStore(s.documentLogger).RoomMetrics(roomID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(resp)
	case "POST":
		// TODO: This needs to be idempotent
		var reqBody SongMetricRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if reqBody.SongID == "" || reqBody.Action == "" || reqBody.UserID == "" {
			http.Error(w, "SongID, Action and UserID are required", http.StatusBadRequest)
			return
		}
		err = storage.NewDocumentStore(s.documentLogger).SongOperation(roomID, reqBody.SongID, reqBody.Action, reqBody.UserID)
		if err != nil {
			switch err {
			case storage.ErrInvalidSongOperation(reqBody.Action):
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) MetricsPlaylistSend(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]
	if roomID == "" {
		http.Error(w, "Missing roomID parameter", http.StatusBadRequest)
		return
	}
	if !storage.NewDocumentStore(s.documentLogger).RoomExist(roomID) {
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
	s.logger.Printf("Received notify request for roomID: %s with body: %+v\n", roomID, reqBody)
	mostLiked, currentQueue, err := storage.NewDocumentStore(s.documentLogger).GetRoomsPlaylist(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := reqBody.Sendmessages(mostLiked, currentQueue)
	if len(resp["failed"].([]FailNotification)) > 0 {
		w.WriteHeader(http.StatusPartialContent)
		json.NewEncoder(w).Encode(resp)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
func (s *Server) MetricsHistory(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]
	if roomID == "" {
		http.Error(w, "Missing roomID parameter", http.StatusBadRequest)
		return
	}
	if !storage.NewDocumentStore(s.documentLogger).RoomExist(roomID) {
		http.Error(w, fmt.Sprintf("Room with ID [ %s ] does not exist", roomID), http.StatusNotFound)
		return
	}
	resp, err := storage.NewDocumentStore(s.documentLogger).QueueHistory(roomID)
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

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

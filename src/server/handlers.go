package server

import (
	"BeatBus/storage"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	log.Printf("Received SignUp request: %+v\n", reqBody)
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
func Refresh(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if !strings.HasPrefix(token, "Bearer ") {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")
	fmt.Println("Received Refresh request with token:", token)
	w.WriteHeader(http.StatusOK)

}

// Rooms
func JoinRoom(w http.ResponseWriter, r *http.Request) {}
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
		w.WriteHeader(http.StatusCreated)
		log.Printf("Received CreateRoom request: %+v\n", reqBody)

	case "PUT":
		// Update room settings
	case "DELETE":
		// Delete a room
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func RoomState(w http.ResponseWriter, r *http.Request) {}

// Queue
func QueuesPlaylist(w http.ResponseWriter, r *http.Request) {}

// Metrics
func Metrics(w http.ResponseWriter, r *http.Request)             {}
func MetricsPlaylistSend(w http.ResponseWriter, r *http.Request) {}
func MetricsHistory(w http.ResponseWriter, r *http.Request)      {}

// Handlers

// Middleware
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path, "executing logMiddleware")
		next.ServeHTTP(w, r)
		log.Println(r.URL.Path, "executing logMiddleware again")
	})
}

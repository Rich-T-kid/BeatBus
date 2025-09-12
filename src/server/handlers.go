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
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func RoomState(w http.ResponseWriter, r *http.Request) {}

// Queue
func QueuesPlaylist(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

		// Add song to queue
	case "GET":
		// Get current queue
	case "PUT":
		err := jwtValidation(*r)
		if err != nil {
			http.Error(w, "[Invalid Token] "+err.Error(), http.StatusUnauthorized)
			return
		}
		// Update queue (e.g., reorder songs)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Metrics
func Metrics(w http.ResponseWriter, r *http.Request) {}
func MetricsPlaylistSend(w http.ResponseWriter, r *http.Request) {
	err := jwtValidation(*r)
	if err != nil {
		http.Error(w, "[Invalid Token] "+err.Error(), http.StatusUnauthorized)
		return
	}
}
func MetricsHistory(w http.ResponseWriter, r *http.Request) {}

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

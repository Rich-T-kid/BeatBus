package server

import (
	"net/http"
)

// Authentication
func SignUp(w http.ResponseWriter, r *http.Request) {}
func LogIn(w http.ResponseWriter, r *http.Request)  {}

// Rooms
func JoinRoom(w http.ResponseWriter, r *http.Request)  {}
func Rooms(w http.ResponseWriter, r *http.Request)     {}
func RoomState(w http.ResponseWriter, r *http.Request) {}

// Queue
func QueuesPlaylist(w http.ResponseWriter, r *http.Request) {}

// Metrics
func Metrics(w http.ResponseWriter, r *http.Request)             {}
func MetricsPlaylistSend(w http.ResponseWriter, r *http.Request) {}
func MetricsHistory(w http.ResponseWriter, r *http.Request)      {}

package server

import (
	"BeatBus/internal"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var cfg = internal.GetConfig()

type Server struct {
	port   string
	logger *log.Logger
}

func NewServer() *Server {
	return &Server{
		port:   ":" + cfg.Port,
		logger: log.New(os.Stdout, "[Server] : ", log.Lshortfile|log.Ltime),
	}
}
func (s *Server) registerMiddleware(r *mux.Router, middleware []mux.MiddlewareFunc) *mux.Router {
	for _, m := range middleware {
		r.Use(m)
	}
	return r
}

func (s *Server) StartServer() error {
	middleware := []mux.MiddlewareFunc{
		logMiddleware,
	}
	router := s.registerRoutes()
	router = s.registerMiddleware(router, middleware)
	s.logger.Println(" Starting server on port", s.port)
	http.ListenAndServe(s.port, router)
	return nil
}

func (s *Server) registerRoutes() *mux.Router {
	router := mux.NewRouter()
	// Authentication
	router.HandleFunc("/signUp", SignUp).Methods("POST")
	router.HandleFunc("/logIn", LogIn).Methods("POST")

	// Rooms
	router.HandleFunc("/rooms/{roomId}", JoinRoom).Methods("GET")
	router.HandleFunc("/rooms", Rooms).Methods("POST", "PUT", "DELETE")
	router.HandleFunc("/rooms/{roomid}/state", RoomState).Methods("GET")

	// Queue
	router.HandleFunc("/queues/{roomID}/playlist", QueuesPlaylist).Methods("POST", "GET", "PUT")

	// Metrics
	router.HandleFunc("/metrics/{roomID}", Metrics).Methods("GET", "POST")
	router.HandleFunc("/metrics/{roomID}/playlist/send", MetricsPlaylistSend).Methods("POST")
	router.HandleFunc("/metrics/{roomID}/history", MetricsHistory).Methods("GET")

	return router
}

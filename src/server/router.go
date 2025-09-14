package server

import (
	"BeatBus/internal"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var cfg = internal.GetConfig()

type Server struct {
	port           string
	documentLogger *log.Logger
	cacheLogger    *log.Logger
	logger         *log.Logger
}

func NewServer() *Server {
	logFile, err := os.OpenFile(cfg.OutputFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file (%s) due to error: %v", cfg.OutputFileName, err)
	}

	// Create MultiWriter to write to both stdout and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	return &Server{
		port:           ":" + cfg.Port,
		logger:         log.New(multiWriter, "[Server] : ", log.Lshortfile|log.Ltime),
		documentLogger: log.New(multiWriter, "[DocumentStore] : ", log.Lshortfile|log.Ltime),
		cacheLogger:    log.New(multiWriter, "[Cache] : ", log.Lshortfile|log.Ltime),
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
		Cors,
	}
	router := s.registerRoutes()
	router = s.registerMiddleware(router, middleware)
	s.logger.Println(" Starting server on port", s.port)
	http.ListenAndServe(s.port, router)
	return nil
}

func (s *Server) registerRoutes() *mux.Router {
	router := mux.NewRouter()

	//Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}).Methods("GET")

	// Authentication
	router.HandleFunc("/signUp", s.SignUp).Methods("POST")
	router.HandleFunc("/login", s.LogIn).Methods("POST")

	// Rooms
	router.HandleFunc("/rooms/{roomId}", s.JoinRoom).Methods("GET")
	router.HandleFunc("/rooms", s.Rooms).Methods("POST", "PUT", "DELETE")
	router.HandleFunc("/rooms/{roomid}/state", s.RoomState).Methods("GET")

	// Queue
	router.HandleFunc("/queues/{roomID}/playlist", s.QueuesPlaylist).Methods("POST", "GET", "PUT")

	// Metrics
	router.HandleFunc("/metrics/{roomID}", s.Metrics).Methods("GET", "POST")
	router.HandleFunc("/metrics/{roomID}/playlist/send", s.MetricsPlaylistSend).Methods("POST")
	router.HandleFunc("/metrics/{roomID}/history", s.MetricsHistory).Methods("GET")

	return router
}

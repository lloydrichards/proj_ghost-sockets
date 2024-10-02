package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	_ "github.com/joho/godotenv/autoload"

	"server/internal/database"
)

type State struct {
	X      int  `json:"x"`
	Y      int  `json:"y"`
}
type Client struct {
	Username string `json:"username"`
	Online bool `json:"online"`
	State
}

type Server struct {
	port  int
	conns map[*websocket.Conn]bool
	hub   map[uuid.UUID]Client
	db    database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:  port,
		conns: make(map[*websocket.Conn]bool),
		hub:   make(map[uuid.UUID]Client),
		db:    database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

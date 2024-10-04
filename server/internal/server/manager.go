package server

import (
	"errors"
	"log"
	"net/http"
	"server/internal/database"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/rand"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Manager struct {
	Clients ClientList
	db      database.Service
	sync.RWMutex

	handlers map[string]EventHandler
}

func NewManager(db *database.Service) *Manager {
	m := &Manager{
		Clients:  make(ClientList),
		db:       *db,
		handlers: make(map[string]EventHandler),
	}

	m.setupHandlers()

	return m

}

func (m *Manager) setupHandlers() {
	m.handlers["update_position"] = UpdatePosition
}

func (m *Manager) initiateWSConnection(c *gin.Context) {
	username := c.Query("username")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("%s, error while Upgrading websocket connection\n", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	log.Printf("New connection from %s\n", username)

	client := NewClient(username, conn, m)

	m.addClient(client)

	// Go Routines
	go client.readMsgs()
	go client.writeMsg()

}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.db.CreateUser(database.User{
		Name:  client.username,
		Color: rand.Intn(10),
	})

	m.db.CreateSession(database.Session{
		ID:       client.id,
		UserName: client.username,
	})

	m.Clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.db.UpdateSession(client.id)

	if _, ok := m.Clients[client]; ok {
		client.conn.Close()
		delete(m.Clients, client)
	}
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	handler, ok := m.handlers[event.Type]
	if !ok {
		return errors.New("no handler for event type")
	}
	err := handler(event, c)
	if err != nil {
		return err
	}
	return nil
}

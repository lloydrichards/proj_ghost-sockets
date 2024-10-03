package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	sync.RWMutex

	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		Clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}

	m.setupHandlers()

	return m

}

func (m *Manager) setupHandlers() {
	m.handlers["update_position"] = UpdatePosition
}

func UpdatePosition(event Event, c *Client) error {

	var update UpdatePositionEvent
	err := json.Unmarshal(event.Payload, &update)
	if err != nil {
		return err
	}
	c.state.X = update.X
	c.state.Y = update.Y
	log.Printf("Update: %s ->    x %d   y %d", c.username, update.X, update.Y)

	for client := range c.manager.Clients {
		Broadcast(event, client)
	}
	return nil
}

func Broadcast(event Event, c *Client) error {
	payloadJson := make(map[string]interface{})

	for client := range c.manager.Clients {
		payloadJson[client.id.String()] = map[string]interface{}{
			"username": client.username,
			"state":    client.state,
		}
	}

	json, err := json.Marshal(payloadJson)
	if err != nil {
		return err
	}
	c.egress <- Event{
		Type:    "broadcast",
		Payload: json,
	}

	return nil
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

func (m *Manager) handleWSConnect(c *gin.Context) {
	username := c.Query("username")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// panic(err)
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

	m.Clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.Clients[client]; ok {
		client.conn.Close()
		delete(m.Clients, client)
	}
}

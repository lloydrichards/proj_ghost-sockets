package server

import (
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
}

func NewManager() *Manager {
	return &Manager{
		Clients: make(ClientList),
	}
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

	client := NewClient(conn, m)
	m.addClient(client)

	// Go Routines
	go client.readMsgs()
	go client.writeMsgs()

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

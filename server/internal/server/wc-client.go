package server

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type State struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Client struct {
	id       uuid.UUID
	username string
	state    State

	// websocket connection
	conn    *websocket.Conn
	manager *Manager

	// channels for communication
	egress chan Event
}

func NewClient(username string, conn *websocket.Conn, manager *Manager) *Client {
	id := uuid.New()
	return &Client{
		id:       id,
		username: username,
		conn:     conn,
		manager:  manager,
		state:    State{X: 0, Y: 0},

		egress: make(chan Event),
	}
}

func (c *Client) readMsgs() {
	defer func() {
		// cleanup connection
		c.manager.removeClient(c)
	}()

	for {
		_, payload, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading Msg: %v", err)
			}
			break
		}

		var request Event

		err = json.Unmarshal(payload, &request)
		if err != nil {
			log.Printf("error unmarshalling Msg: %v", err)
			continue
		}

		err = c.manager.routeEvent(request, c)
		if err != nil {
			log.Printf("error routing Msg: %v", err)
			continue
		}
	}
}

func (c *Client) writeMsg() {
	defer func() {
		// cleanup connection
		c.manager.removeClient(c)
	}()

	for {
		select {
		case msg, ok := <-c.egress:
			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, nil)
				if err != nil {
					log.Printf("connection closed: %v", err)
				}
				return
			}
			data, err := json.Marshal(msg)
			if err != nil {
				log.Printf("failed to marshal Msg: %v", err)
				return
			}

			err = c.conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("failed to writing Msg: %v", err)
			}
		}
	}
}

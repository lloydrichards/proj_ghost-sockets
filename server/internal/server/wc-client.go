package server

import (
	"log"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	conn    *websocket.Conn
	manager *Manager

	// channels for communication
	egress chan []byte
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		conn:    conn,
		manager: manager,
		egress:  make(chan []byte),
	}
}

func (c *Client) readMsgs() {
	defer func() {
		// cleanup connection
		c.manager.removeClient(c)
	}()

	for {
		msgType, payload, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading Msg: %v", err)
			}
			break
		}

		for wsClient := range c.manager.Clients {
			wsClient.egress <- payload
		}
		log.Println(msgType)
		log.Println(payload)
	}
}

func (c *Client) writeMsgs() {
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
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("failed to writing Msg: %v", err)
			}
			log.Println("msg sent")
		}
	}
}

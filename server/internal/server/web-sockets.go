package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) handleWSConnect(c *gin.Context) {
	username := c.Query("username")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// panic(err)
		log.Printf("%s, error while Upgrading websocket connection\n", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	client := s.initClient(conn, username)

	s.readMsgLoop(conn, client.ID)

	defer s.removeClient(conn, client.ID)

}

func (s *Server) readMsgLoop(wc *websocket.Conn, userId uuid.UUID) {
	for {
		var pos State
		err := wc.ReadJSON(&pos)
		if err != nil {
			log.Printf("error while reading message: %s\n", err.Error())
			wc.Close()
			delete(s.conns, wc)
			delete(s.hub, userId)
			break
		}

		// update the position of the user
		s.hub[userId] = Client{
			ID:       userId,
			Username: s.hub[userId].Username,
			Online:   true,
			State: State{
				X: pos.X,
				Y: pos.Y,
			},
		}

		// broadcast the users
		s.broadcastHub()
	}
}

func (s *Server) broadcastHub() {
	log.Printf("broadcasting to %d connections...", len(s.conns))
	for conn := range s.conns {
		err := conn.WriteJSON(s.hub)
		if err != nil {
			log.Printf("error while sending message to connection: %s\n", err.Error())
			conn.Close()
			delete(s.conns, conn)
		}
	}
}

func (s *Server) initClient(wc *websocket.Conn, username string) Client {
	userId := uuid.New()

	log.Printf("new connection from %s (%s)\n", username, userId)
	s.conns[wc] = true
	client := Client{
		Username: username,
		Online:   true,
		State: State{
			X: 0,
			Y: 0,
		},
	}
	s.hub[userId] = client

	return client
}

func (s *Server) removeClient(wc *websocket.Conn, userId uuid.UUID) {
	wc.Close()
	delete(s.conns, wc)
	delete(s.hub, userId)
}

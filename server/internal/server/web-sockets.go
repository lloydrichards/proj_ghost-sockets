package server

import (
	"log"
	"net/http"

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

func (s *Server) handleWSConnect(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// panic(err)
		log.Printf("%s, error while Upgrading websocket connection\n", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	for {
		// Read message from client
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			// panic(err)
			log.Printf("%s, error while reading message\n", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
			break
		}

		// Echo message back to client
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			// panic(err)
			log.Printf("%s, error while writing message\n", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
			break
		}
	}
}

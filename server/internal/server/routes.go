package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	r.GET("/users", s.getUsersHandler)

	r.GET("/users/:username", s.getUserHandler)

	r.PATCH("/users/:username", s.updateUserHandler)

	r.GET("/ws", s.manager.initiateWSConnection)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

type User struct {
	Username      string `json:"username"`
	Color         string `json:"color"`
	Mood          string `json:"mood"`
	IsActive      bool   `json:"is_active"`
	LastSessionId string `json:"last_session_id"`
}

func (s *Server) getUsersHandler(c *gin.Context) {
	clients, err := s.db.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	users := []User{}
	for _, client := range clients {
		session, err := s.db.GetLatestSession(client.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, User{
			Username:      session.UserName,
			Color:         client.Color,
			Mood:          client.Mood,
			IsActive:      session.IsActive,
			LastSessionId: session.ID.String(),
		})

	}
	c.JSON(http.StatusOK, users)
}

func (s *Server) getUserHandler(c *gin.Context) {

	username := c.Param("username")

	client, err := s.db.GetUser(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	session, err := s.db.GetLatestSession(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := User{
		Username:      session.UserName,
		Color:         client.Color,
		Mood:          client.Mood,
		IsActive:      session.IsActive,
		LastSessionId: session.ID.String(),
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) updateUserHandler(c *gin.Context) {
	username := c.Param("username")

	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := s.db.GetUser(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// overwrite the fields that are not empty
	if user.Color != "" {
		client.Color = user.Color
	}
	if user.Mood != "" {
		client.Mood = user.Mood
	}

	if err := s.db.UpdateUser(client); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, client)
}

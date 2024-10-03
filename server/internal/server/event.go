package server

import "encoding/json"

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventUpdatePosition = "update_position"
)

type UpdatePositionEvent struct {
	X int `json:"x"`
	Y int `json:"y"`
}

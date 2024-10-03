package server

import (
	"encoding/json"
	"log"
	"math"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventUpdatePosition = "update_position"
)

type UpdatePositionEvent struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Delta int `json:"delta"`
}

func UpdatePosition(event Event, c *Client) error {

	var update UpdatePositionEvent
	err := json.Unmarshal(event.Payload, &update)
	if err != nil {
		return err
	}

	log.Printf("Update: %s ->    x %d   y %d   delta %d", c.username, update.X, update.Y, update.Delta)

	prevPos := Position{X: c.state.X, Y: c.state.Y}
	curPos := Position{X: float64(update.X), Y: float64(update.Y)}

	c.state.X = curPos.X
	c.state.Y = curPos.Y
	vx, vy := velocity(prevPos, curPos, float64(update.Delta))
	ang := angle(prevPos, curPos)
	spd := speed(vx, vy)
	acc := acceleration(c.state.Spd, spd, float64(update.Delta))

	c.state.Vx = vx
	c.state.Vy = vy
	c.state.Ang = ang
	c.state.Spd = spd
	c.state.Acc = acc

	for client := range c.manager.Clients {
		broadcastState(client)
	}
	return nil
}

type Position struct {
	X, Y float64
}

func displacement(prev, curr Position) (dx, dy float64) {
	dx = curr.X - prev.X
	dy = curr.Y - prev.Y
	return
}

func velocity(prev, curr Position, deltaTime float64) (vx, vy float64) {
	dx, dy := displacement(prev, curr)
	vx = dx / deltaTime
	vy = dy / deltaTime
	return
}

func speed(vx, vy float64) float64 {
	return math.Sqrt(vx*vx + vy*vy)
}

func angle(prev, curr Position) float64 {
	dx, dy := displacement(prev, curr)
	return math.Atan2(dy, dx)
}

func acceleration(v1, v2 float64, deltaTime float64) float64 {
	return (v2 - v1) / deltaTime
}

func broadcastState(c *Client) error {
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

package websocket

import "encoding/json"

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type CreateRoomEvent struct {
	RoomName string
	RoomSize int
}

type EventHandler func(event Event, c *Client) error

const (
	// EventStatusUpdate = "notify_user_status"
	// EventSendMessage    = "send_message"
	EventRegisterRoom   = "register_room"
	EventUnregisterRoom = "unregister_room"
	EventJoinRoom       = "join_room"
	EventLeaveRoom      = "leave_room"
	// EventStartQuiz      = "start_quiz"
)

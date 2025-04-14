package websocket

import "encoding/json"

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"info"`
}

type CreateRoomEvent struct {
	RoomName string `json:"room_name"`
	RoomSize int    `json:"room_size"`
}

type JoinRoomEvent struct {
	RoomID string `json:"room_id"`
}

type LeaveRoomEvent struct {
	RoomID string `json:"room_id"`
}

const (
	// EventStatusUpdate = "notify_user_status"
	// EventSendMessage    = "send_message"
	EventCreateRoom = "create_room"
	EventJoinRoom   = "join_room"
	EventLeaveRoom  = "leave_room"
	// EventStartQuiz      = "start_quiz"
)

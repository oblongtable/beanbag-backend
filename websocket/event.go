package websocket

import "encoding/json"

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"info"`
}

type CreateRoomEvent struct {
	RoomName string `json:"room_name"`
	RoomSize int    `json:"room_size"`
	UserName string `json:"username"`
}

type JoinRoomEvent struct {
	RoomID string `json:"room_id"`
	Name   string `json:"name"`
}

type LeaveRoomEvent struct {
	RoomID string `json:"room_id"`
}

type StartQuizEvent struct {
	RoomID string `json:"room_id"`
}

type SubmitAnswerEvent struct {
	AnswerIndex int `json:"answer_index"`
}

const (
	// EventStatusUpdate = "notify_user_status"
	// EventSendMessage    = "send_message"
	EventCreateRoom   = "create_room"
	EventJoinRoom     = "join_room"
	EventLeaveRoom    = "leave_room"
	EventStartQuiz    = "start_quiz"
	EventForwardQuiz  = "quiz_forward" // New event for moving the quiz forward
	EventSubmitAnswer = "submit_answer"
)

package websocket

import "encoding/json"

type BaseMessage struct {
	Type string `json:"type"`
}

type EventCallbackMessage struct {
	BaseMessage
	IsSuccess bool            `json:"success"`
	Message   string          `json:"message"`
	Info      json.RawMessage `json:"info"`
}

func NewEventCallbackMessage() *EventCallbackMessage {
	return &EventCallbackMessage{
		BaseMessage: BaseMessage{},
		IsSuccess:   true,
		Message:     "",
		Info:        nil,
	}
}

type UserInfoMessage struct {
	BaseMessage
	User UserInfo `json:"user"`
}

type UserInfo struct {
	ID       string `json:"user_id"`
	Username string `json:"user_name"`
}

type UserStatusMessages struct {
	BaseMessage
	Users []*UserInfo
}

type RoomInfoMessage struct {
	BaseMessage
	Room RoomInfo
}

type RoomInfo struct {
	ID        string      `json:"room_id"`
	Name      string      `json:"room_name"`
	Size      int         `json:"room_size"`
	UsersInfo []*UserInfo `json:"users_info"`
}

type RoomInfoMessages struct {
	BaseMessage
	Rooms []*RoomInfo
}

const (
	// MessageStatusUpdateUserAll = "status_update_user_all"
	// MessageStatusUpdateRoomAll = "status_update_room_all"
	// MessageStatusUpdateUser = "status_update_user"
	MessageRoomStatusUpdate = "room_status_update"

	MessageUserJoinRoomUpdate  = "user_join_room_update"
	MessageUserLeaveRoomUpdate = "user_leave_room_update"

	MessageCreateRoom = "create_room_callback"
	MessageJoinRoom   = "join_room_callback"
	MessageLeaveRoom  = "leave_room_callback"
)

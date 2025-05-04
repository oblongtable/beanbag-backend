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

type RoomInfo struct {
	BaseMessage
	ID        string      `json:"room_id"`
	Name      string      `json:"room_name"`
	Size      int         `json:"room_size"`
	UsersInfo []*UserInfo `json:"users_info"`
	HostID    string      `json:"host_id"` // Add HostID field
}

type RoomInfoMessages struct {
	BaseMessage
	Rooms []*RoomInfo `json:"rooms_info"`
}

type Serialisable interface {
	RoomInfo | UserInfo | EventCallbackMessage
}

const (
	MessageRoomStatusUpdate = "room_status_update"

	MessageUserJoinRoomUpdate  = "user_join_room_update"
	MessageUserLeaveRoomUpdate = "user_leave_room_update"

	MessageCreateRoom = "create_room_callback"
	MessageJoinRoom   = "join_room_callback"
	MessageLeaveRoom  = "leave_room_callback"
)

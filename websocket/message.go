package websocket

type BaseMessage struct {
	Type string `json:"type"`
}

type UserStatusMessage struct {
	BaseMessage
	User UserStatus
}

type UserStatus struct {
	ID       string
	Username string
	IsAlive  bool
}

type UserStatusMessages struct {
	BaseMessage
	Users []UserStatus
}

type RoomStatusMessage struct {
	BaseMessage
	Room RoomStatus
}

type RoomStatus struct {
	ID      string
	Name    string
	Size    int
	IsAlive bool
}

type RoomStatusMessages struct {
	BaseMessage
	Rooms []RoomStatus
}

const (
	MessageStatusUpdateUserAll = "status_update_user_all"
	MessageStatusUpdateRoomAll = "status_update_room_all"
	MessageStatusUpdateUser    = "status_update_user"
	MessageStatusUpdateRoom    = "status_update_room"
)

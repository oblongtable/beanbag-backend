package websocket

type BaseMessage struct {
	Type string `json:"type"`
}

type UserStatusMessage struct {
	BaseMessage
	User struct {
		ID       string
		Username string
		IsAlive  bool
	}
}

const (
	MessageStatusUpdate = "status_update"
)

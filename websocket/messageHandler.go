package websocket

import (
	"encoding/json"
	"log"
)

func NewEventCallbackMessage() *EventCallbackMessage {
	return &EventCallbackMessage{
		BaseMessage: BaseMessage{},
		IsSuccess:   true,
		Message:     "",
		Info:        nil,
	}
}

func SendEventCallback[S Serialisable](c *Client, msg_type string, isSuccess bool, msg string, ser *S) {
	evtCbMsg := &EventCallbackMessage{
		BaseMessage: BaseMessage{msg_type},
		IsSuccess:   isSuccess,
		Message:     msg,
		Info:        nil,
	}

	// Serialise Info
	if isSuccess {
		if jsonBytes, err := json.Marshal(ser); err != nil {
			isSuccess = false
			log.Printf("SendEventCallback: Failed to marshal message (%v)\n", err)
			return

		} else {
			evtCbMsg.Info = jsonBytes
		}
	}

	// Serialise callback message and send it
	if strmsg, err := json.Marshal(evtCbMsg); err == nil {
		log.Printf("Marshaled: %s", strmsg)
		c.Send <- strmsg
	} else {
		log.Printf("Failed to marshal: %v", err)
	}
}

func NotifyUserRoomStatus(r *Room, c *Client, msg_type string) error {
	roomInfo := &RoomInfo{
		BaseMessage: BaseMessage{msg_type},
		ID:          r.ID,
		Name:        r.Name,
		Size:        r.Size,
		UsersInfo:   make([]*UserInfo, 0),
		HostID:      r.Host.ID, // Populate HostID
	}

	for cli := range r.Clients {
		userInfo := &UserInfo{
			ID:       cli.ID,
			Username: cli.Username,
		}
		roomInfo.UsersInfo = append(roomInfo.UsersInfo, userInfo)
	}

	if strmsg, err := json.Marshal(roomInfo); err == nil {
		log.Printf("Sending room status update: %s", strmsg)
		c.Send <- strmsg
	} else {
		log.Printf("Failed to marshal: %v", err)
	}

	return nil
}

func NotifyUserRoomUpdate(r *Room, c *Client, msg_type string) error {
	var msg UserInfoMessage
	msg.Type = msg_type
	msg.User.ID = r.ID
	msg.User.Username = r.Name

	if strmsg, err := json.Marshal(&msg); err == nil {
		for cli := range r.Clients {
			if cli == c {
				continue
			}
			cli.Send <- strmsg
		}
	} else {
		log.Printf("Failed to marshal: %v", err)
	}

	return nil
}

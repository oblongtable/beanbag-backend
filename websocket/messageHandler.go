package websocket

import (
	"encoding/json"
	"log"
)

func SendEventCallback(c *Client, evtCbMsg *EventCallbackMessage) {
	if strmsg, err := json.Marshal(evtCbMsg); err == nil {
		c.Send <- strmsg
	} else {
		log.Printf("failed to marshal: %v", err)
	}
}

func NotifyUserRoomStatus(r *Room, c *Client) error {
	var msgs RoomInfoMessages
	msgs.Type = MessageRoomStatusUpdate

	roomInfo := RoomInfo{
		ID:        r.ID,
		Name:      r.Name,
		Size:      r.Size,
		UsersInfo: make([]*UserInfo, 1),
	}

	for cli := range r.Clients {
		userInfo := &UserInfo{
			ID:       cli.ID,
			Username: cli.Username,
		}
		roomInfo.UsersInfo = append(roomInfo.UsersInfo, userInfo)
	}

	if strmsg, err := json.Marshal(roomInfo); err == nil {
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

	if strmsg, err := json.Marshal(msg); err == nil {
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

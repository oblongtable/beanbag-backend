package websocket

import (
	"encoding/json"
	"log"
	"sort"
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
		SenderID:    c.ID,
	}

	// Create a slice of client-timestamp pairs
	type clientTimestamp struct {
		Client    *Client
		Timestamp int64 // Assuming timestamp is int64 based on common usage in Go maps
	}
	var clientsWithTimestamps []clientTimestamp
	for cli, ts := range r.Clients {
		clientsWithTimestamps = append(clientsWithTimestamps, clientTimestamp{Client: cli, Timestamp: ts})
	}

	// Sort the slice by timestamp
	sort.SliceStable(clientsWithTimestamps, func(i, j int) bool {
		return clientsWithTimestamps[i].Timestamp < clientsWithTimestamps[j].Timestamp
	})

	// Build the UsersInfo list from the sorted slice
	for _, ct := range clientsWithTimestamps {
		userInfo := &UserInfo{
			ID:       ct.Client.ID,
			Username: ct.Client.Username,
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

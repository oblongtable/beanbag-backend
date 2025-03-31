package websocket

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

func NotifyRoomsStatusAll(c *Client) error {
	var msgs RoomStatusMessages
	msgs.Type = MessageStatusUpdateRoomAll

	for _, r := range c.Wssvr.Rooms {
		var room RoomStatus
		room.ID = r.ID
		room.Name = r.Name
		room.MaxNumUsers = r.MaxNumUsers
		room.IsAlive = true
		msgs.Rooms = append(msgs.Rooms, room)
	}

	if strmsg, err := json.Marshal(msgs); err == nil {
		for cli := range c.Wssvr.Clients {
			cli.Send <- strmsg
		}
	} else {
		log.Printf("failed to marshal: %v", err)
	}

	return nil
}

func NotifyUsersStatusAll(c *Client) error {
	var msgs UserStatusMessages
	msgs.Type = MessageStatusUpdateRoomAll

	for _, r := range c.Wssvr.Rooms {
		var user UserStatus
		user.ID = r.ID
		user.Username = r.Name
		user.IsAlive = true
		msgs.Users = append(msgs.Users, user)
	}

	if strmsg, err := json.Marshal(msgs); err == nil {
		for cli := range c.Wssvr.Clients {
			cli.Send <- strmsg
		}
	} else {
		log.Printf("failed to marshal: %v", err)
	}

	return nil
}

func NotifyRoomsStatus(r *Room, isAlive bool) error {

	var msg RoomStatusMessage
	msg.Type = MessageStatusUpdateRoom
	msg.Room.ID = r.ID
	msg.Room.Name = r.Name
	msg.Room.MaxNumUsers = r.MaxNumUsers
	msg.Room.IsAlive = isAlive

	if strmsg, err := json.Marshal(msg); err == nil {
		for cli := range r.Leader.Wssvr.Clients {
			cli.Send <- strmsg
		}
	} else {
		log.Printf("failed to marshal: %v", err)
	}

	return nil
}

func NotifyClientsStatus(c *Client, isAlive bool) error {

	var msg UserStatusMessage
	msg.Type = MessageStatusUpdateUser
	msg.User.ID = uuid.New().String()
	msg.User.Username = "foo"
	msg.User.IsAlive = isAlive

	if strmsg, err := json.Marshal(msg); err == nil {
		for cli, _ := range c.Wssvr.Clients {
			cli.Send <- strmsg
		}
	} else {
		log.Printf("failed to marshal: %v", err)
	}

	return nil
}

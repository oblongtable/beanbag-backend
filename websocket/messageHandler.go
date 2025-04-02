package websocket

import (
	"encoding/json"
	"log"
)

func NotifyRoomsStatusAll(c *Client) error {
	var msgs RoomStatusMessages
	msgs.Type = MessageStatusUpdateRoomAll

	for _, r := range c.Wssvr.Rooms {
		var room RoomStatus
		room.ID = r.ID
		room.Name = r.Name
		room.Size = r.Size
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

	for cli := range c.Wssvr.Clients {
		var user UserStatus
		user.ID = cli.ID
		user.Username = cli.Username
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
	msg.Room.Size = r.Size
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
	msg.User.ID = c.ID
	msg.User.Username = c.Username
	msg.User.IsAlive = isAlive

	if strmsg, err := json.Marshal(msg); err == nil {
		for cli := range c.Wssvr.Clients {
			cli.Send <- strmsg
		}
	} else {
		log.Printf("failed to marshal: %v", err)
	}

	return nil
}

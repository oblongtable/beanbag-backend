package websocket

import (
	"encoding/json"
	"fmt"
	"log"
)

/*
@Expected JSON:

	{
		"event":<EVENT_TYPE>,
		"payload":{
			"RoomName": <ROOM_NAME>,
			"MaxNumUsers": <MAX_NUM_USERS>
		}
	}
*/
func CreateRoomEventHandler(event Event, c *Client) error {
	fmt.Println(event)
	var crevt CreateRoomEvent
	if err := json.Unmarshal(event.Payload, &crevt); err != nil {
		log.Printf("Create room failed: %v", err)
		return nil
	}
	if crevt.RoomSize > MAX_ROOM_SIZE {
		log.Printf("Create room failed: Room size exceeded (%v)", crevt.RoomSize)
		return nil
	}
	room := NewRoom(crevt.RoomName, crevt.RoomSize, c)
	c.Wssvr.RegisterRoom <- room
	return nil
}

func JoinRoomEventHandler(event Event, c *Client) error {
	fmt.Println(event)
	var jrevt JoinRoomEvent
	if err := json.Unmarshal(event.Payload, &jrevt); err != nil {
		log.Printf("Join room failed: %v", err)
		return nil
	}

	room, ok := c.Wssvr.Rooms[jrevt.RoomID]
	if !ok {
		log.Printf("Join room failed: Room doesn't exist")
		return nil
	}

	c.Wssvr.JoinRoom <- &ClientRoomPair{cli: c, room: room}

	return nil
}

func LeaveRoomEventHandler(event Event, c *Client) error {
	// TODO
	return nil
}

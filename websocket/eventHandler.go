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
func CreateRoom(event Event, c *Client) error {
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

func JoinRoom(event Event, c *Client) error {
	// TODO
	return nil
}

func LeaveRoom(event Event, c *Client) error {
	// TODO
	return nil
}

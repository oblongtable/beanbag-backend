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

	room := NewRoom(crevt.RoomName, crevt.MaxNumUsers, c)
	c.Wssvr.Rooms[room.ID] = room
	NotifyRoomsStatus(room, true)
	return nil
}

package websocket

/*
@Expected JSON:

	{
		"event":<EVENT_TYPE>,
		"payload":{
			"RoomName": <ROOM_NAME>,
			"RoomSize": <MAX_NUM_USERS>
		}
	}
*/
func CreateRoomEventHandler(cliEvt *ClientEvent) error {
	cliEvt.Requester.Wssvr.RegisterRoom <- cliEvt
	return nil
}

func JoinRoomEventHandler(cliEvt *ClientEvent) error {
	cliEvt.Requester.Wssvr.JoinRoom <- cliEvt
	return nil
}

func LeaveRoomEventHandler(cliEvt *ClientEvent) error {
	cliEvt.Requester.Wssvr.LeaveRoom <- cliEvt
	return nil
}

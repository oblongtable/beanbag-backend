package websocket

// func SendMessage(event Event, c *Client) error {
// 	fmt.Println(event)
// 	message, err := event.Payload.MarshalJSON()
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	c.conn.WriteMessage(websocket.TextMessage, message)
// 	return nil
// }

/*
@Expected JSON:

	{
		"event":<EVENT_TYPE>,
		"username":<NAME>,
		"state":<STATUS>
	}
*/

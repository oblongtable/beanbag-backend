package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	pongWait = 60 * time.Second

	pingInterval = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Client struct {
	ID       string // UUID
	Username string // Not sure how to get it upon init, set to dummy for now
	RoomID   string // Room Joined
	Conn     *websocket.Conn

	Wssvr *WebSocServer

	// Buffered channel of outbound messages.
	Send chan []byte
}

func (c Client) String() string {
	return fmt.Sprintf("Client {ID:\"%s\", Username:\"%s\", RoomID:\"%s\"}", c.ID, c.Username, c.RoomID)
}

func NewClient(conn *websocket.Conn, wssvr *WebSocServer) (c *Client) {
	c = &Client{
		ID:       uuid.New().String(),
		Username: "foo",
		Conn:     conn,
		Wssvr:    wssvr,
		Send:     make(chan []byte, 512),
	}

	wssvr.Register <- c
	return c
}

func (c *Client) PongHandler(pongMsg string) error {
	log.Println("pong")
	return c.Conn.SetReadDeadline(time.Now().Add(pongWait))
}

func (c *Client) ReadMessage() {
	defer func() {
		c.Wssvr.Unregister <- c
		fmt.Println("ReadMessage ERR: Client removed")
	}()
	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		fmt.Println(err)
		return
	}
	c.Conn.SetReadLimit(int64(maxMessageSize))
	c.Conn.SetPongHandler(c.PongHandler)

	for {
		_, evtJson, err := c.Conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Unexpected Close error")

			}
			log.Println(err)
			fmt.Println(evtJson)
			break
		}
		fmt.Println(evtJson)
		var evt Event
		if err := json.Unmarshal(evtJson, &evt); err != nil {
			log.Printf("error marshalling event: %v", err)
			break
		}
		if err := c.Wssvr.RouteEvent(evt, c); err != nil {
			log.Println("error handling message: ", err)
			break
		}
	}
}

func (c *Client) WriteMessage() {
	defer func() {
		c.Wssvr.Unregister <- c
		fmt.Println("WriteMessage ERR: Client removed")
	}()
	ticker := time.NewTicker(pingInterval)
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				if err := c.Conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed")
				}
				return
			}
			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("failed to send message: %v", err)
			}
			log.Println("message sent")
		case <-ticker.C:
			log.Println("Ping")
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Println("writeMessage Error:", err)
				return
			}
		}
	}
}

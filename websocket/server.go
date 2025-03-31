package websocket

import (
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type RoomList map[string]*Room

type EventHandlerList map[string]EventHandler

type WebSocServer struct {
	Clients  ClientList
	Rooms    RoomList
	Handlers EventHandlerList

	Register   chan *Client
	Unregister chan *Client
	// broadcast:  make(chan *Message, 5)
	// Mu sync.RWMutex
}

// Upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebSockServer() (wssvr *WebSocServer) {
	wssvr = &WebSocServer{
		Clients:    make(ClientList),
		Rooms:      make(RoomList),
		Handlers:   make(EventHandlerList),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
	wssvr.SetupEventHandlers()
	return wssvr
}

func (wssvr *WebSocServer) SetupEventHandlers() {
	wssvr.Handlers[EventCreateRoom] = CreateRoom
}

func (wssvr *WebSocServer) RouteEvent(evt Event, c *Client) error {
	if handler, ok := wssvr.Handlers[evt.Type]; ok {
		if err := handler(evt, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("there is no such event")
	}
}

func (wssvr *WebSocServer) AddClient(c *Client) {
	wssvr.Clients[c] = true
	NotifyClientsStatus(c, true)
}

func (wssvr *WebSocServer) RemoveClient(c *Client) {
	if _, ok := wssvr.Clients[c]; !ok {
		return
	}

	delete(wssvr.Clients, c)
	c.Conn.Close()
	NotifyClientsStatus(c, false)
}

func (wssvr *WebSocServer) Run() {
	for {
		select {

		case cli := <-wssvr.Register:
			wssvr.AddClient(cli)

		case cli := <-wssvr.Unregister:
			wssvr.RemoveClient(cli)
		}
	}
}

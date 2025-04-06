package websocket

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type RoomList map[string]*Room

type ClientRoomPair struct {
	cli  *Client
	room *Room
}

type EventHandlerList map[string]EventHandler

type WebSocServer struct {
	Clients  ClientList
	Rooms    RoomList
	Handlers EventHandlerList

	Register       chan *Client
	Unregister     chan *Client
	RegisterRoom   chan *Room
	UnregisterRoom chan *Room
	JoinRoom       chan *ClientRoomPair
	LeaveRoom      chan *ClientRoomPair
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
		Clients:        make(ClientList),
		Rooms:          make(RoomList),
		Handlers:       make(EventHandlerList),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		RegisterRoom:   make(chan *Room),
		UnregisterRoom: make(chan *Room),
		JoinRoom:       make(chan *ClientRoomPair),
		LeaveRoom:      make(chan *ClientRoomPair),
	}
	wssvr.SetupEventHandlers()
	go wssvr.Run()
	return wssvr
}

func (wssvr *WebSocServer) SetupEventHandlers() {
	wssvr.Handlers[EventCreateRoom] = CreateRoomEventHandler
	wssvr.Handlers[EventJoinRoom] = JoinRoomEventHandler
	wssvr.Handlers[EventLeaveRoom] = LeaveRoomEventHandler
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

	if room, ok := wssvr.Rooms[c.RoomID]; ok {
		wssvr.RemoveRoom(room)
	}

	delete(wssvr.Clients, c)
	c.Conn.Close()
	NotifyClientsStatus(c, false)
}

func (wssvr *WebSocServer) AddRoom(room *Room) {
	cli := room.Leader
	cli.RoomID = room.ID
	cli.Wssvr.Rooms[room.ID] = room
	NotifyRoomsStatus(room, true)
}

func (wssvr *WebSocServer) RemoveRoom(room *Room) {
	cli := room.Leader
	cli.RoomID = ""
	delete(cli.Wssvr.Rooms, room.ID)
	NotifyRoomsStatus(room, false)
}

func (wssvr *WebSocServer) JoinRoomF(crp *ClientRoomPair) {
	cli := crp.cli
	room := crp.room

	if _, ok := room.Clients[cli]; ok {
		log.Printf("Join room failed: You are already in the room")
		return
	} else if len(room.Clients) >= room.Size {
		log.Printf("Join room failed: Room is full")
		return
	}

	cli.RoomID = room.ID
	room.Clients[cli] = true
}

func (wssvr *WebSocServer) LeaveRoomF(crp *ClientRoomPair) {
	cli := crp.cli
	room := crp.room

	if _, ok := room.Clients[cli]; !ok {
		log.Printf("Leave room failed: You are not in the room")
		return
	}

	cli.RoomID = ""
	room.Clients[cli] = false
}

func (wssvr *WebSocServer) Run() {
	for {
		select {

		case cli := <-wssvr.Register:
			wssvr.AddClient(cli)

		case cli := <-wssvr.Unregister:
			wssvr.RemoveClient(cli)

		case room := <-wssvr.RegisterRoom:
			wssvr.AddRoom(room)

		case room := <-wssvr.UnregisterRoom:
			wssvr.RemoveRoom(room)

		case crp := <-wssvr.JoinRoom:
			wssvr.JoinRoomF(crp)

		case crp := <-wssvr.LeaveRoom:
			wssvr.LeaveRoomF(crp)
		}
	}
}

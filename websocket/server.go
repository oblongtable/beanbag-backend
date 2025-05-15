package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]int64

type RoomList map[string]*Room

type EventHandler func(cliEvt *ClientEvent) error

type EventHandlerList map[string]EventHandler

type ClientEvent struct {
	Requester *Client
	EventInfo *Event
}

type WebSocServer struct {
	Clients  ClientList
	Rooms    RoomList
	Handlers EventHandlerList

	Register       chan *Client
	Unregister     chan *Client
	RegisterRoom   chan *ClientEvent
	UnregisterRoom chan *Room
	JoinRoom       chan *ClientEvent
	LeaveRoom      chan *ClientEvent
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
		RegisterRoom:   make(chan *ClientEvent),
		UnregisterRoom: make(chan *Room),
		JoinRoom:       make(chan *ClientEvent),
		LeaveRoom:      make(chan *ClientEvent),
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

func (wssvr *WebSocServer) RouteEvent(evt *Event, c *Client) error {
	var cliEvt ClientEvent
	cliEvt.Requester = c
	cliEvt.EventInfo = evt

	if handler, ok := wssvr.Handlers[evt.Type]; ok {
		if err := handler(&cliEvt); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("there is no such event")
	}
}

func (wssvr *WebSocServer) AddClient(c *Client) {

	wssvr.Clients[c] = -1
}

func (wssvr *WebSocServer) RemoveClient(c *Client) {
	if _, ok := wssvr.Clients[c]; !ok {
		return
	}

	if room, ok := wssvr.Rooms[c.RoomID]; ok {
		room.Leave <- c
	}

	delete(wssvr.Clients, c)
	c.Conn.Close()
}

func (wssvr *WebSocServer) AddRoom(cliEvt *ClientEvent) {
	var crevt CreateRoomEvent
	var msg string
	var roomInfo RoomInfo

	isSuccess := false
	cli := cliEvt.Requester
	jsonRaw := cliEvt.EventInfo.Payload
	if err := json.Unmarshal(jsonRaw, &crevt); err != nil {
		msg = fmt.Sprintf("Create room failed: %v", err)
		log.Println(msg)

	} else if crevt.RoomSize > MAX_ROOM_SIZE {
		msg = fmt.Sprintf("Create room failed: Room size cannot be larger than %d", MAX_ROOM_SIZE)
		log.Println(msg)

	} else {
		cli.Username = crevt.UserName
		room := NewRoom(crevt.RoomName, crevt.RoomSize, cli)
		if room == nil {
			msg = "Create room failed: Server overloaded, please try again later."
		} else {
			isSuccess = true

			cli.RoomID = room.ID
			cli.Wssvr.Rooms[room.ID] = room

			roomInfo.ID = room.ID
			roomInfo.Name = room.Name
			roomInfo.Size = room.Size
			roomInfo.UsersInfo = make([]*UserInfo, 0)
			roomInfo.UsersInfo = append(roomInfo.UsersInfo, &UserInfo{
				ID:       cli.ID,
				Username: cli.Username,
			})
			roomInfo.HostID = cli.ID
			roomInfo.SenderID = cli.ID

			msg = "Create room Success"
			log.Println(msg)
		}
	}

	// Message callback
	SendEventCallback(cli, MessageCreateRoom, isSuccess, msg, &roomInfo)
}

func (wssvr *WebSocServer) RemoveRoom(room *Room) {
	cli := room.Host
	cli.RoomID = ""
	delete(cli.Wssvr.Rooms, room.ID)
}

func (wssvr *WebSocServer) JoinRoomF(cliEvt *ClientEvent) {
	var jrevt JoinRoomEvent
	var msg string
	var roomInfo RoomInfo

	isSuccess := true
	cli := cliEvt.Requester
	jsonRaw := cliEvt.EventInfo.Payload
	if err := json.Unmarshal(jsonRaw, &jrevt); err != nil {
		isSuccess = false
		msg = fmt.Sprintf("Join room failed: %v", err)
		log.Println(msg)

	} else if len(cli.RoomID) > 0 {
		isSuccess = false
		msg = "Join room failed: You have already joined a room"
		log.Println(msg)

	} else if room, ok := wssvr.Rooms[jrevt.RoomID]; !ok {
		isSuccess = false
		msg = "Join room failed: Room not found"
		log.Println(msg)

	} else if len(room.Clients) >= room.Size+1 {
		isSuccess = false
		msg = "Join room failed: Room is full"
		log.Println(msg)

	} else {
		cli.Username = jrevt.Name

		msg = "Join room Success"
		log.Println(msg)

		roomInfo.ID = room.ID
		roomInfo.Name = room.Name
		roomInfo.Size = room.Size

		var userInfo []*UserInfo

		// Populate with players already in the lobby
		for roomCli := range room.Clients {
			userInfo = append(userInfo, &UserInfo{
				ID:       roomCli.ID,
				Username: roomCli.Username,
			})
		}

		// Add myself in
		userInfo = append(userInfo, &UserInfo{
			ID:       cli.ID,
			Username: cli.Username,
		})
		roomInfo.UsersInfo = userInfo
		roomInfo.HostID = room.Host.ID
		roomInfo.SenderID = cli.ID

		room.Join <- cli

	}

	// Message callback
	SendEventCallback(cli, MessageJoinRoom, isSuccess, msg, &roomInfo)
}

func (wssvr *WebSocServer) LeaveRoomF(cliEvt *ClientEvent) {
	var jrevt LeaveRoomEvent
	var msg string
	var roomInfo RoomInfo

	isSuccess := true
	cli := cliEvt.Requester
	jsonRaw := cliEvt.EventInfo.Payload
	if err := json.Unmarshal(jsonRaw, &jrevt); err != nil {
		isSuccess = false
		msg = fmt.Sprintf("Leave room failed: %v", err)
		log.Println(msg)

	} else if len(cli.RoomID) < 1 {
		isSuccess = false
		msg = "Leave room failed: You haven't joined a room"
		log.Println(msg)

	} else if room, ok := wssvr.Rooms[jrevt.RoomID]; !ok {
		isSuccess = false
		msg = "Leave room failed: Room not found"
		log.Println(msg)

	} else {
		msg = "Leave room Success"
		log.Println(msg)

		room.Leave <- cli
	}

	// Message callback
	SendEventCallback(cli, MessageLeaveRoom, isSuccess, msg, &roomInfo)
}

func (wssvr *WebSocServer) Run() {
	for {
		select {

		case cli := <-wssvr.Register:
			wssvr.AddClient(cli)

		case cli := <-wssvr.Unregister:
			wssvr.RemoveClient(cli)

		case cliEvt := <-wssvr.RegisterRoom:
			wssvr.AddRoom(cliEvt)

		case room := <-wssvr.UnregisterRoom:
			wssvr.RemoveRoom(room)

		case cliEvt := <-wssvr.JoinRoom:
			wssvr.JoinRoomF(cliEvt)

		case cliEvt := <-wssvr.LeaveRoom:
			wssvr.LeaveRoomF(cliEvt)
		}
	}
}

package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/oblongtable/beanbag-backend/internal/game"
)

type ClientList map[*Client]bool

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
	Games    *game.GameService // Add GameService

	Register       chan *Client
	Unregister     chan *Client
	RegisterRoom   chan *ClientEvent
	UnregisterRoom chan *Room
	JoinRoom       chan *ClientEvent
	LeaveRoom      chan *ClientEvent

	StartQuiz    chan *ClientEvent
	ForwardQuiz  chan *ClientEvent // New channel for advancing the quiz
	SubmitAnswer chan *ClientEvent // New channel for submitting answers
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
		Games:          game.NewService(), // Initialize GameService
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		RegisterRoom:   make(chan *ClientEvent),
		UnregisterRoom: make(chan *Room),
		JoinRoom:       make(chan *ClientEvent),
		LeaveRoom:      make(chan *ClientEvent),

		StartQuiz:    make(chan *ClientEvent),
		ForwardQuiz:  make(chan *ClientEvent), // Initialize the new channel
		SubmitAnswer: make(chan *ClientEvent), // Initialize the new channel
	}
	wssvr.SetupEventHandlers()
	go wssvr.Run()
	return wssvr
}

func (wssvr *WebSocServer) SetupEventHandlers() {
	wssvr.Handlers[EventCreateRoom] = CreateRoomEventHandler
	wssvr.Handlers[EventJoinRoom] = JoinRoomEventHandler
	wssvr.Handlers[EventLeaveRoom] = LeaveRoomEventHandler
	wssvr.Handlers[EventStartQuiz] = StartQuizEventHandler
	wssvr.Handlers[EventForwardQuiz] = ForwardQuizEventHandler   // Register the new handler
	wssvr.Handlers[EventSubmitAnswer] = SubmitAnswerEventHandler // Register the new handler
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

	wssvr.Clients[c] = true
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
				Username: cli.Username,
				Role:     RoleCreator.String(),
			})
			roomInfo.SenderID = cli.ID
			roomInfo.IsHost = true // Set IsHost to true for the creator

			msg = "Create room Success"
			log.Println(msg)
		}
	}

	// Message callback
	SendEventCallback(cli, MessageCreateRoom, isSuccess, msg, &roomInfo)
}

func (wssvr *WebSocServer) RemoveRoom(room *Room) {
	cli := room.Creator
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

	} else if len(room.Participants) >= room.Size+1 {
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
		roomInfo.UsersInfo = room.GetSortedUserInfo()
		roomInfo.SenderID = cli.ID

		// Add myself to user info
		var clientRole Role
		if (room.ParticipantCount()) == 1 {
			clientRole = RoleHost
		} else {
			clientRole = RolePlayer
		}

		roomInfo.UsersInfo = append(roomInfo.UsersInfo, &UserInfo{
			Username: cli.Username,
			Role:     clientRole.String(),
		})

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

func (wssvr *WebSocServer) StartQuizF(cliEvt *ClientEvent) {
	var jrevt StartQuizEvent
	var msg string

	isSuccess := true
	cli := cliEvt.Requester
	jsonRaw := cliEvt.EventInfo.Payload
	if err := json.Unmarshal(jsonRaw, &jrevt); err != nil {
		isSuccess = false
		msg = fmt.Sprintf("Start Quiz failed: %v", err)
		log.Println(msg)

	} else {
		log.Println("Start Quiz Event received")
		log.Println(msg)
		room, ok := wssvr.Rooms[cli.RoomID]
		if !ok {
			isSuccess = false
			msg = "Start Quiz failed: Room not found"
			log.Println(msg)
		}
		if room.Host.ID != cli.ID {
			isSuccess = false
			msg = "Start Quiz failed: Only the room creator can start the quiz"
			log.Println(msg)
		} else {
			// Define the broadcast function for this specific room
			broadcastFunc := func(msgType string, payload interface{}) {
				evt := &Event{
					Type: msgType,
				}

				jsonPayload, err := json.Marshal(payload)
				if err != nil {
					log.Printf("Error marshaling broadcast payload for type %s: %v", msgType, err)
					return
				}
				evt.Payload = jsonPayload

				strmsg, err := json.Marshal(evt)
				if err != nil {
					log.Printf("Error marshaling broadcast event for type %s: %v", msgType, err)
					return
				}

				for _, participant := range room.Participants {
					select {
					case participant.Client.Send <- strmsg:
						// Message sent successfully
					default:
						log.Printf("Failed to send message to client %s in room %s", participant.Client.ID, room.ID)
					}
				}
			}

			// Prepare initial player info from room participants
			initialPlayers := make([]game.InitialPlayerInfo, 0, len(room.Participants))
			for _, pDetail := range room.Participants {
				initialPlayers = append(initialPlayers, game.InitialPlayerInfo{
					ID:       pDetail.Client.ID,
					Username: pDetail.Client.Username,
				})
			}

			// Call CreateGame here, passing the broadcast function and initial players
			game, err := wssvr.Games.CreateGame(room.ID, room.Creator.ID, room.Host.ID, broadcastFunc, initialPlayers)
			if err != nil {
				isSuccess = false
				msg = fmt.Sprintf("Start Quiz failed: Failed to create game: %v", err)
				log.Println(msg)
			} else {

				SendEventCallback(cli, MessageQuizStart, isSuccess, msg, &RoomInfo{})

				// Call StartGame here
				err = wssvr.Games.StartGame(game.ID, room.Host.ID)
				if err != nil {
					isSuccess = false
					msg = fmt.Sprintf("Start Quiz failed: Failed to start game: %v", err)
					log.Println(msg)
				} else {
					msg = "Start Quiz Success"
					log.Println(msg)
				}
			}
		}

	}
}

func (wssvr *WebSocServer) ForwardQuizF(cliEvt *ClientEvent) {
	var msg string
	isSuccess := true
	cli := cliEvt.Requester

	room, ok := wssvr.Rooms[cli.RoomID]
	if !ok {
		isSuccess = false
		msg = "Forward Quiz failed: Room not found"
		log.Println(msg)
	} else if room.Host.ID != cli.ID {
		isSuccess = false
		msg = "Forward Quiz failed: Only the room host can advance the quiz"
		log.Println(msg)
	} else {
		log.Println("Forward Quiz Event received")
		err := wssvr.Games.NextAction(room.ID, cli.ID) // Use the existing NextAction method
		if err != nil {
			isSuccess = false
			msg = fmt.Sprintf("Forward Quiz failed: %v", err)
			log.Println(msg)
		} else {
			msg = "Forward Quiz Success"
			log.Println(msg)
		}
	}
	// Define MessageQuizForward in message.go and use it here
	SendEventCallback(cli, MessageQuizForward, isSuccess, msg, &BaseMessage{})

}

func (wssvr *WebSocServer) SubmitAnswerF(cliEvt *ClientEvent) {
	var saEvt SubmitAnswerEvent
	var msg string
	isSuccess := true
	cli := cliEvt.Requester
	jsonRaw := cliEvt.EventInfo.Payload

	if err := json.Unmarshal(jsonRaw, &saEvt); err != nil {
		isSuccess = false
		msg = fmt.Sprintf("Submit Answer failed: %v", err)
		log.Println(msg)
	} else {
		log.Printf("Submit Answer Event received for answer index: %d", saEvt.AnswerIndex)
		room, ok := cli.Wssvr.Rooms[cli.RoomID]
		if !ok {
			isSuccess = false
			msg = "Submit Answer failed: Room not found"
			log.Println(msg)
		} else {
			err := wssvr.Games.HandleAnswer(room.ID, cli.ID, saEvt.AnswerIndex)
			if err != nil {
				isSuccess = false
				msg = fmt.Sprintf("Submit Answer failed: %v", err)
				log.Println(msg)
			} else {
				msg = "Submit Answer Success"
				log.Println(msg)
			}
		}
	}
	SendEventCallback(cli, MessageSubmitAnswer, isSuccess, msg, &BaseMessage{})
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

		case cliEvt := <-wssvr.StartQuiz:
			wssvr.StartQuizF(cliEvt)

		case cliEvt := <-wssvr.ForwardQuiz: // Handle the new ForwardQuiz event
			wssvr.ForwardQuizF(cliEvt)

		case cliEvt := <-wssvr.SubmitAnswer: // Handle the new SubmitAnswer event
			wssvr.SubmitAnswerF(cliEvt)
		}
	}
}

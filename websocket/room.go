package websocket

import "github.com/google/uuid"

type Room struct {
	ID          string
	Name        string
	MaxNumUsers int
	Leader      *Client
	Clients     ClientList
	Register    chan *Client
	Unregister  chan *Client
	Boardcast   chan *Client
}

func NewRoom(name string, maxNumUsers int, leader *Client) (r *Room) {
	r = &Room{
		ID:          uuid.New().String(),
		Name:        name,
		MaxNumUsers: maxNumUsers,
		Leader:      leader,
		Clients:     make(ClientList, maxNumUsers),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Boardcast:   make(chan *Client),
	}

	return r
}

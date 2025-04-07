package websocket

import (
	"fmt"

	"github.com/google/uuid"
)

const MAX_ROOM_SIZE = 8

type Room struct {
	ID      string // UUID
	Name    string // Customised
	Size    int    // Min: 1, Max: 8
	Leader  *Client
	Clients ClientList
	Join    chan *Client
	Leave   chan *Client
	Disband chan *Client
}

func (r Room) String() string {
	return fmt.Sprintf("Room {ID:\"%s\", Name:\"%s\", Size:%d, Leader:%s, Clients:%v}",
		r.ID, r.Name, r.Size, r.Leader, r.Clients)
}

func NewRoom(name string, size int, leader *Client) (r *Room) {
	r = &Room{
		ID:      uuid.New().String(),
		Name:    name,
		Size:    size,
		Leader:  leader,
		Clients: make(ClientList, size),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Disband: make(chan *Client),
	}
	go r.Run()
	return r
}

func (r *Room) JoinRoom(c *Client) {
	// Check for max room size
	if len(r.Clients) >= r.Size {
		return
	}

	// Check if user is already in the room
	if _, ok := r.Clients[c]; ok {
		return
	}
	r.Clients[c] = true
}

func (r *Room) LeaveRoom(c *Client) {
	// Check if user is not in the room
	if _, ok := r.Clients[c]; !ok {
		return
	}
	// Check if user is the room leader
	if c == r.Leader {
		r.Disband <- c
		return
	}
	r.Clients[c] = false
}

func (r *Room) DisbandRoom(c *Client) {
	// TODO
}

func (r *Room) Run() {
	defer func() {
		r.Leader.Wssvr.UnregisterRoom <- r
	}()
	// Room leader must join the room
	r.Join <- r.Leader

	for {
		select {
		case cli := <-r.Join:
			r.JoinRoom(cli)

		case cli := <-r.Leave:
			r.LeaveRoom(cli)

		case cli := <-r.Disband:
			r.DisbandRoom(cli)
			return
		}
	}
}

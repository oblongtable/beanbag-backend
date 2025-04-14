package websocket

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

const MAX_ROOM_SIZE = 8

type Room struct {
	ID      string // UUID
	Name    string // Customised
	Size    int    // Min: 1, Max: 8
	IsAlive bool
	Host    *Client
	Clients ClientList
	Join    chan *Client
	Leave   chan *Client
}

func (r Room) String() string {
	return fmt.Sprintf("Room {ID:\"%s\", Name:\"%s\", Size:%d, Host:%s, Clients:%v}",
		r.ID, r.Name, r.Size, r.Host, r.Clients)
}

func NewRoom(name string, size int, host *Client) (r *Room) {
	r = &Room{
		ID:      uuid.New().String(),
		Name:    name,
		Size:    size,
		Host:    host,
		Clients: make(ClientList, size),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
	}
	go r.Run()
	return r
}

func (r *Room) JoinRoom(c *Client) {
	c.RoomID = r.ID
	r.Clients[c] = true

	NotifyUserRoomUpdate(r, c, MessageJoinRoom)
	NotifyUserRoomStatus(r, c)
}

func (r *Room) LeaveRoom(c *Client) {
	c.RoomID = ""
	delete(r.Clients, c)

	// Check if user is the room host
	if c == r.Host {
		r.IsAlive = false
	} else {
		NotifyUserRoomUpdate(r, c, MessageLeaveRoom)
	}
}

func (r *Room) Run() {
	defer func() {
		r.Host.Wssvr.UnregisterRoom <- r
	}()

	r.IsAlive = true

	for r.IsAlive {
		select {
		case cli := <-r.Join:
			r.JoinRoom(cli)

		case cli := <-r.Leave:
			r.LeaveRoom(cli)

		}
	}
	log.Printf("Room %s removed.", r.ID)
}

package websocket

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

const MAX_ROOM_SIZE = 20

type Room struct {
	ID      string // UUID
	Name    string // Customised
	Size    int    // Min: 1, Max: 20
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
		ID:      GenerateRandomCode(4),
		Name:    name,
		Size:    size,
		Host:    host,
		Clients: make(ClientList, size+1),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
	}

	// Add the host to the clients list immediately
	r.Clients[host] = time.Now().UnixMilli()

	// Re-generate room ID if already exist such ID for 5 times
	exist := true
	for i := 0; i < 5; i++ {
		if _, ok := host.Wssvr.Rooms[r.ID]; !ok {
			exist = false
			break
		}
		r.ID = GenerateRandomCode(4)
	}

	// Return NULL pointer if it still fails to generate unique code
	if exist {
		return nil
	}

	go r.Run()
	return r
}

func (r *Room) JoinRoom(c *Client) {
	c.RoomID = r.ID
	r.Clients[c] = time.Now().UnixMilli()

	// Notify all clients in the room of the updated user list
	for client := range r.Clients {
		log.Printf("Notify %s of room status update", client.Username)
		NotifyUserRoomStatus(r, client, MessageRoomStatusUpdate)
	}
}

func (r *Room) LeaveRoom(c *Client) {
	c.RoomID = ""
	delete(r.Clients, c)

	// Check if user is the room host
	if c == r.Host {
		r.IsAlive = false
	} else {
		// Notify all clients in the room of the updated user list
		for client := range r.Clients {
			log.Printf("Notify %s of room status update", client.Username)
			NotifyUserRoomStatus(r, client, MessageRoomStatusUpdate)
		}
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

func GenerateRandomCode(length int) string {

	const asciiBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	b := make([]byte, length)

	for i := range b {
		randomIndex := r.Intn(len(asciiBytes))
		b[i] = asciiBytes[randomIndex]
	}

	return string(b)
}

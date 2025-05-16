package websocket

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"sync"
	"time"
)

const MAX_ROOM_SIZE = 20

type ParticipantsDetail struct {
	client   *Client
	role     Role
	joinedAt time.Time
}

type Room struct {
	ID           string // UUID
	Name         string // Customised
	Size         int    // Min: 1, Max: 20
	IsAlive      bool
	Creator      *Client
	Participants map[string]ParticipantsDetail
	Join         chan *Client
	Leave        chan *Client
	mu           sync.Mutex
}

func (r *Room) String() string {
	return fmt.Sprintf("Room {ID:\"%s\", Name:\"%s\", Size:%d, Creator:%s, Participants:%v}",
		r.ID, r.Name, r.Size, r.Creator, r.Participants)
}

func NewRoom(name string, size int, creator *Client) (r *Room) {
	r = &Room{
		ID:           GenerateRandomCode(4),
		Name:         name,
		Size:         size,
		Creator:      creator,
		Participants: make(map[string]ParticipantsDetail),
		Join:         make(chan *Client),
		Leave:        make(chan *Client),
		mu:           sync.Mutex{},
	}

	// Add the host to the clients list immediately
	r.Participants[creator.ID] = ParticipantsDetail{
		client:   creator,
		role:     RoleCreator,
		joinedAt: time.Now(),
	}

	// Re-generate room ID if already exist such ID for 5 times
	exist := true
	for range 5 {
		if _, ok := creator.Wssvr.Rooms[r.ID]; !ok {
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

	r.mu.Lock()

	c.RoomID = r.ID

	// Only the creator is in the room so I am the host
	if len(r.Participants) == 1 {

		r.Participants[c.ID] = ParticipantsDetail{
			client:   c,
			role:     RoleHost,
			joinedAt: time.Now(),
		}

	} else {

		r.Participants[c.ID] = ParticipantsDetail{
			client:   c,
			role:     RolePlayer,
			joinedAt: time.Now(),
		}

	}

	r.mu.Unlock()

	// Notify all clients in the room of the updated user list
	log.Printf("Length of participants: %d", len(r.Participants))
	for id, pd := range r.Participants {
		log.Printf("Notify %s (%s) of room status update", pd.client.Username, id)
		NotifyUserRoomStatus(r, pd.client, r.GetSortedUserInfo(), MessageRoomStatusUpdate)
	}
}

func (r *Room) LeaveRoom(c *Client) {
	r.mu.Lock()

	c.RoomID = ""
	delete(r.Participants, c.ID)

	r.mu.Unlock()

	// Check if user is the room host
	if c == r.Creator {
		r.IsAlive = false

		// Notify that the room has died

	} else {
		// Notify all clients in the room of the updated user list
		for id, pd := range r.Participants {
			log.Printf("Notify %s (%s) of room status update", pd.client.Username, id)
			NotifyUserRoomStatus(r, pd.client, r.GetSortedUserInfo(), MessageRoomStatusUpdate)
		}
	}
}

func (r *Room) Run() {
	defer func() {
		r.Creator.Wssvr.UnregisterRoom <- r
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

func (r *Room) ParticipantCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	count := len(r.Participants)
	return count
}

func (r *Room) GetSortedUserInfo() []*UserInfo {
	r.mu.Lock()
	defer r.mu.Unlock()

	var clients []ParticipantsDetail
	for _, pd := range r.Participants {
		clients = append(clients, pd)
	}

	sort.Slice(clients, func(i, j int) bool {
		return clients[i].joinedAt.Before(clients[j].joinedAt)
	})

	userInfo := make([]*UserInfo, 0)
	for _, pd := range clients {
		userInfo = append(userInfo, &UserInfo{
			ID:       pd.client.ID,
			Username: pd.client.Username,
			Role:     pd.role.String(),
		})
	}

	return userInfo
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

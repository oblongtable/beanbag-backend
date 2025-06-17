package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"sync"
	"time"
)

const MAX_ROOM_SIZE = 20

type ParticipantsDetail struct {
	Client   *Client
	Role     Role
	joinedAt time.Time
}

type Room struct {
	ID           string // UUID
	Name         string // Customised
	Size         int    // Min: 1, Max: 20
	IsAlive      bool
	Creator      *Client
	Host         *Client
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
		Client:   creator,
		Role:     RoleCreator,
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

		r.Host = c
		r.Participants[c.ID] = ParticipantsDetail{
			Client:   c,
			Role:     RoleHost,
			joinedAt: time.Now(),
		}

	} else {

		r.Participants[c.ID] = ParticipantsDetail{
			Client:   c,
			Role:     RolePlayer,
			joinedAt: time.Now(),
		}

	}

	r.mu.Unlock()

	// Notify all clients in the room of the updated user list
	log.Printf("Length of participants: %d", len(r.Participants))
	for id, pd := range r.Participants {
		log.Printf("Notify %s (%s) of room status update", pd.Client.Username, id)
		NotifyUserRoomStatus(r, pd.Client, r.GetSortedUserInfo(), MessageRoomStatusUpdate)
	}
}

func (r *Room) LeaveRoom(c *Client) {
	r.mu.Lock()

	leavingParticipantDetail, exists := r.Participants[c.ID]
	log.Printf(" %s", c.ID)
	isLeavingHost := false
	log.Printf("Participant %s (%s) is leaving room %s %s", c.Username, c.ID, r.ID, leavingParticipantDetail.Role)
	if exists && leavingParticipantDetail.Role == RoleHost {
		isLeavingHost = true
	}

	c.RoomID = ""
	delete(r.Participants, c.ID)

	r.mu.Unlock()

	// If the leaving participant is the creator, shut down the room
	if c == r.Creator {
		r.IsAlive = false
		// Notify all clients in the room that the room is shutting down
		shutdownMessage := BaseMessage{Type: MessageRoomShutdown}
		message, _ := json.Marshal(shutdownMessage)
		for _, pd := range r.Participants {
			log.Printf("Notify %s (%s) of room closure", pd.Client.Username, pd.Client.ID)
			select {
			case pd.Client.Send <- message:
			default:
				close(pd.Client.Send)
				// It's safe to delete here as the loop is finishing for this participant
				delete(r.Participants, pd.Client.ID)
			}
		}
	} else if isLeavingHost { // If the leaving participant is the host (and not the creator)
		// If there are other participants (besides the creator), transfer host role
		if len(r.Participants) > 0 {
			var nextHostID string
			var earliestJoinTime time.Time

			// Initialize with a time far in the future
			earliestJoinTime = time.Now().Add(24 * time.Hour)

			for id, pd := range r.Participants {
				// Exclude the leaving host and the creator
				if pd.Client.ID != c.ID && pd.Client.ID != r.Creator.ID {
					if pd.joinedAt.Before(earliestJoinTime) {
						earliestJoinTime = pd.joinedAt
						nextHostID = id
					}
				}
			}

			if nextHostID != "" {
				// Assign the host role to the next participant
				pd := r.Participants[nextHostID]
				pd.Role = RoleHost
				r.Participants[nextHostID] = pd
				r.Host = pd.Client

				log.Printf("Host role transferred to %s (%s) in room %s", pd.Client.Username, nextHostID, r.ID)

				// Notify all remaining clients of the updated user list and role change
				for _, pd := range r.Participants {
					log.Printf("Notify %s (%s) of room status update with new host", pd.Client.Username, pd.Client.ID)
					NotifyUserRoomStatus(r, pd.Client, r.GetSortedUserInfo(), MessageRoomStatusUpdate)
				}
			} else {
				// No other participants besides the creator, room stays open, notify creator
				log.Printf("Host left, no other participants to transfer role to. Notifying creator %s (%s) in room %s", r.Creator.Username, r.Creator.ID, r.ID)
				NotifyUserRoomStatus(r, r.Creator, r.GetSortedUserInfo(), MessageRoomStatusUpdate)
			}
		} else {
			// Should not happen if the creator is still in the room, but as a fallback
			// No other participants, shut down the room
			r.IsAlive = false
			// Notify all clients in the room that the room is shutting down
			shutdownMessage := BaseMessage{Type: MessageRoomShutdown}
			message, _ := json.Marshal(shutdownMessage)
			for _, pd := range r.Participants {
				log.Printf("Notify %s (%s) of room closure", pd.Client.Username, pd.Client.ID)
				select {
				case pd.Client.Send <- message:
				default:
					close(pd.Client.Send)
					// It's safe to delete here as the loop is finishing for this participant
					delete(r.Participants, pd.Client.ID)
				}
			}
		}
	} else { // User is not the creator and not the host, just remove them and notify others of updated user list
		// Notify all clients in the room of the updated user list
		for _, pd := range r.Participants {
			log.Printf("Notify %s (%s) of room status update", pd.Client.Username, pd.Client.ID)
			NotifyUserRoomStatus(r, pd.Client, r.GetSortedUserInfo(), MessageRoomStatusUpdate)
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
			Username: pd.Client.Username,
			Role:     pd.Role.String(),
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

func Broadcast(message []byte) {

}

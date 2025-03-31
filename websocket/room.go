package websocket

type Room struct {
	ID         string
	Name       string
	Register   chan *Client
	Unregister chan *Client
	Boardcast  chan *Client
}

func NewRoom(name string) (r *Room) {
	r = &Room{
		Name:       name,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Boardcast:  make(chan *Client),
	}

	return r
}

package websocket

// Role represents the role of a client in a room.
type Role string

const (
	// RoleCreator signifies the client who created the room.
	RoleCreator Role = "CREATOR"
	// RoleHost signifies a client with hosting privileges.
	RoleHost Role = "HOST"
	// RolePlayer signifies a regular player in the room.
	RolePlayer Role = "PLAYER"
)

// IsValid checks if the role is one of the predefined valid roles.
func (r Role) IsValid() bool {
	switch r {
	case RoleCreator, RoleHost, RolePlayer:
		return true
	}
	return false
}

func (r Role) String() string {
	return string(r)
}
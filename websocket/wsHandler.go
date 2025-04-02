package websocket

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Endpoint handler
func (wssvr *WebSocServer) ServeWs(ctx *gin.Context) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	c := NewClient(conn, wssvr)

	// Update client the current users and rooms' statuses
	// (Might cause performance issue)
	NotifyRoomsStatusAll(c)
	NotifyUsersStatusAll(c)

	go c.ReadMessage()
	go c.WriteMessage()
}

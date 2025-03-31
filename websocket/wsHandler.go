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

	go c.ReadMessage()
	go c.WriteMessage()
}

package handlers

import (
	"time"

	// Import local packages
	"github.com/ChatGPT-Hackers/go-server/types"
	"github.com/ChatGPT-Hackers/go-server/utils"

	"github.com/gin-gonic/gin"
)

// // # Client routes
func ClientRegister(c *gin.Context) {
	// Make websocket connection
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	// Add connection to the pool
	connection := &types.Connection{
		Id:              utils.GenerateId(),
		Ws:              ws,
		Heartbeat:       time.Now().Unix(),
		LastMessageTime: time.Now().Unix(),
	}
	connectionsMu.Lock()
	connections = append(connections, connection)
	connectionsMu.Unlock()
	// Send connection id to the client
	err = connection.Ws.WriteJSON(types.Message{
		Id:      connection.Id,
		Message: "Connection id",
	})
	if err != nil {
		return
	}
	// Close the connection when it is no longer needed
	defer connection.Ws.Close()
}

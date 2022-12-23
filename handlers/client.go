package handlers

import (
	"time"

	// Import local packages
	"github.com/ChatGPT-Hackers/ChatGPT-API-server/types"
	"github.com/ChatGPT-Hackers/ChatGPT-API-server/utils"

	"github.com/gin-gonic/gin"
)

// // # Client routes
func Client_register(c *gin.Context) {
	// Make websocket connection
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	// Generate connection id
	id := utils.GenerateId()
	// Send connection id
	err = ws.WriteJSON(types.Message{
		Id:      id,
		Message: "Connection id",
	})
	if err != nil {
		return
	}
	// Wait for client to send connection id
	for {
		// Read message
		var message types.Message
		err = ws.ReadJSON(&message)
		if err != nil {
			return
		}
		// Check if the message is the connection id
		if message.Id == id {
			break
		}
	}
	// Add connection to the pool
	connection := &types.Connection{
		Id:              id,
		Ws:              ws,
		LastMessageTime: time.Now(),
		Heartbeat:       time.Now(),
	}
	connectionPool.Set(connection)
	// Debug
	println("New connection:", connection.Id)
}

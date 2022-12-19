package handlers

import (
	"github.com/ChatGPT-Hackers/go-server/types"
	"github.com/ChatGPT-Hackers/go-server/utils"
	"github.com/gin-gonic/gin"
)

// // # API routes
func ApiAsk(c *gin.Context) {
	// Get request
	var request types.ChatGptRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request",
		})
		return
	}
	// Get connection with oldest last message time
	var connection *types.Connection
	connectionsMu.RLock()
	for _, conn := range connections {
		if connection == nil || conn.LastMessageTime < connection.LastMessageTime {
			connection = conn
		}
	}
	connectionsMu.RUnlock()
	// If Id is not set, generate a new one
	if request.Id == "" {
		request.Id = utils.GenerateId()
	}
	// If parent id is not set, generate a new one
	if request.ParentId == "" {
		request.ParentId = utils.GenerateId()
	}
	// Send request to the client
	err = connection.Ws.WriteJSON(request)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to send request to the client",
		})
		return
	}
	// Wait for response
	for {
		// Read response
		var response types.ChatGptResponse
		err = connection.Ws.ReadJSON(&response)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Failed to read response from the client",
			})
			return
		}
		// Check if the response is for the request
		if response.Id == request.Id {
			c.JSON(200, gin.H{
				"response": response,
			})
			return
		}
	}
}

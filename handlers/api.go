package handlers

import (
	"encoding/json"
	"time"

	"github.com/ChatGPT-Hackers/ChatGPT-API-server/types"
	"github.com/ChatGPT-Hackers/ChatGPT-API-server/utils"
	"github.com/gin-gonic/gin"
)

// // # API routes
func API_ask(c *gin.Context) {
	// Get request
	var request types.ChatGptRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request",
		})
		return
	}
	// Check if "Authorization" in Header
	if c.Request.Header["Authorization"] == nil {
		c.JSON(401, gin.H{
			"error": "API key not provided",
		})
		return
	}
	// Check if API key is valid
	verified, err := utils.VerifyToken(c.Request.Header["Authorization"][0])
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to verify API key",
		})
		return
	}
	if !verified {
		c.JSON(401, gin.H{
			"error": "Invalid API key",
		})
		return
	}
	// If Id is not set, generate a new one
	if request.MessageId == "" {
		request.MessageId = utils.GenerateId()
	}
	// If parent id is not set, generate a new one
	if request.ParentId == "" {
		request.ParentId = utils.GenerateId()
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to convert request to json",
		})
		return
	}
	var connection *types.Connection
	// Check conversation id
	connectionPool.Mu.RLock()
	// Check number of connections
	if len(connectionPool.Connections) == 0 {
		c.JSON(503, gin.H{
			"error": "No available clients",
		})
		return
	}
	connectionPool.Mu.RUnlock()
	if request.ConversationId == "" {
		// Retry 3 times before giving up
		var succeeded bool = false
		for i := 0; i < 3; i++ {
			// Find connection with the lowest load and where heartbeat is after last message time
			connectionPool.Mu.RLock()
			for _, conn := range connectionPool.Connections {
				if connection == nil || conn.LastMessageTime.Before(connection.LastMessageTime) {
					if conn.Heartbeat.After(conn.LastMessageTime) {
						connection = conn
					}
				}
			}
			connectionPool.Mu.RUnlock()
			// Check if connection was found
			if connection == nil {
				c.JSON(503, gin.H{
					"error": "No available clients",
				})
				return
			}
			// Ping before sending request
			var pingSucceeded bool = ping(connection.Id)
			if !pingSucceeded {
				// Ping failed. Try again
				connectionPool.Delete(connection.Id)
				succeeded = false
				connection = nil
				continue
			} else {
				succeeded = true
				break
			}
		}
		if !succeeded {
			// Delete connection
			c.JSON(503, gin.H{
				"error": "Ping failed",
			})
			return
		}
	} else {
		// Check if conversation exists
		conversation, ok := conversationPool.Get(request.ConversationId)
		if !ok {
			// Error
			c.JSON(500, gin.H{
				"error": "Conversation doesn't exists",
			})
			return
		} else {
			// Get connectionId of the conversation
			connectionId := conversation.ConnectionId
			// Check if connection exists
			connection, ok = connectionPool.Get(connectionId)
			if !ok {
				// Error
				c.JSON(500, gin.H{
					"error": "Connection no longer exists",
				})
				return
			}
		}
		// Ping before sending request
		if !ping(connection.Id) {
			c.JSON(503, gin.H{
				"error": "Ping failed",
			})
			return
		}
	}
	message := types.Message{
		Id:      utils.GenerateId(),
		Message: "ChatGptRequest",
		// Convert request to json
		Data: string(jsonRequest),
	}
	err = connection.Ws.WriteJSON(message)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to send request to the client",
		})
		// Delete connection
		connectionPool.Delete(connection.Id)
		return
	}
	// Set last message time
	connection.LastMessageTime = time.Now()
	// Wait for response with a timeout
	for {
		// Read message
		var receive types.Message
		connection.Ws.SetReadDeadline(time.Now().Add(120 * time.Second))
		err = connection.Ws.ReadJSON(&receive)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Failed to read response from the client",
				"err":   err.Error(),
			})
			// Delete connection
			connectionPool.Delete(connection.Id)
			return
		}
		// Check if the message is the response
		if receive.Id == message.Id {
			// Convert response to ChatGptResponse
			var response types.ChatGptResponse
			err = json.Unmarshal([]byte(receive.Data), &response)
			if err != nil {
				c.JSON(500, gin.H{
					"error":    "Failed to convert response to ChatGptResponse",
					"response": receive,
				})
				return
			}
			// Add conversation to pool
			conversation := &types.Conversation{
				Id:           response.ConversationId,
				ConnectionId: connection.Id,
			}
			conversationPool.Set(conversation)
			// Send response
			c.JSON(200, response)
			// Heartbeat
			connection.Heartbeat = time.Now()
			return
		} else {
			// Error
			c.JSON(500, gin.H{
				"error": "Failed to find response from the client",
			})
			return
		}

	}

}

func API_getConnections(c *gin.Context) {
	// Get connections
	var connections []*types.Connection
	connectionPool.Mu.RLock()
	for _, connection := range connectionPool.Connections {
		connections = append(connections, connection)
	}
	connectionPool.Mu.RUnlock()
	// Send connections
	c.JSON(200, gin.H{
		"connections": connections,
	})
}

func ping(connection_id string) bool {
	// Get connection
	connection, ok := connectionPool.Get(connection_id)
	// Send "ping" to the connection
	if ok {
		id := utils.GenerateId()
		send := types.Message{
			Id:      id,
			Message: "ping",
		}
		connection.Ws.SetReadDeadline(time.Now().Add(5 * time.Second))
		err := connection.Ws.WriteJSON(send)
		if err != nil {
			return false
		}
		// Wait for response with a timeout
		for {
			// Read message
			var receive types.Message
			err = connection.Ws.ReadJSON(&receive)
			if err != nil {
				return false
			}
			// Check if the message is the response
			if receive.Id == send.Id {
				return true
			} else {
				// Error
				return false
			}
		}
	}
	return false
}

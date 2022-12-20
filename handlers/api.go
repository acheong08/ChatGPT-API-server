package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ChatGPT-Hackers/go-server/types"
	"github.com/ChatGPT-Hackers/go-server/utils"
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
	// If Id is not set, generate a new one
	if request.MessageId == "" {
		request.MessageId = utils.GenerateId()
	}
	// If parent id is not set, generate a new one
	if request.ParentId == "" {
		request.ParentId = utils.GenerateId()
	}
	// Get connection with the lowest load
	var connection *types.Connection
	connectionPool.Mu.RLock()
	for _, conn := range connectionPool.Connections {
		if connection == nil || conn.LastMessageTime.Before(connection.LastMessageTime) {
			connection = conn
		}
	}
	connectionPool.Mu.RUnlock()
	// Do not send request if connection currently has a request
	if connection.LastMessageTime.After(connection.Heartbeat) {
		c.JSON(503, gin.H{
			"error": "No available clients",
		})
		return
	}
	// Send request to the client
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to convert request to json",
		})
		return
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
		return
	}
	// Wait for response
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
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
		if response.Id == message.Id {
			c.JSON(200, gin.H{
				"response": response,
			})
			// Heartbeat
			connection.Heartbeat = time.Now()
			return
		}
		select {
		case <-ctx.Done():
			c.JSON(504, gin.H{
				"error": "Timed out waiting for response from the client",
			})
			return
		default:
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

func API_connectionPing(c *gin.Context) {
	// Get connection id
	id := c.Param("connection_id")
	// Get connection
	connectionPool.Mu.RLock()
	connection, ok := connectionPool.Connections[id]
	connectionPool.Mu.RUnlock()
	// Send "ping" to the connection
	if ok {
		send := types.Message{
			Id:      utils.GenerateId(),
			Message: "ping",
		}
		err := connection.Ws.WriteJSON(send)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Failed to send ping to the client",
			})
			return
		}
		// Wait for response
		for {
			// Read message
			var receive types.Message
			err = connection.Ws.ReadJSON(&receive)
			if err != nil {
				return
			}
			// Check if the message is the connection id
			if receive.Id == send.Id {
				c.JSON(200, gin.H{
					"message": receive,
				})
				// Heartbeat
				connection.Heartbeat = time.Now()
				return
			} else {
				// Return incorrect message
				c.JSON(500, gin.H{
					"error":    "Incorrect message",
					"expected": send,
					"received": receive,
				})
				return
			}
		}
	} else {
		c.JSON(404, gin.H{
			"error": "Connection not found",
		})
		return
	}
}

package main

import (
	"sync"
	"time"

	// Import local packages
	"github.com/ChatGPT-Hackers/go-server/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Message struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

type ChatGptRequest struct {
	Id             string `json:"id"`
	ConversationId string `json:"conversation_id"`
	ParentId       string `json:"parent_id"`
	ConnectionId   string `json:"connection_id"`
	Content        string `json:"content"`
}

type ChatGptResponse struct {
	Id             string `json:"id"`
	ConversationId string `json:"conversation_id"`
}

type Connection struct {
	// The websocket connection.
	Ws *websocket.Conn
	// The connecton id.
	Id string
	// Last heartbeat time.
	Heartbeat int64
	// Last message time.
	LastMessageTime int64
}

var (
	// The websocket upgrader.
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	// Connection pool.
	connections []*Connection
	// Mutex to synchronize access to the connections slice.
	connectionsMu sync.RWMutex
)

func main() {
	router := gin.Default()
	//// # Add routes
	// Register new client connection
	router.GET("/client/register", clientRegister)
	router.POST("/api/ask", apiAsk)
}

// // # API routes
func apiAsk(c *gin.Context) {
	// Get the message
	var message Message
	err := c.BindJSON(&message)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid message",
		})
		return
	}
	// Find the client with the oldest last message time
	var oldestConnection *Connection
	connectionsMu.RLock()
	for _, connection := range connections {
		if oldestConnection == nil || connection.LastMessageTime < oldestConnection.LastMessageTime {
			oldestConnection = connection
		}
	}
	connectionsMu.RUnlock()
	// Send the message to the client
	err = oldestConnection.Ws.WriteJSON(message)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to send message",
		})
		return
	}
	// Update the last message time
	oldestConnection.LastMessageTime = time.Now().Unix()
	// Get the response
	_, response, err := oldestConnection.Ws.ReadMessage()
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to get response",
		})
		return
	}
	// Send the response
	c.JSON(200, gin.H{
		"response": string(response),
	})
}

// // # Client routes
func clientRegister(c *gin.Context) {
	// Make websocket connection
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	// Add connection to the pool
	connection := &Connection{
		Id:              utils.GenerateId(),
		Ws:              ws,
		Heartbeat:       time.Now().Unix(),
		LastMessageTime: time.Now().Unix(),
	}
	connectionsMu.Lock()
	connections = append(connections, connection)
	connectionsMu.Unlock()
	// Send connection id to the client
	err = connection.Ws.WriteJSON(Message{
		Id:      connection.Id,
		Message: "Connection id",
	})
	if err != nil {
		return
	}
	// Close the connection when it is no longer needed
	defer connection.Ws.Close()
}

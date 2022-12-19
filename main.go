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
}

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

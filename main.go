package main

import (
	_ "fmt"
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

type User struct {
	Id              string `json:"id"`
	LastMessageTime int64  `json:"last_message_time"`
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
	connections = make(map[int]*Connection)
)

func main() {
	router := gin.Default()

	//// # Add routes
	// Register new client connection
	router.GET("/client/register", client_register)
}

func client_register(c *gin.Context) {
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
	connections[len(connections)] = connection
	// Send connection id to the client
	connection.Ws.WriteJSON(Message{
		Id:      connection.Id,
		Message: "Connection id",
	})
}

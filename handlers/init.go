package handlers

import (
	"github.com/gorilla/websocket"

	"github.com/ChatGPT-Hackers/go-server/types"
)

var (
	// The websocket upgrader.
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)
var connectionPool = types.NewConnectionPool()

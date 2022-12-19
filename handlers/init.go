package handlers

import (
	"sync"

	"github.com/gorilla/websocket"

	"github.com/ChatGPT-Hackers/go-server/types"
)

var (
	// The websocket upgrader.
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	// Mutex to synchronize access to the connection pool.
	connectionPoolMu sync.RWMutex
)
var connectionPool = types.NewConnectionPool()

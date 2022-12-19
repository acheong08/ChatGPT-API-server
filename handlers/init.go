package handlers

import (
	"sync"

	"github.com/ChatGPT-Hackers/go-server/types"
	"github.com/gorilla/websocket"
)

var (
	// The websocket upgrader.
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	// Connection pool.
	connections []*types.Connection
	// Mutex to synchronize access to the connections slice.
	connectionsMu sync.RWMutex
)

package handlers

import (
	"net/http"

	"github.com/ChatGPT-Hackers/ChatGPT-API-server/types"
	"github.com/gorilla/websocket"
)

var (
	// The websocket upgrader.
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

var connectionPool = types.NewConnectionPool()
var conversationPool = types.NewConversationPool()

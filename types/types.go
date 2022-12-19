package types

import "github.com/gorilla/websocket"

type Message struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

type ChatGptResponse struct {
	Id             string `json:"id"`
	ResponseId     string `json:"response_id"`
	ConversationId string `json:"conversation_id"`
	Content        string `json:"content"`
}

type ChatGptRequest struct {
	Id             string `json:"id"`
	MessageId      string `json:"message_id"`
	ConversationId string `json:"conversation_id"`
	ParentId       string `json:"parent_id"`
	ConnectionId   string `json:"connection_id"`
	Content        string `json:"content"`
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

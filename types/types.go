package types

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Id      string `json:"id"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type ChatGptResponse struct {
	Id             string `json:"id"`
	ResponseId     string `json:"response_id"`
	ConversationId string `json:"conversation_id"`
	Content        string `json:"content"`
	Error          string `json:"error"`
}

type ChatGptRequest struct {
	MessageId      string `json:"message_id"`
	ConversationId string `json:"conversation_id"`
	ParentId       string `json:"parent_id"`
	Content        string `json:"content"`
}

type Connection struct {
	// The websocket connection.
	Ws *websocket.Conn
	// The connecton id.
	Id string
	// Last heartbeat time.
	Heartbeat time.Time
	// Last message time.
	LastMessageTime time.Time
}

type ConnectionPool struct {
	Connections map[string]*Connection
	Mu          sync.RWMutex
}

func (p *ConnectionPool) Get(id string) (*Connection, bool) {
	p.Mu.RLock()
	defer p.Mu.RUnlock()
	conn, ok := p.Connections[id]
	if conn == nil {
		ok = false
	}
	return conn, ok
}

func (p *ConnectionPool) Set(conn *Connection) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	p.Connections[conn.Id] = conn
}

func (p *ConnectionPool) Delete(id string) error {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	delete(p.Connections, id)
	return nil
}

func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		Connections: make(map[string]*Connection),
	}
}

type Conversation struct {
	Id           string `json:"id"`
	ConnectionId string `json:"connection_id"`
}

type ConversationPool struct {
	Conversations map[string]*Conversation
	Mu            sync.RWMutex
}

func (p *ConversationPool) Get(id string) (*Conversation, bool) {
	p.Mu.RLock()
	defer p.Mu.RUnlock()
	conversation, ok := p.Conversations[id]
	return conversation, ok
}

func (p *ConversationPool) Set(conversation *Conversation) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	p.Conversations[conversation.Id] = conversation
}

func (p *ConversationPool) Delete(id string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	delete(p.Conversations, id)
}

func NewConversationPool() *ConversationPool {
	return &ConversationPool{
		Conversations: make(map[string]*Conversation),
	}
}

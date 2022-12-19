package types

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

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
	return conn, ok
}

func (p *ConnectionPool) Set(conn *Connection) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	p.Connections[conn.Id] = conn
}

func (p *ConnectionPool) Delete(id string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	delete(p.Connections, id)
}

func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		Connections: make(map[string]*Connection),
	}
}

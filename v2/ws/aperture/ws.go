package aperture

import (
	"fmt"
	"sync"
)

type TopicCollection struct {
	list map[string]map[*Conn]struct{}
	mu   sync.RWMutex
}

type WebSocket struct {
	Open          func(conn *Conn)
	Message       func(message any, conn *Conn)
	Close         func(conn *Conn, code string, reason string)
	OnPublish     func(topic string, message any, conn *Conn)
	IdleTimeout   int
	PrivateAccess bool
	Description   string
	GetSequre     func() bool

	topicCollections TopicCollection
	topicDocs        map[string]any
	// docs             []string
}

func (ws *WebSocket) Publish(topic string, message any) {
	ws.topicCollections.mu.RLock()
	defer ws.topicCollections.mu.RUnlock()

	clients, exists := ws.topicCollections.list[topic]

	if !exists {
		return
	}

	for conn := range clients {
		conn.Publish(topic, message)
	}
}

func (ws *WebSocket) Subscribe(c *Conn, topic string) {
	ws.topicCollections.mu.Lock()
	defer ws.topicCollections.mu.Unlock()

	if _, exists := ws.topicCollections.list[topic]; !exists {
		ws.topicCollections.list[topic] = map[*Conn]struct{}{}
	}

	ws.topicCollections.list[topic][c] = struct{}{}

	fmt.Println("subscribe+", topic, len(ws.topicCollections.list[topic]))
}

func (ws *WebSocket) Unsubscribe(c *Conn, topic string) {
	ws.topicCollections.mu.Lock()
	defer ws.topicCollections.mu.Unlock()

	delete(ws.topicCollections.list[topic], c)
}

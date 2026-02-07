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
	Open    func(conn *Conn)
	Message func(message any, conn *Conn)
	Close   func(conn *Conn, code string, reason string)

	OnPublish func(topic string, message any, conn *Conn)

	IdleTimeout   int
	PrivateAccess bool

	Description string
	GetSequre   func() bool

	TopicsCollection TopicCollection
	docs             []string
}

func (ws *WebSocket) Publish(topic string, message any) {
	ws.TopicsCollection.mu.RLock()
	defer ws.TopicsCollection.mu.RUnlock()

	clients, exists := ws.TopicsCollection.list[topic]

	if !exists {
		return
	}

	for conn := range clients {
		conn.Publish(topic, message)
	}
}

func (ws *WebSocket) Subscribe(c *Conn, topic string) {
	ws.TopicsCollection.mu.Lock()
	defer ws.TopicsCollection.mu.Unlock()

	if _, exists := ws.TopicsCollection.list[topic]; !exists {
		ws.TopicsCollection.list[topic] = map[*Conn]struct{}{}
	}

	ws.TopicsCollection.list[topic][c] = struct{}{}

	fmt.Println("subscribe+", topic, len(ws.TopicsCollection.list[topic]))
}

func (ws *WebSocket) Unsubscribe(c *Conn, topic string) {
	ws.TopicsCollection.mu.Lock()
	defer ws.TopicsCollection.mu.Unlock()

	delete(ws.TopicsCollection.list[topic], c)
}

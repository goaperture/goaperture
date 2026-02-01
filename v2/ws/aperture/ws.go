package aperture

import "sync"

type TopicCollection struct {
	list map[string]map[*Conn]struct{}
	mu   sync.RWMutex
}

type WebSocket struct {
	Open          func(conn *Conn)
	Message       func(conn *Conn, message string)
	Close         func(conn *Conn, code string, reason string)
	IdleTimeout   int
	PrivateAccess bool
	Description   string
	GetSequre     func() bool

	TopicsCollection TopicCollection
	docs             []string
}

func (ws *WebSocket) Publish(topic string, messate string) {
	ws.TopicsCollection.mu.RLock()
	defer ws.TopicsCollection.mu.RUnlock()

	clients, exists := ws.TopicsCollection.list[topic]

	if !exists {
		return
	}

	for conn := range clients {
		conn.Send(messate)
	}
}

func (ws *WebSocket) Subscribe(c *Conn, topic string) {
	ws.TopicsCollection.mu.Lock()
	defer ws.TopicsCollection.mu.Unlock()

	ws.TopicsCollection.list[topic][c] = struct{}{}
}

func (ws *WebSocket) Unsubscribe(c *Conn, topic string) {
	ws.TopicsCollection.mu.Lock()
	defer ws.TopicsCollection.mu.Unlock()

	delete(ws.TopicsCollection.list[topic], c)
}

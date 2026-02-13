package aperture

import (
	"fmt"
	"strings"
	"sync"
)

type TopicCollection struct {
	list map[string]map[*Conn]struct{}
	mu   sync.RWMutex
}

type WebSocket struct {
	Open             func(conn *Conn)
	Message          func(message any, conn *Conn)
	Close            func(conn *Conn, code string, reason string)
	OnPublish        func(topic string, message any, conn *Conn) bool
	IdleTimeout      int
	PrivateAccess    bool
	Description      string
	GetSequre        func() bool
	topicCollections TopicCollection
	jsonTopics       map[string]jsonTopic
}

func (ws *WebSocket) Publish(topic string, message any) {
	ws.topicCollections.mu.RLock()
	defer ws.topicCollections.mu.RUnlock()

	for prefix, clients := range ws.topicCollections.list {
		if strings.HasPrefix(topic, prefix) {
			for conn := range clients {
				if err := conn.Publish(topic, message); err != nil {
					delete(ws.topicCollections.list[prefix], conn)
					fmt.Println("remove subscriber from ", prefix)
				}
			}
		}
	}
}

func (ws *WebSocket) GetTopicLen(topic string) int {
	ws.topicCollections.mu.RLock()
	defer ws.topicCollections.mu.RUnlock()

	result := 0

	for prefix, clients := range ws.topicCollections.list {
		if strings.HasPrefix(topic, prefix) {
			result += len(clients)
		}
	}

	return result
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

	fmt.Println("unsubscribe-", topic, len(ws.topicCollections.list[topic]))
}

func (ws *WebSocket) handlePublish(topic string, message any, client *Conn) {
	var sendToClient = true
	if ws.OnPublish != nil {
		if !ws.OnPublish(topic, message, client) {
			sendToClient = false
		}
	}

	for prefix, jsonTopic := range ws.jsonTopics {
		if strings.HasPrefix(topic, prefix) {
			if !jsonTopic.handle(topic, message) {
				sendToClient = false
			}
		}
	}

	if !sendToClient {
		return
	}

	// Отправляем реальным клиентам
	for prefix, connections := range ws.topicCollections.list {
		if strings.HasPrefix(topic, prefix) {
			for conn := range connections {
				if conn != client {
					conn.Publish(topic, message)
				}
			}
		}
	}
}

func (ws *WebSocket) getTopicDocs() map[string]any {
	var result = make(map[string]any)

	for prefix, jsonTopic := range ws.jsonTopics {
		result[prefix] = jsonTopic.PrepareCall()
	}

	return result
}

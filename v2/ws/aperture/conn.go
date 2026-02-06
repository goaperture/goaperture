package aperture

type Conn struct {
	Id     int
	Send   func(message SocketData) error
	topics map[string]struct{}

	ws *WebSocket
}

func (c *Conn) Subscribe(topic string) error {
	c.ws.Subscribe(c, topic)
	c.topics[topic] = struct{}{}

	return nil
}

func (c *Conn) Unsubscribe(topic string) error {
	c.ws.Unsubscribe(c, topic)
	delete(c.topics, topic)

	return nil
}

func (c *Conn) Publish(topic string, message string) error {
	return c.Send(SocketData{
		Message: message,
		Topic:   topic,
	})
}

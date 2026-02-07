package aperture

import "encoding/json"

type Topic[Message any] struct {
	ws        *WebSocket
	prefix    string
	onPublish func(topic string, message Message) error
	prepare   Message
}

func CreateTopic[Message any](ws *WebSocket, prefix string, onPublish func(topic string, message Message) error) Topic[Message] {
	topic := Topic[Message]{
		ws:        ws,
		prefix:    prefix,
		onPublish: onPublish,
	}

	if ws.topicDocs == nil {
		ws.topicDocs = make(map[string]any)
	}

	ws.topicDocs[prefix] = topic.prepare

	return topic
}

func (t *Topic[Message]) Publish(topic string, message Message) error {
	data, err := json.Marshal(message)

	if err != nil {
		return err
	}

	t.ws.Publish(t.prefix+topic, string(data))
	return nil
}

// func (t *Topic[Message]) Prepare(test Message) {
// 	t.prepare = test
// }

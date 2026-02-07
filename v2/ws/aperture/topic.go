package aperture

import "encoding/json"

type Topic[Message any] struct {
	ws        *WebSocket
	pefix     string
	onPublish func(topic string, message Message) error
	prepare   Message
}

func CreateTopic[Message any](ws *WebSocket, pefix string, onPublish func(topic string, message Message) error) Topic[Message] {
	return Topic[Message]{
		ws:        ws,
		pefix:     pefix,
		onPublish: onPublish,
	}
}

func (t *Topic[Message]) Publish(topic string, message Message) error {
	data, err := json.Marshal(message)

	if err != nil {
		return err
	}

	t.ws.Publish(t.pefix+topic, string(data))
	return nil
}

// func (t *Topic[Message]) Prepare(test Message) {
// 	t.prepare = test
// }

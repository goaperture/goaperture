package aperture

import "encoding/json"

type Topic[Message any] struct {
	ws        *WebSocket
	Name      string
	OnMessage func(message Message) error
}

func (t *Topic[Message]) Publish(message Message) error {
	data, err := json.Marshal(message)

	if err != nil {
		return err
	}

	t.ws.Publish(t.Name, string(data))
	return nil
}

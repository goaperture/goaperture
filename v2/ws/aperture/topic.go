package aperture

import "encoding/json"

type Topic[Message any] struct {
	ws        *WebSocket
	Prefix    string
	OnPublish func(topic string, message Message) bool
	onJoin    func(topic string)
	onLeave   func(topic string, reasone string)
	prepare   Message
}

type jsonTopic struct {
	prepare any
	handle  func(topic string, message any) bool
}

func CreateTopic[Message any](ws *WebSocket, topic Topic[Message]) Topic[Message] {
	topic.ws = ws

	if ws.jsonTopics == nil {
		ws.jsonTopics = make(map[string]jsonTopic)
	}

	ws.jsonTopics[topic.Prefix] = jsonTopic{
		prepare: topic.prepare,
		handle: func(key string, message any) bool {
			if topic.OnPublish == nil {
				return true
			}

			strMessage, _ := json.Marshal(message)
			var jsonMessage Message
			json.Unmarshal(strMessage, &jsonMessage)

			return topic.OnPublish(key, jsonMessage)

		},
	}

	return topic
}

func (t *Topic[Message]) Publish(topic string, message Message) {
	data, _ := json.Marshal(message)
	t.ws.Publish(t.Prefix+topic, string(data))
}

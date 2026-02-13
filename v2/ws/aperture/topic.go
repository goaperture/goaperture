package aperture

import (
	"encoding/json"

	"github.com/goaperture/goaperture/v2/api/collector"
)

type CL[Message any] = collector.Collector[Message, Message]

type Topic[Message any] struct {
	ws              *WebSocket
	Prefix          string
	OnPublish       func(topic string, message Message) bool
	onJoin          func(topic string)
	onLeave         func(topic string, reasone string)
	prepareDataItem Message
	Prepare         func(collector *CL[Message])
}

type jsonTopic struct {
	PrepareCall func() any
	handle      func(topic string, message any) bool
}

func CreateTopic[Message any](ws *WebSocket, topic Topic[Message]) Topic[Message] {
	topic.ws = ws

	if ws.jsonTopics == nil {
		ws.jsonTopics = make(map[string]jsonTopic)
	}

	ws.jsonTopics[topic.Prefix] = jsonTopic{
		PrepareCall: func() any { // #TODO - типы топиков обработать нормально
			if topic.Prepare == nil {
				return []any{topic.prepareDataItem}
			}

			var cll = collector.Collector[Message, Message]{}
			topic.Prepare(&cll)
			return cll.GetDump().Outputs
		},
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

func (t *Topic[Message]) Publish(suffix string, message Message) {
	t.ws.Publish(t.Prefix+suffix, message)
}

func (t *Topic[Messge]) GetLen(suffix string) int {
	return t.ws.GetTopicLen(t.Prefix + suffix)
}

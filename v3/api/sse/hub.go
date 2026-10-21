package sse

import "sync"

// hub — реестр подписчиков по ключам стрима.
// Один ключ -> множество каналов (по одному на открытое соединение).
type hub struct {
	mu   sync.RWMutex
	subs map[string]map[chan []byte]struct{}
}

var defaultHub = &hub{
	subs: make(map[string]map[chan []byte]struct{}),
}

// subscribe регистрирует канал подписчика для ключа и возвращает его.
// Канал буферизированный: медленный клиент не должен блокировать рассылку.
func subscribe(key string) chan []byte {
	ch := make(chan []byte, 16)

	defaultHub.mu.Lock()
	defer defaultHub.mu.Unlock()

	if defaultHub.subs[key] == nil {
		defaultHub.subs[key] = make(map[chan []byte]struct{})
	}

	defaultHub.subs[key][ch] = struct{}{}

	return ch
}

// unsubscribe удаляет подписчика из реестра.
// Канал НЕ закрывается: Publish копирует список каналов под RLock и шлёт
// уже без блокировки, поэтому close дал бы panic при concurrent рассылке.
// Канал соберёт GC, когда горутина-подписчик завершится.
func unsubscribe(key string, ch chan []byte) {
	defaultHub.mu.Lock()
	defer defaultHub.mu.Unlock()

	if set := defaultHub.subs[key]; set != nil {
		delete(set, ch)

		if len(set) == 0 {
			delete(defaultHub.subs, key)
		}
	}
}

// Publish копирует ОДНО сообщение всем клиентам с одинаковым ключом.
// data должна быть уже сериализована в JSON.
func Publish(key string, data []byte) {
	defaultHub.mu.RLock()

	set := defaultHub.subs[key]
	targets := make([]chan []byte, 0, len(set))
	for ch := range set {
		targets = append(targets, ch)
	}

	defaultHub.mu.RUnlock()

	for _, ch := range targets {
		select {
		case ch <- data:
		default: // канал полон — клиент медленный, пропускаем
		}
	}
}

package sse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func Run(w http.ResponseWriter, r *http.Request, result any, key string) bool {
	acceptHeader := r.Header.Get("Accept")

	if !strings.Contains(acceptHeader, "text/event-stream") {
		return false
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		return false
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	ch := subscribe(key)
	defer unsubscribe(key, ch)

	if err := send(w, flusher, result); err != nil {
		// log.Println("error", err)
		return false
	}

	for {
		select {
		case <-r.Context().Done():
			log.Println("Клиент отключился")
			return false

		case msg := <-ch:
			if err := send(w, flusher, msg); err != nil {
				// log.Println("error", err)
			}
		}
	}
}

func send(w http.ResponseWriter, flusher http.Flusher, result any) error {
	var jsonData []byte

	switch v := result.(type) {
	case []byte:
		// уже сериализовано (сообщение из хаба) — пишем как есть,
		// иначе json.Marshal([]byte) дал бы base64
		jsonData = v
	default:
		var err error
		jsonData, err = json.Marshal(result)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(w, "data: %s\n\n", jsonData)

	flusher.Flush()

	return nil
}

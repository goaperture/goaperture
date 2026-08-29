package sse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Message struct {
	Status    string    `json:"status"`
	Value     int       `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func Run(w http.ResponseWriter, r *http.Request, result any) bool {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return false
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	send(w, flusher, result)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	counter := 0

	for {
		select {
		case <-r.Context().Done():
			log.Println("Клиент отключился")
			return false

		case <-ticker.C:
			counter++

			msg := Message{
				Status:    "active",
				Value:     counter,
				Timestamp: time.Now(),
			}
			err := send(w, flusher, msg)
			if err != nil {
				log.Println("error", err)
				continue
			}

		}
	}
}

func send(w http.ResponseWriter, flusher http.Flusher, result any) error {
	jsonData, err := json.Marshal(result)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "data: %s\n\n", jsonData)

	flusher.Flush()

	return nil
}

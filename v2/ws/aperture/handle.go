package aperture

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/goaperture/goaperture/v2/api/auth"
	"github.com/goaperture/goaperture/v2/api/client"
	"github.com/goaperture/goaperture/v2/exception"
)

type SocketInput struct {
	Message   any      `json:"message"`
	Subscribe []string `json:"subscribe,omitempty"`
	Topic     string   `json:"topic,omitempty"`
}

type SocketData struct {
	Message any    `json:"message"`
	Topic   string `json:"topic,omitempty"`
}

type WebSockets map[string]SocketSwitch

type SocketSwitch struct {
	Handler       func(secret auth.XSecret) func(w http.ResponseWriter, r *http.Request)
	PrivateAccess bool
	Description   string
	Sequre        bool
	TopicDocs     map[string]any
}

func (ws *WebSockets) BindHanders(server *http.ServeMux, secret auth.XSecret) {
	for path, route := range *ws {
		server.HandleFunc(path, route.Handler(secret))
	}
}

func Handle(ws *WebSocket) SocketSwitch {
	isSequre := false

	if ws.GetSequre != nil {
		isSequre = ws.GetSequre()
	}

	ws.topicCollections.list = make(map[string]map[*Conn]struct{})

	return SocketSwitch{
		Handler:       createHandler(ws, isSequre),
		PrivateAccess: ws.PrivateAccess,
		Description:   ws.Description,
		Sequre:        isSequre,
		TopicDocs:     ws.getTopicDocs(),
	}
}

func createHandler(ws *WebSocket, isSequre bool) func(secret auth.XSecret) func(w http.ResponseWriter, r *http.Request) {
	return func(secret auth.XSecret) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			defer exception.Catch(&w)

			jwt, exists := auth.ParseAccessToken(r)

			if ws.PrivateAccess {
				accessKey := auth.GetAccessKeyFromUrl(r.Pattern)
				if !exists {
					exception.NotAccess(accessKey)
				}

				payload := auth.GetPayloadFromJwt[auth.TempPayload](jwt, secret)
				payload.Permissions.CheckX(accessKey)
			}

			conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
				InsecureSkipVerify: !isSequre,
			})
			if err != nil {
				log.Println("Ошибка Accept:", err)
				return
			}

			IdleTimeout := ws.IdleTimeout

			if IdleTimeout == 0 {
				IdleTimeout = 3
			}

			ctx, cancel := context.WithTimeout(r.Context(), time.Minute*time.Duration(IdleTimeout))
			defer cancel()

			if exists {
				ctx = client.WithToken(ctx, jwt)
			}

			var client = Conn{
				ws: ws,
				Send: func(message SocketData) error {
					data, _ := json.Marshal(message)

					return conn.Write(ctx, websocket.MessageText, data)
				},
			}

			defer func() {
				conn.Close(websocket.StatusInternalError, "close connection")

				if ws.Close != nil {
					ws.Close(&client, "code", "R")
				}
			}()

			if ws.Open != nil {
				ws.Open(&client)
			}

			for {
				messageType, data, err := conn.Read(ctx)
				if err != nil {
					log.Println("Ошибка чтения:", err)
					break
				}

				if messageType != websocket.MessageText {
					continue
				}

				var socketData SocketInput
				err = json.Unmarshal(data, &socketData)
				if err != nil {
					continue
				}

				if len(socketData.Subscribe) != 0 {
					for _, topic := range socketData.Subscribe {
						ws.Subscribe(&client, topic)
					}
					continue
				}

				if socketData.Topic != "" {
					ws.handlePublish(socketData.Topic, socketData.Message, &client)
					continue
				}

				if ws.Message != nil {
					ws.Message(socketData.Message, &client)
				}
			}
		}
	}
}

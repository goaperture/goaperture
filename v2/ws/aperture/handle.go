package aperture

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/goaperture/goaperture/v2/api/auth"
	"github.com/goaperture/goaperture/v2/api/client"
	"github.com/goaperture/goaperture/v2/exception"
)

func Handle(ws *WebSocket) SocketSwitch {
	isSequre := false

	if ws.GetSequre != nil {
		isSequre = ws.GetSequre()
	}

	return SocketSwitch{
		Handler: func(secret auth.XSecret) func(w http.ResponseWriter, r *http.Request) {
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

				ctx, cancel := context.WithTimeout(r.Context(), time.Minute*time.Duration(ws.IdleTimeout))
				defer cancel()

				if exists {
					ctx = client.WithToken(ctx, jwt)
				}

				var client = Conn{
					ws: ws,
					Send: func(message string) error {
						return conn.Write(ctx, websocket.MessageText, []byte(message))
					},
				}

				defer func() {
					conn.Close(websocket.StatusInternalError, "close connection")
					ws.Close(&client, "code", "R")
				}()

				ws.Open(&client)

				for {
					_, data, err := conn.Read(ctx)
					if err != nil {
						log.Println("Ошибка чтения:", err)
						break
					}

					ws.Message(&client, string(data))
				}
			}
		},
		PrivateAccess: ws.PrivateAccess,
		Description:   ws.Description,
		Sequre:        isSequre,
	}
}

type WebSockets map[string]SocketSwitch

type SocketSwitch struct {
	Handler       func(secret auth.XSecret) func(w http.ResponseWriter, r *http.Request)
	PrivateAccess bool
	Description   string
	Sequre        bool
}

func (ws *WebSockets) BindHanders(server *http.ServeMux, secret auth.XSecret) {
	for path, route := range *ws {
		server.HandleFunc(path, route.Handler(secret))
	}
}

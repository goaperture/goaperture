package aperture

import (
	"net/http"

	ws "github.com/goaperture/goaperture/v2/ws/aperture"
)

func (a *Api[P]) Middleware(m func(next http.Handler) http.Handler) *Api[P] {
	a.middleware = &m

	return a
}

func (a *Api[P]) SetPort(port int) *Api[P] {
	a.Port = port
	return a
}

func (a *Api[P]) SetToken(token string) *Api[P] {
	a.Token = token
	return a
}

func (a *Api[P]) Routes(routes *Routes) *Api[P] {
	a.routes = *routes
	return a
}

func (a *Api[P]) WebSockets(ws *ws.WebSockets) *Api[P] {
	a.ws = ws
	return a
}

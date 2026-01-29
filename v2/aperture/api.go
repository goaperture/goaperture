package aperture

import (
	"fmt"
	"net/http"

	"github.com/goaperture/goaperture/v2/auth"
)

type Payload any

type Api[P Payload] struct {
	Port       int
	routes     Routes
	Token      string
	Auth       *auth.Auth[P]
	Middleware *func(next http.Handler) http.Handler
}

func (a *Api[P]) ListenAndServe() {
	server := http.NewServeMux()

	for path, route := range a.routes {
		server.HandleFunc(path, route.Handler)
	}

	if a.Auth != nil {
		a.Auth.BindHanders(server)
	}

	if a.Token != "" {
		server.HandleFunc(DOC_URL, docHandle(a))
	}

	http.ListenAndServe(fmt.Sprintf(":%d", a.Port), a.wrap(server))
}

func (a *Api[P]) wrap(server *http.ServeMux) http.Handler {
	if a.Middleware != nil {
		return (*a.Middleware)(server)
	}

	return server
}

func (a *Api[P]) Run(routes *Routes) {
	a.routes = *routes
	a.ListenAndServe()
}

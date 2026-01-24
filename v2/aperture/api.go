package aperture

import (
	"fmt"
	"net/http"

	"github.com/goaperture/goaperture/v2/doc"
)

type Payload any

type Secret struct {
	Token string
	Key   struct {
		Public  string
		Private string
	}
	SecretKey string
}

type Api[P Payload] struct {
	Port       int
	Routes     Routes
	Payload    *P
	Secret     Secret
	Middleware *func(next http.Handler) http.Handler
}

func (a *Api[P]) Run() {
	server := http.NewServeMux()

	for path, route := range a.Routes {
		server.HandleFunc(path, route.Handler)
	}

	if a.Secret.Token != "" {
		server.HandleFunc("__doc__", doc.Handle(a))
	}

	http.ListenAndServe(fmt.Sprintf(":%d", a.Port), a.wrap(server))
}

func (a *Api[P]) wrap(server *http.ServeMux) http.Handler {
	if a.Middleware != nil {
		return (*a.Middleware)(server)
	}

	return server
}

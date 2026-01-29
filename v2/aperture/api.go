package aperture

import (
	"fmt"
	"net/http"

	"github.com/goaperture/goaperture/v2/auth"
	"github.com/goaperture/goaperture/v2/metrics"
)

type Payload any

type Api[P Payload] struct {
	Port       int
	routes     Routes
	Token      string
	Auth       *auth.Auth[P]
	middleware *func(next http.Handler) http.Handler
	Metrics    bool
}

func (a *Api[P]) Run() {
	server := http.NewServeMux()

	for path, route := range a.routes {
		server.HandleFunc(path, route.Handler(a.Auth.GetSecret()))
	}

	if a.Auth != nil {
		a.Auth.BindHanders(server)
	}

	if a.Metrics {
		metrics.BindHanders(server)
	}

	if a.Token != "" {
		server.HandleFunc(DOC_URL, docHandle(a))
	}

	http.ListenAndServe(fmt.Sprintf(":%d", a.Port), a.wrap(server))
}

func (a *Api[P]) wrap(server *http.ServeMux) http.Handler {
	if a.middleware != nil {
		return (*a.middleware)(server)
	}

	return server
}

package types

import (
	"net/http"

	"github.com/goaperture/goaperture/lib/aperture"
)

type Input interface {
	any
}

type Payload interface {
	any
}

type Request struct {
	Request  *http.Request
	Responce *http.ResponseWriter
	secret   string
}

type DocInput struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Responce struct {
	Data  any    `json:"data"`
	Error *Error `json:"error,omitempty"`
}

type Route[T Input, P Payload] struct {
	Path    string
	Handler func(T, aperture.Client[P]) (any, error)
	Test    func(func(T))
}

type Aperture struct {
	Mux        *http.ServeMux
	Middleware func(next http.Handler) http.Handler
	GetSecret  func() string
}

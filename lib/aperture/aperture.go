package aperture

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/goaperture/goaperture/lib/params"
)

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Responce struct {
	Data  any    `json:"data"`
	Error *Error `json:"error,omitempty"`
}

type Input interface {
	any
}

type Route[T Input, P Payload] struct {
	Path    string
	Handler func(T, Client[P]) (any, error)
	Test    func(func(T))
}

type Aperture[P Payload] struct {
	Mux        *http.ServeMux
	Middleware func(next http.Handler) http.Handler
	GetSecret  func() string
}

func NewServer[P Payload](secret *string) *Aperture[P] {
	return &Aperture[P]{
		Mux: http.NewServeMux(),
		Middleware: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		},
		GetSecret: func() string {
			if secret != nil {
				return *secret
			}

			return "OP3$1$$VF555EJX6GTOBDCP5HPA.ZA7A@CN29k.Kj"
		},
	}
}

func NewRoute[I Input, P Payload](api *Aperture[P], route Route[I, P]) {
	api.Mux.HandleFunc(route.Path, invoke(route.Handler, true, api.GetSecret()))
	newDoc(route)
}

func (api *Aperture[P]) Run(port int, token *string, clients []P) error {
	if token != nil {
		api.Mux.HandleFunc("/__doc__", invoke(docHandler(*token, clients), false, api.GetSecret()))
	}
	return http.ListenAndServe(":"+strconv.Itoa(port), api.Middleware(api.Mux))
}

func invoke[I Input, P Payload](method func(I, Client[P]) (any, error), wrap bool, secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(Responce{
					Error: &Error{Message: fmt.Sprint(r), Code: 401},
				})
			}
		}()

		var props = params.GetInput[I](r)
		data, err := method(props, NewClient[P](r, &w, secret, false))

		var ResponceErr *Error
		if err != nil {
			ResponceErr = &Error{Message: err.Error(), Code: 400}
		}

		w.Header().Set("Content-Type", "application/json")

		if !wrap {
			json.NewEncoder(w).Encode(data)
			return
		}

		result := Responce{
			Data:  data,
			Error: ResponceErr,
		}

		json.NewEncoder(w).Encode(result)
	}
}

package aperture

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

type Aperture struct {
	Mux        *http.ServeMux
	Middleware func(next http.Handler) http.Handler
	GetSecret  func() string
}

func NewServer() *Aperture {
	return &Aperture{
		Mux: http.NewServeMux(),
		Middleware: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		},
		GetSecret: func() string { return "OP3$1$$VF555EJX6GTOBDCP5HPA.ZA7A@CN29k.Kj" },
	}
}

func NewRoute[I Input, P Payload](api *Aperture, route Route[I, P]) {
	api.Mux.HandleFunc(route.Path, invoke(route.Handler, true, api.GetSecret()))
	newDoc(route)
}

func (api *Aperture) Run(port int, token *string) error {
	api.Mux.HandleFunc("/__doc__", invoke(docHandler(token), false, api.GetSecret()))
	return http.ListenAndServe(":"+strconv.Itoa(port), api.Middleware(api.Mux))
}

func invoke[I Input, P Payload](method func(I, Client[P]) (any, error), wrap bool, secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var props I

		switch r.Method {
		case http.MethodGet:
			getParamsToStruct(r.URL.Query(), &props)
		case http.MethodPost:
			json.NewDecoder(r.Body).Decode(&props)
		default:
			log.Println("Неизвестный метод")
		}

		defer func() {
			if r := recover(); r != nil {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(Responce{
					Error: &Error{Message: fmt.Sprint(r), Code: 401},
				})
			}
		}()

		data, err := method(props, NewClient[P](r, &w, secret))

		var ResponceErr *Error
		if err != nil {
			ResponceErr = &Error{Message: err.Error(), Code: 400}
		}

		if !wrap {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)
			return
		}

		result := Responce{
			Data:  data,
			Error: ResponceErr,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

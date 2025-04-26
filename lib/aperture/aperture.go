package aperture

import (
	"encoding/json"
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

type Route[T Input] struct {
	Path    string
	Handler func(T) (any, error)
	Test    func(func(T))
}

type Aperture struct {
	Mux        *http.ServeMux
	Middleware func(next http.Handler) http.Handler
}

func NewServer() *Aperture {
	return &Aperture{
		Mux: http.NewServeMux(),
		Middleware: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		},
	}
}

func NewRoute[I Input](api *Aperture, route Route[I]) {
	api.Mux.HandleFunc(route.Path, invoke(route.Handler, true))
	newDoc(route)
}

func (api *Aperture) Run(port int, token *string) error {
	api.Mux.HandleFunc("/__doc__", invoke(docHandler(token), false))
	return http.ListenAndServe(":"+strconv.Itoa(port), api.Middleware(api.Mux))
}

func invoke[I Input](method func(I) (any, error), wrap bool) func(w http.ResponseWriter, r *http.Request) {
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

		data, err := method(props)

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

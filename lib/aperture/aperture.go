package aperture

import (
	"encoding/json"
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

func NewRoute[I Input](route Route[I]) {
	http.HandleFunc(route.Path, invoke(route.Handler))
	newDoc(route)
}

func Run(port int, token *string) error {
	http.HandleFunc("/__doc__", invoke(docHandler(token)))
	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func invoke[I Input](method func(I) (any, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var props I

		switch r.Method {
		case http.MethodGet:
			getParamsToStruct(r.URL.Query(), &props)
		case http.MethodPost:
			json.NewDecoder(r.Body).Decode(props)
		}

		data, err := method(props)

		var ResponceErr *Error
		if err != nil {
			ResponceErr = &Error{Message: err.Error(), Code: 400}
		}

		result := Responce{
			Data:  data,
			Error: ResponceErr,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

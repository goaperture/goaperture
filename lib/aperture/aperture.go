package aperture

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/goaperture/goaperture/lib/aperture/doc"
	"github.com/goaperture/goaperture/lib/aperture/types"
)

func NewServer() *types.Aperture {
	return &types.Aperture{
		Mux: http.NewServeMux(),
		Middleware: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		},
		GetSecret: func() string { return "OP3$1$$VF555EJX6GTOBDCP5HPA.ZA7A@CN29k.Kj" },
	}
}

func New[I types.Input, P types.Payload]() func(api *types.Aperture, route types.Route[I, P]) {
	var x = doc.New[I, P]()

	return func(api *types.Aperture, route types.Route[I, P]) {
		api.Mux.HandleFunc(route.Path, invoke(route.Handler, true, api.GetSecret()))
		x.Add(route)
	}
}

func (api *types.Aperture) Run(port int, token *string, clients *[]types.Payload) error {
	api.Mux.HandleFunc("/__doc__", invoke(doc.docHandler(token, clients), false, api.GetSecret()))
	return http.ListenAndServe(":"+strconv.Itoa(port), api.Middleware(api.Mux))
}

func invoke[I types.Input, P types.Payload](method func(I, types.Client[P]) (any, error), wrap bool, secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var props I

		switch r.Method {
		case http.MethodGet:
			getParamsToStruct(r.URL.Query(), &props)
		default:
			json.NewDecoder(r.Body).Decode(&props)
		}

		defer func() {
			if r := recover(); r != nil {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(types.Responce{
					Error: &types.Error{Message: fmt.Sprint(r), Code: 401},
				})
			}
		}()

		data, err := method(props, NewClient[P](r, &w, secret))

		var ResponceErr *types.Error
		if err != nil {
			ResponceErr = &types.Error{Message: err.Error(), Code: 400}
		}

		w.Header().Set("Content-Type", "application/json")

		if !wrap {
			json.NewEncoder(w).Encode(data)
			return
		}

		result := types.Responce{
			Data:  data,
			Error: ResponceErr,
		}

		json.NewEncoder(w).Encode(result)
	}
}

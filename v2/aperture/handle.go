package aperture

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/goaperture/goaperture/lib/params"
)

type Switch struct {
	Handler    func(w http.ResponseWriter, r *http.Request)
	DirectCall func(input any) any
}

func Handle[I Input, O Output](route Route[I, O]) Switch {
	return Switch{
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var input = params.GetInput[I](r)

			var data = route.Handler(context.Background(), input)

			w.Header().Set("Content-Type", "application/json")
			result := Responce{
				Data: data,
			}

			json.NewEncoder(w).Encode(result)
		},
		DirectCall: func(input any) any {
			if v, ok := input.(I); ok {
				var data = route.Handler(context.Background(), v)
				return data
			}

			return nil
		},
	}
}

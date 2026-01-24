package doc

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/goaperture/goaperture/lib/aperture"
	"github.com/goaperture/goaperture/lib/params"
)

type DocInput struct {
}

type RouteHandler func(w http.ResponseWriter, r *http.Request)

func Handle() RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		var input = params.GetInput[DocInput](r)

		var data = route.Handler(context.Background(), input)

		w.Header().Set("Content-Type", "application/json")
		result := aperture.Responce{
			Data: data,
		}

		json.NewEncoder(w).Encode(result)
	}
}

func getDocs() {}

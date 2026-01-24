package aperture

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/goaperture/goaperture/lib/aperture"
	"github.com/goaperture/goaperture/lib/params"
)

type DocInput struct {
	Token string `json:"token"`
}

type RouteHandler func(w http.ResponseWriter, r *http.Request)

const (
	DOC_URL = "__doc__"
)

func docHandle[P Payload](api *Api[P]) RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		input := params.GetInput[DocInput](r)

		if input.Token != api.Secret.Token {
			http.Error(w, "invalid token", 401)
			return
		}

		data := getDocs(api.Routes)

		result := aperture.Responce{
			Data: data,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func getDocs(routes Routes) []aperture.DocOutput {
	var result = []aperture.DocOutput{}

	for path, route := range routes {
		dump := route.PrepareCall()

		result = append(result, aperture.DocOutput{
			Url:         path,
			Input:       dump.Inputs,
			Output:      dump.Outputs,
			Exceptions:  dump.Errors,
			Description: dump.Description,
			Pathprops:   strings.Split(path, "/"),
		})
	}

	return result
}

package aperture

import (
	"encoding/json"
	"net/http"

	"github.com/goaperture/goaperture/v2/api/auth"
	"github.com/goaperture/goaperture/v2/api/auth/auth_paths"
	"github.com/goaperture/goaperture/v2/api/params"
	"github.com/goaperture/goaperture/v2/exception"
	"github.com/goaperture/goaperture/v2/responce"
	"github.com/goaperture/goaperture/v2/ws/aperture"
)

type DocOutput struct {
	Url         string   `json:"url"`
	Type        string   `json:"type"`
	Method      string   `json:"method"`
	Input       any      `json:"inputs,omitempty"`
	Output      any      `json:"outputs,omitempty"`
	Errors      []string `json:"errors,omitempty"`
	Description string   `json:"description"`
	AccessKey   string   `json:"accessKey,omitempty"`
}

type DocResult struct {
	Schema  []DocOutput `json:"schema"`
	Version int         `json:"version"`
}

type DocInput struct {
	Token   string `json:"token"`
	Version int    `json:"version"`
}

type RouteHandler func(w http.ResponseWriter, r *http.Request)

const (
	DOC_URL = "/__doc__"
)

func docHandle[P Payload](api *Api[P]) RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		defer exception.Catch(&w)

		input := params.GetInput[DocInput](r)

		if input.Token != api.Token {
			http.Error(w, "invalid token", 401)
			return
		}

		schema := getDocs(api.routes, api.ws)

		if api.Auth != nil {
			schema = append(schema, getAuthDocs()...)
		}

		w.Header().Set("Content-Type", "application/json")

		if input.Version == 2 {
			json.NewEncoder(w).Encode(convertToV2(&schema))
			return
		}

		json.NewEncoder(w).Encode(DocResult{
			Schema:  schema,
			Version: 3,
		})
	}
}

func getDocs(routes Routes, ws *aperture.WebSockets) []DocOutput {
	var result = []DocOutput{}

	for path, route := range routes {
		dump := route.PrepareCall()

		var accessKey string

		if route.PrivateAccess {
			accessKey = auth.GetAccessKeyFromUrl(path)
		}

		var method = "post"
		if route.Method != "" {
			method = route.Method
		}

		result = append(result, DocOutput{
			Url:         path,
			Type:        "rest",
			Input:       dump.Inputs,
			Output:      dump.Outputs,
			Errors:      dump.Errors,
			Description: route.Description,
			AccessKey:   accessKey,
			Method:      method,
		})
	}

	if ws == nil {
		return result
	}

	for path, socket := range *ws {
		var accessKey string

		if socket.PrivateAccess {
			accessKey = auth.GetAccessKeyFromUrl(path)
		}

		shema := "ws"
		if socket.Sequre {
			shema = "wss"
		}

		result = append(result, DocOutput{
			Url:         path,
			Type:        "ws",
			Description: socket.Description,
			AccessKey:   accessKey,
			Method:      shema,
		})
	}

	return result
}

func getAuthDocs() []DocOutput {
	return []DocOutput{
		{
			Url:         auth_paths.LOGIN,
			Type:        "rest",
			Description: "Получить Access Token",
			Method:      "POST",
			Input: []any{
				auth.LoginInput{},
			},
			Output: []any{
				auth.LoginOutput{},
			},
		},
		{
			Url:         auth_paths.LOGOUT,
			Type:        "rest",
			Description: "Заблокировать Access Token",
			Method:      "POST",
			Output: []any{
				responce.SuccessType{},
			},
		},
		{
			Url:         auth_paths.REFRESH,
			Type:        "rest",
			Description: "Обновить Access Token",
			Method:      "POST",
			Output: []any{
				auth.LoginOutput{},
			},
		},
	}
}

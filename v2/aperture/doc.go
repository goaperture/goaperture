package aperture

import (
	"encoding/json"
	"net/http"

	"github.com/goaperture/goaperture/v2/auth"
	"github.com/goaperture/goaperture/v2/auth/auth_paths"
	"github.com/goaperture/goaperture/v2/exception"
	"github.com/goaperture/goaperture/v2/params"
	"github.com/goaperture/goaperture/v2/responce"
)

type DocOutput struct {
	Url         string   `json:"url"`
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
	Token string `json:"token"`
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

		data := getDocs(api.routes)

		if api.Auth != nil {
			data = append(data, getAuthDocs()...)
		}

		result := Responce{
			Data: data,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func getDocs(routes Routes) []DocOutput {
	var result = []DocOutput{}

	for path, route := range routes {
		dump := route.PrepareCall()

		var accessKey string

		if route.PrivateAccess {
			accessKey = auth.GetAccessKeyFromUrl(path)
		}

		result = append(result, DocOutput{
			Url:         path,
			Input:       dump.Inputs,
			Output:      dump.Outputs,
			Errors:      dump.Errors,
			Description: route.Description,
			AccessKey:   accessKey,
			Method:      route.Method,
		})
	}

	return result
}

func getAuthDocs() []DocOutput {
	return []DocOutput{
		{
			Url:         auth_paths.LOGIN,
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
			Description: "Заблокировать Access Token",
			Method:      "POST",
			Output: []any{
				responce.SuccessType{},
			},
		},
		{
			Url:         auth_paths.REFRESH,
			Description: "Обновить Access Token",
			Method:      "POST",
			Output: []any{
				auth.LoginOutput{},
			},
		},
	}
}

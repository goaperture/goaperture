package aperture

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/goaperture/goaperture/v2/auth"
	"github.com/goaperture/goaperture/v2/auth/auth_paths"
	"github.com/goaperture/goaperture/v2/exception"
	"github.com/goaperture/goaperture/v2/params"
)

type DocOutput struct {
	Url         string          `json:"url"`
	Version     string          `json:"version"`
	Method      string          `json:"method"`
	Alias       string          `json:"alias"`
	Input       any             `json:"inputType,omitempty"`
	Output      any             `json:"outputType,omitempty"`
	Errors      []string        `json:"errors,omitempty"`
	Description string          `json:"description"`
	AccessKey   auth.Permission `json:"accessKey,omitempty"`
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

		result = append(result, DocOutput{
			Url:         path,
			Input:       dump.Inputs,
			Output:      dump.Outputs,
			Errors:      dump.Errors,
			Description: dump.Description,
			Alias:       getAliasFromUrl(path),
			AccessKey:   "",
			Version:     "v2",
			Method:      "POST",
		})
	}

	return result
}

func getAliasFromUrl(url string) string {
	return strings.ToLower(strings.ReplaceAll(url, "/", "_"))
}

func getAuthDocs() []DocOutput {
	return []DocOutput{
		{
			Url:         auth_paths.LOGIN,
			Alias:       getAliasFromUrl(auth_paths.LOGIN),
			Description: "Получить Access Token",
			Version:     "v1",
			Method:      "POST",
			Input: []any{
				auth.LoginInput{},
			},
		},
		{
			Url:         auth_paths.LOGOUT,
			Alias:       getAliasFromUrl(auth_paths.LOGOUT),
			Description: "Заблокировать Access Token",
			Version:     "v1",
			Method:      "POST",
		},
		{
			Url:         auth_paths.REFRESH,
			Alias:       getAliasFromUrl(auth_paths.REFRESH),
			Description: "Обновить Access Token",
			Version:     "v1",
			Method:      "POST",
		},
	}
}

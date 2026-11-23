package aperture

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goaperture/goaperture/v2/api/auth"
	"github.com/goaperture/goaperture/v2/api/auth/auth_paths"
	"github.com/goaperture/goaperture/v2/api/params"
	"github.com/goaperture/goaperture/v2/api/tstree"
	"github.com/goaperture/goaperture/v2/exception"
	"github.com/goaperture/goaperture/v2/response"
	"github.com/goaperture/goaperture/v2/ws/aperture"
)

type DocOutput struct {
	Url                  string         `json:"url"`
	Type                 string         `json:"type"`
	Method               string         `json:"method"`
	Input                any            `json:"inputs,omitempty"`
	Output               any            `json:"outputs,omitempty"`
	Errors               []string       `json:"errors,omitempty"`
	Description          string         `json:"description"`
	AccessKey            string         `json:"accessKey,omitempty"`
	AccessKeyDescription string         `json:"accessKeyDescription,omitempty"`
	TopicDocs            map[string]any `json:"topics,omitempty"`
	WithPagination       bool           `json:"pagination,omitempty"`
	Types                Types          `json:"types,omitempty"`
	TS                   *tstree.Result `json:"ts,omitempty"`
}

type DocResult struct {
	Schema  []DocOutput    `json:"schema"`
	Version int            `json:"version"`
	Token   string         `json:"token,omitempty"`
	Named   []tstree.Named `json:"named,omitempty"`
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

		reg := tstree.NewRegistry()

		schema := getDocs(api.routes, api.ws, api.AccessPrefix, api.Auth.GetDevToken(), reg)

		if api.Auth.Login != nil {
			schema = append(schema, getAuthDocs()...)
		}

		w.Header().Set("Content-Type", "application/json")

		token := api.Auth.GetDevToken()

		if input.Version == 2 {
			json.NewEncoder(w).Encode(convertToV2(&schema, token, reg))
			return
		}

		json.NewEncoder(w).Encode(DocResult{
			Schema:  schema,
			Version: 3,
			Token:   token,
			Named:   reg.Named(),
		})
	}
}

func getDocs(routes Routes, ws *aperture.WebSockets, accessPrefix string, devToken string, reg *tstree.Registry) []DocOutput {
	var result = []DocOutput{}

	fmt.Println("DOC: ")
	for path, route := range routes {
		fmt.Println("\t", route.Description)
		dump := route.PrepareCall(devToken, reg)

		var accessKey string
		var accessDescription string

		if route.PrivateAccess {
			accessKey = auth.GetAccessKeyFromUrl(path, accessPrefix)
			accessDescription = route.AccessDescription
		}

		var method = "query"
		if route.Method != "" {
			method = route.Method
		}

		alias, _ := getAlias(path)

		result = append(result, DocOutput{
			Url:    path,
			Type:   "rest",
			Input:  dump.Inputs,
			Output: dump.Outputs,
			Errors: dump.Errors,
			TS:     dump.TS,

			// Description — описание эндпоинта.
			// Укажите целевое действие в инфинитиве (например: получить..., создать..., обновить...).
			Description:          route.Description,
			AccessKey:            accessKey,
			AccessKeyDescription: accessDescription,
			Method:               method,
			WithPagination:       dump.WithPagination,
			Types:                route.Types.GetWithDefault(alias),
		})
	}

	if ws == nil {
		return result
	}

	for path, socket := range *ws {
		var accessKey string

		if socket.PrivateAccess {
			accessKey = auth.GetAccessKeyFromUrl(path, accessPrefix)
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
			TopicDocs:   socket.TopicDocs,
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
			Types: Types{
				Input:  "LoginInput",
				Output: "LoginSuccess",
			},
			Method: "POST",
			Input: []any{
				auth.LoginInput{},
			},
			Output: []any{
				auth.LoginOutput{},
			},
		},
		{
			Url: auth_paths.LOGOUT,
			Types: Types{
				Input:  "LogoutInput",
				Output: "LogoutSuccess",
			},
			Type:        "rest",
			Description: "Заблокировать Access Token",
			Method:      "POST",
			Output: []any{
				response.SuccessType{},
			},
		},
		{
			Url: auth_paths.REFRESH,
			Types: Types{
				Input:  "TokenRefresh",
				Output: "TokenRefreshSuccess",
			},
			Type:        "rest",
			Description: "Обновить Access Token",
			Method:      "POST",
			Output: []any{
				auth.LoginOutput{},
			},
		},
	}
}

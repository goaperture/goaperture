package doc

import (
	"errors"
	"strings"

	"github.com/goaperture/goaperture/lib/aperture/types"
)

type Doc[T types.Input, P types.Payload] struct {
	Routes []types.TestItem[P]
}

func New[T types.Input, P types.Payload]() Doc[T, P] {
	return Doc[T, P]{
		Routes: []types.TestItem[P]{},
	}
}

func (docs Doc[T, P]) Add(route types.Route[T, P]) {
	docs.Routes = append(docs.Routes, types.TestItem[P]{
		Path: route.Path,
		Test: func(client types.Client[P]) types.TestData {
			data := types.TestData{}

			route.Test(func(input T) {
				data.Inputs = append(data.Inputs, input)

				output, err := route.Handler(input, client)
				if err == nil {
					data.Outputs = append(data.Outputs, output)
				}

			})

			return data
		},
	})
}

type DocOutput struct {
	Method     string   `json:"method"`
	Url        string   `json:"url"`
	Alias      string   `json:"alias"`
	Input      any      `json:"inputType"`
	Output     any      `json:"outputType"`
	Pathprops  []string `json:"pathProps"`
	Exceptions []string `json:"exceptions"`
}

type DocResult struct {
	Schema  any `json:"schema"`
	Version int `json:"version"`
}

func docHandler(token *string, clients *[]types.Payload) func(input types.DocInput, client types.Client[types.Payload]) (any, error) {
	return func(input types.DocInput, client types.Client[types.Payload]) (any, error) {

		if token != nil && input.Token != *token {
			return nil, errors.New("invalid token")
		}

		schema := []DocOutput{}

		for _, test := range routes {
			data := test.Test()
			alias := getAlias(test.Path)

			schema = append(schema, DocOutput{
				Url:        test.Path,
				Method:     "post",
				Alias:      alias,
				Input:      map[string]any{alias + "__TYPE__": data.Inputs},
				Output:     map[string]any{alias + "__TYPE__": data.Outputs},
				Pathprops:  []string{},
				Exceptions: []string{},
			})
		}

		result := DocResult{
			Schema:  schema,
			Version: 2,
		}

		return result, nil
	}
}

func getAlias(path string) string {
	alias := ""
	nextUp := false

	for index, char := range strings.TrimPrefix(path, "/") {
		if index == 0 || nextUp {
			alias += strings.ToUpper(string(char))
			nextUp = false
			continue
		}
		if char == '/' || char == '-' {
			nextUp = true
			continue
		}

		nextUp = false
		alias += string(char)
	}

	return alias
}

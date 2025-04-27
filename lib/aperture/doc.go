package aperture

import (
	"errors"
	"strings"
)

type TestData struct {
	Inputs  []any
	Outputs []any
}

type TestItem struct {
	Path string
	Test func() TestData
}

type DocInput struct {
	Token string `json:"token"`
}

var routes = []TestItem{}

func newDoc[T Input, C client](route Route[T, C]) {
	routes = append(routes, TestItem{
		Path: route.Path,
		Test: func() TestData {
			data := TestData{}

			route.Test(func(input T) {
				data.Inputs = append(data.Inputs, input)
				output, err := route.Handler(input, nil)
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

func docHandler(token *string) func(input DocInput, client *client) (any, error) {
	return func(input DocInput, client *client) (any, error) {

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

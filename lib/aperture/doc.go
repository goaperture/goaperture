package aperture

import (
	"errors"
	"fmt"
	"strings"
)

type TestData struct {
	Inputs     []any
	Outputs    []any
	Exceptions []string
}

type TestItem struct {
	Path string
	Test func(secret string) TestData
}

type DocInput struct {
	Token string `json:"token"`
}

var routes = []TestItem{}

func newDoc[T Input, P Payload](route Route[T, P]) {
	routes = append(routes, TestItem{
		Path: route.Path,
		Test: func(secret string) TestData {
			data := TestData{}

			route.Test(func(input T) {
				defer func() {
					if r := recover(); r != nil {
						data.Exceptions = append(
							data.Exceptions,
							fmt.Sprintf("panic: %v", r),
						)
					}
				}()

				data.Inputs = append(data.Inputs, input)
				output, err := route.Handler(input, NewClient[P](nil, nil, secret, true))
				if err == nil {
					data.Outputs = append(data.Outputs, output)
				}
			})

			return data
		},
	})
}

type DocOutput struct {
	Url        string   `json:"url"`
	Version    string   `json:"version"`
	Alias      string   `json:"alias"`
	Method     string   `json:"method"`
	Input      any      `json:"inputType"`
	Output     any      `json:"outputType"`
	Pathprops  []string `json:"pathProps"`
	Exceptions []string `json:"exceptions"`
}

type DocResult struct {
	Schema  any `json:"schema"`
	Version int `json:"version"`
}

func docHandler[P Payload](token string, testClients []P) func(input DocInput, client Client[P]) (any, error) {
	return func(input DocInput, client Client[P]) (any, error) {

		if input.Token != token {
			return nil, errors.New("invalid token")
		}

		schema := []DocOutput{}

		for _, test := range routes {
			func() {
				alias, version := getAlias(test.Path)

				data := test.Test("sdkjflskdjflkjl")

				schema = append(schema, DocOutput{
					Url:        test.Path,
					Version:    version,
					Method:     "post",
					Alias:      alias,
					Input:      map[string]any{alias + "__TYPE__": data.Inputs},
					Output:     map[string]any{alias + "__TYPE__": data.Outputs},
					Pathprops:  []string{},
					Exceptions: []string{},
				})
			}()

		}

		result := DocResult{
			Schema:  schema,
			Version: 2,
		}

		return result, nil
	}
}

func getAlias(path string) (string, string) {
	alias := ""
	nextUp := false

	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	version := parts[0]
	url := strings.Join(parts[1:], "/")

	for index, char := range url {
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

	return alias, version
}

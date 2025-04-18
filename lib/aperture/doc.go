package aperture

import (
	"errors"
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

func newDoc[T Input](route Route[T]) {
	routes = append(routes, TestItem{
		Path: route.Path,
		Test: func() TestData {
			data := TestData{}

			route.Test(func(input T) {
				data.Inputs = append(data.Inputs, input)
				output, err := route.Handler(input)
				if err == nil {
					data.Outputs = append(data.Outputs, output)
				}
			})

			return data
		},
	})
}

type DocOutput struct {
	Path string   `json:"path"`
	Data TestData `json:"data"`
}

func docHandler(token *string) func(input DocInput) (any, error) {
	return func(input DocInput) (any, error) {
		if token != nil && input.Token != *token {
			return nil, errors.New("Invalid token")
		}

		result := []DocOutput{}

		for _, test := range routes {
			result = append(result, DocOutput{
				Path: test.Path,
				Data: test.Test(),
			})
		}

		return result, nil
	}
}

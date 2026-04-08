package aperture

import (
	"strings"

	"github.com/goaperture/goaperture/v2/responce"
)

type DocOutputV2 struct {
	Url         string         `json:"url"`
	Version     string         `json:"version"`
	Alias       string         `json:"alias"`
	Method      string         `json:"method"`
	Input       any            `json:"inputType"`
	Output      any            `json:"outputType"`
	Pathprops   []string       `json:"pathProps"`
	Exceptions  []string       `json:"exceptions"`
	Description string         `json:"description"`
	Topics      map[string]any `json:"topics,omitempty"`
}

type DocResultV2 struct {
	Schema  []DocOutputV2 `json:"schema"`
	Version int           `json:"version"`
}

func convertToV2(doc *[]DocOutput) DocResultV2 {
	schema := []DocOutputV2{}

	for _, route := range *doc {
		alias, version := getAlias(route.Url)

		var input any
		var output any

		if route.Input != nil {
			input = map[string]any{
				alias + "Input___TYPE__": route.Input,
			}

		}

		if route.Method == "ws" {
			input = map[string]any{
				alias + "Topic___TYPE__": route.TopicDocs,
			}
		}

		if route.Output != nil {
			if route.WithPagination {
				output = map[string][]responce.Responce{
					alias + "Output___TYPE__": getOutputWithPagination(route.Output),
				}
			} else {
				output = map[string]any{
					alias + "Output___TYPE__": route.Output,
				}
			}
		}

		schema = append(schema, DocOutputV2{
			Url:         route.Url,
			Version:     version,
			Alias:       alias,
			Method:      route.Method,
			Input:       input,
			Output:      output,
			Exceptions:  route.Errors,
			Description: route.Description,
			Topics:      route.TopicDocs,
		})
	}

	return DocResultV2{
		Schema:  schema,
		Version: 2,
	}
}

func getOutputWithPagination(output any) []responce.Responce {
	var result []responce.Responce

	if array, ok := output.([]any); ok {
		for _, val := range array {
			result = append(result, responce.Responce{
				Data:       val,
				Pagination: responce.Pagination{},
			})
		}
	}

	return result
}

func getAlias(path string) (string, string) {
	alias := ""
	nextUp := false

	var versionIndex = 0

	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if parts[0] == "ws" {
		versionIndex = 1
	}

	version := parts[versionIndex]
	url := strings.Join(parts[versionIndex+1:], "/")

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

package collector

import (
	"context"
	"fmt"
)

type Collector[I any, O any] struct {
	Method      string
	Description string
	Handler     func(context.Context, I) O
	Inputs      []I
	Outputs     []any
	Errors      []string
}

func (c *Collector[I, O]) Execute(input I) {
	defer func() {
		if r := recover(); r != nil {
			c.Errors = append(c.Errors, fmt.Sprintf("%s", r))
		}
	}()

	c.Inputs = append(c.Inputs, input)
	output := c.Handler(context.Background(), input)
	c.Outputs = append(c.Outputs, output)
}

func (c *Collector[I, O]) Expect(output any) {
	c.Outputs = append(c.Outputs, output)
}

type RouteDump struct {
	Method      string
	Description string
	Inputs      []any
	Outputs     []any
	Errors      []string
}

func (c *Collector[I, O]) GetDump() RouteDump {
	return RouteDump{
		Method:      c.Method,
		Description: c.Description,
		Inputs:      convertToAny(c.Inputs),
		Outputs:     convertToAny(c.Outputs),
		Errors:      c.Errors,
	}
}

func convertToAny[T any](input []T) []any {
	result := make([]any, len(input))

	for key, value := range input {
		result[key] = value
	}

	return result
}

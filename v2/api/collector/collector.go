package collector

import (
	"context"
	"fmt"

	"github.com/goaperture/goaperture/v2/api/client"
)

type Collector[I any, O any] struct {
	Method         string
	Description    string
	Handler        func(context.Context, I) O
	Inputs         []I
	Outputs        []any
	Errors         []string
	WithPagination bool
}

func (c *Collector[I, O]) Execute(input I) *Collector[I, O] {
	defer func() {
		if r := recover(); r != nil {
			c.Errors = append(c.Errors, fmt.Sprintf("%s", r))
		}
	}()

	ctx := client.WithPagination(context.Background())

	c.Inputs = append(c.Inputs, input)
	output := c.Handler(ctx, input)
	c.Outputs = append(c.Outputs, output)
	c.WithPagination = client.GetPagination(ctx).Use

	return c
}

func (c *Collector[I, O]) Entry(input I) *Collector[I, O] {
	c.Inputs = append(c.Inputs, input)

	return c
}

func (c *Collector[I, O]) Expect(output any) *Collector[I, O] {
	c.Outputs = append(c.Outputs, output)

	return c
}

func (c *Collector[I, O]) ExpectArray(output any) *Collector[I, O] {
	c.Outputs = append(c.Outputs, []any{output})

	return c
}

//
// DUMP
//

type RouteDump struct {
	Method         string
	Description    string
	AccessKey      string
	Inputs         []any
	Outputs        []any
	Errors         []string
	WithPagination bool
}

func (c *Collector[I, O]) GetDump() RouteDump {
	return RouteDump{
		Method:         c.Method,
		Description:    c.Description,
		Inputs:         convertToAny(c.Inputs),
		Outputs:        convertToAny(c.Outputs),
		Errors:         c.Errors,
		WithPagination: c.WithPagination,
	}
}

func convertToAny[T any](input []T) []any {
	result := make([]any, len(input))

	for key, value := range input {
		result[key] = value
	}

	return result
}

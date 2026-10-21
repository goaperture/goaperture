package sse

import (
	"context"
)

type SSEKey struct{}

type RequestSSE struct {
	Key string
	Use bool
}

func With(ctx context.Context) context.Context {
	return context.WithValue(ctx, SSEKey{}, new(RequestSSE))
}

func Get(ctx context.Context) *RequestSSE {
	val := ctx.Value(SSEKey{})
	if req, ok := val.(*RequestSSE); ok {
		return req
	}

	return nil
}

package client

import (
	"context"
	"net/http"
)

type requestkey struct{}

func WithRequest(ctx context.Context, request *http.Request) context.Context {
	return context.WithValue(ctx, requestkey{}, request)
}

func GetRequest(ctx context.Context) *http.Request {
	val := ctx.Value(requestkey{})
	if req, ok := val.(*http.Request); ok {
		return req
	}
	return nil
}

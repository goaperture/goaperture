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

type responsekey struct{}

func WithResponse(ctx context.Context, response *http.ResponseWriter) context.Context {
	return context.WithValue(ctx, responsekey{}, response)
}

func GetResponse(ctx context.Context) http.ResponseWriter {
	val := ctx.Value(responsekey{})
	if req, ok := val.(*http.ResponseWriter); ok {
		return *req
	}
	return nil
}

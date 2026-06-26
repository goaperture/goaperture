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

type responcekey struct{}

func WithResponce(ctx context.Context, responce *http.ResponseWriter) context.Context {
	return context.WithValue(ctx, responcekey{}, responce)
}

func GetResponce(ctx context.Context) *http.ResponseWriter {
	val := ctx.Value(responcekey{})
	if req, ok := val.(*http.ResponseWriter); ok {
		return req
	}
	return nil
}

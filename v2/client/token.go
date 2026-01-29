package client

import "context"

type tokenkey struct{}

func WithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenkey{}, &token)
}

func GetToken(ctx context.Context) string {
	return ctx.Value(tokenkey{}).(string)
}

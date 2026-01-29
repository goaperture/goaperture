package client

import (
	"context"

	"github.com/goaperture/goaperture/v2/auth"
)

type Client[P any] struct {
	Token   string
	payload *P
	read    bool
}

type key struct{}

func With[P any](ctx context.Context, payload Client[P]) context.Context {
	return context.WithValue(ctx, key{}, &payload)
}

func Get[P any](ctx context.Context) *Client[P] {
	return ctx.Value(key{}).(*Client[P])
}

func New[P any](token string) Client[P] {
	return Client[P]{Token: token}
}

func (c *Client[P]) GetPayload() *P {
	if c.payload != nil {
		return c.payload
	}

	if c.read {
		return nil
	}

	c.payload = auth.GetPayloadFromJwt[P](c.Token)
	c.read = true

	return c.payload
}

package aperture

import (
	"context"

	"github.com/goaperture/goaperture/v3/api/auth"
	"github.com/goaperture/goaperture/v3/api/client"
	"github.com/goaperture/goaperture/v3/exception"
)

func (a *Api[P]) GetClientPayload(ctx context.Context) *P {
	if a.Auth == nil {
		return nil
	}

	token := client.GetToken(ctx)
	if token == nil {
		return nil
	}

	payload := auth.GetPayloadFromJwt[P](*token, a.Auth.GetSecret())
	return payload
}

func (a *Api[P]) GetClientPayloadX(ctx context.Context) *P {
	payload := a.GetClientPayload(ctx)

	if payload == nil {
		exception.NotAccess("Не удалось получить доступ")
	}

	return payload
}

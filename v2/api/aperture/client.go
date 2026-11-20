package aperture

import (
	"context"
	"fmt"

	"github.com/goaperture/goaperture/v2/api/auth"
	"github.com/goaperture/goaperture/v2/api/client"
	"github.com/goaperture/goaperture/v2/exception"
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

func (a *Api[P]) CreateClientAccessToken(ctx context.Context, id auth.ID) (string, error) {
	if a.Auth == nil {
		return "", fmt.Errorf("Auth is not enabled")
	}

	w := a.Response(ctx)
	return a.Auth.BindRefreshAndGetAccessTokens(&w, id), nil
}

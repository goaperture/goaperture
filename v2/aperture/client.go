package aperture

import (
	"context"

	"github.com/goaperture/goaperture/v2/auth"
	"github.com/goaperture/goaperture/v2/client"
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

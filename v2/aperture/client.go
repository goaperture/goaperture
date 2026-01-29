package aperture

import (
	"context"

	"github.com/goaperture/goaperture/v2/auth"
	"github.com/goaperture/goaperture/v2/client"
)

func (a *Api[P]) GetClientPayload(ctx context.Context) *P {
	token := client.GetToken(ctx)
	if token == nil {
		return nil
	}

	payload := auth.GetPayloadFromJwt[P](*token)
	return payload
}

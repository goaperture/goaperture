package aperture

import (
	"context"
	"net/http"

	"github.com/goaperture/goaperture/v2/client"
)

func (a *Api[P]) GetClient(ctx context.Context) *P {
	return client.Get[P](ctx).GetPayload()
}

func GetBearerToken(r *http.Request) {

}

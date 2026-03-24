package aperture

import (
	"context"
	"net/http"

	"github.com/goaperture/goaperture/v2/api/client"
)

func (a *Api[P]) Request(ctx context.Context) *http.Request {
	return client.GetRequest(ctx)
}

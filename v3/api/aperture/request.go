package aperture

import (
	"context"
	"net/http"

	"github.com/goaperture/goaperture/v3/api/client"
)

func (a *Api[P]) Request(ctx context.Context) *http.Request {
	return client.GetRequest(ctx)
}

func (a *Api[P]) Response(ctx context.Context) http.ResponseWriter {
	return client.GetResponse(ctx)
}

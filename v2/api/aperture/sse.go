package aperture

import (
	"context"

	"github.com/goaperture/goaperture/v2/api/sse"
)

func (a *Api[P]) SetStreamKey(ctx context.Context, key string) {
	request := sse.Get(ctx)

	if request == nil {
		return
	}

	request.Use = true
	request.Key = key
}

package aperture

import (
	"context"

	"github.com/goaperture/goaperture/v2/api/collector"
)

type Input interface {
	any
}
type Output interface {
	any
}

type Handler[I Input, O Output] = func(ctx context.Context, input I) O

type T[P Input] interface {
	Execute(input P)
	Expect(r any)
}

type CL[P Input, O Output] = collector.Collector[P, O]
type Prepare[P Input, O Output] = func(collector *CL[P, O])

type Route[I Input, O Output] struct {
	Method        string
	Handler       Handler[I, O]
	PrivateAccess bool
	Description   string
	Prepare       Prepare[I, O]
}

func GetPayload[P any](ctx context.Context) (*P, bool) {
	return nil, true
}

type ApiServer interface {
	UseRoutes(r *Routes)
}

type Routes map[string]Switch

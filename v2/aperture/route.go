package aperture

import (
	"context"
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

type Prepare[P Input] = func(controller T[P])

type Route[I Input, O Output] struct {
	Method        string
	Handler       Handler[I, O]
	PrivateAccess bool
	Description   string
	Prepare       Prepare[I]
}

func GetPayload[P any](ctx context.Context) (*P, bool) {
	return nil, true
}

type Routes = map[string]Switch

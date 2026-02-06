package Hello

import (
	"context"

	"github.com/goaperture/goaperture/v2/api/aperture"
)

type PrivateWorldInput struct{}
type PrivateWorldOutput interface{ any }

var PrivateWorld = aperture.Route[PrivateWorldInput, PrivateWorldOutput]{
	Description:   "Тут нужен приватный доступ",
	PrivateAccess: true,
	Handler:       PrivateWorldHandler,
	Prepare: func(cl *aperture.CL[PrivateWorldInput, PrivateWorldOutput]) {
		cl.Execute(PrivateWorldInput{})
	},
}

func PrivateWorldHandler(ctx context.Context, input PrivateWorldInput) PrivateWorldOutput {
	return "Private world"
}

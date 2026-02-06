package Hello

import (
	"context"

	"github.com/goaperture/goaperture/v2/api/aperture"
)

type HelloWorldInput struct{}
type HelloWorldOutput interface{ any }

var HelloWorld = aperture.Route[HelloWorldInput, HelloWorldOutput]{
	Description: "HW",
	Handler:     HelloWorldHandler,
	Prepare: func(cl *aperture.CL[HelloWorldInput, HelloWorldOutput]) {
		cl.Execute(HelloWorldInput{})
	},
}

func HelloWorldHandler(ctx context.Context, input HelloWorldInput) HelloWorldOutput {
	return "Hello World"
}

package aperture

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/goaperture/goaperture/v2/auth"
	"github.com/goaperture/goaperture/v2/client"
	"github.com/goaperture/goaperture/v2/collector"
	"github.com/goaperture/goaperture/v2/exception"
	"github.com/goaperture/goaperture/v2/params"
)

type TempPayload struct {
	auth.Permissions `json:"permissions"`
}

type Switch struct {
	Handler       func(secret auth.XSecret) func(w http.ResponseWriter, r *http.Request)
	DirectCall    func(input any) any
	PrepareCall   func() collector.RouteDump
	PrivateAccess bool
	Description   string
	Method        string
}

func Handle[I Input, O Output](route Route[I, O]) Switch {
	return Switch{
		Handler: func(secret auth.XSecret) func(w http.ResponseWriter, r *http.Request) {
			return func(w http.ResponseWriter, r *http.Request) {
				defer exception.Catch(&w)

				jwt, exists := auth.ParseAccessToken(r)

				if route.PrivateAccess {
					accessKey := auth.GetAccessKeyFromUrl(r.Pattern)
					if !exists {
						exception.NotAccess(accessKey)
					}

					payload := auth.GetPayloadFromJwt[TempPayload](jwt, secret)
					payload.Permissions.CheckX(accessKey)
				}

				ctx := r.Context()
				if exists {
					ctx = client.WithToken(ctx, jwt)
				}

				var input = params.GetInput[I](r)

				var data = route.Handler(ctx, input)

				w.Header().Set("Content-Type", "application/json")
				result := Responce{
					Data: data,
				}

				json.NewEncoder(w).Encode(result)
			}
		},
		DirectCall: func(input any) any {
			if v, ok := input.(I); ok {
				var data = route.Handler(context.Background(), v)
				return data
			}

			return nil
		},
		PrepareCall: func() collector.RouteDump {

			var cll = collector.Collector[I, O]{
				Handler: route.Handler,
			}

			route.Prepare(&cll)

			return cll.GetDump()
		},
		PrivateAccess: route.PrivateAccess,
		Description:   route.Description,
		Method:        route.Method,
	}
}

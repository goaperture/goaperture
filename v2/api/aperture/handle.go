package aperture

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/goaperture/goaperture/v2/api/auth"
	"github.com/goaperture/goaperture/v2/api/client"
	"github.com/goaperture/goaperture/v2/api/collector"
	"github.com/goaperture/goaperture/v2/api/params"
	"github.com/goaperture/goaperture/v2/exception"
)

func Handle[I Input, O Output](route *Route[I, O]) Switch {
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

					payload := auth.GetPayloadFromJwt[auth.TempPayload](jwt, secret)
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

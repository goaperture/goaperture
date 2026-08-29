package aperture

import (
	"encoding/json"
	"net/http"

	"github.com/goaperture/goaperture/v2/api/auth"
	"github.com/goaperture/goaperture/v2/api/client"
	"github.com/goaperture/goaperture/v2/api/params"
	"github.com/goaperture/goaperture/v2/exception"
)

type HTTPHandler func(w http.ResponseWriter, r *http.Request)
type ConfigHandler func(secret auth.XSecret, accessPrefix string) HTTPHandler

func mainHandler[I Input, O Output](route *Route[I, O]) ConfigHandler {
	return func(secret auth.XSecret, accessPrefix string) HTTPHandler {

		return func(w http.ResponseWriter, r *http.Request) {
			defer exception.Catch(&w)

			jwt, exists := auth.ParseAccessToken(r)

			if route.PrivateAccess {
				accessKey := auth.GetAccessKeyFromUrl(r.Pattern, accessPrefix)
				if !exists {
					exception.NotAccess(accessKey)
				}

				payload := auth.GetPayloadFromJwt[auth.TempPayload](jwt, secret)
				payload.Permissions.CheckX(accessKey)
			}

			ctx := client.WithRequest(r.Context(), r)
			ctx = client.WithResponce(ctx, &w)
			ctx = client.WithPagination(ctx)

			if exists {
				ctx = client.WithToken(ctx, jwt)
			}

			var input = params.GetInput[I](r)

			var data = route.Handler(ctx, input)

			w.Header().Set("Content-Type", "application/json")

			pagination := client.GetPagination(ctx).Export()

			result := Responce{
				Data:       data,
				Pagination: pagination,
			}

			json.NewEncoder(w).Encode(result)

		}
	}
}

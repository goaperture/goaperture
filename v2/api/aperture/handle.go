package aperture

import (
	"context"

	"github.com/goaperture/goaperture/v2/api/client"
	"github.com/goaperture/goaperture/v2/api/collector"
)

func Handle[I Input, O Output](route *Route[I, O]) Switch {
	return Switch{
		Handler: mainHandler(route),
		DirectCall: func(input any) any {
			if v, ok := input.(I); ok {
				var data = route.Handler(context.Background(), v)
				return data
			}

			return nil
		},
		PrepareCall: func(token string) collector.RouteDump {
			ctx := client.WithToken(context.TODO(), token)

			var cll = collector.Collector[I, O]{
				Context: ctx,
				Handler: route.Handler,
			}

			route.Prepare(&cll)

			return cll.GetDump()
		},
		PrivateAccess: route.PrivateAccess,
		Description:   route.Description,
		Method:        route.Method,
		Types:         route.Types,
	}
}

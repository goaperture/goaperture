package aperture

import (
	"net/http"

	"github.com/goaperture/goaperture/v2/api/auth"
	"github.com/goaperture/goaperture/v2/api/collector"
)

type Switch struct {
	Handler       func(secret auth.XSecret) func(w http.ResponseWriter, r *http.Request)
	DirectCall    func(input any) any
	PrepareCall   func() collector.RouteDump
	PrivateAccess bool
	Description   string
	Method        string
}

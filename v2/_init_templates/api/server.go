package api

import (
	"github.com/goaperture/goaperture/v2/api/aperture"
)

var Server = aperture.Api[Payload]{
	Port:    3003,
	Auth:    &Auth,
	Metrics: true,
}

package aperture

import (
	"github.com/goaperture/goaperture/v2/api/collector"
	"github.com/goaperture/goaperture/v2/api/tstree"
)

type Switch struct {
	Handler           ConfigHandler
	DirectCall        func(input any) any
	PrepareCall       func(token string, reg *tstree.Registry) collector.RouteDump
	PrivateAccess     bool
	AccessDescription string
	// Description — описание эндпоинта.
	// Укажите целевое действие в инфинитиве (например: получить..., создать..., обновить...).
	Description string
	Method      string
	Types       Types
}

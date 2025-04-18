package api

import (
	"github.com/goaperture/goaperture/test/api/routes"
	"github.com/goaperture/goaperture/test/api/routes/filter"
	"github.com/goaperture/goaperture/test/api/routes/filter/connections"
	"github.com/goaperture/goaperture/test/api/routes/home"
	"github.com/goaperture/goaperture/test/api/routes/user"

	api "github.com/goaperture/goaperture/lib/aperture"
)

func Serve(port int, token *string) error {
	api.NewRoute(api.Route[connections.GetConnectionsInput]{
		Handler: connections.GetConnections,
		Path:    "/filter/connections/getconnections",
		Test:    connections.GetConnectionsTest,
	})
	api.NewRoute(api.Route[connections.GetTypesInput]{
		Handler: connections.GetTypes,
		Path:    "/filter/connections/gettypes",
		Test:    connections.GetTypesTest,
	})
	api.NewRoute(api.Route[filter.GetFilterInfoInput]{
		Handler: filter.GetFilterInfo,
		Path:    "/filter/getfilterinfo",
		Test:    filter.GetFilterInfoTest,
	})
	api.NewRoute(api.Route[home.HomeInput]{
		Handler: home.Home,
		Path:    "/home/home",
		Test:    home.HomeTest,
	})
	api.NewRoute(api.Route[routes.IndexInput]{
		Handler: routes.Index,
		Path:    "/index",
		Test:    routes.IndexTest,
	})
	api.NewRoute(api.Route[user.GetInfoInput]{
		Handler: user.GetInfo,
		Path:    "/user/getinfo",
		Test:    user.GetInfoTest,
	})
	return api.Run(port, token)
}

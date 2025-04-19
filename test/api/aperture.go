package api

import (
	"net/http"

	"github.com/goaperture/goaperture/test/api/routes"
	"github.com/goaperture/goaperture/test/api/routes/filter"
	"github.com/goaperture/goaperture/test/api/routes/filter/connections"
	"github.com/goaperture/goaperture/test/api/routes/home"
	"github.com/goaperture/goaperture/test/api/routes/user"

	api "github.com/goaperture/goaperture/lib/aperture"
)

type Server struct {
	aperture api.Aperture
}

func NewServer() *Server {
	return &Server{
		aperture: *api.NewServer(),
	}
}

func (server *Server) Middleware(middleware func(next http.Handler) http.Handler) {
	server.aperture.Middleware = middleware
}

func (server *Server) Run(port int, token *string) error {
	api.NewRoute(&server.aperture, api.Route[connections.GetConnectionsInput]{
		Handler: connections.GetConnections,
		Path:    "/filter/connections/getconnections",
		Test:    connections.GetConnectionsTest,
	})
	api.NewRoute(&server.aperture, api.Route[connections.GetTypesInput]{
		Handler: connections.GetTypes,
		Path:    "/filter/connections/gettypes",
		Test:    connections.GetTypesTest,
	})
	api.NewRoute(&server.aperture, api.Route[filter.GetFilterInfoInput]{
		Handler: filter.GetFilterInfo,
		Path:    "/filter/getfilterinfo",
		Test:    filter.GetFilterInfoTest,
	})
	api.NewRoute(&server.aperture, api.Route[home.HomeInput]{
		Handler: home.Home,
		Path:    "/home/home",
		Test:    home.HomeTest,
	})
	api.NewRoute(&server.aperture, api.Route[routes.IndexInput]{
		Handler: routes.Index,
		Path:    "/index",
		Test:    routes.IndexTest,
	})
	api.NewRoute(&server.aperture, api.Route[user.GetInfoInput]{
		Handler: user.GetInfo,
		Path:    "/user/getinfo",
		Test:    user.GetInfoTest,
	})

	return server.aperture.Run(port, token)
}

func Serve(port int, token *string) {
	if err := NewServer().Run(port, token); err != nil {
		panic(err)
	}
}

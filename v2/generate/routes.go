package generate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dave/jennifer/jen"
)

func GenerateRoutes(app, outPath string) {
	routesPath := filepath.Join(outPath, "routes")

	if _, err := os.Stat(routesPath); os.IsNotExist(err) {
		return
	}

	target := filepath.Join(routesPath, "routes.go")

	var routes []FileRoute
	getRoutesFrom(routesPath, &routes, true)

	// имя пакета, в котором будет сгенерирован файл
	f := jen.NewFile("routes")

	// основной импорт
	f.ImportAlias(
		"github.com/goaperture/goaperture/v2/api/aperture",
		"api",
	)

	// импорты всех роутов
	for _, route := range routes {
		importPath := app + "/" + route.Import
		f.ImportName(importPath, route.Package)
	}

	// собираем map[string]aperture.Handle
	routesDict := jen.Dict{}

	for _, route := range routes {
		path := app + "/" + route.Import
		cutUrl := route.Url[len(routesPath):]

		routesDict[jen.Lit(cutUrl)] = jen.Qual(
			"github.com/goaperture/goaperture/v2/api/aperture",
			"Handle",
		).Call(
			jen.Op("&").Qual(path, route.Route),
		)
	}

	// var Routes = aperture.Routes{ ... }
	f.Var().Id("Routes").Op("=").Qual(
		"github.com/goaperture/goaperture/v2/api/aperture",
		"Routes",
	).Values(
		routesDict,
	)

	if err := f.Save(target); err != nil {
		panic(fmt.Sprintf("error saving file: %v", err))
	}
}

func GenerateWebsockets(app, outPath string) {
	routesPath := filepath.Join(outPath, "ws")

	if _, err := os.Stat(routesPath); os.IsNotExist(err) {
		return
	}

	target := filepath.Join(routesPath, "websockets.go")

	var routes []FileRoute
	getRoutesFrom(routesPath, &routes, true)

	f := jen.NewFile("ws")
	// основной импорт
	f.ImportAlias(
		"github.com/goaperture/goaperture/v2/ws/aperture",
		"sockets",
	)

	// собираем map[string]aperture.Handle
	routesDict := jen.Dict{}

	for _, route := range routes {
		path := app + "/" + route.Import
		cutUrl := "/ws" + route.Url[len(routesPath):]

		routesDict[jen.Lit(cutUrl)] = jen.Qual(
			"github.com/goaperture/goaperture/v2/ws/aperture",
			"Handle",
		).Call(
			jen.Op("&").Qual(path, route.Route),
		)
	}

	// var Routes = aperture.Routes{ ... }
	f.Var().Id("Routes").Op("=").Qual(
		"github.com/goaperture/goaperture/v2/ws/aperture",
		"WebSockets",
	).Values(
		routesDict,
	)

	if err := f.Save(target); err != nil {
		panic(fmt.Sprintf("error saving file: %v", err))
	}
}

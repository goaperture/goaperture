package generate

import (
	"fmt"
	"path/filepath"

	"github.com/dave/jennifer/jen"
)

const (
	targetFile = "aperture.go"
)

func GenerateRoutes(app, outPath string) {

	target := filepath.Join(outPath, targetFile)
	// if _, err := os.Stat(target); err == nil {
	// 	os.Remove(target)
	// }

	var routes []FileRoute
	getRoutesFrom(outPath, &routes, true)

	// имя пакета, в котором будет сгенерирован файл
	f := jen.NewFile("routes")

	// основной импорт
	f.ImportAlias(
		"github.com/goaperture/goaperture/v2/aperture",
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
		cutUrl := route.Url[len(outPath):]

		routesDict[jen.Lit(cutUrl)] = jen.Qual(
			"github.com/goaperture/goaperture/v2/aperture",
			"Handle",
		).Call(
			jen.Qual(path, route.Route),
		)
	}

	// var Routes = aperture.Routes{ ... }
	f.Var().Id("Routes").Op("=").Qual(
		"github.com/goaperture/goaperture/v2/aperture",
		"Routes",
	).Values(
		routesDict,
	)

	if err := f.Save(target); err != nil {
		panic(fmt.Sprintf("error saving file: %v", err))
	}
}

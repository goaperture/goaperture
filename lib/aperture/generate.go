package aperture

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/dave/jennifer/jen"
)

type FileRoute struct {
	Import  string
	Package string
	Route   string
	Test    string
	Type    string
	Url     string
}

func Generate(app string, routes string, _package string) {
	list := []FileRoute{}
	getRoutesFrom(routes, &list)

	f := jen.NewFile(_package)

	f.ImportAlias("github.com/goaperture/goaperture/lib/aperture", "api")

	for _, route := range list {
		path := app + "/" + route.Import
		f.ImportName(path, route.Package)
	}

	serveBody := []jen.Code{}

	for _, route := range list {
		path := app + "/" + route.Import
		cutUrl := route.Url[len(routes):]

		serveBody = append(serveBody, jen.Qual("github.com/goaperture/goaperture/lib/aperture", "NewRoute").Call(
			jen.Qual("github.com/goaperture/goaperture/lib/aperture", "Route").Index(jen.Qual(path, route.Type)).Values(jen.Dict{
				jen.Id("Path"):    jen.Lit(cutUrl),
				jen.Id("Handler"): jen.Qual(path, route.Route),
				jen.Id("Test"):    jen.Qual(path, route.Test),
			}),
		))
	}

	serveBody = append(serveBody, jen.Return(
		jen.Qual("github.com/goaperture/goaperture/lib/aperture", "Run").Call(
			jen.Id("port"),
			jen.Id("token"),
		),
	))

	f.Func().Id("Serve").Params(
		jen.Id("port").Int(),
		jen.Id("token").Op("*").String(),
	).Error().Block(serveBody...)

	if _, err := os.Stat(_package); os.IsNotExist(err) {
		os.Mkdir(_package, 0777)
	}
	if err := f.Save(_package + "/aperture.go"); err != nil {
		panic(fmt.Sprintf("Error saving file: %v", err))
	}
}

func getRoutesFrom(path string, routes *[]FileRoute) {
	content, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Не удалось прочитать ", path)
		return
	}

	for _, entry := range content {
		full := path + "/" + entry.Name()
		if entry.IsDir() {
			getRoutesFrom(full, routes)
			continue
		}

		if !strings.HasSuffix(full, ".go") {
			continue
		}

		name, _ := strings.CutSuffix(entry.Name(), ".go")

		exp := strings.Split(path, "/")
		routePackage := exp[len(exp)-1]
		routeMethod := strings.ToUpper(name[:1]) + name[1:]
		routeTest := routeMethod + "Test"
		routeInputType := routeMethod + "Input"
		routeUrl := getPrettyPath(path + "/" + name)

		*routes = append(*routes, FileRoute{
			Import:  path,
			Route:   routeMethod,
			Package: routePackage,
			Test:    routeTest,
			Type:    routeInputType,
			Url:     routeUrl,
		})
	}
}

func getPrettyPath(path string) string {
	var builder strings.Builder
	pt := false

	for _, ch := range path {
		if unicode.IsUpper(ch) {
			if !pt {
				builder.WriteRune('-')
			}
			builder.WriteRune(unicode.ToLower(ch))

			continue
		}

		pt = ch == '/'

		builder.WriteRune(ch)
	}

	return builder.String()
}

package generate

import (
	"errors"
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

func Schema(app string, routes string, _package string) {
	list := []FileRoute{}
	getRoutesFrom(routes, &list)

	f := jen.NewFile(_package)

	f.ImportAlias("github.com/goaperture/goaperture/lib/aperture", "api")
	f.ImportAlias(app+"/"+_package+"/config", "config")

	for _, route := range list {
		path := app + "/" + route.Import
		f.ImportName(path, route.Package)
	}

	serveBody := []jen.Code{}

	// serveBody = append(serveBody,
	// 	jen.Id("server").Op(":=").Qual("github.com/goaperture/goaperture/lib/aperture", "NewServer").Call(),
	// )

	for _, route := range list {
		path := app + "/" + route.Import
		cutUrl := route.Url[len(routes):]

		serveBody = append(serveBody, jen.Qual("github.com/goaperture/goaperture/lib/aperture", "NewRoute").Call(
			jen.Op("&").Id("server").Dot("aperture"),
			jen.Qual("github.com/goaperture/goaperture/lib/aperture", "Route").Types(jen.Qual(path, route.Type), jen.Qual(app+"/"+_package+"/config", "Payload")).Values(jen.Dict{
				jen.Id("Path"):    jen.Lit(cutUrl),
				jen.Id("Handler"): jen.Qual(path, route.Route),
				jen.Id("Test"):    jen.Qual(path, route.Test),
			}),
		))
	}

	serveBody = append(serveBody, jen.Return(
		jen.Id("server").Dot("aperture").Dot("Run").Call(
			jen.Id("port"),
			jen.Id("token"),
		),
	))

	f.Func().Params(jen.Id("server").Op("*").Id("Server")).
		Id("Run").Params(
		jen.Id("port").Int(),
		jen.Id("token").Op("*").String(),
	).Error().Block(serveBody...)

	// ---

	f.Type().Id("Server").Struct(
		jen.Id("aperture").Qual("github.com/goaperture/goaperture/lib/aperture", "Aperture"),
	)

	f.Func().Id("NewServer").Params().Op("*").Id("Server").Block(
		jen.Return(jen.Op("&").Id("Server").Values(jen.Dict{
			jen.Id("aperture"): jen.Op("*").Qual("github.com/goaperture/goaperture/lib/aperture", "NewServer").Call(),
		})),
	)

	f.Func().Params(jen.Id("server").Op("*").Id("Server")).
		Id("Middleware").
		Params(jen.Id("middleware").Func().Params(jen.Id("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler"))).
		Block(
			jen.Id("server").Dot("aperture").Dot("Middleware").Op("=").Id("middleware"),
		)

	f.Func().Id("Serve").Params(
		jen.Id("port").Int(),
		jen.Id("token").Op("*").String(),
	).Block(
		jen.If(
			jen.Err().Op(":=").Id("NewServer").Call().Dot("Run").Call(jen.Id("port"), jen.Id("token")),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Id("panic").Call(jen.Err()),
		),
	)

	// ----

	if _, err := os.Stat(_package); os.IsNotExist(err) {
		os.MkdirAll(_package, 0755)
	}
	if err := f.Save(_package + "/aperture.go"); err != nil {
		panic(fmt.Sprintf("Error saving file: %v", err))
	}

	fmt.Println(">", _package+"/config/config.go")
	_, err := os.Stat(_package + "/config/config.go")
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println(">>> generate", _package)
		Config(_package)
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

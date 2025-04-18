package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/goaperture/goaperture/lib/aperture"
)

func Run() {
	app, routes, output := commands()
	if app == nil || routes == nil || output == nil {
		panic("Нужны аргументы")
	}
	aperture.Generate(*app, *routes, *output)
	fmt.Println("OK")
}

func commands() (app *string, routes *string, output *string) {
	generate := flag.NewFlagSet("generate", flag.ExitOnError)

	app = generate.String("module", "app", "Название модуля (go.mod module)")
	routes = generate.String("routes", "api/routes", "Папка с маршрутами")
	output = generate.String("output", "api", "Папка для сохранения")

	generate.Parse(os.Args[1:])

	return
}

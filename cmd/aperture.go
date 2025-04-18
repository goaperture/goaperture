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
	cmd := flag.NewFlagSet("generate", flag.ExitOnError)

	app = cmd.String("module", "app", "Название модуля (go.mod module)")
	routes = cmd.String("routes", "api/routes", "Папка с маршрутами")
	output = cmd.String("output", "api", "Папка для сохранения")

	cmd.Parse(os.Args[1:])

	return
}

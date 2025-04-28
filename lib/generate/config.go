package generate

import (
	"fmt"
	"os"

	"github.com/dave/jennifer/jen"
)

func Config(_package string) {
	f := jen.NewFile("config")

	f.ImportAlias("github.com/goaperture/goaperture/lib/aperture", "aperture")

	f.Type().Id("Payload").Struct(
		jen.Id("Id").String().Tag(map[string]string{"json": "id"}),
		jen.Id("Name").String().Tag(map[string]string{"json": "name"}),
		jen.Id("Email").String().Tag(map[string]string{"json": "email"}),
		jen.Id("Avatar").String().Tag(map[string]string{"json": "avatar"}),
		jen.Id("Permissions").Qual("github.com/goaperture/goaperture/lib/aperture", "Permissions").Tag(map[string]string{"json": "permissions"}),
	)

	f.Type().Id("Client").Op("=").Qual("github.com/goaperture/goaperture/lib/aperture", "Client").Types(jen.Id("Payload"))

	if err := os.MkdirAll(_package+"/config", 0755); err != nil {
		panic(err)
	}

	if err := f.Save(_package + "/config/config.go"); err != nil {
		panic(fmt.Sprintf("Error saving file: %v", err))
	}
}

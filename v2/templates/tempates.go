package templates

import (
	"go/format"
	"log"
	"strings"
	"text/template"
)

type FileData struct {
	Pkg         string
	Name        string
	Description string
	Sequre      bool
}

func GetRouteCode(pkg, name, description string, sequre bool) string {
	var code = `package {{.Pkg}}
	import (
		"context"
		"github.com/goaperture/goaperture/v2/api/aperture"
	)
	type {{.Name}}Input struct { }
	type {{.Name}}Output interface {any}
	var {{.Name}} = aperture.Route[{{.Name}}Input, {{.Name}}Output]{
		Description:   "{{.Description}}",
		Handler:       {{.Name}}Handler,
		{{if .Sequre}}PrivateAccess: true,{{end}}
		Prepare: func(cl *aperture.CL[{{.Name}}Input, {{.Name}}Output]) {
			cl.Execute({{.Name}}Input{})
		},
	}
	func {{.Name}}Handler(ctx context.Context, input {{.Name}}Input) {{.Name}}Output {
		return "Hello from {{.Name}}"
	}
	`
	var result strings.Builder
	tmp, err := template.New("route").Parse(code)
	if err != nil {
		log.Fatal(err)
	}

	tmp.Execute(&result, FileData{pkg, name, description, sequre})

	formatted, err := format.Source([]byte(result.String()))
	if err != nil {
		return result.String()
	}

	return string(formatted)
}

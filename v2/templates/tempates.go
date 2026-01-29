package templates

import (
	"go/format"
	"strings"
	"text/template"
)

type FileData struct {
	Pkg         string
	Name        string
	Description string
}

func GetRouteCode(pkg, name, description string) string {
	var code = `package {{.Pkg}}
	import (
		"context"
		"github.com/goaperture/goaperture/v2/aperture"
	)
	type {{.Name}}Input struct {}
	type {{.Name}}Output interface {any}
	var {{.Name}} = aperture.Route[{{.Name}}Input, {{.Name}}Output]{
		Description:   "{{.Description}}",
		Handler:       {{.Name}}Handler,
		Prepare: func(cl *aperture.CL[{{.Name}}Input, {{.Name}}Output]) {
			cl.Execute({{.Name}}Input{})
		},
	}
	func {{.Name}}Handler(ctx context.Context, input {{.Name}}Input) {{.Name}}Output {
		return "Hello from {{.Name}}"
	}
	`
	var result strings.Builder
	template.Must(template.New("route").Parse(code)).Execute(&result, FileData{pkg, name, description})

	formatted, err := format.Source([]byte(result.String()))
	if err != nil {
		return result.String()
	}

	return string(formatted)
}

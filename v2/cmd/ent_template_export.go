package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var entTemplateExport = &cobra.Command{
	Use:   "ent-templates-export",
	Short: "Перейдите в нужную папку в выполните команду чтобы создать ent template",
	Run: func(cmd *cobra.Command, args []string) {

		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		imprt := filepath.Base(dir)
		fmt.Println(dir, imprt)

		paginationCode := `
{{ define "import/additional" }}
    {{- if hasField $ "Config" }}
        stdsql "database/sql"
    {{- end }}
	"github.com/goaperture/goaperture/v2/api/client"
	"github.com/goaperture/goaperture/v2/responce"
{{ end }}



{{ define "query/additional/methods" }}

// Paginate добавляет LIMIT и OFFSET в запрос.
func (q *{{ $.Name }}Query) Paginate(ctx context.Context, page, size int) *{{ $.Name }}Query {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size

	total, _ := q.Clone().Offset(0).Limit(-1).Count(ctx)
	client.SetPagination(ctx, responce.Pagination{
		Page:  page,
		Size:  size,
		Total: total,
	})

	return q.Limit(size).Offset(offset)
}
{{ end }}

`

		if err := os.WriteFile(filepath.Join(dir, "pagination.tmpl"), []byte(paginationCode), 0644); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Отлично - теперь можно добавить в файл `generate.go`")
		fmt.Printf("--template ./%s\n\n\n", imprt)
	},
}

func init() {
	rootCmd.AddCommand(entTemplateExport)

}

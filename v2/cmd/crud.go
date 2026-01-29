package cmd

import (
	"github.com/goaperture/goaperture/v2/templates"
	"github.com/spf13/cobra"
)

type Method struct {
	Name        string
	Description string
}

var crudCmd = &cobra.Command{
	Use:   "crud",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		path, name, description := args[0], args[1], args[2]

		methods := []Method{
			{
				Name:        "Create" + name,
				Description: "Создать " + description,
			},
			{
				Name:        "Get" + name,
				Description: "Получить " + description,
			},
			{
				Name:        "Update" + name,
				Description: "Обновить " + description,
			},
			{
				Name:        "Delete" + name,
				Description: "Удалить " + description,
			},
		}

		for _, method := range methods {
			templates.CreateRoute(path, method.Name, method.Description, "")
		}
	},
}

func init() {
	rootCmd.AddCommand(crudCmd)
}

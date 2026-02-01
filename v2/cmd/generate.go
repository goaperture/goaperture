package cmd

import (
	"github.com/goaperture/goaperture/v2/generate"
	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) {
	app, _ := cmd.Flags().GetString("app")
	path, _ := cmd.Flags().GetString("path")

	generate.GenerateRoutes(app, path)
	generate.GenerateWebsockets(app, path)
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: run,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("app", "a", "app", "Название модуля (go.mod module)")
	generateCmd.Flags().StringP("path", "p", "api", "Папка с маршрутами (Обычно api) - там должна быть папка routes или ws")
}

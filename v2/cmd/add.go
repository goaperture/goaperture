package cmd

import (
	"log"
	"os"

	"github.com/goaperture/goaperture/v2/templates"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add-route",
	Short: "Перейдите в нужную папку в выполните команду чтобы создать роут",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if len(args) > 0 {
			name = args[0]
		}
		description, _ := cmd.Flags().GetString("description")
		if len(args) > 1 {
			description = args[1]
		}

		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		sequre, _ := cmd.Flags().GetBool("sequre")

		templates.CreateRoute(dir, name, description, sequre)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringP("name", "n", "Hello", "Название роута (обязательный первый аргумент)")
	addCmd.Flags().StringP("description", "d", "hello world", "Описание роута (обязательный второй аргумент)")
	addCmd.Flags().BoolP("sequre", "s", false, "Нужен доступ")

	addCmd.Flags().Bool("test", false, "Вызывать метод при тестировании")

}

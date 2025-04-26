package cmd

import (
	"github.com/goaperture/goaperture/lib/aperture"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Генерирует файл aperture.go на основе указанных директорий и файлов.",
	Long: `
Команда generate анализирует переданные директории и файлы, после чего автоматически 
создаёт или обновляет файл aperture.go. Этот файл объединяет структуру и 
содержимое указанных источников, упрощая настройку проекта, инициализацию 
зависимостей или генерацию шаблонного кода. Поддерживает гибкую конфигурацию 
через аргументы командной строки.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		app, _ := cmd.Flags().GetString("app")
		routes, _ := cmd.Flags().GetString("routes")
		output, _ := cmd.Flags().GetString("output")

		aperture.Generate(app, routes, output)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringP("app", "a", "app", "Название модуля (go.mod module)")
	generateCmd.Flags().StringP("routes", "r", "api/routes", "Папка с маршрутами")
	generateCmd.Flags().StringP("output", "o", "api", "Папка для сохранения файла aperture.go")
}

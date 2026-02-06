package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed _init/*
var templateFS embed.FS

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Инициализация нового проекта",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetDir := "."
		if len(args) > 0 {
			targetDir = args[0]
		}

		fmt.Printf("🚀 Инициализация проекта в: %s\n", targetDir)

		// Рекурсивно проходим по встроенным файлам
		err := fs.WalkDir(templateFS, "templates", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// Вычисляем путь назначения (убираем префикс "templates")
			relPath, _ := filepath.Rel("templates", path)
			if relPath == "." {
				return nil
			}

			targetPath := filepath.Join(targetDir, relPath)

			if d.IsDir() {
				// Создаем папку
				return os.MkdirAll(targetPath, 0755)
			} else {
				// Читаем файл из бинарника
				data, err := templateFS.ReadFile(path)
				if err != nil {
					return err
				}
				// Записываем файл на диск
				fmt.Printf("  Создаю файл: %s\n", targetPath)
				return os.WriteFile(targetPath, data, 0644)
			}
		})

		if err != nil {
			fmt.Printf("❌ Ошибка: %v\n", err)
			return
		}

		fmt.Println("✅ Готово!")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

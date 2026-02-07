package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed templates/*
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

		fmt.Printf("Инициализация проекта в: %s\n", targetDir)

		err := fs.WalkDir(templateFS, "templates", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			relPath, _ := filepath.Rel("templates", path)
			if relPath == "." {
				return nil
			}

			targetPath := filepath.Join(targetDir, relPath)

			if d.IsDir() {
				return os.MkdirAll(targetPath, 0755)
			} else {
				data, err := templateFS.ReadFile(path)
				if err != nil {
					return err
				}
				return os.WriteFile(targetPath, data, 0644)
			}
		})

		if err != nil {
			fmt.Printf("❌ Ошибка: %v\n", err)
			return
		}

		if targetDir != "." {
			fmt.Printf("cd %s\n", targetDir)
		}
		fmt.Println("go mod init app")
		fmt.Println("go mod tidy")
		fmt.Println("a2 generate")

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/goaperture/goaperture/v2/templates"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		path, name, description := args[0], args[1], args[2]
		createRoute(path, name, description)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func createRoute(path, name, description string) {
	absPath, _ := filepath.Abs(path)
	pkg := filepath.Base(absPath)

	fileContent := templates.GetRouteCode(pkg, name, description)

	file := filepath.Join(absPath, name+".go")

	os.WriteFile(file, []byte(fileContent), 0777)
}

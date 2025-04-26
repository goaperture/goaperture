/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Инициализировать старовый проект",
	Long:  `Инициализирует стартовый простой проект`,
	Run: func(cmd *cobra.Command, args []string) {
		folder, _ := cmd.Flags().GetString("folder")

		clone := exec.Command("sh", "-c", "git clone git@github.com:goaperture/templates.git "+folder)
		output, _ := clone.CombinedOutput()
		fmt.Println(string(output))

		os.RemoveAll(folder + "/.git")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("folder", "f", ".", "Папка в которой создать проект")
}

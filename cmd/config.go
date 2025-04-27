package cmd

import (
	"github.com/goaperture/goaperture/lib/generate"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "generate config",
	Long:  `generate config`,
	Run: func(cmd *cobra.Command, args []string) {
		ouput, _ := cmd.Flags().GetString("output")
		generate.Config(ouput)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringP("output", "o", "api", "Папка где находится файла aperture.go")
}

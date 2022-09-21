package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"nacos_check/app"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "查看本地配置文件路径",
	Run: func(cmd *cobra.Command, args []string) {
		configfile := app.GetConfigFilePath()
		fmt.Println("本地配置文件路径:", configfile)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "查看版本",
	Run: func(cmd *cobra.Command, args []string) {
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version: 0.7.7")
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "查看本地配置文件路径",
	Run: func(cmd *cobra.Command, args []string) {
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		configfile := getConfigFilePath()
		fmt.Println("本地配置文件路径:", configfile)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
}

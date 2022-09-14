package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"nacos_check/app"
	"path/filepath"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "查看本地配置文件路径",
	Run: func(cmd *cobra.Command, args []string) {
		homedir, err := app.HomeDir()
		if err != nil {
			fmt.Println("获取系统家目录获取异常")
			homedir = "."
		}
		configfile := filepath.Join(homedir, ".nacos_url")
		fmt.Println("本地配置文件路径:", configfile)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

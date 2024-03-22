package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"nacos-check/internal/nacos"
)

var (
	backupmode bool
	format     string
	timestamp  bool
)

var configsyncCmd = &cobra.Command{
	Use:   "config-sync",
	Short: "实时同步nacos配置到本地",
	Run: func(cmd *cobra.Command, args []string) {
		sync := nacos.NewSync(backupmode, timestamp, format)
		sync.RunSync()
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if !backupmode && timestamp {
			log.Fatalln("使用timestamp必须同时启用--backup")
		}
	},
}

func init() {
	configsyncCmd.Flags().BoolVarP(&backupmode, "backup", "b", false, "启用备份模式,配置文件同步本地后文件自带时间")
	configsyncCmd.Flags().StringVarP(&format, "format", "f", "2006-01-0215:04:05", "备份模式备份时间戳格式化风格，只可修改排版，不可以修改内容")
	configsyncCmd.Flags().BoolVarP(&timestamp, "timestamp", "t", false, "备份模式下使用时间戳后缀")
	rootCmd.AddCommand(configsyncCmd)
}

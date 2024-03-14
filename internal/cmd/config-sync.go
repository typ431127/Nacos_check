package cmd

import (
	"github.com/spf13/cobra"
	"nacos-check/internal/nacos"
)

var configsyncCmd = &cobra.Command{
	Use:   "config-sync",
	Short: "实时同步nacos配置到本地",
	Run: func(cmd *cobra.Command, args []string) {
		sync := nacos.NewSync()
		sync.RunSync()
	},
}

func init() {
	rootCmd.AddCommand(configsyncCmd)
}

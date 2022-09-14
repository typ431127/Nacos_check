package cmd

import (
	"github.com/spf13/cobra"
	"nacos_check/nacos"
	"nacos_check/web"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "开启web api Prometheus http_sd_configs",
	Run: func(cmd *cobra.Command, args []string) {
		web.Runwebserver()
	},
}

func init() {
	webCmd.Flags().StringVarP(&nacos.Port, "port", "p", ":8099", "web 端口")
	rootCmd.AddCommand(webCmd)
}

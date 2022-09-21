package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"nacos_check/app"
	"nacos_check/nacos"
	"os"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "nacos_check",
	Short: "Nacos",
	Long:  `Nacos`,
	Run: func(cmd *cobra.Command, args []string) {
		app.NacosRun()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		app.NacosInit()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&nacos.Nacosurl, "url", "u", "http://dev-k8s-nacos:8848", "Nacos地址")
	rootCmd.Flags().StringVarP(&nacos.Writefile, "write", "o", "", "导出json文件, prometheus 自动发现文件路径")
	rootCmd.Flags().StringVarP(&nacos.Ipfile, "ipfile", "i", "salt_ip.json", "ip解析文件")
	rootCmd.Flags().StringVarP(&nacos.Findstr, "find", "f", "", "查找服务")
	rootCmd.Flags().BoolVarP(&nacos.ExportJson, "json", "", false, "输出json")
	rootCmd.Flags().BoolVarP(&nacos.Watch, "watch", "w", false, "监控服务")
	rootCmd.Flags().DurationVarP(&nacos.Second, "second", "s", 5*time.Second, "监控服务间隔刷新时间")
	rootCmd.PersistentFlags().StringToStringVarP(&nacos.AddLable, "lable", "l", map[string]string{}, "添加标签 -l env=dev,pro=java")
}

func Execute() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

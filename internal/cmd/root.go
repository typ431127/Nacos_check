package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"nacos-check/internal/config"
	"nacos-check/internal/core"
	"os"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "nacos-check",
	Short: "Nacos",
	Long:  `Nacos`,
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case config.EXPORTJSON:
			jsondata, err := config.Nacos.GetJson("byte", false)
			if err != nil {
				fmt.Println("获取json发生错误")
				os.Exit(2)
			}
			data := make([]byte, 0)
			var check bool
			if data, check = jsondata.([]byte); !check {
				fmt.Println("转换失败")
				os.Exit(2)
			}
			fmt.Println(string(data))
			os.Exit(0)
		case config.WRITEFILE != "":
			config.Nacos.WriteFile()
		default:
			if config.STDOUT == "table" {
				fmt.Println("Nacos:", config.NACOSURL)
			} else {
				fmt.Printf("[Nacos](%s)\n", config.NACOSURL)
			}
			if config.WATCH {
				fmt.Printf("监控模式 刷新时间:%s/次\n", config.SECOND)
				for {
					config.Nacos.GetNacosInstance()
					config.Nacos.Render()
					time.Sleep(config.SECOND)
				}
			}
			config.Nacos.Render()
		}
		config.Nacos.Client.CloseIdleConnections()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		core.PreFunc()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&config.NACOSURL, "url", "u", "http://dev-k8s-nacos:8848", "Nacos地址")
	rootCmd.PersistentFlags().StringVarP(&config.USERNAME, "username", "", "nacos", "账户")
	rootCmd.PersistentFlags().StringVarP(&config.PASSWORD, "password", "", "", "密码")
	rootCmd.PersistentFlags().StringVarP(&config.CONTEXTPATH, "contextpath", "", "/nacos", "server.servlet.contextPath")
	rootCmd.Flags().StringVarP(&config.NAMESPACE, "namespace", "", "", "指定命名空间ID 多个: id1,id2,id3")
	rootCmd.PersistentFlags().StringVarP(&config.GROUP, "group", "", "DEFAULT_GROUP", "指定分组 多个分组 group1,group2")
	rootCmd.Flags().StringVarP(&config.WRITEFILE, "write", "o", "", "导出json文件, prometheus 自动发现文件路径")
	rootCmd.Flags().StringVarP(&config.IPFILE, "ipfile", "i", "salt_ip.json", "ip解析文件")
	rootCmd.Flags().StringVarP(&config.FIND, "find", "f", "", "查找服务")
	rootCmd.Flags().BoolVarP(&config.CLUSTER, "cluster", "", false, "全集群查找")
	rootCmd.Flags().BoolVarP(&config.EXPORTJSON, "json", "", false, "输出json")
	rootCmd.Flags().BoolVarP(&config.WATCH, "watch", "w", false, "监控服务")
	rootCmd.Flags().DurationVarP(&config.SECOND, "second", "s", 5*time.Second, "监控服务间隔刷新时间")
	rootCmd.Flags().StringVarP(&config.STDOUT, "stdout", "", "table", "输出类型 table / markdown")
	rootCmd.PersistentFlags().StringToStringVarP(&config.ADDLABEL, "lable", "l", map[string]string{}, "添加标签 -l env=dev,pro=java")
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

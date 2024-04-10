package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"nacos-check/internal/nacos"
	"nacos-check/pkg/fmtd"
	"os"
	"time"
)

var Nacos *nacos.Nacos

var rootCmd = &cobra.Command{
	Use:   "nacos-check",
	Short: "Nacos",
	Long:  `Nacos`,
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case nacos.EXPORTJSON:
			jsondata, err := Nacos.GetJson("byte", false)
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
		case nacos.WRITEFILE != "":
			Nacos.WriteFile()
		default:
			if nacos.STDOUT == "table" {
				fmt.Println("Nacos:", nacos.NACOSURL)
			} else {
				fmt.Printf("[Nacos](%s)\n", nacos.NACOSURL)
			}
			if nacos.WATCH {
				fmt.Printf("监控模式 刷新时间:%s/次\n", nacos.SECOND)
				for {
					Nacos.GetNacosInstance()
					Nacos.Render()
					time.Sleep(nacos.SECOND)
				}
			}
			Nacos.Render()
		}
		Nacos.Client.CloseIdleConnections()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		preFunc()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&nacos.NACOSURL, "url", "u", "http://dev-k8s-nacos:8848", "Nacos地址")
	rootCmd.PersistentFlags().StringVarP(&nacos.USERNAME, "username", "", "nacos", "账户")
	rootCmd.PersistentFlags().StringVarP(&nacos.PASSWORD, "password", "", "", "密码")
	rootCmd.PersistentFlags().StringVarP(&nacos.CONTEXTPATH, "contextpath", "", "/nacos", "server.servlet.contextPath")
	rootCmd.Flags().StringVarP(&nacos.NAMESPACE, "namespace", "", "", "指定命名空间ID 多个: id1,id2,id3")
	rootCmd.PersistentFlags().StringVarP(&nacos.GROUP, "group", "", "DEFAULT_GROUP", "指定分组 多个分组 group1,group2")
	rootCmd.PersistentFlags().StringVarP(&nacos.FileConfig.ConfigFile, "config", "", "", "指定配置文件路径")
	rootCmd.Flags().StringVarP(&nacos.WRITEFILE, "write", "o", "", "导出json文件, prometheus 自动发现文件路径")
	rootCmd.Flags().StringVarP(&nacos.IPFILE, "ipfile", "i", "salt_ip.json", "ip解析文件")
	rootCmd.Flags().StringVarP(&nacos.FIND, "find", "f", "", "查找服务")
	rootCmd.Flags().BoolVarP(&nacos.CLUSTER, "cluster", "", false, "全集群查找")
	rootCmd.Flags().BoolVarP(&nacos.EXPORTJSON, "json", "", false, "输出json")
	rootCmd.Flags().BoolVarP(&nacos.WATCH, "watch", "w", false, "监控服务")
	rootCmd.Flags().DurationVarP(&nacos.SECOND, "second", "s", 5*time.Second, "监控服务间隔刷新时间")
	rootCmd.Flags().StringVarP(&nacos.STDOUT, "stdout", "", "table", "输出类型 table / markdown")
	rootCmd.PersistentFlags().StringToStringVarP(&nacos.ADDLABEL, "lable", "l", map[string]string{}, "添加标签 -l env=dev,pro=java")
}

func Execute() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if err := rootCmd.Execute(); err != nil {
		fmtd.Fatalln(err)
	}
}

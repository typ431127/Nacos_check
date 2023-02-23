package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"nacos-check/internal/config"
	"net/url"
	"os"
	"time"
)

var Nacos config.Nacos
var rootCmd = &cobra.Command{
	Use:   "nacos-check",
	Short: "Nacos",
	Long:  `Nacos`,
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case config.EXPORTJSON:
			jsondata, err := Nacos.GetJson("byte")
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
			Nacos.WriteFile()
		default:
			fmt.Println("Nacos:", config.NACOSURL)
			if config.WATCH {
				fmt.Printf("监控模式 刷新时间:%s/次\n", config.SECOND)
				for {
					Nacos.GetNacosInstance()
					Nacos.TableRender()
					time.Sleep(config.SECOND)
				}
			}
			Nacos.TableRender()
		}
		Nacos.Client.CloseIdleConnections()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		NacosFilePathLoad()
		u, err := url.Parse(config.NACOSURL)
		if err != nil {
			fmt.Println("url解析错误!")
			os.Exit(config.EXITCODE)
		}
		Nacos.DefaultUlr = config.NACOSURL
		Nacos.Host = u.Host
		Nacos.Scheme = u.Scheme
		Nacos.Port = u.Port()
		Nacos.GetNameSpace()
		Nacos.GetNacosInstance()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&config.NACOSURL, "url", "u", "http://dev-k8s-nacos:8848", "Nacos地址")
	rootCmd.Flags().StringVarP(&config.WRITEFILE, "write", "o", "", "导出json文件, prometheus 自动发现文件路径")
	rootCmd.Flags().StringVarP(&config.IPFILE, "ipfile", "i", "salt_ip.json", "ip解析文件")
	rootCmd.Flags().StringVarP(&config.FIND, "find", "f", "", "查找服务")
	rootCmd.Flags().BoolVarP(&config.EXPORTJSON, "json", "", false, "输出json")
	rootCmd.Flags().BoolVarP(&config.WATCH, "watch", "w", false, "监控服务")
	rootCmd.Flags().DurationVarP(&config.SECOND, "second", "s", 5*time.Second, "监控服务间隔刷新时间")
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

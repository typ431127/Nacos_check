package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"nacos-check/internal/config"
	"os"
	"strconv"
)

var v2upgrade bool
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "集群状态",
	Run: func(cmd *cobra.Command, args []string) {
		config.CLUSTER = true
		if v2upgrade {
			config.V2UPGRADE = true
		}
		Nacos.GetNacosInstance()
		tablecluser := tablewriter.NewWriter(os.Stdout)
		tablecluser.SetHeader([]string{"节点", "端口", "状态", "版本", "刷新时间", "健康实例", "异常实例"})
		for _, key := range Nacos.Clusterdata {
			tablecluser.Append([]string{
				key.Ip,
				key.Port,
				key.State,
				key.Version,
				key.LastRefreshTime,
				strconv.Itoa(len(key.HealthInstance)),
				strconv.Itoa(len(key.UnHealthInstance)),
			})
		}
		leader := gjson.Get(Nacos.Cluster, "servers.0.extendInfo.raftMetaData.metaDataMap.naming_instance_metadata.leader")
		fmt.Printf("Nacos集群状态: (数量:%d)\n集群Master: %s\n", tablecluser.NumLines(), leader)
		tablecluser.Render()
		if v2upgrade {
			tablecluser := tablewriter.NewWriter(os.Stdout)
			tablecluser.SetHeader([]string{"节点", "双写", "v2服务", "v2实例", "v1服务", "v1实例", "upgraded", "全部V2"})
			for _, key := range Nacos.Clusterdata {
				tablecluser.Append([]string{
					key.Ip,
					strconv.FormatBool(key.V2Upgrade.IsDoubleWriteEnabled),
					strconv.FormatInt(key.V2Upgrade.ServiceCountV2, 10),
					strconv.FormatInt(key.V2Upgrade.InstanceCountV2, 10),
					strconv.FormatInt(key.V2Upgrade.ServiceCountV1, 10),
					strconv.FormatInt(key.V2Upgrade.InstanceCountV1, 10),
					strconv.FormatBool(key.V2Upgrade.Upgraded),
					strconv.FormatBool(key.V2Upgrade.IsAll20XVersion),
				})
			}
			fmt.Printf("v2版本升级接口详情\n")
			tablecluser.Render()
		}
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterCmd.Flags().BoolVarP(&v2upgrade, "v2upgrade", "v", false, "v2升级接口状态")
}

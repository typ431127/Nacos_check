package cmd

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"nacos-check/internal/nacos"
	"time"
)

var Refreshtime time.Duration
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "开启web api Prometheus http_sd_configs",
	Run: func(cmd *cobra.Command, args []string) {
		Webserver()
	},
}

func init() {
	webCmd.Flags().BoolVarP(&nacos.CLUSTER, "cluster", "", false, "启用全集群查找")
	webCmd.Flags().StringVarP(&nacos.NAMESPACE, "namespace", "", "", "指定命名空间ID 多个: id1,id2,id3")
	webCmd.Flags().DurationVarP(&Refreshtime, "refresh", "", time.Second*3600, "Token刷新时间,默认3600")
	webCmd.Flags().StringVarP(&nacos.WEBPORT, "port", "p", ":8099", "web 端口")
	rootCmd.AddCommand(webCmd)
}

func response(c *gin.Context) {
	if c.Request.RequestURI == "/health" {
		c.JSON(200, gin.H{"status": true})
		return
	}
	if c.Request.RequestURI == "/favicon.ico" {
		c.JSON(404, "404")
		return
	}
	result, err := Nacos.GetJson("json", true)
	if err != nil {
		c.JSON(500, []string{})
		return
	}
	c.JSON(200, result)
}

func Webserver() {
	fmt.Println("Nacos:", nacos.NACOSURL)
	gin.SetMode(gin.DebugMode)
	RefreshToken()
	r := gin.Default()
	v1 := r.Group("/")
	v1.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "OPTIONS", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	{
		v1.GET("/*route", response)
	}
	err := r.Run(nacos.WEBPORT)
	if err != nil {
		fmt.Println(err)
	}
}

func RefreshToken() {
	if len(nacos.USERNAME) != 0 && len(nacos.PASSWORD) != 0 {
		go func() {
			for {
				Nacos.WithAuth()
				time.Sleep(Refreshtime)
			}
		}()
	}
}

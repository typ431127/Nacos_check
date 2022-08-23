package web

import (
	"github.com/gin-gonic/gin"
	"nacos_check/nacos"
)

func response(c *gin.Context) {
	var d *nacos.Nacos
	d = &nacos.Na
	result := d.GetJson("json")
	c.JSON(200, result)
}
func Runwebserver() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	v1 := r.Group("/")
	{
		v1.GET("/*action", response)
	}
	r.Run(nacos.Port)
}

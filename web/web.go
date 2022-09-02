package web

import (
	"github.com/gin-gonic/gin"
	"nacos_check/nacos"
)

func response(c *gin.Context) {
	var d *nacos.Nacos
	d = &nacos.Na
	if c.Request.RequestURI == "/health" {
		c.JSON(200, gin.H{"status": true})
		return
	}
	if c.Request.RequestURI == "/favicon.ico" {
		c.JSON(404, "404")
		return
	}
	result, err := d.GetJson("json")
	if err != nil {
		c.JSON(500, []string{})
		return
	}
	c.JSON(200, result)
}

func Runwebserver() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	v1 := r.Group("/")
	{
		v1.GET("/*route", response)
	}
	r.Run(nacos.Port)
}

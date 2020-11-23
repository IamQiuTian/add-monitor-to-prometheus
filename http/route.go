package http

import (
	"automonitor/g"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

var r = gin.Default()

func Run()  {
	r.LoadHTMLGlob("templates/*")

	r.GET("/", Index)
	r.GET("/node_exporter", Node_exporter)
	r.GET("/port_exporter", Port_exporter)
	r.GET("/get_exporter", Get_exporter)
	r.GET("/ssl_exporter", SSL_exporter)
	r.GET("/delete", DeleteHtml)
	r.POST("/addnode_exporter", Addnode_exporter)
	r.POST("/addport_exporter",Addport_exporter)
	r.POST("/addget_exporter",Addget_exporter)
	r.POST("/addssl_exporter",Addssl_exporter)
	r.POST("/deleteservice", DeleteService)


	listen := fmt.Sprintf("%s:%s", g.CF.Config.HttpListen.Host, g.CF.Config.HttpListen.Port)

	log.Fatal(r.Run(listen))
}

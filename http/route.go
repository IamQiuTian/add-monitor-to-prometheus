package http

import (
	"add-monitor-to-prometheus/g"
	"add-monitor-to-prometheus/utils"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

var r = gin.Default()

func Run() {
	r.LoadHTMLGlob("templates/*")
	r.GET("/", Login)
	r.POST("/auth", Auth)
	authorized := r.Group("/")
	authorized.Use(utils.JWTAuthMiddleware())
	{
		authorized.GET("/index", Index)
		authorized.GET("/node_exporter", Node_exporter)
		authorized.GET("/port_exporter", Port_exporter)
		authorized.GET("/mysqld_exporter", Mysqld_exporter)
		authorized.GET("/mongodb_exporter", Mongodb_exporter)
		authorized.GET("/ssl_exporter", SSL_exporter)
		authorized.GET("/domain_exporter", Domain_exporter)
		authorized.GET("/delete", DeleteHtml)
		authorized.POST("/addnode_exporter", Addnode_exporter)
		authorized.POST("/addport_exporter", Addport_exporter)
		authorized.POST("/addmysqld_exporter", Addmysqld_exporter)
		authorized.POST("/addmongodb_exporter", Addmongodb_exporter)
		authorized.POST("/addssl_exporter", Addssl_exporter)
		authorized.POST("/adddomain_exporter", Adddomain_exporter)
		authorized.POST("/deleteservice", DeleteService)
		authorized.POST("/deleteauth", DeleteAuth)
	}

	listen := fmt.Sprintf("%s:%s", g.CF.Config.HttpListen.Host, g.CF.Config.HttpListen.Port)

	r.Use(Cors())
	log.Fatal(r.Run(listen))
}

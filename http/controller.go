package http

import (
	"add-monitor-to-prometheus/consul"
	"add-monitor-to-prometheus/g"
	"add-monitor-to-prometheus/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func init() {
	log.SetPrefix("[DEBUG]")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", nil)
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", nil)
}

func DeleteHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "delete.tmpl", nil)
}

func Node_exporter(c *gin.Context) {
	c.HTML(http.StatusOK, "node_exporter.tmpl", nil)
}

func Port_exporter(c *gin.Context) {
	c.HTML(http.StatusOK, "port_exporter.tmpl", nil)
}

func Mysqld_exporter(c *gin.Context) {
	c.HTML(http.StatusOK, "mysqld_exporter.tmpl", nil)
}

func Mongodb_exporter(c *gin.Context) {
	c.HTML(http.StatusOK, "mongodb_exporter.tmpl", nil)
}

func SSL_exporter(c *gin.Context) {
	c.HTML(http.StatusOK, "ssl_exporter.tmpl", nil)
}

func Domain_exporter(c *gin.Context) {
	c.HTML(http.StatusOK, "domain_exporter.tmpl", nil)
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		}
		c.Next()
	}
}

func Auth(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	var userjson *g.Userjson
	err = json.Unmarshal(data, &userjson)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	if userjson.Username == g.CF.Config.AuthorizaTion.Username && userjson.Password == g.CF.Config.AuthorizaTion.Password {
		authStr, err := c.Cookie("Authorization")
		if err == nil {
			mc, err := utils.ParseToken(authStr)
			if err != nil {
				c.JSON(http.StatusNotExtended, gin.H{
					"code": 2005,
					"msg":  "无效的Token",
				})
				c.Abort()
				return
			}
			c.Set("username", mc.Username)
		}

		tokenString, _ := utils.GenToken(userjson.Username)
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "Authorization",
			Value:    tokenString,
			Path:     "/",
			Domain:   "",
			MaxAge:   3600,
			HttpOnly: false,
		})
		c.JSON(http.StatusOK, gin.H{
			"code": 2000,
			"msg":  "success",
			"data": gin.H{"Authorization": tokenString},
		})
		return
	}

	c.JSON(http.StatusNotExtended, gin.H{
		"code": 2002,
		"msg":  "鉴权失败",
	})
	return
}

func DeleteAuth(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	var userjson *g.Userjson
	err = json.Unmarshal(data, &userjson)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	tokenString, _ := utils.GenToken(userjson.Username)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Path:     "/",
		Domain:   "",
		MaxAge:   -1,
		HttpOnly: false,
	})
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "success",
		"data": gin.H{"Authorization": tokenString},
	})
	return
	c.JSON(http.StatusNotExtended, gin.H{
		"code": 2002,
		"msg":  "删除失败",
	})
	return
}

func Addnode_exporter(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	var consulDataJson *g.ConsulDataJson
	err = json.Unmarshal(data, &consulDataJson)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}
	_, _, err = consul.ReadService(consulDataJson.Hostname)
	if err == nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": "service already exists",
		})
		log.Println("service already exists")
		return
	}

	checkUrl := consulDataJson.Address
	checkPort := consulDataJson.Port
	var consulData = g.ConsulData{
		ID:      consulDataJson.Hostname,
		Name:    consulDataJson.Group,
		Address: consulDataJson.Address,
		Port:    consulDataJson.Port,
		Tags:    g.CF.Config.Tags.Node_exporter,
	}

	err = consul.WriteService(consulData, checkUrl, checkPort)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	address, port, err := consul.ReadService(consulData.ID)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": fmt.Sprintf("%s: %n", address, port),
	})
}

func Addport_exporter(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	var consulDataJson *g.ConsulDataJson
	err = json.Unmarshal(data, &consulDataJson)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}
	ID := fmt.Sprintf("%s-port-check-%v", consulDataJson.Hostname, consulDataJson.Port)
	_, _, err = consul.ReadService(ID)
	if err == nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": "service already exists",
		})
		log.Println("service already exists")
		return
	}

	var consulData = g.ConsulData{
		ID:      ID,
		Name:    consulDataJson.Group,
		Address: consulDataJson.Address,
		Port:    consulDataJson.Port,
		Tags:    g.CF.Config.Tags.Port_exporter,
	}

	checkUrl := ""
	checkPort := 0
	err = consul.WriteService(consulData, checkUrl, checkPort)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	address, port, err := consul.ReadService(consulData.ID)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": fmt.Sprintf("%s: %n", address, port),
	})
}

func Addmysqld_exporter(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	var consulDataJson *g.ConsulDataJson
	err = json.Unmarshal(data, &consulDataJson)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	ID := strings.Replace(consulDataJson.Hostname, ".", "_", -1) + "-mysqld-check"

	_, _, err = consul.ReadService(ID)
	if err == nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": "service already exists",
		})
		log.Println("service already exists")
		return
	}

	var consulData = g.ConsulData{
		ID:      ID,
		Name:    consulDataJson.Group,
		Address: consulDataJson.Address,
		Port:    consulDataJson.Port,
		Tags:    g.CF.Config.Tags.Mysqld_exporter,
	}

	checkUrl := ""
	checkPort := 0
	err = consul.WriteService(consulData, checkUrl, checkPort)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	address, port, err := consul.ReadService(consulData.ID)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": fmt.Sprintf("%s: %n", address, port),
	})
}

func Addmongodb_exporter(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	var consulDataJson *g.ConsulDataJson
	err = json.Unmarshal(data, &consulDataJson)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	ID := strings.Replace(consulDataJson.Hostname, ".", "_", -1) + "-mongodb-check"

	_, _, err = consul.ReadService(ID)
	if err == nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": "service already exists",
		})
		log.Println("service already exists")
		return
	}

	var consulData = g.ConsulData{
		ID:      ID,
		Name:    consulDataJson.Group,
		Address: consulDataJson.Address,
		Port:    consulDataJson.Port,
		Tags:    g.CF.Config.Tags.Mongodb_exporter,
	}

	checkUrl := ""
	checkPort := 0
	err = consul.WriteService(consulData, checkUrl, checkPort)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	address, port, err := consul.ReadService(consulData.ID)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": fmt.Sprintf("%s: %n", address, port),
	})
}

func Addssl_exporter(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	var consulDataJson *g.ConsulDataJson
	err = json.Unmarshal(data, &consulDataJson)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	ID := strings.Replace(consulDataJson.Address, ".", "_", -1) + "-ssl-check"

	_, _, err = consul.ReadService(ID)
	if err == nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": "service already exists",
		})
		log.Println("service already exists")
		return
	}

	var consulData = g.ConsulData{
		ID:      ID,
		Name:    consulDataJson.Group,
		Address: consulDataJson.Address,
		Port:    443,
		Tags:    g.CF.Config.Tags.SSL_exporter,
	}

	checkUrl := ""
	checkPort := 0
	err = consul.WriteService(consulData, checkUrl, checkPort)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	address, port, err := consul.ReadService(consulData.ID)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": fmt.Sprintf("%s: %n", address, port),
	})
}

func Adddomain_exporter(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	var consulDataJson *g.ConsulDataJson
	err = json.Unmarshal(data, &consulDataJson)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	ID := strings.Replace(consulDataJson.Address, ".", "_", -1) + "-domain-check"

	_, _, err = consul.ReadService(ID)
	if err == nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": "service already exists",
		})
		log.Println("service already exists")
		return
	}

	var consulData = g.ConsulData{
		ID:      ID,
		Name:    consulDataJson.Group,
		Address: consulDataJson.Address,
		Port:    80,
		Tags:    g.CF.Config.Tags.Domain_exporter,
	}

	checkUrl := ""
	checkPort := 0
	err = consul.WriteService(consulData, checkUrl, checkPort)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	address, port, err := consul.ReadService(consulData.ID)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": fmt.Sprintf("%s: %n", address, port),
	})
}

func DeleteService(c *gin.Context) {
	item := c.PostForm("message")
	_, _, err := consul.ReadService(item)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
	}

	err = consul.DeleteService(item)
	if err != nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "ok",
	})
}

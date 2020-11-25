package http

import (
	"add-monitor-to-prometheu/consul"
	"add-monitor-to-prometheu/g"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)


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

func Get_exporter(c *gin.Context) {
	c.HTML(http.StatusOK, "get_exporter.tmpl", nil)
}

func SSL_exporter(c *gin.Context) {
	c.HTML(http.StatusOK, "ssl_exporter.tmpl", nil)
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
	var consulData =  g.ConsulData {
		ID: consulDataJson.Hostname,
		Name: consulDataJson.Group,
		Address: consulDataJson.Address,
		Port: consulDataJson.Port,
		Tags: g.CF.Config.Tags.Node_exporter ,
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
		"success":  fmt.Sprintf("%s: %n", address, port),
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


	var consulData =  g.ConsulData {
		ID: ID,
		Name: consulDataJson.Group,
		Address: consulDataJson.Address,
		Port: consulDataJson.Port,
		Tags: g.CF.Config.Tags.Port_exporter ,
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
		"success":  fmt.Sprintf("%s: %n", address, port),
	})
}


func Addget_exporter(c *gin.Context) {
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


	ID := strings.Replace(consulDataJson.Address, ".", "_",-1 ) + "-httpget-check"

	_, _, err = consul.ReadService(ID)
	if err == nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": "service already exists",
		})
		log.Println("service already exists")
		return
	}


	var consulData =  g.ConsulData {
		ID: ID,
		Name: consulDataJson.Group,
		Address: consulDataJson.Address,
		Port: consulDataJson.Port,
		Tags: g.CF.Config.Tags.Get_exporter ,
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
		"success":  fmt.Sprintf("%s: %n", address, port),
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


	ID := strings.Replace(consulDataJson.Address, ".", "_",-1 ) + "-ssl-check"

	_, _, err = consul.ReadService(ID)
	if err == nil {
		c.JSON(http.StatusNotExtended, gin.H{
			"error": "service already exists",
		})
		log.Println("service already exists")
		return
	}


	var consulData =  g.ConsulData {
		ID: ID,
		Name: consulDataJson.Group,
		Address: consulDataJson.Address,
		Port: 443,
		Tags: g.CF.Config.Tags.SSL_exporter ,
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
		"success":  fmt.Sprintf("%s: %n", address, port),
	})
}



func DeleteService(c *gin.Context) {
	item := c.PostForm("message")
	_,_,err := consul.ReadService(item)
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
		"success":  "ok",
	})
}


package main

import (
	"automonitor/g"
	"automonitor/http"
	"flag"
	"log"
)

func init() {
	ConfigFile := flag.String("conf", "cfg.yaml", "Config file for this listener and ldap configs")
	flag.Parse()

	if err := g.ReadYamlConfig(*ConfigFile); err != nil {
		log.Fatal(err)
	}
}

func main()  {
	http.Run()
}

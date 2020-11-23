package g

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var CF *Cfg

type Cfg struct {
	Config config `yaml:"config"`
}

type config struct {
	ConsulServerList []string  `yaml:"consulServerList"`
	Tags tags `yaml:"tags"`
	HttpListen httpListen `yaml:"httpListen"`
}


type tags struct {
	Node_exporter []string `yaml:"node_exporter"`
	Port_exporter []string `yaml:"port_exporter"`
	Get_exporter []string `yaml:"get_exporter"`
	SSL_exporter []string `yaml:"ssl_exporter"`
}

type httpListen struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}


func ReadYamlConfig(path string)  (error){
	CF = new(Cfg)
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = yaml.Unmarshal(yamlFile, &CF)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}


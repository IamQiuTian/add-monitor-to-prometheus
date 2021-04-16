package g

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var CF *Cfg

type Cfg struct {
	Config config `yaml:"config"`
}

type config struct {
	ConsulServerList []string      `yaml:"consulServerList"`
	Tags             tags          `yaml:"tags"`
	HttpListen       httpListen    `yaml:"httpListen"`
	AuthorizaTion    authorizaTion `yaml:"authorizaTion"`
}

type tags struct {
	Node_exporter    []string `yaml:"node_exporter"`
	Port_exporter    []string `yaml:"port_exporter"`
	Mysqld_exporter  []string `yaml:"mysqld_exporter"`
	Mongodb_exporter []string `yaml:"mongodb_exporter"`
	SSL_exporter     []string `yaml:"ssl_exporter"`
	Domain_exporter  []string `yaml:"domain_exporter"`
}

type httpListen struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type authorizaTion struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func ReadYamlConfig(path string) error {
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

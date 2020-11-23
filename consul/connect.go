package consul

import (
	consulapi "github.com/hashicorp/consul/api"
)

var config = *consulapi.DefaultConfig()

func InitConfig(ipAddress string) (client *consulapi.Client, err error) {
	config.Address = ipAddress
	client, err = consulapi.NewClient(&config)
	if err != nil {
		return client, err
	}
	return
}
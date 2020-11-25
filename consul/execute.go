package consul

import (
	"add-monitor-to-prometheus/g"
	"errors"
	consulapi "github.com/hashicorp/consul/api"
	"log"
	"fmt"
)


func ReadService(serviceid string) (address string, port int, err error) {
	for _, ipaddr := range g.CF.Config.ConsulServerList {
		client,err := InitConfig(ipaddr)
		if err != nil {
			log.Println(err)
			continue
		}
		srv, _, err := client.Agent().Service(serviceid, nil)
		if err !=nil {
			log.Println(err)
			continue
		}
		return srv.Address, srv.Port, nil
	}
	return "",int(0), errors.New("Not Fond Service")
}

func WriteService(consuldata g.ConsulData, checkUrl string, checkPort int) error {
	for _, ipaddr := range g.CF.Config.ConsulServerList {
		client, err := InitConfig(ipaddr)
		if err != nil {
			log.Println(err)
			continue
		}

		//注册服务参数
		registration := new(consulapi.AgentServiceRegistration)
		registration.ID = consuldata.ID
		registration.Name = consuldata.Name
		registration.Address = consuldata.Address
		registration.Port = consuldata.Port
		registration.Tags = consuldata.Tags

		var check = new(consulapi.AgentServiceCheck)
		if checkUrl != "" {
			check.HTTP = fmt.Sprintf("http://%s:%d/metrics", checkUrl, checkPort)
			check.Timeout = "15s"
			check.Interval = "35s"

			registration.Check = check
		}

		err = client.Agent().ServiceRegister(registration)
		if err != nil {
			log.Println(err)
			continue
		}
		return nil
	}
	return errors.New("registration failed")
}

func DeleteService(seviceid string) error {
	for _, ipaddr := range g.CF.Config.ConsulServerList {
		client, err := InitConfig(ipaddr)
		if err != nil {
			log.Println(err)
			continue
		}

		err = client.Agent().ServiceDeregister(seviceid)
		if err != nil {
			log.Println(err)
			continue
		}
		_, _, err = ReadService(seviceid)
		if err != nil {
			return nil
		}
	}
	return errors.New("Delete failed or no such service")
}

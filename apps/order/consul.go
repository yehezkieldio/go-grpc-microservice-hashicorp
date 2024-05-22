package main

import (
	"os"

	capi "github.com/hashicorp/consul/api"
)

var Consul *capi.Client

func InitConsul() {
	consulConfig := capi.DefaultConfig()
	consulConfig.Address = "consul:8500"
	Consul, _ = capi.NewClient(consulConfig)
}

func getHostname() (hostname string) {
	hostname, _ = os.Hostname()
	return
}

func RegisterOrderService() {
	serviceId := "order"
	port := 50051
	address := getHostname()

	registration := &capi.AgentServiceRegistration{
		ID:      serviceId,
		Name:    serviceId,
		Address: address,
		Port:    port,
		Check: &capi.AgentServiceCheck{
			HTTP:     "http://" + address + ":" + string(rune(port)) + "/check",
			Interval: "10s",
		},
	}

	registrationError := Consul.Agent().ServiceRegister(registration)
	if registrationError != nil {
		panic(registrationError)
	}
}

func RegisterInventoryService() {
	serviceId := "inventory"
	port := 50051
	address := getHostname()

	registration := &capi.AgentServiceRegistration{
		ID:      serviceId,
		Name:    serviceId,
		Address: address,
		Port:    port,
		Check: &capi.AgentServiceCheck{
			GRPC:       address + ":" + string(rune(port)),
			GRPCUseTLS: false,
			Interval:   "10s",
		},
	}

	registrationError := Consul.Agent().ServiceRegister(registration)
	if registrationError != nil {
		panic(registrationError)
	}
}

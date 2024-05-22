package main

import (
	"fmt"
	"os"
	"strconv"

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

func getPort() (port string) {
 port = os.Getenv("PORT")
 if len(port) == 0 {
  port = "8080"
 }
 port = ":" + port
 return
}

func RegisterOrderService() {
	serviceId := "order"
	port, _ := strconv.Atoi(getPort()[1:len(getPort())])
	address := getHostname()

	fmt.Printf("Registering order service with address: %s and port: %d\n", address, port)

	registration := &capi.AgentServiceRegistration{
		ID:      serviceId,
		Name:    serviceId,
		Address: address,
		Port:    port,
		Check: &capi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/check", address, port),
			Interval: "10s",
		},
	}

	registrationError := Consul.Agent().ServiceRegister(registration)
	if registrationError != nil {
		fmt.Println("Error registering service with consul")
		panic(registrationError)
	}

	fmt.Printf("Registered order service with ID: %s\n", serviceId)
}

func RegisterInventoryService() {
	serviceId := "inventory"
	port := 50051
	address := "inventory"

	fmt.Printf("Registering inventory service with address: %s and port: %d\n", address, port)

	registration := &capi.AgentServiceRegistration{
		ID:      serviceId,
		Name:    serviceId,
		Address: address,
		Port:    port,
		Check: &capi.AgentServiceCheck{
			GRPC: "inventory:50051",
			Interval:   "10s",
		},
	}

	registrationError := Consul.Agent().ServiceRegister(registration)
	if registrationError != nil {
		fmt.Println("Error registering service with consul")
		panic(registrationError)
	}

	fmt.Printf("Registered inventory service with ID: %s\n", serviceId)
}

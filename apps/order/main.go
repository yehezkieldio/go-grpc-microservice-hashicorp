package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	pb "go-grpc-microservice-hashicorp/gen"

	capi "github.com/hashicorp/consul/api"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getPort() (port string) {
	port = os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	port = ":" + port
	return
}

func getHostname() (hostname string) {
	hostname, _ = os.Hostname()
	return
}

func main() {
	// TODO: Read the server address from the environment variable or a configuration file
	consulConfig := capi.DefaultConfig()
	consulConfig.Address = "consul:8500"
	consul, _ := capi.NewClient(consulConfig)


	serviceId := "order"
	port, _ := strconv.Atoi(getPort()[1:len(getPort())])
	fmt.Println(port)
	address := getHostname()
	fmt.Println(address)

	registration := &capi.AgentServiceRegistration{
		ID:      serviceId,
		Name:   serviceId,
		Address: address,
		Port:   port,
		Check: &capi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/check", address, port),
			Interval: "10s",
		},
	}

	regError := consul.Agent().ServiceRegister(registration)
	if regError != nil {
		fmt.Println(regError)
	}

	// FIXME: Replace with your own server address
	serverAddr := "dns:///inventory:50051"

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connection, err := grpc.DialContext(ctx, serverAddr, opts...)
	if err != nil {
		log.Fatalln(err)
	}
	defer connection.Close()

	client := pb.NewInventoryClient(connection)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})
	e.GET("/check", func(c echo.Context) error {
		return c.String(200, "Health Check OK!")
	})
	e.GET("/order/:item_id", func(c echo.Context) error {
		itemID := c.Param("item_id")
		res, err := client.GetInventory(ctx, &pb.InventoryRequest{
			ItemId: itemID,
		})
		if err != nil {
			return c.String(500, err.Error())
		}

		return c.String(200, fmt.Sprintf("Item ID: %s, Quantity: %d", res.ItemId, res.Quantity))
	})
	e.Logger.Fatal(e.Start(":8080"))
}

package main

import (
	"context"
	"fmt"
	"log"

	pb "go-grpc-microservice-hashicorp/gen"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	InitConsul()
	RegisterInventoryService()
	RegisterOrderService()

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

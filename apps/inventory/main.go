package main

import (
	"context"
	pb "go-grpc-microservice-hashicorp/gen"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedInventoryServer
}

func (s *server) GetInventory(ctx context.Context, in *pb.InventoryRequest) (*pb.InventoryResponse, error) {
	log.Printf(("Received: %v"), in.ItemId)

	return &pb.InventoryResponse{ItemId: in.ItemId, Quantity: 100}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterInventoryServer(s, &server{})
	log.Println("Starting gRPC server on :50051")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

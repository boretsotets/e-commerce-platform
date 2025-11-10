package main

import (
	"fmt"
	"log"
	"net"

	client "github.com/boretsotets/e-commerce-platform/order-service/cmd"
	"github.com/boretsotets/e-commerce-platform/order-service/cmd/service"
	"github.com/boretsotets/e-commerce-platform/order-service/internal/repository"
	pb "github.com/boretsotets/e-commerce-platform/order-service/pkg/api"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	productClient := client.NewProductClient("localhost:50051")

	repo := repository.NewInmemOrderRepo()
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, &service.OrderServer{Repo: repo, ProductClient: productClient})

	fmt.Println("gRPC server listening on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

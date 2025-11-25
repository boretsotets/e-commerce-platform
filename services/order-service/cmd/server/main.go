package main

import (
	"fmt"
	"log"
	"net"

	service "github.com/boretsotets/e-commerce-platform/order-service/internal/application/usecase"
	client "github.com/boretsotets/e-commerce-platform/order-service/internal/infra/clients/productclient"
	repository "github.com/boretsotets/e-commerce-platform/order-service/internal/infra/persistence"
	controller "github.com/boretsotets/e-commerce-platform/order-service/internal/infra/transport/grpc"
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
	orderservice := service.NewOrderServiceServer(repo, productClient)
	pb.RegisterOrderServiceServer(s, &controller.OrderServer{Service: *orderservice})

	fmt.Println("gRPC server listening on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

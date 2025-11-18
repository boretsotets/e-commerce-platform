package main

import (
	"fmt"
	"log"
	"net"

	service "github.com/boretsotets/e-commerce-platform/product-service/internal/application/usecase"
	repository "github.com/boretsotets/e-commerce-platform/product-service/internal/infra/persistence"
	server "github.com/boretsotets/e-commerce-platform/product-service/internal/infra/transport/grpc"
	pb "github.com/boretsotets/e-commerce-platform/product-service/pkg/api"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	repo := repository.NewInmemProductRepo()
	s := grpc.NewServer()
	service := service.NewProductService(repo)

	pb.RegisterProductServiceServer(s, &server.ProductServer{Service: *service})

	fmt.Println("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

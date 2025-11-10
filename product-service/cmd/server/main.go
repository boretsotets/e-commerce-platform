package main

import (
	"fmt"
	"log"
	"net"

	"github.com/boretsotets/e-commerce-platform/product-service/internal/repository"
	productgetter "github.com/boretsotets/e-commerce-platform/product-service/internal/service"
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
	pb.RegisterProductServiceServer(s, &productgetter.ProductServer{Repo: repo})

	fmt.Println("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package app

import (
	"fmt"
	"log"
	"net"

	"github.com/boretsotets/e-commerce-platform/product-service/internal/application/usecase"
	"github.com/boretsotets/e-commerce-platform/product-service/internal/infra/db"
	"github.com/boretsotets/e-commerce-platform/product-service/internal/infra/persistence"
	server "github.com/boretsotets/e-commerce-platform/product-service/internal/infra/transport/grpc"
	pb "github.com/boretsotets/e-commerce-platform/product-service/pkg/api"
	"google.golang.org/grpc"
)

func Run() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db, err := db.NewPostgres()
	if err != nil {

	}
	repo := persistence.NewProductRepo(db)
	s := grpc.NewServer()
	service := usecase.NewProductService(repo)

	pb.RegisterProductServiceServer(s, &server.ProductServer{Service: *service})

	fmt.Println("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

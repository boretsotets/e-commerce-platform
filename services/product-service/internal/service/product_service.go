package service

import (
	"context"

	"github.com/boretsotets/e-commerce-platform/product-service/internal/repository"
	pb "github.com/boretsotets/e-commerce-platform/product-service/pkg/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductServer struct {
	pb.UnimplementedProductServiceServer
	Repo repository.ProductRepository
}

func (s *ProductServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	response, err := s.Repo.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id:    int64(response.Id),
			Name:  response.Name,
			Price: response.Price,
			Stock: response.Stock,
		},
	}, nil
}

func (s *ProductServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	response, err := s.Repo.InsertNewProduct(ctx, req.Name, req.Price, req.Stock)
	if err != nil {
		return nil, err
	}

	return &pb.CreateProductResponse{
		Product: &pb.Product{
			Id:    int64(response.Id),
			Name:  response.Name,
			Price: response.Price,
			Stock: response.Stock,
		},
	}, nil
}

func (s *ProductServer) BatchChangeStock(ctx context.Context, req *pb.BatchChangeStockRequest) (*emptypb.Empty, error) {
	err := s.Repo.BatchChangeStock(ctx, req.Items)
	return nil, err
}

func (s *ProductServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	response, err := s.Repo.GetList(ctx, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	return &pb.ListProductsResponse{
		Products: response,
	}, nil
}

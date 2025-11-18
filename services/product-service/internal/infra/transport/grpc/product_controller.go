package server

import (
	"context"

	service "github.com/boretsotets/e-commerce-platform/product-service/internal/application/usecase"
	"github.com/boretsotets/e-commerce-platform/product-service/internal/domain/models"
	pb "github.com/boretsotets/e-commerce-platform/product-service/pkg/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductServer struct {
	pb.UnimplementedProductServiceServer
	Service service.ProductService
}

func (s *ProductServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	product, err := s.Service.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id:    product.ID,
			Name:  product.Name,
			Price: product.Price,
			Stock: product.Stock,
		},
	}, nil
}

func (s *ProductServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	product, err := s.Service.CreateProduct(ctx, req.Name, req.Price, req.Stock)
	if err != nil {
		return nil, err
	}
	return &pb.CreateProductResponse{
		Product: &pb.Product{
			Id:    product.ID,
			Name:  product.Name,
			Price: product.Price,
			Stock: product.Stock,
		},
	}, nil
}

func (s *ProductServer) BatchChangeStock(ctx context.Context, req *pb.BatchChangeStockRequest) (*emptypb.Empty, error) {
	modelsStockChange := PbtomodelStockChangeItem(req.Items)
	err := s.Service.BatchChangeStock(ctx, modelsStockChange)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func PbtomodelStockChangeItem(pblist []*pb.StockChangeItem) []*models.StockChangeItem {
	res := make([]*models.StockChangeItem, len(pblist))
	for i := range pblist {
		res[i].ProductID = pblist[i].ProductId
		res[i].Delta = pblist[i].Delta
	}
	return res
}

func (s *ProductServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	productList, err := s.Service.ListProducts(ctx, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}
	return &pb.ListProductsResponse{
		Products: PbtomodelProductList(productList),
	}, nil
}

func PbtomodelProductList(input []*models.Product) []*pb.Product {
	res := make([]*pb.Product, len(input))
	for i := range input {
		res[i].Id = input[i].ID
		res[i].Name = input[i].Name
		res[i].Price = input[i].Price
		res[i].Stock = input[i].Stock
	}
	return res
}

func (s *ProductServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*emptypb.Empty, error) {
	err := s.Service.DeleteProduct(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

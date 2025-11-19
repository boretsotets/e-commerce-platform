package service

import (
	"context"
	"errors"

	"github.com/boretsotets/e-commerce-platform/product-service/internal/domain"
	"github.com/boretsotets/e-commerce-platform/product-service/internal/domain/models"
)

type ProductService struct {
	Repo domain.ProductRepository
}

func NewProductService(Repo domain.ProductRepository) *ProductService {
	return &ProductService{Repo: Repo}
}

func (s *ProductService) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	response, err := s.Repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.Product{
		ID:    int64(response.ID),
		Name:  response.Name,
		Price: response.Price,
		Stock: response.Stock,
	}, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, ProductName string, ProductPrice float64, ProductStock int32) (*models.Product, error) {
	productExsist := s.Repo.CheckProductExsistence(ctx, ProductName)
	if productExsist {
		return nil, errors.New("Product with this ProductName already exists")
	}
	if ProductPrice <= 0 || ProductStock < 0 {
		return nil, errors.New("Price should be > 0 and stock shold be >= 0")
	}

	response, err := s.Repo.InsertNewProduct(ctx, ProductName, ProductPrice, ProductStock)
	if err != nil {
		return nil, err
	}

	return &models.Product{
		ID:    int64(response.ID),
		Name:  response.Name,
		Price: response.Price,
		Stock: response.Stock,
	}, nil
}

func (s *ProductService) BatchChangeStock(ctx context.Context, List []*models.StockChangeItem) error {
	err := s.Repo.BatchChangeStock(ctx, List)
	return err
}

func (s *ProductService) ListProducts(ctx context.Context, Page int32, Limit int32) ([]*models.Product, error) {
	response, err := s.Repo.GetList(ctx, Page, Limit)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, ProductID int64) error {
	return nil
}

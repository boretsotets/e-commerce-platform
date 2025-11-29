package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/boretsotets/e-commerce-platform/product-service/internal/domain"
	"github.com/boretsotets/e-commerce-platform/product-service/internal/domain/models"
)

var (
	ErrInvalidInput   = errors.New("invalid input")
	ErrNotFound       = errors.New("not found")
	ErrConflict       = errors.New("conflict")
	ErrNotEnoughStock = errors.New("not enough stock")
)

type ProductService struct {
	Repo domain.ProductRepository
}

func NewProductService(Repo domain.ProductRepository) *ProductService {
	return &ProductService{Repo: Repo}
}

func (s *ProductService) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	if id < 1 {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, "id must be positive")
	}
	response, err := s.Repo.GetById(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, fmt.Errorf("%w: %v", ErrNotFound, "product not found")
		}
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
	if ProductPrice <= 0 {
		return nil, fmt.Errorf("%w: price must be > 0, got %f", ErrInvalidInput, ProductPrice)
	}
	if ProductStock < 0 {
		return nil, fmt.Errorf("%w: stock must be > 0, got %d", ErrInvalidInput, ProductStock)
	}

	productExsist, err := s.Repo.CheckProductExsistence(ctx, ProductName)
	if productExsist {
		return nil, fmt.Errorf("%w: product with product name %s already exists", ErrConflict, ProductName)
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
	for _, v := range List {
		product, err := s.Repo.GetById(ctx, v.ProductID)
		if err != nil {
			return err
		}
		if product.Stock < v.Delta {
			return fmt.Errorf("%w: product %d requires %d, only %d available",
				ErrNotEnoughStock, v.ProductID, v.Delta, product.Stock)
		}
	}
	err := s.Repo.BatchChangeStock(ctx, List)
	return err
}

func (s *ProductService) ListProducts(ctx context.Context, Page int32, Limit int32) ([]*models.Product, error) {
	if Page < 1 {
		return nil, fmt.Errorf("%w: page must be > 0, got %d", ErrInvalidInput, Page)
	}
	if Limit < 1 {
		return nil, fmt.Errorf("%w: limit must be > 0, got %d", ErrInvalidInput, Limit)
	}

	response, err := s.Repo.GetList(ctx, Page, Limit)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, ProductID int64) error {
	// error if deleting product that is in active order
	if ProductID < 1 {
		return fmt.Errorf("%w: %v", ErrInvalidInput, "invalid id")
	}
	err := s.Repo.DeleteProduct(ctx, ProductID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return fmt.Errorf("%w: %v", ErrNotFound, "product not found")
		}
		return err
	}
	return nil
}

func (s *ProductService) ApplyOrderEvent(ctx context.Context, evt models.EventOrderCreated) error {
	return nil
}

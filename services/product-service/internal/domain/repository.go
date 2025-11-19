package domain

import (
	"context"

	"github.com/boretsotets/e-commerce-platform/product-service/internal/domain/models"
)

type ProductRepository interface {
	GetById(ctx context.Context, id int64) (*models.Product, error)
	CheckProductExsistence(ctx context.Context, name string) bool
	InsertNewProduct(ctx context.Context, name string, price float64, stock int32) (*models.Product, error)
	UpdateStock(ctx context.Context, id int64, delta int32) (int32, error)
	BatchChangeStock(ctx context.Context, items []*models.StockChangeItem) error
	GetList(ctx context.Context, page int32, limit int32) ([]*models.Product, error)
	DeleteProduct(ctx context.Context, productID int64) error
}

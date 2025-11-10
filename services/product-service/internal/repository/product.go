package repository

import (
	"context"

	product "github.com/boretsotets/e-commerce-platform/product-service/pkg/api"
)

type Product struct {
	ID    int64
	Name  string
	Price float64
	Stock int32
}

type ProductRepository interface {
	GetById(ctx context.Context, id int64) (*product.Product, error)
	InsertNewProduct(ctx context.Context, name string, price float64, stock int32) (*product.Product, error)
	UpdateStock(ctx context.Context, id int64, delta int32) (int32, error)
	BatchChangeStock(ctx context.Context, items []*product.StockChangeItem) error
	GetList(ctx context.Context, page int32, limit int32) ([]*product.Product, error)
}

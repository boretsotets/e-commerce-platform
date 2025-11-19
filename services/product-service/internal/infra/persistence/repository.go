package repository

import (
	"context"
	"errors"

	"github.com/boretsotets/e-commerce-platform/product-service/internal/domain/models"
	"gorm.io/gorm"
)

type InmemProductRepo struct {
	data []*models.Product
}

type ProductRepo struct {
	Repo *gorm.DB
}

func NewProductRepo(repo *gorm.DB) *ProductRepo {
	return &ProductRepo{Repo: repo}
}

func NewInmemProductRepo() *InmemProductRepo {
	return &InmemProductRepo{
		data: []*models.Product{
			{ID: 1, Name: "Pen", Price: 1.5, Stock: 100},
			{ID: 2, Name: "Notebook", Price: 3.2, Stock: 50},
		},
	}
}

func (r *ProductRepo) GetById(ctx context.Context, id int64) (*models.Product, error) {
	var p models.Product
	err := r.Repo.WithContext(ctx).First(&p, id).Error
	return &p, err
}

func (r *InmemProductRepo) GetById(ctx context.Context, id int64) (*models.Product, error) {
	if id >= int64(len(r.data)) {
		return nil, errors.New("not found")
	}
	return r.data[id], nil
}

func (r *InmemProductRepo) CheckProductExsistence(ctx context.Context, name string) bool {
	for i := range r.data {
		if r.data[i].Name == name {
			return true
		}
	}
	return false
}

func (r *InmemProductRepo) GetList(ctx context.Context, page int32, limit int32) ([]*models.Product, error) {
	start := page * limit
	if start+limit >= int32(len(r.data)) {
		return nil, errors.New("Indexes out of range")
	}
	return r.data[start : start+limit], nil
}

func (r *InmemProductRepo) InsertNewProduct(ctx context.Context, name string, price float64, stock int32) (*models.Product, error) {
	newProduct := models.Product{
		ID:    int64(len(r.data) + 1),
		Name:  name,
		Price: price,
		Stock: stock,
	}
	r.data = append(r.data, &newProduct)
	return &newProduct, nil
}

func (r *InmemProductRepo) UpdateStock(ctx context.Context, id int64, delta int32) (int32, error) {
	r.data[id].Stock += delta
	return r.data[id].Stock, nil
}

func (r *InmemProductRepo) BatchChangeStock(ctx context.Context, items []*models.StockChangeItem) error {
	for _, i := range items {
		delta := r.data[i.ProductID].Stock - i.Delta
		if delta < 0 {
			return errors.New("Not enough stock left")
		}
		r.data[i.ProductID].Stock = delta
	}

	return nil
}

func (r *InmemProductRepo) DeleteProduct(ctx context.Context, productID int64) error {

	return nil
}

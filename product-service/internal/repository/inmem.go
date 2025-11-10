package repository

import (
	"context"
	"errors"

	product "github.com/boretsotets/e-commerce-platform/product-service/pkg/api"
)

type InmemProductRepo struct {
	data []*product.Product
}

func NewInmemProductRepo() *InmemProductRepo {
	return &InmemProductRepo{
		data: []*product.Product{
			{Id: 1, Name: "Pen", Price: 1.5, Stock: 100},
			{Id: 2, Name: "Notebook", Price: 3.2, Stock: 50},
		},
	}
}

func (r *InmemProductRepo) GetById(ctx context.Context, id int64) (*product.Product, error) {
	if id >= int64(len(r.data)) {
		return nil, errors.New("not found")
	}
	return r.data[id], nil
}

func (r *InmemProductRepo) GetList(ctx context.Context, page int32, limit int32) ([]*product.Product, error) {
	start := page * limit
	if start+limit >= int32(len(r.data)) {
		return nil, errors.New("Indexes out of range")
	}
	return r.data[start : start+limit], nil
}

func (r *InmemProductRepo) InsertNewProduct(ctx context.Context, name string, price float64, stock int32) (*product.Product, error) {
	newProduct := product.Product{
		Id:    int64(len(r.data) + 1),
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

func (r *InmemProductRepo) BatchChangeStock(ctx context.Context, items []*product.StockChangeItem) error {
	for _, i := range items {
		r.data[i.ProductId].Stock += i.Delta
	}

	return nil
}

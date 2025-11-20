package persistence

import (
	"context"

	"github.com/boretsotets/e-commerce-platform/product-service/internal/domain/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepo struct {
	Repo *gorm.DB
}

func NewProductRepo(repo *gorm.DB) *ProductRepo {
	return &ProductRepo{Repo: repo}
}

func (r *ProductRepo) GetById(ctx context.Context, id int64) (*models.Product, error) {
	var p models.Product
	err := r.Repo.WithContext(ctx).First(&p, id).Error
	return &p, err
}

func (r *ProductRepo) CheckProductExsistence(ctx context.Context, name string) (bool, error) {
	var p models.Product
	err := r.Repo.WithContext(ctx).Where("name = ?", name).First(&p).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ProductRepo) GetList(ctx context.Context, offset int32, limit int32) ([]*models.Product, error) {
	var products []*models.Product
	err := r.Repo.WithContext(ctx).Offset(int(offset)).Limit(int(limit)).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepo) InsertNewProduct(ctx context.Context, name string, price float64, stock int32) (*models.Product, error) {
	record := &models.Product{
		Name:  name,
		Price: price,
		Stock: stock,
	}

	// check if this returns to record
	err := r.Repo.Create(record).Error
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (r *ProductRepo) UpdateStock(ctx context.Context, id int64, delta int32) (int32, error) {
	var p models.Product
	err := r.Repo.WithContext(ctx).
		Where("productID = ?", id).
		Clauses(clause.Returning{}).
		Update("stock", delta).
		Scan(&p).
		Error
	if err != nil {
		return 0, err
	}
	return p.Stock, nil

}

func (r *ProductRepo) BatchChangeStock(ctx context.Context, items []*models.StockChangeItem) error {
	for _, item := range items {
		err := r.Repo.WithContext(ctx).
			Where("productID = ?", item.ProductID).
			Update("stock", gorm.Expr("stock + ?", item.Delta)).
			Error

		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ProductRepo) DeleteProduct(ctx context.Context, productID int64) error {
	var p models.Product
	err := r.Repo.WithContext(ctx).Find(&p, productID).Error
	if err != nil {
		return err
	}
	err = r.Repo.Delete(&p).Error
	if err != nil {
		return err
	}
	return nil
}

package persistence

import (
	"context"
	"log"
	"testing"

	"github.com/boretsotets/e-commerce-platform/product-service/internal/infra/db"
)

func TestGetByID(t *testing.T) {
	ctx := context.Background()

	db, err := db.NewTestPostgres()
	if err != nil {
		log.Fatalf("error connecting to database: %v\n", err)
	}

	repo := NewProductRepo(db)

	// 1. тест создания продукта
	product, err := repo.InsertNewProduct(ctx, "Карандаш", 10.5, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if product.ID == 1 && product.Name == "Карандаш" && product.Price == 10.5 && product.Stock == 100 {
	} else {
		t.Errorf("inserted product does not match input")
	}

	// 2. тест получения продукта по ID
	product, err = repo.GetById(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

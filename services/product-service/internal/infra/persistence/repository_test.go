package repository

import (
	"context"
	"testing"
)

func TestGetByID(t *testing.T) {
	ctx := context.Background()

	repo := NewInmemProductRepo()

	// тест создания продукта
	product, err := repo.InsertNewProduct(ctx, "Карандаш", 10.5, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if product.ID == 1 && product.Name == "Карандаш" && product.Price == 10.5 && product.Stock == 100 {
	} else {
		t.Errorf("inserted product does not match input")
	}

}

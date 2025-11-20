package usecase

import (
	"testing"

	"github.com/boretsotets/e-commerce-platform/product-service/internal/domain/models"
)

type mockProductRepo struct {
	product []*models.Product
}

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name         string
		productID    int64
		productName  string
		productPrice float64
		productStock int32
		productInfo  models.Product
		wantError    bool
	}{
		{
			name:         "Успешное создание продукта",
			productID:    1,
			productName:  "Ручка",
			productPrice: 150.5,
			productStock: 15,
			productInfo: models.Product{
				ID:    1,
				Name:  "Ручка",
				Price: 150.5,
				Stock: 15,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mockProductRepo{}

			service := NewProductService(mockRepo)
		})
	}

}

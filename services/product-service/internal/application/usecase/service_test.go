package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/boretsotets/e-commerce-platform/product-service/internal/domain/models"
	mocks "github.com/boretsotets/e-commerce-platform/product-service/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestService_CreateProduct(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockProductRepository(ctrl)
	srv := NewProductService(mockRepo)

	expectedProduct := &models.Product{
		ID:    1,
		Name:  "Ink",
		Price: 30,
		Stock: 70,
	}

	mockRepo.EXPECT().CheckProductExsistence(ctx, expectedProduct.Name).Return(false, nil)
	mockRepo.
		EXPECT().
		InsertNewProduct(ctx, expectedProduct.Name, expectedProduct.Price, expectedProduct.Stock).
		Return(expectedProduct, nil)

	result, err := srv.CreateProduct(ctx, "Ink", 30, 70)

	require.NoError(t, err)
	require.Equal(t, expectedProduct, result)
}

func TestService_CreateProduct_DuplicateName(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockProductRepository(ctrl)
	srv := NewProductService(mockRepo)

	t.Run("duplicate name error", func(t *testing.T) {
		mockRepo.EXPECT().CheckProductExsistence(ctx, "Ink").Return(true, nil)

		_, err := srv.CreateProduct(ctx, "Ink", float64(15), int32(15))

		require.Error(t, err)
		require.Contains(t, err.Error(), "Product with this ProductName already exists")
	})
}

func TestService_CreateProduct_InvalidPrice(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockProductRepository(ctrl)
	srv := NewProductService(mockRepo)

	_, err := srv.CreateProduct(ctx, "invalid product", 0, 15)

	require.Error(t, err)
	require.Contains(t, err.Error(), "Price should be > 0 and stock shold be >= 0")
}

func TestService_CreateProduct_InvalidStock(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockProductRepository(ctrl)
	srv := NewProductService(mockRepo)

	_, err := srv.CreateProduct(ctx, "invalid product", 20, -1)

	require.Error(t, err)
	require.Contains(t, err.Error(), "Price should be > 0 and stock shold be >= 0")
}

func TestService_GetProduct(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockProductRepository(ctrl)
	srv := NewProductService(mockRepo)

	expectedProduct := &models.Product{
		ID:    1,
		Name:  "Ink",
		Price: 30,
		Stock: 70,
	}

	t.Run("get by id", func(t *testing.T) {
		mockRepo.EXPECT().GetById(ctx, int64(1)).Return(expectedProduct, nil)
		result, err := srv.GetProduct(ctx, int64(1))

		require.NoError(t, err)
		require.Equal(t, expectedProduct, result)
	})
}

func TestService_GetProduct_InvalidID(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockProductRepository(ctrl)
	srv := NewProductService(mockRepo)

	_, err := srv.GetProduct(ctx, 0)

	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid id")
}

func TestService_BatchChangeStock(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockProductRepository(ctrl)
	srv := NewProductService(mockRepo)

	one := &models.Product{
		ID:    1,
		Name:  "one",
		Price: 10,
		Stock: 30,
	}

	two := &models.Product{
		ID:    1,
		Name:  "two",
		Price: 10,
		Stock: 30,
	}

	t.Run("insert product one", func(t *testing.T) {
		mockRepo.EXPECT().CheckProductExsistence(ctx, "one").Return(false, nil)
		mockRepo.EXPECT().InsertNewProduct(ctx, "one", float64(10), int32(30)).Return(&models.Product{}, nil)
		srv.CreateProduct(ctx, "one", 10, 30)
	})

	t.Run("insert product two", func(t *testing.T) {
		mockRepo.EXPECT().CheckProductExsistence(ctx, "two").Return(false, nil)
		mockRepo.EXPECT().InsertNewProduct(ctx, "two", float64(10), int32(30)).Return(&models.Product{}, nil)
		srv.CreateProduct(ctx, "two", 10, 30)
	})

	input := []*models.StockChangeItem{
		{
			ProductID: one.ID,
			Delta:     10,
		},
		{
			ProductID: two.ID,
			Delta:     20,
		},
	}

	t.Run("batch change stock success", func(t *testing.T) {
		mockRepo.EXPECT().GetById(ctx, one.ID).Return(one, nil)
		mockRepo.EXPECT().GetById(ctx, two.ID).Return(two, nil)

		mockRepo.EXPECT().BatchChangeStock(ctx, input).Return(nil)

		err := srv.BatchChangeStock(ctx, input)
		require.NoError(t, err)
	})

	inputError := []*models.StockChangeItem{
		{
			ProductID: one.ID,
			Delta:     10,
		},
		{
			ProductID: two.ID,
			Delta:     60,
		},
	}
	t.Run("batch change stock error", func(t *testing.T) {
		mockRepo.EXPECT().GetById(ctx, one.ID).Return(one, nil)
		mockRepo.EXPECT().GetById(ctx, two.ID).Return(two, nil)

		err := srv.BatchChangeStock(ctx, inputError)

		require.Error(t, err)
		require.Contains(t, err.Error(), "not enough stock of "+two.Name)
	})

}

func TestService_DeleteProduct(t *testing.T) {

	one := &models.Product{
		ID: 1,
	}

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		mockRepo := mocks.NewMockProductRepository(ctrl)
		srv := NewProductService(mockRepo)

		mockRepo.EXPECT().DeleteProduct(ctx, one.ID).Return(nil)
		err := srv.DeleteProduct(ctx, one.ID)
		require.NoError(t, err)

	})

	t.Run("error negative id", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		mockRepo := mocks.NewMockProductRepository(ctrl)
		srv := NewProductService(mockRepo)

		err := srv.DeleteProduct(ctx, 0)
		require.Error(t, err)
		require.Contains(t, err.Error(), "product ID must be positive")

	})

	t.Run("error id not found", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		mockRepo := mocks.NewMockProductRepository(ctrl)
		srv := NewProductService(mockRepo)

		mockRepo.EXPECT().DeleteProduct(ctx, one.ID).Return(errors.New("no rows"))

		err := srv.DeleteProduct(ctx, one.ID)
		require.Error(t, err)
		require.Contains(t, err.Error(), "product ID not found")

	})

}

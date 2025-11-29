package persistence

import (
	"context"
	"log"
	"testing"

	"github.com/boretsotets/e-commerce-platform/services/order-service/internal/domain/model"
	"github.com/boretsotets/e-commerce-platform/services/order-service/internal/infra/db"
	"github.com/stretchr/testify/require"
)

func TestRepo(t *testing.T) {
	ctx := context.Background()
	db, err := db.NewTestPostgres()
	if err != nil {
		log.Fatalf("error connecting to database: %v\n", err)
	}

	repo := NewOrdertRepo(db)
	items := []*model.OrderItem{
		{
			ProductID: 1,
			Count:     10,
		}, {
			ProductID: 2,
			Count:     15,
		},
	}

	order := &model.Order{
		ClientID:        1,
		Items:           items,
		ShippingAddress: "New York",
	}

	// тест создания заказа
	order, err = repo.RepoCreateOrder(ctx, order.ClientID, order.Items, order.ShippingAddress)
	require.Equal(items, order.Items)
	require.Equal(order.ClientID, int64(1))
	require.Equal(order.ShippingAddress, "New York")
	require.Equal(order.Status, "created")

	// тест получения заказа
	getorder, err := repo.RepoGetOrder(ctx, order.OrderID)
	require.NoError(err)
	require.Equal(order, getorder)

}

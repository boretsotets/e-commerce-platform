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
	require.Equal(t, items, order.Items)
	require.Equal(t, order.ClientID, int64(1))
	require.Equal(t, order.ShippingAddress, "New York")
	require.Equal(t, order.Status, "created")

	// тест получения заказа
	getorder, err := repo.RepoGetOrder(ctx, order.OrderID)
	require.NoError(t, err)
	require.Equal(t, order, getorder)

	// тест изменения заказа
	var orderItemsChanges []*model.OrderItem
	orderItemsChanges = append(orderItemsChanges, &model.OrderItem{
		ProductID: 2,
		Count:     25,
	})

	newOrderItems := []*model.OrderItem{
		{
			ProductID: 1,
			Count:     10,
		}, {
			ProductID: 2,
			Count:     25,
		},
	}

	updatedOrder, err := repo.RepoUpdateOrder(ctx, getorder, orderItemsChanges, "in delivery", "California")
	require.NoError(t, err)
	require.Equal(t, updatedOrder.ClientID, getorder.ClientID)
	require.Equal(t, updatedOrder.OrderID, getorder.OrderID)
	require.Equal(t, updatedOrder.Items, newOrderItems)
	require.Equal(t, updatedOrder.ShippingAddress, "California")
	require.Equal(t, updatedOrder.Status, "in delivery")

	// тест списка заказов
	order2, _ := repo.RepoCreateOrder(ctx, order.ClientID, order.Items, order.ShippingAddress)
	_, _ = repo.RepoCreateOrder(ctx, int64(2), order.Items, order.ShippingAddress)

	orderList, err := repo.RepoListOrders(ctx, int64(1))

	require.NoError(t, err)
	require.Equal(t, len(orderList), 2)
	order.OrderID = 2 // чтобы сработало сравнение ниже
	require.Equal(t, order, orderList[1])
	require.Equal(t, orderList[0], updatedOrder)

	// тест удаления заказа
	err = repo.RepoDeleteOrder(ctx, order2.OrderID)
	require.NoError(t, err)
	_, err = repo.RepoGetOrder(ctx, order2.OrderID)
	require.Error(t, err)
}

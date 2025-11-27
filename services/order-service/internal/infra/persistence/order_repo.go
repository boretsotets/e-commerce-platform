package persistence

import (
	"fmt"

	model "github.com/boretsotets/e-commerce-platform/services/order-service/internal/domain/model"
)

type InmemOrderRepo struct {
	data []*model.Order
}

func NewInmemOrderRepo() *InmemOrderRepo {
	return &InmemOrderRepo{
		data: []*model.Order{
			{ClientID: 1, OrderID: 1, Items: []*model.OrderItem{{ProductID: 1, Count: 2}, {ProductID: 2, Count: 4}}, Status: "created", ShippingAddress: "New York"},
			{ClientID: 2, OrderID: 2, Items: []*model.OrderItem{{ProductID: 1, Count: 3}, {ProductID: 2, Count: 5}}, Status: "shipped", ShippingAddress: "San Francisco"},
		},
	}
}

func (r *InmemOrderRepo) RepoCreateOrder(clientId int64, items []*model.OrderItem, shippingAddress string) (*model.Order, error) {
	newOrder := &model.Order{
		OrderID:         int64(len(r.data)),
		ClientID:        clientId,
		Items:           items,
		Status:          "created",
		ShippingAddress: shippingAddress,
	}
	r.data = append(r.data, newOrder)
	return newOrder, nil
}

func (r *InmemOrderRepo) RepoGetOrder(orderId int64) (*model.Order, error) {
	return r.data[orderId], nil
}
func (r *InmemOrderRepo) RepoUpdateOrder(oldorder *model.Order, items []*model.OrderItem, status string, shippingAddress string) (*model.Order, error) {
	r.data[oldorder.OrderID] = &model.Order{
		OrderID:         oldorder.OrderID,
		ClientID:        oldorder.ClientID,
		Items:           items,
		Status:          status,
		ShippingAddress: shippingAddress,
	}
	return r.data[oldorder.OrderID], nil
}
func (r *InmemOrderRepo) RepoListOrders(clientId int64) ([]*model.Order, error) {
	clientOrders := make([]*model.Order, 0)
	for i := 0; i < len(r.data); i++ {
		if r.data[i].ClientID == clientId {
			clientOrders = append(clientOrders, r.data[i])
		}
	}
	return clientOrders, nil
}
func (r *InmemOrderRepo) RepoDeleteOrder(orderId int64) (bool, error) {
	fmt.Printf("Deleted order %v\n", orderId)
	return true, nil
}

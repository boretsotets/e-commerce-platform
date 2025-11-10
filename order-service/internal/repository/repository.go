package repository

import (
	// "context"
	// "errors"

	"fmt"

	order "github.com/boretsotets/e-commerce-platform/order-service/pkg/api"
)

type InmemOrderRepo struct {
	data []*order.Order
}

type OrderRepository interface {
	RepoCreateOrder(clientId int64, items []*order.OrderItem, shippingAddress string) (*order.Order, error)
	RepoGetOrder(orderId int64) (*order.Order, error)
	RepoUpdateOrder(order *order.Order, items []*order.OrderItem, status string, shippingAddress string) (*order.Order, error)
	RepoListOrders(clientId int64) ([]*order.Order, error)
	RepoDeleteOrder(orderId int64) (bool, error)
}

func NewInmemOrderRepo() *InmemOrderRepo {
	return &InmemOrderRepo{
		data: []*order.Order{
			{ClientId: 1, OrderId: 1, Items: []*order.OrderItem{{ProductId: 1, Count: 2}, {ProductId: 2, Count: 4}}, Status: "created", ShippingAddress: "New York"},
			{ClientId: 2, OrderId: 2, Items: []*order.OrderItem{{ProductId: 1, Count: 3}, {ProductId: 2, Count: 5}}, Status: "shipped", ShippingAddress: "San Francisco"},
		},
	}
}

func (r *InmemOrderRepo) RepoCreateOrder(clientId int64, items []*order.OrderItem, shippingAddress string) (*order.Order, error) {
	newOrder := &order.Order{
		OrderId:         int64(len(r.data)),
		ClientId:        clientId,
		Items:           items,
		Status:          "created",
		ShippingAddress: shippingAddress,
	}
	r.data = append(r.data, newOrder)
	return newOrder, nil
}

func (r *InmemOrderRepo) RepoGetOrder(orderId int64) (*order.Order, error) {
	return r.data[orderId], nil
}
func (r *InmemOrderRepo) RepoUpdateOrder(oldorder *order.Order, items []*order.OrderItem, status string, shippingAddress string) (*order.Order, error) {
	r.data[oldorder.OrderId] = &order.Order{
		OrderId:         oldorder.OrderId,
		ClientId:        oldorder.ClientId,
		Items:           items,
		Status:          status,
		ShippingAddress: shippingAddress,
	}
	return r.data[oldorder.OrderId], nil
}
func (r *InmemOrderRepo) RepoListOrders(clientId int64) ([]*order.Order, error) {
	clientOrders := make([]*order.Order, 0)
	for i := 0; i < len(r.data); i++ {
		if r.data[i].ClientId == clientId {
			clientOrders = append(clientOrders, r.data[i])
		}
	}
	return clientOrders, nil
}
func (r *InmemOrderRepo) RepoDeleteOrder(orderId int64) (bool, error) {
	fmt.Printf("Deleted order %v\n", orderId)
	return true, nil
}

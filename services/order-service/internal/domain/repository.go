package domain

import (
	"github.com/boretsotets/e-commerce-platform/services/order-service/internal/domain/model"
)

type OrderRepository interface {
	RepoCreateOrder(clientId int64, items []*model.OrderItem, shippingAddress string) (*model.Order, error)
	RepoGetOrder(orderId int64) (*model.Order, error)
	RepoUpdateOrder(order *model.Order, items []*model.OrderItem, status string, shippingAddress string) (*model.Order, error)
	RepoListOrders(clientId int64) ([]*model.Order, error)
	RepoDeleteOrder(orderId int64) (bool, error)
}

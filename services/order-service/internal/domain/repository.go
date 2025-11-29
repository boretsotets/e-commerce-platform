package domain

import (
	"context"

	"github.com/boretsotets/e-commerce-platform/services/order-service/internal/domain/model"
)

type OrderRepository interface {
	RepoCreateOrder(ctx context.Context, clientId int64, items []*model.OrderItem, shippingAddress string) (*model.Order, error)
	RepoGetOrder(ctx context.Context, id int64) (*model.Order, error)
	RepoUpdateOrder(ctx context.Context, oldorder *model.Order, items []*model.OrderItem, status string, shippingAddress string) (*model.Order, error)
	RepoListOrders(ctx context.Context, clientId int64) ([]*model.Order, error)
	RepoDeleteOrder(ctx context.Context, orderId int64) error
}

package service

import (
	"context"

	"github.com/boretsotets/e-commerce-platform/services/order-service/internal/domain"
	model "github.com/boretsotets/e-commerce-platform/services/order-service/internal/domain/model"
)

type ProductService interface {
	UpdateStock(ctx context.Context, items []*model.OrderItem) error
}

type OrderService struct {
	Repo          domain.OrderRepository
	ProductClient ProductService
}

func NewOrderServiceServer(repo domain.OrderRepository, productClient ProductService) *OrderService {
	return &OrderService{
		Repo:          repo,
		ProductClient: productClient,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, ClientID int64, Items []*model.OrderItem, ShippingAddress string) (*model.Order, error) {

	err := s.ProductClient.UpdateStock(ctx, Items)
	if err != nil {
		return nil, err
	}
	response, err := s.Repo.RepoCreateOrder(ClientID, Items, ShippingAddress)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *OrderService) UpdateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	oldOrder, err := s.GetOrder(ctx, order.OrderID)
	if err != nil {
		return nil, err
	}
	deltas := CountDeltas(oldOrder.Items, order.Items)

	err = s.ProductClient.UpdateStock(ctx, deltas)
	if err != nil {
		return nil, err
	}
	response, err := s.Repo.RepoUpdateOrder(order, order.Items, order.Status, order.ShippingAddress)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *OrderService) GetOrder(ctx context.Context, OrderID int64) (*model.Order, error) {
	response, err := s.Repo.RepoGetOrder(OrderID)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *OrderService) ListOrders(ctx context.Context, ClientID int64) ([]*model.Order, error) {
	response, err := s.Repo.RepoListOrders(ClientID)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, orderID int64) (bool, error) {
	// need to update client stock
	oldOrder, err := s.GetOrder(ctx, orderID)
	if err != nil {
		return false, err
	}
	deltas := deltasForDelete(oldOrder.Items)
	err = s.ProductClient.UpdateStock(ctx, deltas)
	if err != nil {
		return false, err
	}

	response, err := s.Repo.RepoDeleteOrder(orderID)
	if err != nil {
		return false, err
	}
	return response, nil
}

func CountDeltas(oldItems []*model.OrderItem, items []*model.OrderItem) []*model.OrderItem {
	res := make([]*model.OrderItem, len(items))
	for i, item := range items {
		res[i] = &model.OrderItem{
			ProductID: item.ProductID,
			Count:     items[item.ProductID].Count - oldItems[item.ProductID].Count,
		}
	}
	return res
}

func deltasForDelete(oldItems []*model.OrderItem) []*model.OrderItem {
	res := make([]*model.OrderItem, len(oldItems))
	for i := range oldItems {
		res[i] = &model.OrderItem{
			ProductID: oldItems[i].ProductID,
			Count:     int32(0),
		}
	}
	return res
}

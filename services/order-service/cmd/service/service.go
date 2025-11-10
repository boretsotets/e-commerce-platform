package service

import (
	"context"

	client "github.com/boretsotets/e-commerce-platform/order-service/cmd"
	"github.com/boretsotets/e-commerce-platform/order-service/internal/repository"
	pb "github.com/boretsotets/e-commerce-platform/order-service/pkg/api"
)

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	Repo          repository.OrderRepository
	ProductClient *client.ProductClient
}

func NewOrderServiceServer(repo repository.OrderRepository, productClient *client.ProductClient) *OrderServer {
	return &OrderServer{
		Repo:          repo,
		ProductClient: productClient,
	}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	err := s.ProductClient.UpdateStock(ctx, req.Items)
	if err != nil {
		return nil, err
	}
	response, err := s.Repo.RepoCreateOrder(req.ClientId, req.Items, req.ShippingAddress)
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{
		Order: response,
	}, nil
}

func (s *OrderServer) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	oldOrder, err := s.GetOrder(ctx, &pb.GetOrderRequest{
		OrderId: req.Order.OrderId,
	})
	if err != nil {
		return nil, err
	}
	deltas := CountDeltas(oldOrder.Order.Items, req.Items)

	err = s.ProductClient.UpdateStock(ctx, deltas)
	if err != nil {
		return nil, err
	}
	response, err := s.Repo.RepoUpdateOrder(req.Order, req.Items, req.Status, req.ShippingAddress)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateOrderResponse{
		Order: response,
	}, nil
}

func (s *OrderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	response, err := s.Repo.RepoGetOrder(req.OrderId)
	if err != nil {
		return nil, err
	}
	return &pb.GetOrderResponse{
		Order: response,
	}, nil
}

func (s *OrderServer) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	response, err := s.Repo.RepoListOrders(req.ClientId)
	if err != nil {
		return nil, err
	}
	return &pb.ListOrdersResponse{
		OrdersList: response,
	}, nil
}

func (s *OrderServer) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	// need to update client stock
	oldOrder, err := s.GetOrder(ctx, &pb.GetOrderRequest{
		OrderId: req.OrderId,
	})
	if err != nil {
		return nil, err
	}
	deltas := deltasForDelete(oldOrder.Order.Items)
	err = s.ProductClient.UpdateStock(ctx, deltas)
	if err != nil {
		return nil, err
	}

	response, err := s.Repo.RepoDeleteOrder(req.OrderId)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteOrderResponse{
		Success: response,
	}, nil
}

func CountDeltas(oldItems []*pb.OrderItem, items []*pb.OrderItem) []*pb.OrderItem {
	res := make([]*pb.OrderItem, len(items))
	for i, item := range items {
		res[i] = &pb.OrderItem{
			ProductId: item.ProductId,
			Count:     items[item.ProductId].Count - oldItems[item.ProductId].Count,
		}
	}
	return res
}

func deltasForDelete(oldItems []*pb.OrderItem) []*pb.OrderItem {
	res := make([]*pb.OrderItem, len(oldItems))
	for i := range oldItems {
		res[i] = &pb.OrderItem{
			ProductId: oldItems[i].ProductId,
			Count:     int32(0),
		}
	}
	return res
}

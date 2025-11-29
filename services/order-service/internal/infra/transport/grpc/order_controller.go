package grpc

import (
	usecase "github.com/boretsotets/e-commerce-platform/services/order-service/internal/application/usecase"
	domain "github.com/boretsotets/e-commerce-platform/services/order-service/internal/domain/model"
	pb "github.com/boretsotets/e-commerce-platform/services/order-service/pkg/api"

	"context"
)

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	Service usecase.OrderService
}

func (c *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	clientID := req.ClientId
	shippingAddress := req.ShippingAddress
	DomainItems := OrderItemsToDomainItems(req.Items)
	order, err := c.Service.CreateOrder(ctx, clientID, DomainItems, shippingAddress)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Order: DomainOrderToOrder(order),
	}, nil
}

func OrderItemsToDomainItems(orderitems []*pb.OrderItem) []*domain.OrderItem {
	domainItems := make([]*domain.OrderItem, len(orderitems))
	for i := range domainItems {
		domainItems[i].ProductID = orderitems[i].ProductId
		domainItems[i].Count = orderitems[i].Count
	}
	return domainItems
}

func DomainItemsToOrderItems(domainItems []*domain.OrderItem) []*pb.OrderItem {
	orderItems := make([]*pb.OrderItem, len(domainItems))
	for i := range orderItems {
		orderItems[i].ProductId = domainItems[i].ProductID
		orderItems[i].Count = domainItems[i].Count
	}
	return orderItems
}

func OrderToDomainOrder(pborder *pb.Order) *domain.Order {
	return &domain.Order{
		ClientID:        pborder.ClientId,
		OrderID:         pborder.OrderId,
		Items:           OrderItemsToDomainItems(pborder.Items),
		Status:          pborder.Status,
		ShippingAddress: pborder.ShippingAddress,
	}
}

func DomainOrderToOrder(domainorder *domain.Order) *pb.Order {
	return &pb.Order{
		ClientId:        domainorder.ClientID,
		OrderId:         domainorder.OrderID,
		Items:           DomainItemsToOrderItems(domainorder.Items),
		Status:          domainorder.Status,
		ShippingAddress: domainorder.ShippingAddress,
	}
}

func DomainOrderListToOrderList(domainOrderList []*domain.Order) []*pb.Order {
	pbOrderList := make([]*pb.Order, len(domainOrderList))
	for i := range domainOrderList {
		pbOrderList[i] = DomainOrderToOrder(domainOrderList[i])
	}
	return pbOrderList
}

func (c *OrderServer) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	order, err := c.Service.UpdateOrder(ctx, OrderToDomainOrder(req.Order))
	if err != nil {
		return nil, err
	}
	return &pb.UpdateOrderResponse{
		Order: DomainOrderToOrder(order),
	}, nil
}

func (c *OrderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	domainOrder, err := c.Service.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}
	return &pb.GetOrderResponse{
		Order: DomainOrderToOrder(domainOrder),
	}, nil
}

func (c *OrderServer) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	domainOrderList, err := c.Service.ListOrders(ctx, req.ClientId)
	if err != nil {
		return nil, err
	}
	return &pb.ListOrdersResponse{
		OrdersList: DomainOrderListToOrderList(domainOrderList),
	}, nil
}

func (c *OrderServer) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	err := c.Service.DeleteOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteOrderResponse{
		Success: true,
	}, nil
}

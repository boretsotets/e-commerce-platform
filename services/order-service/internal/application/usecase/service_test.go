package service

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	pb "github.com/boretsotets/e-commerce-platform/order-service/pkg/api"
// )

// type mockProductClient struct {
// 	stock map[int64]int32
// }

// func (m *mockProductClient) UpdateStock(ctx context.Context, items []*pb.OrderItem) error {
// 	for _, i := range items {
// 		delta := m.stock[i.ProductId] - i.Count
// 		if delta < 0 {
// 			return errors.New("Not enough stock left")
// 		}
// 		m.stock[i.ProductId] = delta
// 	}

// 	return nil
// }

// type mockOrderRepo struct {
// 	createdOrders []*pb.Order
// }

// func (m *mockOrderRepo) RepoCreateOrder(clientId int64, items []*pb.OrderItem, shippingAddress string) (*pb.Order, error) {
// 	var newOrder = &pb.Order{
// 		ClientId:        clientId,
// 		Items:           items,
// 		Status:          "created",
// 		ShippingAddress: shippingAddress,
// 	}
// 	m.createdOrders = append(m.createdOrders, newOrder)
// 	return newOrder, nil
// }

// func (m *mockOrderRepo) RepoGetOrder(orderId int64) (*pb.Order, error) {
// 	return nil, nil
// }
// func (m *mockOrderRepo) RepoUpdateOrder(order *pb.Order, items []*pb.OrderItem, status string, shippingAddress string) (*pb.Order, error) {
// 	return nil, nil
// }
// func (m *mockOrderRepo) RepoListOrders(clientId int64) ([]*pb.Order, error) {
// 	return nil, nil
// }
// func (m *mockOrderRepo) RepoDeleteOrder(orderId int64) (bool, error) {
// 	return true, nil
// }

// func TestCreateOrder(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		stock     map[int64]int32
// 		items     []*pb.OrderItem
// 		wantError bool
// 	}{
// 		{
// 			name: "Успешное создание заказа",
// 			stock: map[int64]int32{
// 				1: 10,
// 				2: 5,
// 			},
// 			items: []*pb.OrderItem{
// 				{ProductId: 1, Count: 2},
// 				{ProductId: 2, Count: 1},
// 			},
// 			wantError: false,
// 		},
// 		{
// 			name: "Ошибка: недостаточно товара",
// 			stock: map[int64]int32{
// 				1: 10,
// 				2: 5,
// 			},
// 			items: []*pb.OrderItem{
// 				{ProductId: 1, Count: 10},
// 				{ProductId: 2, Count: 10},
// 			},
// 			wantError: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockClient := mockProductClient{stock: tt.stock}
// 			mockRepo := mockOrderRepo{}

// 			svc := OrderServer{
// 				Repo:          &mockRepo,
// 				ProductClient: &mockClient,
// 			}

// 			resp, err := svc.CreateOrder()
// 		})

// 	}
// }

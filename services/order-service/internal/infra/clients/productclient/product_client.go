package client

import (
	"context"
	"log"

	model "github.com/boretsotets/e-commerce-platform/order-service/internal/domain/model"
	productpb "github.com/boretsotets/e-commerce-platform/product-service/pkg/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductClient struct {
	conn   *grpc.ClientConn
	client productpb.ProductServiceClient
}

func NewProductClient(addr string) *ProductClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error connecting to product-service: %v", err)
	}
	c := productpb.NewProductServiceClient(conn)
	return &ProductClient{conn: conn, client: c}
}

func (p *ProductClient) UpdateStock(ctx context.Context, items []*model.OrderItem) error {
	stockChangeItems := ToStockChange(items)
	_, err := p.client.BatchChangeStock(ctx, &productpb.BatchChangeStockRequest{
		Items: stockChangeItems,
	})
	return err
}

func ToStockChange(orderItems []*model.OrderItem) []*productpb.StockChangeItem {
	res := make([]*productpb.StockChangeItem, len(orderItems))
	for i, item := range orderItems {
		res[i] = &productpb.StockChangeItem{
			ProductId: item.ProductID,
			Delta:     item.Count,
		}
	}
	return res
}

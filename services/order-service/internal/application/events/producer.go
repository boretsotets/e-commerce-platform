package events

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/boretsotets/e-commerce-platform/services/order-service/internal/domain/model"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(broker string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    "order.created",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) PublishOrderCreated(ctx context.Context, order *model.Order) error {
	event := EventOrderCreated{
		OrderID:   int64(order.OrderID),
		CreatedAt: time.Now(),
	}

	for _, item := range order.Items {
		event.Items = append(event.Items, struct {
			ProductID int64 "json:\"product_id\""
			Qty       int32 "json:\"qty\""
		}{
			ProductID: item.ProductID,
			Qty:       item.Count,
		})
	}
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(strconv.Itoa(order.OrderID)),
		Value: data,
	})
}

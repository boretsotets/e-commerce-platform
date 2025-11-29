package events

import (
	"context"
	"encoding/json"
	"log"

	"github.com/boretsotets/e-commerce-platform/product-service/internal/application/usecase"
	"github.com/boretsotets/e-commerce-platform/product-service/internal/domain/models"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	c       *kafka.Consumer
	service *usecase.ProductService
	topics  []string
}

func NewConsumer(brokers, group string, topics []string, svc *usecase.ProductService) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.services": brokers,
		"group.id":           group,
		"auto.offset.reset":  "earliest",
	})
	if err != nil {
		return nil, err
	}

	if err := c.SubscribeTopics(topics, nil); err != nil {
		return nil, err
	}

	return &Consumer{c: c, service: svc, topics: topics}, nil
}

func (cons *Consumer) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			_ = cons.c.Close()
			return
		default:
			ev := cons.c.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				var evt models.EventOrderCreated
				if err := json.Unmarshal(e.Value, &evt); err != nil {
					log.Printf("invalid event json: %v", err)
					// dlq
					continue
				}
				if err := cons.service.ApplyOrderEvent(context.Background(), evt); err != nil {
					log.Printf("apply order event failed: %v", err)
					// retry/dlq
				}
			case kafka.Error:
				log.Printf("kafka error: %v", e)
				// handle broker errors
			default:
				// другие ивенты
			}
		}
	}
}

func (cons *Consumer) Close() {
	_ = cons.c.Close()
}

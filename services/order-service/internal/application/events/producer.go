package events

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	p     *kafka.Producer
	topic string
}

func NewProducer(brokers, topic string) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
	})
	if err != nil {
		return nil, err
	}
	return &Producer{p: p, topic: topic}, nil
}

func (pr *Producer) PublishOrderCreated(ctx context.Context, evt EventOrderCreated) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &pr.topic, Partition: kafka.PartitionAny},
		Key:            []byte(strconv.Itoa(int(evt.OrderID))),
		Value:          data,
	}

	deliveryChan := make(chan kafka.Event, 1)
	err = pr.p.Produce(msg, deliveryChan)
	if err != nil {
		return err
	}

	select {
	case ev := <-deliveryChan:
		m := ev.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			return m.TopicPartition.Error
		}
	case <-time.After(5 * time.Second):
		return fmt.Errorf("kafka delivery timeout")
	}

	close(deliveryChan)
	return nil
}

func (pr *Producer) Close() {
	pr.p.Flush(10000)
	pr.p.Close()
}

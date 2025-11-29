package events

import (
	"time"

	"github.com/boretsotets/e-commerce-platform/services/order-service/internal/domain/model"
)

type EventOrderCreated struct {
	OrderID int64 `json:"order_id"`
	Items   []struct {
		ProductID int64 `json:"product_id"`
		Qty       int32 `json:"qty"`
	} `json:"items"`
	CreatedAt time.Time
}

func CreateEvent(order model.Order) EventOrderCreated {
	event := EventOrderCreated{
		OrderID:   int64(order.OrderID),
		CreatedAt: time.Now(),
	}

	for _, item := range order.Items {
		event.Items = append(event.Items, struct {
			ProductID int64 `json:"product_id"`
			Qty       int32 `json:"qty"`
		}{
			ProductID: item.ProductID,
			Qty:       item.Count,
		})
	}
	return event
}

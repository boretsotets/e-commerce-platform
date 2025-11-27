package events

import "time"

type EventOrderCreated struct {
	OrderID int64 `json:"order_id"`
	Items   []struct {
		ProductID int64 `json:"product_id"`
		Qty       int32 `json:"qty"`
	} `json:"items"`
	CreatedAt time.Time
}

package models

import "time"

type Product struct {
	ID    int64 `gorm:"primaryKey"`
	Name  string
	Price float64
	Stock int32
}

type ProductList struct {
	List []*Product
}

type StockChangeItem struct {
	ProductID int64
	Delta     int32
}

type EventOrderCreated struct {
	OrderID int64 `json:"order_id"`
	Items   []struct {
		ProductID int64 `json:"product_id"`
		Qty       int32 `json:"qty"`
	} `json:"items"`
	CreatedAt time.Time
}

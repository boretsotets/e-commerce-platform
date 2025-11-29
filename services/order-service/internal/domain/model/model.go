package model

type Order struct {
	OrderID         int64 `gorm:"primaryKey"`
	ClientID        int64
	Items           []*OrderItem `gorm:"foreignKey:OrderID"`
	Status          string
	ShippingAddress string
}

type OrderItem struct {
	ProductID int64
	Count     int32
}

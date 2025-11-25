package domain

type Order struct {
	ClientID        int64
	OrderID         int64
	Items           []*OrderItem
	Status          string
	ShippingAddress string
}

type OrderItem struct {
	ProductID int64
	Count     int32
}

package models

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

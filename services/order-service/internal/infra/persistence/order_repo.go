package persistence

import (
	"context"

	model "github.com/boretsotets/e-commerce-platform/services/order-service/internal/domain/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepo struct {
	Repo *gorm.DB
}

func NewOrdertRepo(repo *gorm.DB) *OrderRepo {
	return &OrderRepo{Repo: repo}
}

type OrderItemDB struct {
	ID        int64 `gorm:"primaryKey"`
	OrderID   int64
	ProductID int64
	Count     int32
}

func (r *OrderRepo) RepoCreateOrder(ctx context.Context, clientId int64, items []*model.OrderItem, shippingAddress string) (*model.Order, error) {
	orderDB := &model.Order{
		ClientID:        clientId,
		ShippingAddress: shippingAddress,
		Status:          "created",
	}

	tx := r.Repo.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.WithContext(ctx).Create(&orderDB).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var itemsDB []OrderItemDB
	for _, item := range items {
		itemsDB = append(itemsDB, OrderItemDB{
			OrderID:   orderDB.OrderID,
			ProductID: item.ProductID,
			Count:     item.Count,
		})
	}

	if err := tx.WithContext(ctx).Create(&itemsDB).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return orderDB, nil
}

func (r *OrderRepo) RepoGetOrder(ctx context.Context, id int64) (*model.Order, error) {
	var order model.Order
	err := r.Repo.WithContext(ctx).First(&order, id).Error
	if err != nil {
		return nil, err
	}

	var orderItemsDB []*OrderItemDB
	err = r.Repo.WithContext(ctx).Where("order_id = ?", id).Find(&orderItemsDB).Error
	if err != nil {
		return nil, err
	}

	var orderItems []*model.OrderItem
	for _, v := range orderItemsDB {
		orderItems = append(orderItems, &model.OrderItem{
			ProductID: v.ProductID,
			Count:     v.Count,
		})
	}

	order.Items = orderItems

	return &order, nil
}

func (r *OrderRepo) RepoUpdateOrder(ctx context.Context, oldorder *model.Order, items []*model.OrderItem, status string, shippingAddress string) (*model.Order, error) {
	var order model.Order

	tx := r.Repo.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.WithContext(ctx).
		Where("ID = ?", oldorder.OrderID).
		Model(&order).
		Clauses(clause.Returning{}).
		Update("status", status).
		Update("shippingAddress", shippingAddress).
		Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var orderItemsDB *OrderItemDB
	for _, v := range items {
		if err := tx.WithContext(ctx).
			Where("order_id = ?", oldorder.OrderID).
			Where("product_id = ?", v.ProductID).
			Model(&orderItemsDB).
			Clauses(clause.Returning{}).
			Update("count", v.Count).
			Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		order.Items = append(order.Items, v)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepo) RepoListOrders(ctx context.Context, clientId int64) ([]*model.Order, error) {
	var orders []*model.Order
	if err := r.Repo.WithContext(ctx).Where("ID = ?", clientId).Find(&orders).Error; err != nil {
		return nil, err
	}

	var orderItemDB OrderItemDB
	for i, v := range orders {
		if err := r.Repo.WithContext(ctx).Where("order_id = ?", v.OrderID).Find(&orderItemDB).Error; err != nil {
			return nil, err
		}
		orders[i].Items = append(orders[i].Items, &model.OrderItem{
			ProductID: orderItemDB.ProductID,
			Count:     orderItemDB.Count,
		})
	}

	return orders, nil
}

func (r *OrderRepo) RepoDeleteOrder(ctx context.Context, orderId int64) error {
	var order model.Order
	err := r.Repo.WithContext(ctx).Find(&order, orderId).Error
	if err != nil {
		return err
	}

	tx := r.Repo.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.WithContext(ctx).Delete(&order).Error; err != nil {
		tx.Rollback()
		return err
	}

	var orderItem *OrderItemDB
	if err := tx.WithContext(ctx).Where("order_id = ?", orderId).Delete(&orderItem).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

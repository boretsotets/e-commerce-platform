package persistence

import (
	"context"
	"fmt"
	"time"

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

type OrderDB struct {
	OrderID         int64 `gorm:"primaryKey"`
	ClientID        int64
	Status          string
	ShippingAddress string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (r *OrderRepo) RepoCreateOrder(ctx context.Context, clientId int64, items []*model.OrderItem, shippingAddress string) (*model.Order, error) {
	orderDB := &OrderDB{
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

	return &model.Order{
		OrderID:         orderDB.OrderID,
		ClientID:        orderDB.ClientID,
		Items:           items,
		Status:          orderDB.Status,
		ShippingAddress: orderDB.ShippingAddress,
	}, nil
}

func (r *OrderRepo) RepoGetOrder(ctx context.Context, id int64) (*model.Order, error) {
	var order OrderDB
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

	return &model.Order{
		OrderID:         order.OrderID,
		ClientID:        order.ClientID,
		Items:           orderItems,
		Status:          order.Status,
		ShippingAddress: order.ShippingAddress,
	}, nil
}

func (r *OrderRepo) RepoUpdateOrder(ctx context.Context, oldorder *model.Order, items []*model.OrderItem, status string, shippingAddress string) (*model.Order, error) {
	tx := r.Repo.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var orderDB OrderDB
	if err := tx.WithContext(ctx).
		Where("order_id = ?", oldorder.OrderID).
		Model(&orderDB).
		Clauses(clause.Returning{}).
		Update("status", status).
		Update("shippingAddress", shippingAddress).
		Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var orderItemDB *OrderItemDB
	caseStmt := "CASE product_id "
	ids := []int{}

	for _, item := range items {
		caseStmt += fmt.Sprintf("WHEN %d THEN %d ", item.ProductID, item.Count)
		ids = append(ids, int(item.ProductID))
	}
	caseStmt += "ELSE count END"

	if err := tx.WithContext(ctx).
		Model(&orderItemDB).
		Where("order_id = ?", oldorder.OrderID).
		Where("product_id IN ?", ids).
		Update("count", gorm.Expr(caseStmt)).
		Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	var orderItemsDB []*OrderItemDB
	if err := r.Repo.WithContext(ctx).Where("order_id = ?", oldorder.OrderID).Find(&orderItemsDB).Error; err != nil {
		return nil, err
	}

	order := &model.Order{
		OrderID:         orderDB.OrderID,
		ClientID:        orderDB.ClientID,
		Status:          orderDB.Status,
		ShippingAddress: orderDB.ShippingAddress,
	}
	for _, v := range orderItemsDB {
		order.Items = append(order.Items, &model.OrderItem{
			ProductID: v.ProductID,
			Count:     v.Count,
		})
	}

	return order, nil
}

func (r *OrderRepo) RepoListOrders(ctx context.Context, clientId int64) ([]*model.Order, error) {
	var ordersDB []*OrderDB
	if err := r.Repo.WithContext(ctx).Where("client_id = ?", clientId).Find(&ordersDB).Error; err != nil {
		return nil, err
	}
	var orders []*model.Order

	var orderItemsDB []*OrderItemDB
	for i, v := range ordersDB {
		if err := r.Repo.WithContext(ctx).Where("order_id = ?", v.OrderID).Find(&orderItemsDB).Error; err != nil {
			return nil, err
		}
		orders = append(orders, &model.Order{
			OrderID:         v.OrderID,
			ClientID:        v.ClientID,
			ShippingAddress: v.ShippingAddress,
			Status:          v.Status,
		})
		for _, val := range orderItemsDB {
			orders[i].Items = append(orders[i].Items, &model.OrderItem{
				ProductID: val.ProductID,
				Count:     val.Count,
			})

		}
	}

	return orders, nil
}

func (r *OrderRepo) RepoDeleteOrder(ctx context.Context, orderId int64) error {
	var order OrderDB
	err := r.Repo.WithContext(ctx).Find(&order, orderId).Error
	if err != nil {
		return err
	}

	tx := r.Repo.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var orderItem *OrderItemDB
	if err := tx.WithContext(ctx).Where("order_id = ?", orderId).Delete(&orderItem).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.WithContext(ctx).Delete(&order).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

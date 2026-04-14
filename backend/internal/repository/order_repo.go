package repository

import (
	github.com/ethanqian1990/wechat-mall-backend/internal/model
	github.com/gin-gonic/gin
	gorm  gorm.io/gorm
)

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (r *OrderRepository) FindByUserID(userID, status string, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := DB.Model(&model.Order{}).Where('user_id = ?', userID)
	if status != '' {
		query = query.Where('status = ?', status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order('created_at DESC').Offset(offset).Limit(pageSize).Find(&orders).Error

	return orders, total, err
}

func (r *OrderRepository) FindByIDAndUserID(id, userID string) (*model.Order, error) {
	var order model.Order
	err := DB.First(&order, 'id = ? AND user_id = ?', id, userID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &order, err
}

func (r *OrderRepository) FindByToken(token string) (*model.Order, error) {
	var order model.Order
	err := DB.First(&order, 'client_order_token = ?', token).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &order, err
}

func (r *OrderRepository) Create(order *model.Order) error {
	return DB.Create(order).Error
}

func (r *OrderRepository) UpdateStatus(id, status string) error {
	updates := gin.H{'status': status}
	switch status {
	case 'PENDING_SHIP':
		updates['pay_at'] = gorm.Expr('NOW()')
	case 'SHIPPED':
		updates['ship_at'] = gorm.Expr('NOW()')
	case 'FINISHED':
		updates['receive_at'] = gorm.Expr('NOW()')
	case 'CANCELED':
		updates['cancel_at'] = gorm.Expr('NOW()')
	}
	return DB.Model(&model.Order{}).Where('id = ?', id).Updates(updates).Error
}

func (r *OrderRepository) Update(id string, updates map[string]interface{}) error {
	return DB.Model(&model.Order{}).Where('id = ?', id).Updates(updates).Error
}

// OrderItemRepository
type OrderItemRepository struct{}

func NewOrderItemRepository() *OrderItemRepository {
	return &OrderItemRepository{}
}

func (r *OrderItemRepository) FindByOrderID(orderID string) ([]model.OrderItem, error) {
	var items []model.OrderItem
	err := DB.Where('order_id = ?', orderID).Find(&items).Error
	return items, err
}

func (r *OrderItemRepository) BatchCreate(items []model.OrderItem) error {
	return DB.Create(&items).Error
}
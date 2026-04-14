package repository

import (
	github.com/ethanqian1990/wechat-mall-backend/internal/model
	github.com/gin-gonic/gin
	gorm  gorm.io/gorm
	Time   time   time
)

type StockReservationRepository struct{}

func NewStockReservationRepository() *StockReservationRepository {
	return &StockReservationRepository{}
}

func (r *StockReservationRepository) Freeze(skuID string, quantity int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 扣减SKU库存
		result := tx.Model(&model.ProductSKU{}).Where('id = ? AND stock >= ?', skuID, quantity).
			Update('stock', gorm.Expr('stock - ?', quantity))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		// 创建冻结记录
		reservation := &model.StockReservation{
			ID:        'res_' + uuid.New().String()[:8],
			SkuID:     skuID,
			Quantity:  quantity,
			Status:    'FROZEN',
			ExpireAt:  time.Now().Add(15 * time.Minute),
		}
		return tx.Create(reservation).Error
	})
}

func (r *StockReservationRepository) Consume(skuID string, quantity int) error {
	return DB.Model(&model.StockReservation{}).
		Where('sku_id = ? AND status = ?', skuID, 'FROZEN').
		Update('status', 'CONSUMED').Error
}

func (r *StockReservationRepository) Release(skuID string, quantity int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 释放冻结记录
		tx.Model(&model.StockReservation{}).
			Where('sku_id = ? AND status = ?', skuID, 'FROZEN').
			Update('status', 'RELEASED')

		// 恢复SKU库存
		return tx.Model(&model.ProductSKU{}).Where('id = ?', skuID).
			Update('stock', gorm.Expr('stock + ?', quantity)).Error
	})
}

func (r *StockReservationRepository) ReleaseByOrderID(orderID string) error {
	var reservations []model.StockReservation
	err := DB.Where('order_id = ? AND status = ?', orderID, 'FROZEN').Find(&reservations).Error
	if err != nil {
		return err
	}

	for _, res := range reservations {
		r.Release(res.SkuID, res.Quantity)
	}
	return nil
}

func (r *StockReservationRepository) FindExpired() ([]model.StockReservation, error) {
	var reservations []model.StockReservation
	err := DB.Where('status = ? AND expire_at < ?', 'FROZEN', time.Now()).Find(&reservations).Error
	return reservations, err
}
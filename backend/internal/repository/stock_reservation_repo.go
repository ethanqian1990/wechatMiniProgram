package repository

import (
	github.com/ethanqian1990/wechat-mall-backend/internal/model
)

type StockReservationRepository struct{}

func NewStockReservationRepository() *StockReservationRepository {
	return &StockReservationRepository{}
}

func (r *StockReservationRepository) Freeze(skuID string, quantity int) error {
	// TODO: 冻结库存
	return nil
}

func (r *StockReservationRepository) Consume(skuID string, quantity int) error {
	// TODO: 扣减冻结库存
	return nil
}

func (r *StockReservationRepository) Release(skuID string, quantity int) error {
	// TODO: 释放冻结库存
	return nil
}

func (r *StockReservationRepository) ReleaseByOrderID(orderID string) error {
	// TODO: 释放订单的所有冻结库存
	return nil
}

func (r *StockReservationRepository) FindExpired() ([]model.StockReservation, error) {
	// TODO: 查找过期的冻结记录
	return nil, nil
}
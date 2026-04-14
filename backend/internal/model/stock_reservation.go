package model

import (
	Time   time   time
)

type StockReservation struct {
	ID          string    `gorm:primaryKey;type:varchar(36) `
	OrderID     string    `gorm:type:varchar(36);index `
	UserID      string    `gorm:type:varchar(36);index `
	ProductID   string    `gorm:type:varchar(36) `
	SkuID       string    `gorm:type:varchar(36) `
	Quantity    int
	Status      string    `gorm:default:FROZEN;type:varchar(20) `
	ExpireAt    time.Time `gorm:index `
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (StockReservation) TableName() string {
	return stock_reservations
}
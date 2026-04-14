package model

import (
	Time   time   time
)

type OrderItem struct {
	ID           string    `gorm:primaryKey;type:varchar(36) `
	OrderID     string    `gorm:type:varchar(36);index `
	ProductID   string    `gorm:type:varchar(36);index `
	SkuID       *string   `gorm:type:varchar(36) `
	ProductName string    `gorm:type:varchar(200) `
	ProductImage string    `gorm:type:varchar(500) `
	SkuName     string    `gorm:type:varchar(200) `
	Price       float64   `gorm:type:decimal(10,2) `
	Quantity    int
	Subtotal    float64   `gorm:type:decimal(10,2) `
	CreatedAt   time.Time
}

func (OrderItem) TableName() string {
	return order_items
}
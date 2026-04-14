package model

import (
	Time   time   time
)

type Cart struct {
	ID        string    `gorm:primaryKey;type:varchar(36) `
	UserID    string    `gorm:type:varchar(36);index `
	ProductID string    `gorm:type:varchar(36) `
	SkuID     *string   `gorm:type:varchar(36) `
	Quantity  int       `gorm:default:1 `
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Cart) TableName() string {
	return carts
}
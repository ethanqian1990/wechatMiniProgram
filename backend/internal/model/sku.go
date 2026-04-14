package model

import (
	Time   time   time
)

type ProductSKU struct {
	ID            string    `gorm:primaryKey;type:varchar(36) `
	ProductID     string    `gorm:type:varchar(36);index `
	SkuName       string    `gorm:type:varchar(200) `
	Specs         string    `gorm:type:json `
	Price         float64   `gorm:type:decimal(10,2) `
	OriginalPrice *float64  `gorm:type:decimal(10,2) `
	Stock         int       `gorm:default:0 `
	Image         string    `gorm:type:varchar(500) `
	IsEnabled     int       `gorm:default:1 `
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (ProductSKU) TableName() string {
	return product_skus
}
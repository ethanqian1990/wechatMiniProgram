package model

import (
	Time   time   time
)

type Product struct {
	ID               string     `gorm:primaryKey;type:varchar(36) `
	Name             string     `gorm:type:varchar(200) `
	Description      string     `gorm:type:text `
	MainImage        string     `gorm:type:varchar(500) `
	Images           string     `gorm:type:json `
	CategoryID       string     `gorm:type:varchar(36) `
	Price            float64    `gorm:type:decimal(10,2);default:0 `
	OriginalPrice    *float64   `gorm:type:decimal(10,2) `
	Stock            int        `gorm:default:0 `
	SalesCount       int        `gorm:default:0 `
	IsOnShelf        int        `gorm:default:1;column:is_on_shelf `
	ShowOnHome       int        `gorm:default:1;column:show_on_home `
	AvailableRegions string     `gorm:type:json;column:available_regions `
	Tags             string     `gorm:type:json `
	PromoText        string     `gorm:type:varchar(200) `
	LowStockThreshold int       `gorm:default:10;column:low_stock_threshold `
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time `gorm:index `
}

func (Product) TableName() string {
	return products
}
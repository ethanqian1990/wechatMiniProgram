package model

import (
	Time   time   time
)

type Category struct {
	ID        string    `gorm:primaryKey;type:varchar(36) `
	Name      string    `gorm:type:varchar(50) `
	ParentID  string    `gorm:type:varchar(36) `
	Level     int       `gorm:default:1 `
	Icon      string    `gorm:type:varchar(500) `
	SortOrder int       `gorm:default:0 `
	IsEnabled int       `gorm:default:1 `
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Category) TableName() string {
	return categories
}
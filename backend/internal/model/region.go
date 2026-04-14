package model

import (
	Time   time   time
)

type Region struct {
	Code       string    `gorm:primaryKey;type:varchar(20) `
	ID         string    `gorm:type:varchar(36) `
	Name       string    `gorm:column:region_name;type:varchar(50) `
	IsEnabled  int       `gorm:default:1 `
	SortOrder  int       `gorm:default:0 `
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (Region) TableName() string {
	return regions
}
package model

import (
	Time   time   time
)

type HomeConfig struct {
	ID          string    `gorm:primaryKey;type:varchar(36) `
	RegionCode  string    `gorm:type:varchar(20);index `
	ConfigType  string    `gorm:type:varchar(20) `
	ConfigKey   string    `gorm:type:varchar(50) `
	ConfigValue string    `gorm:type:json `
	SortOrder   int       `gorm:default:0 `
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (HomeConfig) TableName() string {
	return home_configs
}
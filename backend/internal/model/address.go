package model

import (
	Time   time   time
)

type Address struct {
	ID         string    `gorm:primaryKey;type:varchar(36) `
	UserID     string    `gorm:type:varchar(36);index `
	Receiver   string    `gorm:type:varchar(50) `
	Phone      string    `gorm:type:varchar(20) `
	Province   string    `gorm:type:varchar(50) `
	City       string    `gorm:type:varchar(50) `
	District   string    `gorm:type:varchar(50) `
	Detail     string    `gorm:type:varchar(200) `
	IsDefault  int       `gorm:default:0;column:is_default `
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (Address) TableName() string {
	return addresses
}
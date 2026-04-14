package model

import (
	Time   time   time
)

type User struct {
	ID          string     `gorm:primaryKey;type:varchar(36) `
	OpenID      string     `gorm:uniqueIndex;type:varchar(128) `
	Nickname    string     `gorm:type:varchar(100) `
	Avatar      string     `gorm:type:varchar(500) `
	Phone       string     `gorm:type:varchar(20) `
	Status      int        `gorm:default:1 `
	LastLoginAt *time.Time `gorm:column:last_login_at `
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:index `
}

func (User) TableName() string {
	return users
}
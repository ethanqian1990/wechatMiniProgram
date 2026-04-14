package model

import (
	Time   time   time
)

type SearchHistory struct {
	ID        string    `gorm:primaryKey;type:varchar(36) `
	UserID    string    `gorm:type:varchar(36);index `
	Keyword   string    `gorm:type:varchar(100) `
	CreatedAt time.Time
}

func (SearchHistory) TableName() string {
	return search_history
}
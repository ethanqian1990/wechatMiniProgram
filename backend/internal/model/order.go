package model

import (
	Time   time   time
)

type Order struct {
	ID                string     `gorm:primaryKey;type:varchar(36) `
	OrderNo           string     `gorm:uniqueIndex;type:varchar(64) `
	UserID            string     `gorm:type:varchar(36);index `
	TotalAmount       float64    `gorm:type:decimal(10,2) `
	DiscountAmount    float64    `gorm:type:decimal(10,2);default:0 `
	FreightAmount     float64    `gorm:type:decimal(10,2);default:0 `
	FinalAmount       float64    `gorm:type:decimal(10,2) `
	Status            string     `gorm:default:PENDING_PAY;type:varchar(20) `
	Receiver          string     `gorm:type:varchar(50) `
	Phone             string     `gorm:type:varchar(20) `
	Province          string     `gorm:type:varchar(50) `
	City              string     `gorm:type:varchar(50) `
	District          string     `gorm:type:varchar(50) `
	Detail            string     `gorm:type:varchar(200) `
	PayAt             *time.Time `gorm:column:pay_at `
	ShipAt            *time.Time `gorm:column:ship_at `
	ReceiveAt         *time.Time `gorm:column:receive_at `
	CancelAt          *time.Time `gorm:column:cancel_at `
	CancelReason      string     `gorm:column:cancel_reason;type:varchar(200) `
	ClientOrderToken  string     `gorm:column:client_order_token;type:varchar(64) `
	RegionCode        string     `gorm:column:region_code;type:varchar(20) `
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (Order) TableName() string {
	return orders
}
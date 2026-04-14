package model

import (
	"time"
)

type User struct {
	ID          string     `gorm:"primaryKey;type:varchar(36)"`
	OpenID      string     `gorm:"uniqueIndex;type:varchar(128)"`
	Nickname    string     `gorm:"type:varchar(100)"`
	Avatar      string     `gorm:"type:varchar(500)"`
	Phone       string     `gorm:"type:varchar(20)"`
	Status      int        `gorm:"default:1"`
	LastLoginAt *time.Time `gorm:"column:last_login_at"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

type Address struct {
	ID         string    `gorm:"primaryKey;type:varchar(36)"`
	UserID     string    `gorm:"type:varchar(36);index"`
	Receiver   string    `gorm:"type:varchar(50)"`
	Phone      string    `gorm:"type:varchar(20)"`
	Province   string    `gorm:"type:varchar(50)"`
	City       string    `gorm:"type:varchar(50)"`
	District   string    `gorm:"type:varchar(50)"`
	Detail     string    `gorm:"type:varchar(200)"`
	IsDefault  int       `gorm:"default:0;column:is_default"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (Address) TableName() string {
	return "addresses"
}

type Region struct {
	Code       string    `gorm:"primaryKey;type:varchar(20)"`
	Name       string    `gorm:"column:region_name;type:varchar(50)"`
	IsEnabled  int       `gorm:"default:1"`
	SortOrder  int       `gorm:"default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (Region) TableName() string {
	return "regions"
}

type Category struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	Name      string    `gorm:"type:varchar(50)"`
	ParentID  string    `gorm:"type:varchar(36)"`
	Level     int       `gorm:"default:1"`
	SortOrder int       `gorm:"default:0"`
	IsEnabled int       `gorm:"default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Category) TableName() string {
	return "categories"
}

type Product struct {
	ID               string     `gorm:"primaryKey;type:varchar(36)"`
	Name             string     `gorm:"type:varchar(200)"`
	Description      string     `gorm:"type:text"`
	MainImage        string     `gorm:"type:varchar(500)"`
	Images           string     `gorm:"type:json"`
	CategoryID       string     `gorm:"type:varchar(36)"`
	Price            float64    `gorm:"type:decimal(10,2);default:0"`
	OriginalPrice    *float64   `gorm:"type:decimal(10,2)"`
	Stock            int        `gorm:"default:0"`
	SalesCount       int        `gorm:"default:0"`
	IsOnShelf        int        `gorm:"default:1;column:is_on_shelf"`
	ShowOnHome       int        `gorm:"default:1;column:show_on_home"`
	AvailableRegions string     `gorm:"type:json;column:available_regions"`
	Tags             string     `gorm:"type:json"`
	PromoText        string     `gorm:"type:varchar(200)"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time `gorm:"index"`
}

func (Product) TableName() string {
	return "products"
}

type ProductSKU struct {
	ID            string    `gorm:"primaryKey;type:varchar(36)"`
	ProductID     string    `gorm:"type:varchar(36);index"`
	SkuName       string    `gorm:"type:varchar(200)"`
	Specs         string    `gorm:"type:json"`
	Price         float64   `gorm:"type:decimal(10,2)"`
	OriginalPrice *float64  `gorm:"type:decimal(10,2)"`
	Stock         int       `gorm:"default:0"`
	Image         string    `gorm:"type:varchar(500)"`
	IsEnabled     int       `gorm:"default:1"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (ProductSKU) TableName() string {
	return "product_skus"
}

type Cart struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	UserID    string    `gorm:"type:varchar(36);index"`
	ProductID string    `gorm:"type:varchar(36)"`
	SkuID     *string   `gorm:"type:varchar(36)"`
	Quantity  int       `gorm:"default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Cart) TableName() string {
	return "carts"
}
type Order struct {
	ID                string     `gorm:"primaryKey;type:varchar(36)"`
	OrderNo           string     `gorm:"uniqueIndex;type:varchar(64)"`
	UserID            string     `gorm:"type:varchar(36);index"`
	TotalAmount       float64    `gorm:"type:decimal(10,2)"`
	DiscountAmount    float64    `gorm:"type:decimal(10,2);default:0"`
	FreightAmount     float64    `gorm:"type:decimal(10,2);default:0"`
	FinalAmount       float64    `gorm:"type:decimal(10,2)"`
	Status            string     `gorm:"default:PENDING_PAY;type:varchar(20)"`
	Receiver          string     `gorm:"type:varchar(50)"`
	Phone             string     `gorm:"type:varchar(20)"`
	Province          string     `gorm:"type:varchar(50)"`
	City              string     `gorm:"type:varchar(50)"`
	District          string     `gorm:"type:varchar(50)"`
	Detail            string     `gorm:"type:varchar(200)"`
	PayAt             *time.Time `gorm:"column:pay_at"`
	ShipAt            *time.Time `gorm:"column:ship_at"`
	ReceiveAt         *time.Time `gorm:"column:receive_at"`
	CancelAt          *time.Time `gorm:"column:cancel_at"`
	CancelReason      string     `gorm:"column:cancel_reason;type:varchar(200)"`
	ClientOrderToken  string     `gorm:"column:client_order_token;type:varchar(64)"`
	RegionCode        string     `gorm:"column:region_code;type:varchar(20)"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	ID           string    `gorm:"primaryKey;type:varchar(36)"`
	OrderID     string    `gorm:"type:varchar(36);index"`
	ProductID   string    `gorm:"type:varchar(36);index"`
	SkuID       *string   `gorm:"type:varchar(36)"`
	ProductName string    `gorm:"type:varchar(200)"`
	ProductImage string    `gorm:"type:varchar(500)"`
	SkuName     string    `gorm:"type:varchar(200)"`
	Price       float64   `gorm:"type:decimal(10,2)"`
	Quantity    int
	Subtotal    float64   `gorm:"type:decimal(10,2)"`
	CreatedAt   time.Time
}

func (OrderItem) TableName() string {
	return "order_items"
}

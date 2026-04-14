package model

import (
	"time"
)

type User struct {
	ID          string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	OpenID      string     `json:"open_id" gorm:"uniqueIndex;type:varchar(128)"`
	Nickname    string     `json:"nickname" gorm:"type:varchar(100)"`
	Avatar      string     `json:"avatar" gorm:"type:varchar(500)"`
	Phone       string     `json:"phone" gorm:"type:varchar(20)"`
	Status      int        `json:"status" gorm:"default:1"`
	LastLoginAt *time.Time `json:"last_login_at" gorm:"column:last_login_at"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

type Address struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	UserID     string    `json:"user_id" gorm:"type:varchar(36);index"`
	Receiver   string    `json:"receiver" gorm:"type:varchar(50)"`
	Phone      string    `json:"phone" gorm:"type:varchar(20)"`
	Province   string    `json:"province" gorm:"type:varchar(50)"`
	City       string    `json:"city" gorm:"type:varchar(50)"`
	District   string    `json:"district" gorm:"type:varchar(50)"`
	Detail     string    `json:"detail" gorm:"type:varchar(200)"`
	IsDefault  int       `json:"is_default" gorm:"default:0;column:is_default"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (Address) TableName() string {
	return "addresses"
}

type Region struct {
	Code       string    `json:"code" gorm:"primaryKey;type:varchar(20)"`
	Name       string    `json:"name" gorm:"column:region_name;type:varchar(50)"`
	IsEnabled  int       `json:"is_enabled" gorm:"default:1"`
	SortOrder  int       `json:"sort_order" gorm:"default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (Region) TableName() string {
	return "regions"
}

type Category struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name      string    `json:"name" gorm:"type:varchar(50)"`
	ParentID  string    `json:"parent_id" gorm:"type:varchar(36)"`
	Level     int       `json:"level" gorm:"default:1"`
	SortOrder int       `json:"sort_order" gorm:"default:0"`
	IsEnabled int       `json:"is_enabled" gorm:"default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Category) TableName() string {
	return "categories"
}

type Product struct {
	ID               string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name             string     `json:"name" gorm:"type:varchar(200)"`
	Description      string     `json:"description" gorm:"type:text"`
	MainImage        string     `json:"main_image" gorm:"type:varchar(500)"`
	Images           string     `json:"images" gorm:"type:json"`
	CategoryID       string     `json:"category_id" gorm:"type:varchar(36)"`
	Price            float64    `json:"price" gorm:"type:decimal(10,2);default:0"`
	OriginalPrice    *float64   `json:"original_price" gorm:"type:decimal(10,2)"`
	Stock            int        `json:"stock" gorm:"default:0"`
	SalesCount       int        `json:"sales_count" gorm:"default:0"`
	IsOnShelf        int        `json:"is_on_shelf" gorm:"default:1;column:is_on_shelf"`
	ShowOnHome       int        `json:"show_on_home" gorm:"default:1;column:show_on_home"`
	AvailableRegions string     `json:"available_regions" gorm:"type:json;column:available_regions"`
	Tags             string     `json:"tags" gorm:"type:json"`
	PromoText        string     `json:"promo_text" gorm:"type:varchar(200)"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time `gorm:"index"`
}

func (Product) TableName() string {
	return "products"
}

type ProductSKU struct {
	ID            string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ProductID     string    `json:"product_id" gorm:"type:varchar(36);index"`
	SkuName       string    `json:"sku_name" gorm:"type:varchar(200)"`
	Specs         string    `json:"specs" gorm:"type:json"`
	Price         float64   `json:"price" gorm:"type:decimal(10,2)"`
	OriginalPrice *float64  `json:"original_price" gorm:"type:decimal(10,2)"`
	Stock         int       `json:"stock" gorm:"default:0"`
	Image         string    `json:"image" gorm:"type:varchar(500)"`
	IsEnabled     int       `json:"is_enabled" gorm:"default:1"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (ProductSKU) TableName() string {
	return "product_skus"
}

type Cart struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	UserID    string    `json:"user_id" gorm:"type:varchar(36);index"`
	ProductID string    `json:"product_id" gorm:"type:varchar(36)"`
	SkuID     *string   `json:"sku_id" gorm:"type:varchar(36)"`
	Quantity  int       `json:"quantity" gorm:"default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Cart) TableName() string {
	return "carts"
}
type Order struct {
	ID               string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	OrderNo          string     `json:"order_no" gorm:"uniqueIndex;type:varchar(64)"`
	UserID           string     `json:"user_id" gorm:"type:varchar(36);index"`
	TotalAmount      float64    `json:"total_amount" gorm:"type:decimal(10,2)"`
	DiscountAmount   float64    `json:"discount_amount" gorm:"type:decimal(10,2);default:0"`
	FreightAmount    float64    `json:"freight_amount" gorm:"type:decimal(10,2);default:0"`
	FinalAmount      float64    `json:"final_amount" gorm:"type:decimal(10,2)"`
	Status           string     `json:"status" gorm:"default:PENDING_PAY;type:varchar(20)"`
	Receiver         string     `json:"receiver" gorm:"type:varchar(50)"`
	Phone            string     `json:"phone" gorm:"type:varchar(20)"`
	Province         string     `json:"province" gorm:"type:varchar(50)"`
	City             string     `json:"city" gorm:"type:varchar(50)"`
	District         string     `json:"district" gorm:"type:varchar(50)"`
	Detail           string     `json:"detail" gorm:"type:varchar(200)"`
	PayAt            *time.Time `json:"pay_at" gorm:"column:pay_at"`
	ShipAt           *time.Time `json:"ship_at" gorm:"column:ship_at"`
	ReceiveAt        *time.Time `json:"receive_at" gorm:"column:receive_at"`
	CancelAt         *time.Time `json:"cancel_at" gorm:"column:cancel_at"`
	CancelReason     string     `json:"cancel_reason" gorm:"column:cancel_reason;type:varchar(200)"`
	ClientOrderToken string     `json:"client_order_token" gorm:"column:client_order_token;type:varchar(64)"`
	RegionCode       string     `json:"region_code" gorm:"column:region_code;type:varchar(20)"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	ID           string   `json:"id" gorm:"primaryKey;type:varchar(36)"`
	OrderID      string   `json:"order_id" gorm:"type:varchar(36);index"`
	ProductID    string   `json:"product_id" gorm:"type:varchar(36);index"`
	SkuID        *string  `json:"sku_id" gorm:"type:varchar(36)"`
	ProductName  string   `json:"product_name" gorm:"type:varchar(200)"`
	ProductImage string   `json:"product_image" gorm:"type:varchar(500)"`
	SkuName      string   `json:"sku_name" gorm:"type:varchar(200)"`
	Price        float64  `json:"price" gorm:"type:decimal(10,2)"`
	Quantity     int      `json:"quantity"`
	Subtotal     float64  `json:"subtotal" gorm:"type:decimal(10,2)"`
	CreatedAt   time.Time
}

func (OrderItem) TableName() string {
	return "order_items"
}

type SearchHistory struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	UserID    string    `json:"user_id" gorm:"type:varchar(36);index"`
	Keyword   string    `json:"keyword" gorm:"type:varchar(100)"`
	CreatedAt time.Time
}

func (SearchHistory) TableName() string {
	return "search_history"
}

type StockReservation struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	OrderID   string    `json:"order_id" gorm:"type:varchar(36);index"`
	UserID    string    `json:"user_id" gorm:"type:varchar(36);index"`
	ProductID string    `json:"product_id" gorm:"type:varchar(36)"`
	SkuID     string    `json:"sku_id" gorm:"type:varchar(36)"`
	Quantity  int       `json:"quantity"`
	Status    string    `json:"status" gorm:"default:FROZEN;type:varchar(20)"`
	ExpireAt  time.Time `json:"expire_at" gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (StockReservation) TableName() string {
	return "stock_reservations"
}

type HomeConfig struct {
	ID          string   `json:"id" gorm:"primaryKey;type:varchar(36)"`
	RegionCode  string   `json:"region_code" gorm:"column:region_code;type:varchar(20);index"`
	ConfigType  string   `json:"config_type" gorm:"column:config_type;type:varchar(50);index"`
	ConfigValue string   `json:"config_value" gorm:"column:config_value;type:json"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (HomeConfig) TableName() string {
	return "home_configs"
}

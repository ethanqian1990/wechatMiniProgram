package repository

import (
	"fmt"
	"time"

	"github.com/ethanqian1990/wechat-mall-backend/internal/config"
	"github.com/ethanqian1990/wechat-mall-backend/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)

	DB = db
	return nil
}

// User Repository
type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindByOpenID(openID string) (*model.User, error) {
	var user model.User
	err := DB.Where("openid = ?", openID).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	err := DB.First(&user, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) Create(user *model.User) error {
	return DB.Create(user).Error
}

func (r *UserRepository) Update(userID, nickname, avatar string) error {
	updates := map[string]interface{}{}
	if nickname != "" {
		updates["nickname"] = nickname
	}
	if avatar != "" {
		updates["avatar"] = avatar
	}
	return DB.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (r *UserRepository) UpdateLastLogin(userID string) error {
	return DB.Model(&model.User{}).Where("id = ?", userID).Update("last_login_at", time.Now()).Error
}

// Address Repository
type AddressRepository struct{}

func NewAddressRepository() *AddressRepository {
	return &AddressRepository{}
}

func (r *AddressRepository) FindByUserID(userID string) ([]model.Address, error) {
	var addresses []model.Address
	err := DB.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses).Error
	return addresses, err
}

func (r *AddressRepository) FindByIDAndUserID(id, userID string) (*model.Address, error) {
	var address model.Address
	err := DB.First(&address, "id = ? AND user_id = ?", id, userID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &address, err
}

func (r *AddressRepository) FindDefault(userID string) (*model.Address, error) {
	var address model.Address
	err := DB.First(&address, "user_id = ? AND is_default = ?", userID, 1).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &address, err
}

func (r *AddressRepository) Create(address *model.Address) error {
	return DB.Create(address).Error
}

func (r *AddressRepository) Update(id, userID string, updates map[string]interface{}) error {
	return DB.Model(&model.Address{}).Where("id = ? AND user_id = ?", id, userID).Updates(updates).Error
}

func (r *AddressRepository) Delete(id, userID string) error {
	return DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Address{}).Error
}

func (r *AddressRepository) ClearDefault(userID string) error {
	return DB.Model(&model.Address{}).Where("user_id = ?", userID).Update("is_default", 0).Error
}

func (r *AddressRepository) SetDefault(id, userID string) error {
	return DB.Model(&model.Address{}).Where("id = ? AND user_id = ?", id, userID).Update("is_default", 1).Error
}

// Region Repository
type RegionRepository struct{}

func NewRegionRepository() *RegionRepository {
	return &RegionRepository{}
}

func (r *RegionRepository) FindEnabled() ([]model.Region, error) {
	var regions []model.Region
	err := DB.Where("is_enabled = ?", 1).Order("sort_order ASC").Find(&regions).Error
	return regions, err
}

func (r *RegionRepository) FindByCode(code string) (*model.Region, error) {
	var region model.Region
	err := DB.First(&region, "code = ?", code).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &region, err
}

// Category Repository
type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (r *CategoryRepository) FindEnabled() ([]model.Category, error) {
	var categories []model.Category
	err := DB.Where("is_enabled = ?", 1).Order("sort_order ASC").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) FindByID(id string) (*model.Category, error) {
	var category model.Category
	err := DB.First(&category, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &category, err
}

// Product Repository
type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

func (r *ProductRepository) FindByFilter(regionCode, categoryID, keyword string, onShelf, showOnHome int, page, pageSize int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := DB.Model(&model.Product{}).Where("is_on_shelf = ? AND deleted_at IS NULL", 1)
	if regionCode != "" {
		query = query.Where("available_regions LIKE ?", "%"+regionCode+"%")
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Order("sales_count DESC").Offset(offset).Limit(pageSize).Find(&products).Error
	return products, total, err
}

func (r *ProductRepository) FindByID(id string) (*model.Product, error) {
	var product model.Product
	err := DB.First(&product, "id = ? AND deleted_at IS NULL", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &product, err
}

func (r *ProductRepository) UpdateOnShelf(id string, onShelf int) error {
	return DB.Model(&model.Product{}).Where("id = ?", id).Update("is_on_shelf", onShelf).Error
}

func (r *ProductRepository) UpdateShowOnHome(id string, showOnHome int) error {
	return DB.Model(&model.Product{}).Where("id = ?", id).Update("show_on_home", showOnHome).Error
}

// SKU Repository
type SKURepository struct{}

func NewSKURepository() *SKURepository {
	return &SKURepository{}
}

func (r *SKURepository) FindByProductID(productID string) ([]model.ProductSKU, error) {
	var skus []model.ProductSKU
	err := DB.Where("product_id = ? AND is_enabled = ?", productID, 1).Find(&skus).Error
	return skus, err
}

func (r *SKURepository) FindByID(id string) (*model.ProductSKU, error) {
	var sku model.ProductSKU
	err := DB.First(&sku, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &sku, err
}

// Cart Repository
type CartRepository struct{}

func NewCartRepository() *CartRepository {
	return &CartRepository{}
}

func (r *CartRepository) FindByUserID(userID string) ([]model.Cart, error) {
	var carts []model.Cart
	err := DB.Where("user_id = ?", userID).Find(&carts).Error
	return carts, err
}

func (r *CartRepository) FindByUserProductSKU(userID, productID, skuID string) (*model.Cart, error) {
	var cart model.Cart
	query := DB.Where("user_id = ? AND product_id = ?", userID, productID)
	if skuID == "" {
		query = query.Where("sku_id IS NULL")
	} else {
		query = query.Where("sku_id = ?", skuID)
	}
	err := query.First(&cart).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &cart, err
}

func (r *CartRepository) Create(cart *model.Cart) error {
	return DB.Create(cart).Error
}

func (r *CartRepository) UpdateQuantity(id string, quantity int) error {
	return DB.Model(&model.Cart{}).Where("id = ?", id).Update("quantity", quantity).Error
}

func (r *CartRepository) Delete(id string) error {
	return DB.Delete(&model.Cart{}, "id = ?", id).Error
}

func (r *CartRepository) DeleteByUserID(userID string) error {
	return DB.Where("user_id = ?", userID).Delete(&model.Cart{}).Error
}
// Order Repository
type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (r *OrderRepository) FindAll(status string, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := DB.Model(&model.Order{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
	return orders, total, err
}

func (r *OrderRepository) FindByID(id string) (*model.Order, error) {
	var order model.Order
	err := DB.First(&order, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &order, err
}

func (r *OrderRepository) UpdateStatus(id, status string) error {
	return DB.Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *OrderRepository) Update(id string, updates map[string]interface{}) error {
	return DB.Model(&model.Order{}).Where("id = ?", id).Updates(updates).Error
}

type OrderItemRepository struct{}

func NewOrderItemRepository() *OrderItemRepository {
	return &OrderItemRepository{}
}

func (r *OrderItemRepository) FindByOrderID(orderID string) ([]model.OrderItem, error) {
	var items []model.OrderItem
	err := DB.Where("order_id = ?", orderID).Find(&items).Error
	return items, err
}

// Additional methods for Product
func (r *ProductRepository) FindAll(page, pageSize int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64
	DB.Model(&model.Product{}).Count(&total)
	offset := (page - 1) * pageSize
	err := DB.Offset(offset).Limit(pageSize).Find(&products).Error
	return products, total, err
}

func (r *ProductRepository) Create(product *model.Product) error {
	return DB.Create(product).Error
}

func (r *ProductRepository) Delete(id string) error {
	return DB.Delete(&model.Product{}, "id = ?", id).Error
}

// Additional methods for Region
func (r *RegionRepository) FindAll() ([]model.Region, error) {
	var regions []model.Region
	err := DB.Order("sort_order ASC").Find(&regions).Error
	return regions, err
}

func (r *RegionRepository) Create(region *model.Region) error {
	return DB.Create(region).Error
}

// Additional methods for Category
func (r *CategoryRepository) FindAll() ([]model.Category, error) {
	var categories []model.Category
	err := DB.Order("sort_order ASC").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) Create(category *model.Category) error {
	return DB.Create(category).Error
}

func (r *CategoryRepository) Update(id string, updates map[string]interface{}) error {
	return DB.Model(&model.Category{}).Where("id = ?", id).Updates(updates).Error
}

func (r *CategoryRepository) Delete(id string) error {

	return DB.Delete(&model.Category{}, "id = ?", id).Error
}

// Missing methods
func (r *ProductRepository) Update(id string, updates map[string]interface{}) error {
	return DB.Model(&model.Product{}).Where("id = ?", id).Updates(updates).Error
}

func (r *UserRepository) FindAll(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	DB.Model(&model.User{}).Count(&total)
	offset := (page - 1) * pageSize
	err := DB.Offset(offset).Limit(pageSize).Find(&users).Error
	return users, total, err
}

func (r *RegionRepository) Update(code string, updates map[string]interface{}) error {
	return DB.Model(&model.Region{}).Where("code = ?", code).Updates(updates).Error
}

func (r *RegionRepository) Delete(code string) error {
	return DB.Delete(&model.Region{}, "code = ?", code).Error
}

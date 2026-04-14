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

func (r *ProductRepository) FindByFilter(regionCode, categoryID, keyword string, onShelf, showOnHome int, sort string, page, pageSize int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := DB.Model(&model.Product{}).Where("deleted_at IS NULL")
	if regionCode != "" {
		query = query.Where("JSON_CONTAINS(available_regions, ?)", fmt.Sprintf("\"%s\"", regionCode))
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if onShelf == 0 || onShelf == 1 {
		query = query.Where("is_on_shelf = ?", onShelf)
	}
	if showOnHome == 0 || showOnHome == 1 {
		query = query.Where("show_on_home = ?", showOnHome)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize

	orderBy := "sales_count DESC"
	switch sort {
	case "sales_desc":
		orderBy = "sales_count DESC"
	case "created_desc":
		orderBy = "created_at DESC"
	case "price_asc":
		orderBy = "price ASC"
	case "price_desc":
		orderBy = "price DESC"
	}

	err := query.Order(orderBy).Offset(offset).Limit(pageSize).Find(&products).Error
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

func (r *ProductRepository) FindHomeCandidates(regionCode, categoryID string, limit int, excludeIDs []string) ([]model.Product, error) {
	var products []model.Product

	query := DB.Model(&model.Product{}).
		Where("is_on_shelf = ? AND show_on_home = ? AND deleted_at IS NULL", 1, 1)
	if regionCode != "" {
		// available_regions is JSON array of region_code strings
		query = query.Where("JSON_CONTAINS(available_regions, ?)", fmt.Sprintf("\"%s\"", regionCode))
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if len(excludeIDs) > 0 {
		query = query.Where("id NOT IN ?", excludeIDs)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Order("sales_count DESC").Find(&products).Error
	return products, err
}

// HomeConfig Repository
type HomeConfigRepository struct{}

func NewHomeConfigRepository() *HomeConfigRepository {
	return &HomeConfigRepository{}
}

func (r *HomeConfigRepository) FindByRegionAndType(regionCode, configType string) (*model.HomeConfig, error) {
	var cfg model.HomeConfig
	err := DB.First(&cfg, "region_code = ? AND config_type = ?", regionCode, configType).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &cfg, err
}

func (r *HomeConfigRepository) Upsert(regionCode, configType, configValue string) error {
	now := time.Now()
	var existing model.HomeConfig
	err := DB.First(&existing, "region_code = ? AND config_type = ?", regionCode, configType).Error
	if err == nil {
		return DB.Model(&model.HomeConfig{}).
			Where("id = ?", existing.ID).
			Updates(map[string]interface{}{
				"config_value": configValue,
				"updated_at":   now,
			}).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return DB.Create(&model.HomeConfig{
		ID:          "homecfg_" + fmt.Sprintf("%d", now.UnixNano()),
		RegionCode:  regionCode,
		ConfigType:  configType,
		ConfigValue: configValue,
		CreatedAt:   now,
		UpdatedAt:   now,
	}).Error
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


func (r *OrderRepository) FindByUser(userID, status string, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := DB.Model(&model.Order{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
	return orders, total, err
}

func (r *OrderRepository) FindByIDAndUser(id, userID string) (*model.Order, error) {
	var order model.Order
	err := DB.First(&order, "id = ? AND user_id = ?", id, userID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &order, err
}

func (r *OrderRepository) FindByUserAndClientToken(userID, token string) (*model.Order, error) {
	var order model.Order
	err := DB.First(&order, "user_id = ? AND client_order_token = ?", userID, token).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &order, err
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

// FindExpiredOrders 查询过期订单
func (r *OrderRepository) FindExpiredOrders(status string, expireTime time.Time) ([]model.Order, error) {
	var orders []model.Order
	err := DB.Where("status = ? AND created_at < ?", status, expireTime).Find(&orders).Error
	return orders, err
}

// ReleaseByOrderID 释放订单的所有冻结库存
func (r *StockReservationRepository) ReleaseByOrderID(orderID string) error {
	var reservations []model.StockReservation
	err := DB.Where("order_id = ? AND status = ?", orderID, "FROZEN").Find(&reservations).Error
	if err != nil {
		return err
	}

	for _, res := range reservations {
		DB.Model(&model.ProductSKU{}).Where("id = ?", res.SkuID).
			Update("stock", gorm.Expr("stock + ?", res.Quantity))
		DB.Model(&model.StockReservation{}).Where("id = ?", res.ID).Update("status", "RELEASED")
	}
	return nil
}

// FindExpired 查询过期的冻结记录
func (r *StockReservationRepository) FindExpired(now time.Time) ([]model.StockReservation, error) {
	var reservations []model.StockReservation
	err := DB.Where("status = ? AND expire_at < ?", "FROZEN", now).Find(&reservations).Error
	return reservations, err
}

// MarkReleased 标记为已释放
func (r *StockReservationRepository) MarkReleased(id string) error {
	return DB.Model(&model.StockReservation{}).Where("id = ?", id).Update("status", "RELEASED").Error
}

// SearchHistory Repository
type SearchHistoryRepository struct{}

func NewSearchHistoryRepository() *SearchHistoryRepository {
	return &SearchHistoryRepository{}
}

func (r *SearchHistoryRepository) Create(history *model.SearchHistory) error {
	return DB.Create(history).Error
}

func (r *SearchHistoryRepository) FindByUserID(userID string) ([]model.SearchHistory, error) {
	var histories []model.SearchHistory
	err := DB.Where("user_id = ?", userID).Order("created_at DESC").Limit(20).Find(&histories).Error
	return histories, err
}

func (r *SearchHistoryRepository) DeleteByUserAndKeyword(userID, keyword string) error {
	return DB.Where("user_id = ? AND keyword = ?", userID, keyword).Delete(&model.SearchHistory{}).Error
}

func (r *SearchHistoryRepository) DeleteByUserID(userID string) error {
	return DB.Where("user_id = ?", userID).Delete(&model.SearchHistory{}).Error
}

func (r *SearchHistoryRepository) FindHotKeywords(limit int) ([]string, error) {
	var keywords []string
	err := DB.Model(&model.SearchHistory{}).
		Select("keyword, COUNT(*) as count").
		Group("keyword").
		Order("count DESC").
		Limit(limit).
		Pluck("keyword", &keywords).Error
	return keywords, err
}

// StockReservation Repository
type StockReservationRepository struct{}

func NewStockReservationRepository() *StockReservationRepository {
	return &StockReservationRepository{}
}

func (r *StockReservationRepository) MarkConsumedByOrderIDTx(tx *gorm.DB, orderID string) error {
	return tx.Model(&model.StockReservation{}).
		Where("order_id = ? AND status = ?", orderID, "FROZEN").
		Update("status", "CONSUMED").Error
}

// ReleaseStock 释放冻结的库存
func (r *SKURepository) ReleaseStock(skuID string, quantity int) error {
	return DB.Model(&model.ProductSKU{}).Where("id = ?", skuID).
		Update("stock", gorm.Expr("stock + ?", quantity)).Error
}

// FindByOrderNo 根据订单号查找

// FindByIDAndUserID 根据ID和用户ID查找
func (r *OrderRepository) FindByIDAndUserID(id, userID string) (*model.Order, error) {
	var order model.Order
	err := DB.First(&order, "id = ? AND user_id = ?", id, userID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &order, err
}

// Update 更新订单

// ConfirmByOrderID 确认库存冻结
func (r *StockReservationRepository) ConfirmByOrderID(orderID string) error {
	return DB.Model(&model.StockReservation{}).Where("order_id = ? AND status = ?", orderID, "FROZEN").
		Update("status", "CONSUMED").Error
}

func (r *OrderRepository) FindByOrderNo(orderNo string) ([]model.Order, error) {
	var orders []model.Order
	err := DB.Where("order_no = ?", orderNo).Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) UpdateOrder(id string, updates map[string]interface{}) error {
	return DB.Model(&model.Order{}).Where("id = ?", id).Updates(updates).Error
}

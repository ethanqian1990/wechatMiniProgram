package service

import (
	"fmt"
	"time"

	"github.com/ethanqian1990/wechat-mall-backend/internal/middleware"
	"github.com/ethanqian1990/wechat-mall-backend/internal/model"
	"github.com/ethanqian1990/wechat-mall-backend/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

func (s *UserService) WxLogin(code string) (*model.User, string, error) {
	openID := code
	user, err := s.userRepo.FindByOpenID(openID)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		user = &model.User{
			ID:       "u_" + uuid.New().String()[:8],
			OpenID:   openID,
			Nickname: "微信用户",
			Status:   1,
		}
		if err := s.userRepo.Create(user); err != nil {
			return nil, "", err
		}
	}
	s.userRepo.UpdateLastLogin(user.ID)
	token, err := middleware.GenerateToken(user.ID, "user")
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (s *UserService) GetUserInfo(userID string) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *UserService) UpdateProfile(userID string, nickname, avatar string) error {
	return s.userRepo.Update(userID, nickname, avatar)
}

type RegionService struct {
	regionRepo *repository.RegionRepository
}

func NewRegionService() *RegionService {
	return &RegionService{
		regionRepo: repository.NewRegionRepository(),
	}
}

func (s *RegionService) GetEnabledRegions() ([]model.Region, error) {
	return s.regionRepo.FindEnabled()
}

func (s *RegionService) GetRegionByCode(code string) (*model.Region, error) {
	return s.regionRepo.FindByCode(code)
}

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService() *CategoryService {
	return &CategoryService{
		categoryRepo: repository.NewCategoryRepository(),
	}
}

func (s *CategoryService) GetEnabledCategories() ([]model.Category, error) {
	return s.categoryRepo.FindEnabled()
}

func (s *CategoryService) GetCategoryByID(id string) (*model.Category, error) {
	return s.categoryRepo.FindByID(id)
}

type ProductService struct {
	productRepo *repository.ProductRepository
	skuRepo     *repository.SKURepository
}

func NewProductService() *ProductService {
	return &ProductService{
		productRepo: repository.NewProductRepository(),
		skuRepo:     repository.NewSKURepository(),
	}
}

func (s *ProductService) GetProducts(regionCode string, categoryID, keyword string, onShelf, showOnHome int, page, pageSize int) ([]model.Product, int64, error) {
	return s.productRepo.FindByFilter(regionCode, categoryID, keyword, onShelf, showOnHome, "sales_desc", page, pageSize)
}

func (s *ProductService) GetProductsWithSort(regionCode string, categoryID, keyword string, onShelf, showOnHome int, sort string, page, pageSize int) ([]model.Product, int64, error) {
	if sort == "" {
		sort = "sales_desc"
	}
	return s.productRepo.FindByFilter(regionCode, categoryID, keyword, onShelf, showOnHome, sort, page, pageSize)
}

func (s *ProductService) GetProductByID(id string) (*model.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *ProductService) GetProductWithSKUs(id string) (*model.Product, []model.ProductSKU, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	skus, err := s.skuRepo.FindByProductID(id)
	return product, skus, err
}

func (s *ProductService) SetOnShelf(id string, onShelf int) error {
	return s.productRepo.UpdateOnShelf(id, onShelf)
}

func (s *ProductService) SetShowOnHome(id string, showOnHome int) error {
	return s.productRepo.UpdateShowOnHome(id, showOnHome)
}

type CartService struct {
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
	skuRepo     *repository.SKURepository
}

func NewCartService() *CartService {
	return &CartService{
		cartRepo:    repository.NewCartRepository(),
		productRepo: repository.NewProductRepository(),
		skuRepo:     repository.NewSKURepository(),
	}
}

func (s *CartService) GetCartList(userID string) ([]model.Cart, error) {
	return s.cartRepo.FindByUserID(userID)
}

func (s *CartService) AddToCart(userID, productID, skuID string, quantity int) error {
	product, err := s.productRepo.FindByID(productID)
	if err != nil || product == nil {
		return fmt.Errorf("商品不存在")
	}
	if skuID != "" {
		sku, err := s.skuRepo.FindByID(skuID)
		if err != nil || sku == nil || sku.ProductID != productID {
			return fmt.Errorf("SKU不存在")
		}
		if sku.Stock < quantity {
			return fmt.Errorf("库存不足")
		}
	}
	existing, _ := s.cartRepo.FindByUserProductSKU(userID, productID, skuID)
	if existing != nil {
		return s.cartRepo.UpdateQuantity(existing.ID, existing.Quantity+quantity)
	}

	var skuIDPtr *string
	if skuID != "" {
		skuIDPtr = &skuID
	}
	cart := &model.Cart{
		ID:        "cart_" + uuid.New().String()[:8],
		UserID:    userID,
		ProductID: productID,
		SkuID:     skuIDPtr,
		Quantity:  quantity,
	}
	return s.cartRepo.Create(cart)
}

func (s *CartService) UpdateCartQuantity(cartID string, quantity int) error {
	if quantity <= 0 {
		return s.cartRepo.Delete(cartID)
	}
	return s.cartRepo.UpdateQuantity(cartID, quantity)
}

func (s *CartService) DeleteCartItem(cartID string) error {
	return s.cartRepo.Delete(cartID)
}

func (s *CartService) ClearCart(userID string) error {
	return s.cartRepo.DeleteByUserID(userID)
}

type AddressService struct {
	addressRepo *repository.AddressRepository
}

func NewAddressService() *AddressService {
	return &AddressService{
		addressRepo: repository.NewAddressRepository(),
	}
}

func (s *AddressService) GetAddresses(userID string) ([]model.Address, error) {
	return s.addressRepo.FindByUserID(userID)
}

func (s *AddressService) GetAddressByID(id, userID string) (*model.Address, error) {
	return s.addressRepo.FindByIDAndUserID(id, userID)
}

func (s *AddressService) CreateAddress(userID, receiver, phone, province, city, district, detail string, isDefault bool) (*model.Address, error) {
	if isDefault {
		s.addressRepo.ClearDefault(userID)
	}
	address := &model.Address{
		ID:        "addr_" + uuid.New().String()[:8],
		UserID:    userID,
		Receiver:  receiver,
		Phone:     phone,
		Province:  province,
		City:      city,
		District:  district,
		Detail:    detail,
		IsDefault: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if isDefault {
		address.IsDefault = 1
	}
	err := s.addressRepo.Create(address)
	return address, err
}

func (s *AddressService) UpdateAddress(id, userID, receiver, phone, province, city, district, detail string, isDefault bool) error {
	if isDefault {
		s.addressRepo.ClearDefault(userID)
	}
	updates := map[string]interface{}{
		"receiver": receiver,
		"phone":    phone,
		"province": province,
		"city":     city,
		"district": district,
		"detail":   detail,
	}
	if isDefault {
		updates["is_default"] = 1
	}
	return s.addressRepo.Update(id, userID, updates)
}

func (s *AddressService) DeleteAddress(id, userID string) error {
	return s.addressRepo.Delete(id, userID)
}

func (s *AddressService) SetDefault(id, userID string) error {
	s.addressRepo.ClearDefault(userID)
	return s.addressRepo.SetDefault(id, userID)
}

func (s *AddressService) GetDefaultAddress(userID string) (*model.Address, error) {
	return s.addressRepo.FindDefault(userID)
}
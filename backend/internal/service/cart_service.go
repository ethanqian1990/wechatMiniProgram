package service

import (
	github.com/ethanqian1990/wechat-mall-backend/internal/model
	github.com/ethanqian1990/wechat-mall-backend/internal/repository
	github.com/google/uuid
)

type CartService struct {
	cartRepo    *repo.CartRepository
	productRepo *repo.ProductRepository
	skuRepo     *repo.SKURepository
}

func NewCartService() *CartService {
	return &CartService{
		cartRepo:    repo.NewCartRepository(),
		productRepo: repo.NewProductRepository(),
		skuRepo:     repo.NewSKURepository(),
	}
}

// GetCartList 获取购物车列表
func (s *CartService) GetCartList(userID string) ([]model.Cart, error) {
	return s.cartRepo.FindByUserID(userID)
}

// AddToCart 添加到购物车
func (s *CartService) AddToCart(userID, productID, skuID string, quantity int) error {
	// 检查商品是否存在
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New('商品不存在')
	}

	// 检查SKU
	if skuID != '' {
		sku, err := s.skuRepo.FindByID(skuID)
		if err != nil {
			return err
		}
		if sku == nil || sku.ProductID != productID {
			return errors.New('SKU不存在')
		}
		if sku.Stock < quantity {
			return errors.New('库存不足')
		}
	}

	// 检查购物车是否已存在
	existing, err := s.cartRepo.FindByUserProductSKU(userID, productID, skuID)
	if err != nil {
		return err
	}

	if existing != nil {
		// 更新数量
		return s.cartRepo.UpdateQuantity(existing.ID, existing.Quantity+quantity)
	}

	// 新增
	cart := &model.Cart{
		ID:        'cart_' + uuid.New().String()[:8],
		UserID:    userID,
		ProductID: productID,
		SkuID:     &skuID,
		Quantity:  quantity,
	}
	return s.cartRepo.Create(cart)
}

// UpdateCartQuantity 更新购物车数量
func (s *CartService) UpdateCartQuantity(cartID string, quantity int) error {
	if quantity <= 0 {
		return s.cartRepo.Delete(cartID)
	}
	return s.cartRepo.UpdateQuantity(cartID, quantity)
}

// DeleteCartItem 删除购物车项
func (s *CartService) DeleteCartItem(cartID string) error {
	return s.cartRepo.Delete(cartID)
}

// ClearCart 清空购物车
func (s *CartService) ClearCart(userID string) error {
	return s.cartRepo.DeleteByUserID(userID)
}
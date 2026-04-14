package service

import (
	github.com/ethanqian1990/wechat-mall-backend/internal/model
	github.com/ethanqian1990/wechat-mall-backend/internal/repository
)

type ProductService struct {
	productRepo *repo.ProductRepository
	skuRepo     *repo.SKURepository
}

func NewProductService() *ProductService {
	return &ProductService{
		productRepo: repo.NewProductRepository(),
		skuRepo:     repo.NewSKURepository(),
	}
}

// GetProducts 获取商品列表
func (s *ProductService) GetProducts(regionCode string, categoryID, keyword string, onShelf, showOnHome int, page, pageSize int) ([]model.Product, int64, error) {
	return s.productRepo.FindByFilter(regionCode, categoryID, keyword, onShelf, showOnHome, page, pageSize)
}

// GetProductByID 获取商品详情
func (s *ProductService) GetProductByID(id string) (*model.Product, error) {
	return s.productRepo.FindByID(id)
}

// GetProductWithSKUs 获取商品详情(含SKU)
func (s *ProductService) GetProductWithSKUs(id string) (*model.Product, []model.ProductSKU, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	skus, err := s.skuRepo.FindByProductID(id)
	return product, skus, err
}

// SetOnShelf 设置上下架
func (s *ProductService) SetOnShelf(id string, onShelf int) error {
	return s.productRepo.UpdateOnShelf(id, onShelf)
}

// SetShowOnHome 设置首页展示
func (s *ProductService) SetShowOnHome(id string, showOnHome int) error {
	return s.productRepo.UpdateShowOnHome(id, showOnHome)
}
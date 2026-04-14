package repository

import (
	github.com/ethanqian1990/wechat-mall-backend/internal/model
	github.com/gin-gonic/gin
	gorm  gorm.io/gorm
)

type SKURepository struct{}

func NewSKURepository() *SKURepository {
	return &SKURepository{}
}

func (r *SKURepository) FindByProductID(productID string) ([]model.ProductSKU, error) {
	var skus []model.ProductSKU
	err := DB.Where('product_id = ? AND is_enabled = ?', productID, 1).Find(&skus).Error
	return skus, err
}

func (r *SKURepository) FindByID(id string) (*model.ProductSKU, error) {
	var sku model.ProductSKU
	err := DB.First(&sku, 'id = ? AND is_enabled = ?', id, 1).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &sku, err
}

func (r *SKURepository) Create(sku *model.ProductSKU) error {
	return DB.Create(sku).Error
}

func (r *SKURepository) Update(id string, updates map[string]interface{}) error {
	return DB.Model(&model.ProductSKU{}).Where('id = ?', id).Updates(updates).Error
}

func (r *SKURepository) FreezeStock(skuID string, quantity int) error {
	return DB.Model(&model.ProductSKU{}).Where('id = ? AND stock >= ?', skuID, quantity).
		Update('stock', gorm.Expr('stock - ?', quantity)).Error
}

func (r *SKURepository) DecreaseStock(skuID string, quantity int) error {
	return DB.Model(&model.ProductSKU{}).Where('id = ?', skuID).
		Update('stock', gorm.Expr('stock - ?', quantity)).Error
}

func (r *SKURepository) ReleaseStock(skuID string, quantity int) error {
	return DB.Model(&model.ProductSKU{}).Where('id = ?', skuID).
		Update('stock', gorm.Expr('stock + ?', quantity)).Error
}
package repository

import (
	github.com/ethanqian1990/wechat-mall-backend/internal/model
	github.com/gin-gonic/gin
	gorm  gorm.io/gorm
)

type CartRepository struct{}

func NewCartRepository() *CartRepository {
	return &CartRepository{}
}

func (r *CartRepository) FindByUserID(userID string) ([]model.Cart, error) {
	var carts []model.Cart
	err := DB.Where('user_id = ?', userID).Find(&carts).Error
	return carts, err
}

func (r *CartRepository) FindByUserProductSKU(userID, productID, skuID string) (*model.Cart, error) {
	var cart model.Cart
	query := DB.Where('user_id = ? AND product_id = ?', userID, productID)
	if skuID == '' {
		query = query.Where('sku_id IS NULL')
	} else {
		query = query.Where('sku_id = ?', skuID)
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
	return DB.Model(&model.Cart{}).Where('id = ?', id).Update('quantity', quantity).Error
}

func (r *CartRepository) Delete(id string) error {
	return DB.Delete(&model.Cart{}, 'id = ?', id).Error
}

func (r *CartRepository) DeleteByUserID(userID string) error {
	return DB.Where('user_id = ?', userID).Delete(&model.Cart{}).Error
}
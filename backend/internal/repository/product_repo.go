package repository

import (
	github.com/ethanqian1990/wechat-mall-backend/internal/model
	github.com/gin-gonic/gin
	gorm  gorm.io/gorm
	strings  strings
)

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

func (r *ProductRepository) FindByFilter(regionCode, categoryID, keyword string, onShelf, showOnHome int, page, pageSize int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := DB.Model(&model.Product{}).Where('is_on_shelf = ? AND deleted_at IS NULL', 1)

	if regionCode != '' {
		query = query.Where('available_regions LIKE ?', '%'+regionCode+'%')
	}
	if categoryID != '' {
		query = query.Where('category_id = ?', categoryID)
	}
	if keyword != '' {
		query = query.Where('name LIKE ?', '%'+keyword+'%')
	}
	if onShelf > 0 {
		query = query.Where('is_on_shelf = ?', onShelf)
	}
	if showOnHome > 0 {
		query = query.Where('show_on_home = ?', showOnHome)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order('sales_count DESC').Offset(offset).Limit(pageSize).Find(&products).Error

	return products, total, err
}

func (r *ProductRepository) FindByID(id string) (*model.Product, error) {
	var product model.Product
	err := DB.First(&product, 'id = ? AND deleted_at IS NULL', id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &product, err
}

func (r *ProductRepository) Create(product *model.Product) error {
	return DB.Create(product).Error
}

func (r *ProductRepository) Update(id string, updates map[string]interface{}) error {
	return DB.Model(&model.Product{}).Where('id = ?', id).Updates(updates).Error
}

func (r *ProductRepository) Delete(id string) error {
	return DB.Delete(&model.Product{}, 'id = ?', id).Error
}

func (r *ProductRepository) UpdateOnShelf(id string, onShelf int) error {
	return DB.Model(&model.Product{}).Where('id = ?', id).Update('is_on_shelf', onShelf).Error
}

func (r *ProductRepository) UpdateShowOnHome(id string, showOnHome int) error {
	return DB.Model(&model.Product{}).Where('id = ?', id).Update('show_on_home', showOnHome).Error
}

func (r *ProductRepository) UpdateSalesCount(id string, count int) error {
	return DB.Model(&model.Product{}).Where('id = ?', id).UpdateColumn('sales_count', gorm.Expr('sales_count + ?', count)).Error
}
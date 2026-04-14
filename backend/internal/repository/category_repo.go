package repository

import (
	github.com/ethanqian1990/wechat-mall-backend/internal/model
	github.com/gin-gonic/gin
	gorm  gorm.io/gorm
)

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (r *CategoryRepository) FindEnabled() ([]model.Category, error) {
	var categories []model.Category
	err := DB.Where('is_enabled = ?', 1).Order('sort_order ASC').Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) FindByID(id string) (*model.Category, error) {
	var category model.Category
	err := DB.First(&category, 'id = ?', id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &category, err
}

func (r *CategoryRepository) Create(category *model.Category) error {
	return DB.Create(category).Error
}

func (r *CategoryRepository) Update(id string, updates map[string]interface{}) error {
	return DB.Model(&model.Category{}).Where('id = ?', id).Updates(updates).Error
}

func (r *CategoryRepository) Delete(id string) error {
	return DB.Delete(&model.Category{}, 'id = ?', id).Error
}
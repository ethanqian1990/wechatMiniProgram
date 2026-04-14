package repository

import (
	github.com/ethanqian1990/wechat-mall-backend/internal/model
)

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (r *CategoryRepository) FindEnabled() ([]model.Category, error) {
	// TODO: 实现数据库查询
	return nil, nil
}

func (r *CategoryRepository) FindByID(id string) (*model.Category, error) {
	// TODO: 实现数据库查询
	return nil, nil
}

func (r *CategoryRepository) Create(category *model.Category) error {
	// TODO: 实现数据库创建
	return nil
}

func (r *CategoryRepository) Update(id string, updates map[string]interface{}) error {
	// TODO: 实现数据库更新
	return nil
}

func (r *CategoryRepository) Delete(id string) error {
	// TODO: 实现数据库删除
	return nil
}
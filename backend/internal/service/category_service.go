package service

import (
	github.com/ethanqian1990/wechat-mall-backend/internal/model
	github.com/ethanqian1990/wechat-mall-backend/internal/repository
)

type CategoryService struct {
	categoryRepo *repo.CategoryRepository
}

func NewCategoryService() *CategoryService {
	return &CategoryService{
		categoryRepo: repo.NewCategoryRepository(),
	}
}

// GetEnabledCategories 获取启用的分类列表
func (s *CategoryService) GetEnabledCategories() ([]model.Category, error) {
	return s.categoryRepo.FindEnabled()
}

// GetCategoryByID 根据ID获取分类
func (s *CategoryService) GetCategoryByID(id string) (*model.Category, error) {
	return s.categoryRepo.FindByID(id)
}
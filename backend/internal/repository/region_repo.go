package repository

import (
	github.com/ethanqian1990/wechat-mall-backend/internal/model
	github.com/gin-gonic/gin
	gorm  gorm.io/gorm
)

type RegionRepository struct{}

func NewRegionRepository() *RegionRepository {
	return &RegionRepository{}
}

func (r *RegionRepository) FindEnabled() ([]model.Region, error) {
	var regions []model.Region
	err := DB.Where('is_enabled = ?', 1).Order('sort_order ASC').Find(&regions).Error
	return regions, err
}

func (r *RegionRepository) FindByCode(code string) (*model.Region, error) {
	var region model.Region
	err := DB.First(&region, 'code = ?', code).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &region, err
}

func (r *RegionRepository) Create(region *model.Region) error {
	return DB.Create(region).Error
}

func (r *RegionRepository) Update(code string, updates map[string]interface{}) error {
	return DB.Model(&model.Region{}).Where('code = ?', code).Updates(updates).Error
}

func (r *RegionRepository) Delete(code string) error {
	return DB.Delete(&model.Region{}, 'code = ?', code).Error
}
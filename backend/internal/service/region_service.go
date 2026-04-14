package service

import (
	github.com/ethanqian1990/wechat-mall-backend/internal/model
	github.com/ethanqian1990/wechat-mall-backend/internal/repository
)

type RegionService struct {
	regionRepo *repo.RegionRepository
}

func NewRegionService() *RegionService {
	return &RegionService{
		regionRepo: repo.NewRegionRepository(),
	}
}

// GetEnabledRegions 获取启用的区域列表
func (s *RegionService) GetEnabledRegions() ([]model.Region, error) {
	return s.regionRepo.FindEnabled()
}

// GetRegionByCode 根据编码获取区域
func (s *RegionService) GetRegionByCode(code string) (*model.Region, error) {
	return s.regionRepo.FindByCode(code)
}
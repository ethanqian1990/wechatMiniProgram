package service

import (
	github.com/ethanqian1990/wechat-mall-backend/internal/model
	github.com/ethanqian1990/wechat-mall-backend/internal/repository
)

type HomeService struct {
	homeConfigRepo *repo.HomeConfigRepository
	productRepo    *repo.ProductRepository
}

func NewHomeService() *HomeService {
	return &HomeService{
		homeConfigRepo: repo.NewHomeConfigRepository(),
		productRepo:    repo.NewProductRepository(),
	}
}

// GetBanners 获取首页Banner
func (s *HomeService) GetBanners(regionCode string) ([]model.HomeConfig, error) {
	return s.homeConfigRepo.FindByType(regionCode, 'banner')
}

// GetFeatured 获取主推商品
func (s *HomeService) GetFeatured(regionCode, categoryID string) (*model.Product, error) {
	// 获取配置
	config, err := s.homeConfigRepo.FindByKey(regionCode, 'featured', 'mode')
	if err != nil || config == nil {
		// 默认按销量
		products, _, err := s.productRepo.FindByFilter(regionCode, categoryID, '', 1, 1, 1, 1)
		if err != nil || len(products) == 0 {
			return nil, nil
		}
		return &products[0], nil
	}

	// 如果是指定模式，获取指定商品
	featuredConfig, _ := s.homeConfigRepo.FindByKey(regionCode, 'featured', 'product_id')
	if featuredConfig != nil && featuredConfig.ConfigValue != '' {
		return s.productRepo.FindByID(featuredConfig.ConfigValue)
	}

	// 否则按销量
	products, _, err := s.productRepo.FindByFilter(regionCode, categoryID, '', 1, 1, 1, 1)
	if err != nil || len(products) == 0 {
		return nil, nil
	}
	return &products[0], nil
}

// GetHomeProducts 获取首页商品列表
func (s *HomeService) GetHomeProducts(regionCode, categoryID string) ([]model.Product, error) {
	// 获取列表数量配置
	config, _ := s.homeConfigRepo.FindByKey(regionCode, 'list', 'size')
	listSize := 24 // 默认
	if config != nil && config.ConfigValue != '' {
		// 解析配置
	}

	products, _, err := s.productRepo.FindByFilter(regionCode, categoryID, '', 1, 1, 1, listSize)
	if err != nil {
		return nil, err
	}

	// 排除已作为主推的商品
	featured, _ := s.GetFeatured(regionCode, categoryID)
	if featured != nil {
		var filtered []model.Product
		for _, p := range products {
			if p.ID != featured.ID {
				filtered = append(filtered, p)
			}
		}
		return filtered, nil
	}

	return products, nil
}

// GetConfig 获取首页配置
func (s *HomeService) GetConfig(regionCode string) (map[string]string, error) {
	return s.homeConfigRepo.GetConfig(regionCode)
}

// UpdateConfig 更新首页配置
func (s *HomeService) UpdateConfig(regionCode string, configType, configKey, configValue string) error {
	return s.homeConfigRepo.Upsert(regionCode, configType, configKey, configValue)
}
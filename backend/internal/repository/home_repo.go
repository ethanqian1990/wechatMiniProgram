package repository

import (
	github.com/ethanqian1990/wechat-mall-backend/internal/model
	github.com/gin-gonic/gin
	gorm  gorm.io/gorm
	encoding/json encoding/json
)

type HomeConfigRepository struct{}

func NewHomeConfigRepository() *HomeConfigRepository {
	return &HomeConfigRepository{}
}

func (r *HomeConfigRepository) FindByType(regionCode, configType string) ([]model.HomeConfig, error) {
	var configs []model.HomeConfig
	err := DB.Where('region_code = ? AND config_type = ?', regionCode, configType).
		Order('sort_order ASC').Find(&configs).Error
	return configs, err
}

func (r *HomeConfigRepository) FindByKey(regionCode, configType, configKey string) (*model.HomeConfig, error) {
	var config model.HomeConfig
	err := DB.First(&config, 'region_code = ? AND config_type = ? AND config_key = ?',
		regionCode, configType, configKey).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &config, err
}

func (r *HomeConfigRepository) GetConfig(regionCode string) (map[string]string, error) {
	var configs []model.HomeConfig
	err := DB.Where('region_code = ?', regionCode).Find(&configs).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, c := range configs {
		result[c.ConfigKey] = c.ConfigValue
	}
	return result, nil
}

func (r *HomeConfigRepository) Upsert(regionCode, configType, configKey, configValue string) error {
	var existing model.HomeConfig
	err := DB.First(&existing, 'region_code = ? AND config_type = ? AND config_key = ?',
		regionCode, configType, configKey).Error

	if err == gorm.ErrRecordNotFound {
		// 新增
		config := &model.HomeConfig{
			ID:          'hc_' + uuid.New().String()[:8],
			RegionCode:  regionCode,
			ConfigType:  configType,
			ConfigKey:   configKey,
			ConfigValue: configValue,
		}
		return DB.Create(config).Error
	} else if err != nil {
		return err
	}

	// 更新
	return DB.Model(&existing).Update('config_value', configValue).Error
}
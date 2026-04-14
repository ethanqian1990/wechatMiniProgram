package handler

import (
	github.com/gin-gonic/gin
	github.com/ethanqian1990/wechat-mall-backend/pkg/response
	github.com/ethanqian1990/wechat-mall-backend/internal/service
	github.com/ethanqian1990/wechat-mall-backend/internal/model
)

type HomeHandler struct {
	homeService *service.HomeService
}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{
		homeService: service.NewHomeService(),
	}
}

// GetBanners 获取首页Banner
func (h *HomeHandler) GetBanners(c *gin.Context) {
	regionCode := c.Query('region_code')
	if regionCode == '' {
		regionCode = 'EC'
	}

	banners, err := h.homeService.GetBanners(regionCode)
	if err != nil {
		response.Error(c, 500, '获取Banner失败')
		return
	}

	var result []gin.H
	for _, b := range banners {
		result = append(result, gin.H{
			'id':       b.ID,
			'image':    b.ConfigValue,
			'goto':     b.ConfigKey,
		})
	}

	response.Success(c, gin.H{
		'region_code': regionCode,
		'banners':     result,
	})
}

// GetFeatured 获取主推商品
func (h *HomeHandler) GetFeatured(c *gin.Context) {
	regionCode := c.Query('region_code')
	categoryID := c.Query('category_id')
	if regionCode == '' {
		regionCode = 'EC'
	}

	product, err := h.homeService.GetFeatured(regionCode, categoryID)
	if err != nil || product == nil {
		response.Success(c, gin.H{
			'region_code': regionCode,
			'category_id': categoryID,
			'featured':    nil,
		})
		return
	}

	response.Success(c, gin.H{
		'region_code': regionCode,
		'category_id': categoryID,
		'featured': gin.H{
			'id':           product.ID,
			'name':         product.Name,
			'main_image':   product.MainImage,
			'price':        product.Price,
			'original_price': product.OriginalPrice,
		},
	})
}

// GetHomeProducts 获取首页商品列表
func (h *HomeHandler) GetHomeProducts(c *gin.Context) {
	regionCode := c.Query('region_code')
	categoryID := c.Query('category_id')
	if regionCode == '' {
		regionCode = 'EC'
	}

	products, err := h.homeService.GetHomeProducts(regionCode, categoryID)
	if err != nil {
		response.Error(c, 500, '获取商品列表失败')
		return
	}

	var result []gin.H
	for _, p := range products {
		result = append(result, gin.H{
			'id':              p.ID,
			'name':            p.Name,
			'main_image':      p.MainImage,
			'price':           p.Price,
			'original_price':  p.OriginalPrice,
			'sales_count':     p.SalesCount,
			'stock':           p.Stock,
		})
	}

	response.Success(c, gin.H{
		'region_code': regionCode,
		'category_id': categoryID,
		'list':        result,
	})
}

// GetConfig 获取首页配置
func (h *HomeHandler) GetConfig(c *gin.Context) {
	regionCode := c.Query('region_code')
	if regionCode == '' {
		regionCode = 'EC'
	}

	config, err := h.homeService.GetConfig(regionCode)
	if err != nil {
		response.Error(c, 500, '获取配置失败')
		return
	}

	response.Success(c, gin.H{
		'region_code': regionCode,
		'config':      config,
	})
}

// UpdateConfig 更新首页配置
type UpdateConfigRequest struct {
	RegionCode  string `json: region_code binding:required`
	ConfigType  string `json: config_type binding:required`
	ConfigKey   string `json: config_key binding:required`
	ConfigValue string `json: config_value binding:required`
}

func (h *HomeHandler) UpdateConfig(c *gin.Context) {
	var req UpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, '参数错误')
		return
	}

	err := h.homeService.UpdateConfig(req.RegionCode, req.ConfigType, req.ConfigKey, req.ConfigValue)
	if err != nil {
		response.Error(c, 500, '更新配置失败')
		return
	}

	response.Success(c, nil)
}
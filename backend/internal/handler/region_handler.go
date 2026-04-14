package handler

import (
	github.com/gin-gonic/gin
	github.com/ethanqian1990/wechat-mall-backend/pkg/response
	github.com/ethanqian1990/wechat-mall-backend/internal/service
)

type RegionHandler struct {
	regionService *service.RegionService
}

func NewRegionHandler() *RegionHandler {
	return &RegionHandler{
		regionService: service.NewRegionService(),
	}
}

// GetRegions 获取区域列表
func (h *RegionHandler) GetRegions(c *gin.Context) {
	regions, err := h.regionService.GetEnabledRegions()
	if err != nil {
		response.Error(c, 500, '获取区域列表失败')
		return
	}

	var result []gin.H
	for _, r := range regions {
		result = append(result, gin.H{
			'code':   r.Code,
			'name':   r.Name,
			'enabled': r.IsEnabled == 1,
			'sort':   r.SortOrder,
		})
	}

	response.Success(c, gin.H{'list': result})
}
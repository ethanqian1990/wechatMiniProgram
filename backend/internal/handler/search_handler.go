package handler

import (
	github.com/gin-gonic/gin
	github.com/ethanqian1990/wechat-mall-backend/pkg/response
	github.com/ethanqian1990/wechat-mall-backend/internal/service
)

type SearchHandler struct {
	productService *service.ProductService
}

func NewSearchHandler() *SearchHandler {
	return &SearchHandler{
		productService: service.NewProductService(),
	}
}

// GetSuggest 搜索建议
func (h *SearchHandler) GetSuggest(c *gin.Context) {
	keyword := c.Query('keyword')
	if keyword == '' {
		response.Success(c, gin.H{'list': []string{}})
		return
	}

	// TODO: 实现搜索建议逻辑
	response.Success(c, gin.H{'list': []string{keyword}})
}

// GetHotSearch 热门搜索
func (h *SearchHandler) GetHotSearch(c *gin.Context) {
	// TODO: 实现热门搜索
	hotKeywords := []string{'苹果', '大米', '土鸡蛋', '葡萄干'}
	response.Success(c, gin.H{'list': hotKeywords})
}

// GetHistory 搜索历史
func (h *SearchHandler) GetHistory(c *gin.Context) {
	userID := c.GetString('user_id')
	// TODO: 从数据库获取搜索历史
	response.Success(c, gin.H{'list': []string{}})
}

// ClearHistory 清空搜索历史
func (h *SearchHandler) ClearHistory(c *gin.Context) {
	userID := c.GetString('user_id')
	// TODO: 清空搜索历史
	response.Success(c, nil)
}
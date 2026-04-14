package handler

import (
	github.com/gin-gonic/gin
	github.com/ethanqian1990/wechat-mall-backend/pkg/response
	github.com/ethanqian1990/wechat-mall-backend/internal/service
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		categoryService: service.NewCategoryService(),
	}
}

// GetCategories 获取分类列表
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.categoryService.GetEnabledCategories()
	if err != nil {
		response.Error(c, 500, '获取分类列表失败')
		return
	}

	var result []gin.H
	for _, c := range categories {
		result = append(result, gin.H{
			'id':      c.ID,
			'name':    c.Name,
			'enabled': c.IsEnabled == 1,
			'sort':    c.SortOrder,
		})
	}

	response.Success(c, gin.H{'list': result})
}
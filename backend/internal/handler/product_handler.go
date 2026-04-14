package handler

import (
	github.com/gin-gonic/gin
	github.com/ethanqian1990/wechat-mall-backend/pkg/response
	github.com/ethanqian1990/wechat-mall-backend/internal/service
	github.com/ethanqian1990/wechat-mall-backend/internal/model
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		productService: service.NewProductService(),
	}
}

// GetProducts 获取商品列表
func (h *ProductHandler) GetProducts(c *gin.Context) {
	regionCode := c.Query('region_code')
	categoryID := c.Query('category_id')
	keyword := c.Query('keyword')
	onShelf := 1 // 默认查上架
	if v := c.Query('on_shelf'); v != '' {
		fmt.Sscanf(v, '%d', &onShelf)
	}
	showOnHome := -1
	if v := c.Query('show_on_home'); v != '' {
		fmt.Sscanf(v, '%d', &showOnHome)
	}
	page := 1
	if v := c.Query('page'); v != '' {
		fmt.Sscanf(v, '%d', &page)
	}
	pageSize := 20
	if v := c.Query('page_size'); v != '' {
		fmt.Sscanf(v, '%d', &pageSize)
	}

	products, total, err := h.productService.GetProducts(regionCode, categoryID, keyword, onShelf, showOnHome, page, pageSize)
	if err != nil {
		response.Error(c, 500, '获取商品列表失败')
		return
	}

	response.Paginate(c, h.formatProducts(products), int(total), page, pageSize)
}

// GetProduct 获取商品详情
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param('id')

	product, skus, err := h.productService.GetProductWithSKUs(id)
	if err != nil {
		response.Error(c, 500, '获取商品详情失败')
		return
	}
	if product == nil {
		response.Error(c, 404, '商品不存在')
		return
	}

	response.Success(c, gin.H{
		'id':               product.ID,
		'name':             product.Name,
		'description':      product.Description,
		'main_image':       product.MainImage,
		'images':           product.Images,
		'category_id':      product.CategoryID,
		'is_on_shelf':      product.IsOnShelf,
		'show_on_home':     product.ShowOnHome,
		'available_regions': product.AvailableRegions,
		'skus':             h.formatSKUs(skus),
	})
}

func (h *ProductHandler) formatProducts(products []model.Product) []gin.H {
	var result []gin.H
	for _, p := range products {
		result = append(result, gin.H{
			'id':              p.ID,
			'name':            p.Name,
			'main_image':      p.MainImage,
			'price':           p.Price,
			'original_price':  p.OriginalPrice,
			'sales_count':    p.SalesCount,
			'stock':           p.Stock,
			'is_on_shelf':     p.IsOnShelf,
			'show_on_home':    p.ShowOnHome,
		})
	}
	return result
}

func (h *ProductHandler) formatSKUs(skus []model.ProductSKU) []gin.H {
	var result []gin.H
	for _, s := range skus {
		result = append(result, gin.H{
			'id':             s.ID,
			'sku_name':       s.SkuName,
			'price':          s.Price,
			'original_price': s.OriginalPrice,
			'stock':          s.Stock,
		})
	}
	return result
}
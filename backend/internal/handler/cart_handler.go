package handler

import (
	github.com/gin-gonic/gin
	github.com/ethanqian1990/wechat-mall-backend/pkg/response
	github.com/ethanqian1990/wechat-mall-backend/internal/service
)

type CartHandler struct {
	cartService    *service.CartService
	productService *service.ProductService
}

func NewCartHandler() *CartHandler {
	return &CartHandler{
		cartService:    service.NewCartService(),
		productService: service.NewProductService(),
	}
}

// GetCart 获取购物车列表
func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.GetString('user_id')

	carts, err := h.cartService.GetCartList(userID)
	if err != nil {
		response.Error(c, 500, '获取购物车失败')
		return
	}

	// 填充商品信息
	var result []gin.H
	for _, cart := range carts {
		product, _ := h.productService.GetProductByID(cart.ProductID)
		if product != nil {
			result = append(result, gin.H{
				'id':         cart.ID,
				'product_id': cart.ProductID,
				'sku_id':     cart.SkuID,
				'quantity':   cart.Quantity,
				'product': gin.H{
					'name':       product.Name,
					'main_image': product.MainImage,
					'price':      product.Price,
				},
			})
		}
	}

	response.Success(c, gin.H{'list': result})
}

// AddCart 添加到购物车
type AddCartRequest struct {
	ProductID string `json: product_id binding:required`
	SkuID     string `json: sku_id`
	Quantity  int    `json: quantity binding:required;min=1`
}

func (h *CartHandler) AddCart(c *gin.Context) {
	userID := c.GetString('user_id')

	var req AddCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, '参数错误')
		return
	}

	err := h.cartService.AddToCart(userID, req.ProductID, req.SkuID, req.Quantity)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateCart 更新购物车数量
type UpdateCartRequest struct {
	Quantity int `json: quantity binding:required;min=0`
}

func (h *CartHandler) UpdateCart(c *gin.Context) {
	userID := c.GetString('user_id')
	cartID := c.Param('id')

	var req UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, '参数错误')
		return
	}

	err := h.cartService.UpdateCartQuantity(cartID, req.Quantity)
	if err != nil {
		response.Error(c, 500, '更新失败')
		return
	}

	response.Success(c, nil)
}

// DeleteCart 删除购物车项
func (h *CartHandler) DeleteCart(c *gin.Context) {
	cartID := c.Param('id')

	err := h.cartService.DeleteCartItem(cartID)
	if err != nil {
		response.Error(c, 500, '删除失败')
		return
	}

	response.Success(c, nil)
}
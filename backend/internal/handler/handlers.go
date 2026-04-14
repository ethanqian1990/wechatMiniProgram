package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ethanqian1990/wechat-mall-backend/pkg/response"
	"github.com/ethanqian1990/wechat-mall-backend/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

func (h *UserHandler) WxLogin(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	user, token, err := h.userService.WxLogin(req.Code)
	if err != nil {
		response.Error(c, 500, "登录失败", err.Error())
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
		},
	})
}

func (h *UserHandler) GetInfo(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := h.userService.GetUserInfo(userID)
	if err != nil {
		response.Error(c, 500, "获取用户信息失败")
		return
	}
	response.Success(c, gin.H{
		"id":       user.ID,
		"nickname": user.Nickname,
		"avatar":   user.Avatar,
		"phone":    user.Phone,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	var req struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	err := h.userService.UpdateProfile(userID, req.Nickname, req.Avatar)
	if err != nil {
		response.Error(c, 500, "更新失败")
		return
	}
	response.Success(c, nil)
}

type RegionHandler struct {
	regionService *service.RegionService
}

func NewRegionHandler() *RegionHandler {
	return &RegionHandler{
		regionService: service.NewRegionService(),
	}
}

func (h *RegionHandler) GetRegions(c *gin.Context) {
	regions, err := h.regionService.GetEnabledRegions()
	if err != nil {
		response.Error(c, 500, "获取区域列表失败")
		return
	}
	response.Success(c, gin.H{"list": regions})
}

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		categoryService: service.NewCategoryService(),
	}
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.categoryService.GetEnabledCategories()
	if err != nil {
		response.Error(c, 500, "获取分类列表失败")
		return
	}
	response.Success(c, gin.H{"list": categories})
}

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		productService: service.NewProductService(),
	}
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	regionCode := c.Query("region_code")
	categoryID := c.Query("category_id")
	keyword := c.Query("keyword")
	page := 1
	pageSize := 20

	products, total, err := h.productService.GetProducts(regionCode, categoryID, keyword, 1, -1, page, pageSize)
	if err != nil {
		response.Error(c, 500, "获取商品列表失败")
		return
	}
	response.Paginate(c, products, int(total), page, pageSize)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	product, skus, err := h.productService.GetProductWithSKUs(id)
	if err != nil {
		response.Error(c, 500, "获取商品详情失败")
		return
	}
	if product == nil {
		response.Error(c, 404, "商品不存在")
		return
	}
	response.Success(c, gin.H{
		"product": product,
		"skus":    skus,
	})
}

type CartHandler struct {
	cartService *service.CartService
}

func NewCartHandler() *CartHandler {
	return &CartHandler{
		cartService: service.NewCartService(),
	}
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.GetString("user_id")
	carts, err := h.cartService.GetCartList(userID)
	if err != nil {
		response.Error(c, 500, "获取购物车失败")
		return
	}
	response.Success(c, gin.H{"list": carts})
}

func (h *CartHandler) AddCart(c *gin.Context) {
	userID := c.GetString("user_id")
	var req struct {
		ProductID string `json:"product_id" binding:"required"`
		SkuID     string `json:"sku_id"`
		Quantity  int    `json:"quantity" binding:"required;min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	err := h.cartService.AddToCart(userID, req.ProductID, req.SkuID, req.Quantity)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *CartHandler) UpdateCart(c *gin.Context) {
	cartID := c.Param("id")
	var req struct {
		Quantity int `json:"quantity" binding:"required;min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	err := h.cartService.UpdateCartQuantity(cartID, req.Quantity)
	if err != nil {
		response.Error(c, 500, "更新失败")
		return
	}
	response.Success(c, nil)
}

func (h *CartHandler) DeleteCart(c *gin.Context) {
	cartID := c.Param("id")
	err := h.cartService.DeleteCartItem(cartID)
	if err != nil {
		response.Error(c, 500, "删除失败")
		return
	}
	response.Success(c, nil)
}

type AddressHandler struct {
	addressService *service.AddressService
}

func NewAddressHandler() *AddressHandler {
	return &AddressHandler{
		addressService: service.NewAddressService(),
	}
}

func (h *AddressHandler) GetAddresses(c *gin.Context) {
	userID := c.GetString("user_id")
	addresses, err := h.addressService.GetAddresses(userID)
	if err != nil {
		response.Error(c, 500, "获取地址列表失败")
		return
	}
	response.Success(c, gin.H{"list": addresses})
}

func (h *AddressHandler) CreateAddress(c *gin.Context) {
	userID := c.GetString("user_id")
	var req struct {
		Receiver  string `json:"receiver" binding:"required"`
		Phone     string `json:"phone" binding:"required"`
		Province  string `json:"province" binding:"required"`
		City      string `json:"city" binding:"required"`
		District  string `json:"district" binding:"required"`
		Detail    string `json:"detail" binding:"required"`
		IsDefault bool   `json:"is_default"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	address, err := h.addressService.CreateAddress(userID, req.Receiver, req.Phone, req.Province, req.City, req.District, req.Detail, req.IsDefault)
	if err != nil {
		response.Error(c, 500, "创建失败")
		return
	}
	response.Success(c, gin.H{"id": address.ID})
}

func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	userID := c.GetString("user_id")
	addressID := c.Param("id")
	var req struct {
		Receiver  string `json:"receiver" binding:"required"`
		Phone     string `json:"phone" binding:"required"`
		Province  string `json:"province" binding:"required"`
		City      string `json:"city" binding:"required"`
		District  string `json:"district" binding:"required"`
		Detail    string `json:"detail" binding:"required"`
		IsDefault bool   `json:"is_default"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	err := h.addressService.UpdateAddress(addressID, userID, req.Receiver, req.Phone, req.Province, req.City, req.District, req.Detail, req.IsDefault)
	if err != nil {
		response.Error(c, 500, "更新失败")
		return
	}
	response.Success(c, nil)
}

func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	userID := c.GetString("user_id")
	addressID := c.Param("id")
	err := h.addressService.DeleteAddress(addressID, userID)
	if err != nil {
		response.Error(c, 500, "删除失败")
		return
	}
	response.Success(c, nil)
}

func (h *AddressHandler) SetDefault(c *gin.Context) {
	userID := c.GetString("user_id")
	addressID := c.Param("id")
	err := h.addressService.SetDefault(addressID, userID)
	if err != nil {
		response.Error(c, 500, "设置失败")
		return
	}
	response.Success(c, nil)
}

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) GetBanners(c *gin.Context) {
	response.Success(c, gin.H{
		"region_code": c.Query("region_code"),
		"banners":     []gin.H{},
	})
}

func (h *HomeHandler) GetFeatured(c *gin.Context) {
	response.Success(c, gin.H{
		"region_code": c.Query("region_code"),
		"featured":    nil,
	})
}

func (h *HomeHandler) GetHomeProducts(c *gin.Context) {
	response.Success(c, gin.H{
		"region_code": c.Query("region_code"),
		"list":        []gin.H{},
	})
}

func (h *HomeHandler) GetConfig(c *gin.Context) {
	response.Success(c, gin.H{})
}

func (h *HomeHandler) UpdateConfig(c *gin.Context) {
	response.Success(c, nil)
}

type SearchHandler struct{}

func NewSearchHandler() *SearchHandler {
	return &SearchHandler{}
}

func (h *SearchHandler) GetSuggest(c *gin.Context) {
	response.Success(c, gin.H{"list": []string{}})
}

func (h *SearchHandler) GetHotSearch(c *gin.Context) {
	response.Success(c, gin.H{"list": []string{"苹果", "大米", "土鸡蛋"}})
}

func (h *SearchHandler) GetHistory(c *gin.Context) {
	response.Success(c, gin.H{"list": []string{}})
}

func (h *SearchHandler) ClearHistory(c *gin.Context) {
	response.Success(c, nil)
}

type OrderHandler struct{}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	response.Success(c, gin.H{"list": []gin.H{}, "total": 0})
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	response.Success(c, gin.H{})
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	response.Success(c, gin.H{})
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	response.Success(c, nil)
}

func (h *OrderHandler) ConfirmReceive(c *gin.Context) {
	response.Success(c, nil)
}

type PaymentHandler struct{}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	response.Success(c, gin.H{})
}

func (h *PaymentHandler) PayCallback(c *gin.Context) {
	response.Success(c, gin.H{"code": "SUCCESS"})
}

func (h *PaymentHandler) GetPayStatus(c *gin.Context) {
	response.Success(c, gin.H{})
}
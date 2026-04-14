package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ethanqian1990/wechat-mall-backend/pkg/response"
	"github.com/ethanqian1990/wechat-mall-backend/internal/repository"
	"github.com/ethanqian1990/wechat-mall-backend/internal/model"
	"github.com/google/uuid"
)

// ===== 后台商品管理 =====

type AdminProductHandler struct {
	productRepo *repository.ProductRepository
	skuRepo    *repository.SKURepository
}

func NewAdminProductHandler() *AdminProductHandler {
	return &AdminProductHandler{
		productRepo: repository.NewProductRepository(),
		skuRepo:     repository.NewSKURepository(),
	}
}

func (h *AdminProductHandler) GetProducts(c *gin.Context) {
	page := 1
	pageSize := 20

	products, total, err := h.productRepo.FindAll(page, pageSize)
	if err != nil {
		response.Error(c, 500, "获取商品列表失败")
		return
	}
	response.Paginate(c, products, int(total), page, pageSize)
}

func (h *AdminProductHandler) CreateProduct(c *gin.Context) {
	var req struct {
		Name             string   `json:"name" binding:"required"`
		Description      string   `json:"description"`
		MainImage        string   `json:"main_image"`
		Images           []string `json:"images"`
		CategoryID       string   `json:"category_id"`
		Price            float64  `json:"price" binding:"required"`
		OriginalPrice    float64  `json:"original_price"`
		Stock            int      `json:"stock"`
		IsOnShelf        int      `json:"is_on_shelf"`
		ShowOnHome       int      `json:"show_on_home"`
		AvailableRegions []string `json:"available_regions"`
		Tags             []string `json:"tags"`
		PromoText        string   `json:"promo_text"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	product := &model.Product{
		ID:               "p_" + uuid.New().String()[:8],
		Name:             req.Name,
		Description:      req.Description,
		MainImage:        req.MainImage,
		Price:            req.Price,
		Stock:            req.Stock,
		IsOnShelf:        1,
		ShowOnHome:       1,
		AvailableRegions:  toJSON(req.AvailableRegions),
		Tags:             toJSON(req.Tags),
		PromoText:       req.PromoText,
	}

	err := h.productRepo.Create(product)
	if err != nil {
		response.Error(c, 500, "创建失败")
		return
	}
	response.Success(c, gin.H{"id": product.ID})
}

func (h *AdminProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name             string   `json:"name"`
		Description      string   `json:"description"`
		MainImage        string   `json:"main_image"`
		Price            float64  `json:"price"`
		Stock            int      `json:"stock"`
		IsOnShelf        int      `json:"is_on_shelf"`
		ShowOnHome       int      `json:"show_on_home"`
		AvailableRegions []string `json:"available_regions"`
		PromoText        string   `json:"promo_text"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.MainImage != "" {
		updates["main_image"] = req.MainImage
	}
	if req.Price > 0 {
		updates["price"] = req.Price
	}
	if req.Stock >= 0 {
		updates["stock"] = req.Stock
	}
	if req.IsOnShelf > 0 {
		updates["is_on_shelf"] = req.IsOnShelf
	}
	if req.ShowOnHome >= 0 {
		updates["show_on_home"] = req.ShowOnHome
	}
	if req.AvailableRegions != nil {
		updates["available_regions"] = toJSON(req.AvailableRegions)
	}
	if req.PromoText != "" {
		updates["promo_text"] = req.PromoText
	}

	err := h.productRepo.Update(id, updates)
	if err != nil {
		response.Error(c, 500, "更新失败")
		return
	}
	response.Success(c, nil)
}

func (h *AdminProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	err := h.productRepo.Delete(id)
	if err != nil {
		response.Error(c, 500, "删除失败")
		return
	}
	response.Success(c, nil)
}

func (h *AdminProductHandler) SetShelf(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		OnShelf int `json:"on_shelf"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	err := h.productRepo.UpdateOnShelf(id, req.OnShelf)
	if err != nil {
		response.Error(c, 500, "设置失败")
		return
	}
	response.Success(c, nil)
}

func (h *AdminProductHandler) SetHome(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		ShowOnHome int `json:"show_on_home"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	err := h.productRepo.UpdateShowOnHome(id, req.ShowOnHome)
	if err != nil {
		response.Error(c, 500, "设置失败")
		return
	}
	response.Success(c, nil)
}

// ===== 后台订单管理 =====

type AdminOrderHandler struct {
	orderRepo     *repository.OrderRepository
	orderItemRepo *repository.OrderItemRepository
}

func NewAdminOrderHandler() *AdminOrderHandler {
	return &AdminOrderHandler{
		orderRepo:     repository.NewOrderRepository(),
		orderItemRepo: repository.NewOrderItemRepository(),
	}
}

func (h *AdminOrderHandler) GetOrders(c *gin.Context) {
	page := 1
	pageSize := 20
	status := c.Query("status")

	orders, total, err := h.orderRepo.FindAll(status, page, pageSize)
	if err != nil {
		response.Error(c, 500, "获取订单列表失败")
		return
	}
	response.Paginate(c, orders, int(total), page, pageSize)
}

func (h *AdminOrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.orderRepo.FindByID(id)
	if err != nil || order == nil {
		response.Error(c, 404, "订单不存在")
		return
	}
	items, _ := h.orderItemRepo.FindByOrderID(id)
	response.Success(c, gin.H{
		"order": order,
		"items": items,
	})
}

func (h *AdminOrderHandler) ShipOrder(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		ExpressCompany string `json:"express_company"`
		ExpressNo     string `json:"express_no"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	err := h.orderRepo.UpdateStatus(id, "SHIPPED")
	if err != nil {
		response.Error(c, 500, "发货失败")
		return
	}
	response.Success(c, nil)
}

func (h *AdminOrderHandler) UpdatePrice(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		FinalAmount float64 `json:"final_amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	err := h.orderRepo.UpdateOrder(id, map[string]interface{}{
		"final_amount": req.FinalAmount,
	})
	if err != nil {
		response.Error(c, 500, "修改价格失败")
		return
	}
	response.Success(c, nil)
}

// ===== 后台用户管理 =====

type AdminUserHandler struct {
	userRepo *repository.UserRepository
}

func NewAdminUserHandler() *AdminUserHandler {
	return &AdminUserHandler{
		userRepo: repository.NewUserRepository(),
	}
}

func (h *AdminUserHandler) GetUsers(c *gin.Context) {
	page := 1
	pageSize := 20

	users, total, err := h.userRepo.FindAll(page, pageSize)
	if err != nil {
		response.Error(c, 500, "获取用户列表失败")
		return
	}
	response.Paginate(c, users, int(total), page, pageSize)
}

func (h *AdminUserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userRepo.FindByID(id)
	if err != nil || user == nil {
		response.Error(c, 404, "用户不存在")
		return
	}
	response.Success(c, user)
}

// ===== 后台区域管理 =====

type AdminRegionHandler struct {
	regionRepo *repository.RegionRepository
}

func NewAdminRegionHandler() *AdminRegionHandler {
	return &AdminRegionHandler{
		regionRepo: repository.NewRegionRepository(),
	}
}

func (h *AdminRegionHandler) GetRegions(c *gin.Context) {
	regions, err := h.regionRepo.FindAll()
	if err != nil {
		response.Error(c, 500, "获取区域列表失败")
		return
	}
	response.Success(c, gin.H{"list": regions})
}

func (h *AdminRegionHandler) CreateRegion(c *gin.Context) {
	var req struct {
		Code  string `json:"code" binding:"required"`
		Name  string `json:"name" binding:"required"`
		Sort  int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	region := &model.Region{
		Code:      req.Code,
		Name:      req.Name,
		SortOrder: req.Sort,
		IsEnabled: 1,
	}

	err := h.regionRepo.Create(region)
	if err != nil {
		response.Error(c, 500, "创建失败")
		return
	}
	response.Success(c, nil)
}

func (h *AdminRegionHandler) UpdateRegion(c *gin.Context) {
	code := c.Param("code")
	var req struct {
		Name  string `json:"name"`
		Sort  int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["region_name"] = req.Name
	}
	if req.Sort >= 0 {
		updates["sort_order"] = req.Sort
	}

	err := h.regionRepo.Update(code, updates)
	if err != nil {
		response.Error(c, 500, "更新失败")
		return
	}
	response.Success(c, nil)
}

func (h *AdminRegionHandler) DeleteRegion(c *gin.Context) {
	code := c.Param("code")
	err := h.regionRepo.Delete(code)
	if err != nil {
		response.Error(c, 500, "删除失败")
		return
	}
	response.Success(c, nil)
}

// ===== 后台分类管理 =====

type AdminCategoryHandler struct {
	categoryRepo *repository.CategoryRepository
}

func NewAdminCategoryHandler() *AdminCategoryHandler {
	return &AdminCategoryHandler{
		categoryRepo: repository.NewCategoryRepository(),
	}
}

func (h *AdminCategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.categoryRepo.FindAll()
	if err != nil {
		response.Error(c, 500, "获取分类列表失败")
		return
	}
	response.Success(c, gin.H{"list": categories})
}

func (h *AdminCategoryHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		ParentID string `json:"parent_id"`
		Sort     int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	category := &model.Category{
		ID:        "c_" + uuid.New().String()[:8],
		Name:      req.Name,
		ParentID:  req.ParentID,
		SortOrder: req.Sort,
		IsEnabled: 1,
		Level:     1,
	}

	err := h.categoryRepo.Create(category)
	if err != nil {
		response.Error(c, 500, "创建失败")
		return
	}
	response.Success(c, nil)
}

func (h *AdminCategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name string `json:"name"`
		Sort int   `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Sort >= 0 {
		updates["sort_order"] = req.Sort
	}

	err := h.categoryRepo.Update(id, updates)
	if err != nil {
		response.Error(c, 500, "更新失败")
		return
	}
	response.Success(c, nil)
}

func (h *AdminCategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	err := h.categoryRepo.Delete(id)
	if err != nil {
		response.Error(c, 500, "删除失败")
		return
	}
	response.Success(c, nil)
}

// 辅助函数
func toJSON(v interface{}) string {
	if v == nil {
		return "[]"
	}
	return fmt.Sprintf("%v", v)
}
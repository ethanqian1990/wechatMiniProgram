package handler

import (
	github.com/gin-gonic/gin
	github.com/ethanqian1990/wechat-mall-backend/pkg/response
	github.com/ethanqian1990/wechat-mall-backend/internal/service
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		orderService: service.NewOrderService(),
	}
}

// GetOrders 获取订单列表
func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID := c.GetString('user_id')
	status := c.Query('status')
	page := 1
	pageSize := 20

	orders, total, err := h.orderService.GetOrders(userID, status, page, pageSize)
	if err != nil {
		response.Error(c, 500, '获取订单列表失败')
		return
	}

	var result []gin.H
	for _, o := range orders {
		result = append(result, gin.H{
			'id':           o.ID,
			'order_no':     o.OrderNo,
			'status':       o.Status,
			'final_amount': o.FinalAmount,
			'created_at':   o.CreatedAt,
		})
	}

	response.Paginate(c, result, int(total), page, pageSize)
}

// GetOrder 获取订单详情
func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID := c.GetString('user_id')
	orderID := c.Param('id')

	order, items, err := h.orderService.GetOrderByID(orderID, userID)
	if err != nil {
		response.Error(c, 500, '获取订单详情失败')
		return
	}
	if order == nil {
		response.Error(c, 404, '订单不存在')
		return
	}

	var itemResult []gin.H
	for _, item := range items {
		itemResult = append(itemResult, gin.H{
			'product_name':   item.ProductName,
			'product_image':  item.ProductImage,
			'sku_name':       item.SkuName,
			'price':          item.Price,
			'quantity':       item.Quantity,
			'subtotal':       item.Subtotal,
		})
	}

	response.Success(c, gin.H{
		'id':           order.ID,
		'order_no':     order.OrderNo,
		'status':       order.Status,
		'total_amount': order.TotalAmount,
		'final_amount': order.FinalAmount,
		'receiver':     order.Receiver,
		'phone':        order.Phone,
		'address':      order.Province + order.City + order.District + order.Detail,
		'items':        itemResult,
		'created_at':   order.CreatedAt,
	})
}

// CreateOrder 创建订单
type CreateOrderRequest struct {
	RegionCode       string `json: region_code`
	AddressID        string `json: address_id binding:required`
	ClientOrderToken string `json: client_order_token`
	Items            []struct {
		ProductID string `json: product_id`
		SkuID     string `json: sku_id`
		Quantity  int    `json: quantity`
	} `json: items binding:required`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetString('user_id')

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, '参数错误')
		return
	}

	params := service.CreateOrderParams{
		RegionCode:       req.RegionCode,
		AddressID:        req.AddressID,
		ClientOrderToken: req.ClientOrderToken,
		Items:            req.Items,
	}

	order, err := h.orderService.CreateOrder(userID, params)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, gin.H{
		'id':          order.ID,
		'order_no':    order.OrderNo,
		'status':      order.Status,
		'final_amount': order.FinalAmount,
	})
}

// CancelOrder 取消订单
type CancelOrderRequest struct {
	Reason string `json: reason`
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID := c.GetString('user_id')
	orderID := c.Param('id')

	var req CancelOrderRequest
	c.ShouldBindJSON(&req)

	err := h.orderService.CancelOrder(orderID, userID, req.Reason)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

// ConfirmReceive 确认收货
func (h *OrderHandler) ConfirmReceive(c *gin.Context) {
	userID := c.GetString('user_id')
	orderID := c.Param('id')

	err := h.orderService.ConfirmReceive(orderID, userID)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}
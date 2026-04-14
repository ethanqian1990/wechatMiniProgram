package handler

import (
	github.com/gin-gonic/gin
	github.com/ethanqian1990/wechat-mall-backend/pkg/response
	github.com/ethanqian1990/wechat-mall-backend/internal/service
	github.com/ethanqian1990/wechat-mall-backend/internal/repository
	github.com/google/uuid
)

type PaymentHandler struct {
	orderService   *service.OrderService
	orderRepo      *repo.OrderRepository
	stockRepo      *repo.StockReservationRepository
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{
		orderService: service.NewOrderService(),
		orderRepo:     repo.NewOrderRepository(),
		stockRepo:     repo.NewStockReservationRepository(),
	}
}

// CreatePayment 创建支付
type CreatePaymentRequest struct {
	OrderID string `json: order_id binding:required`
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	userID := c.GetString('user_id')

	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, '参数错误')
		return
	}

	order, err := h.orderRepo.FindByIDAndUserID(req.OrderID, userID)
	if err != nil || order == nil {
		response.Error(c, 404, '订单不存在')
		return
	}

	if order.Status != 'PENDING_PAY' {
		response.Error(c, 400, '订单状态不允许支付')
		return
	}

	// TODO: 调用微信统一下单接口
	// 模拟返回支付参数
	payParams := gin.H{
		'appId':      'wx1234567890',
		'timeStamp':  '1234567890',
		'nonceStr':   uuid.New().String()[:16],
		'package':    'prepay_id=wx1234567890',
		'signType':   'RSA',
		'paySign':    'mock_signature',
	}

	response.Success(c, payParams)
}

// PayCallback 支付回调
func (h *PaymentHandler) PayCallback(c *gin.Context) {
	// TODO: 验签并处理支付回调
	// 模拟处理成功
	
	// 1. 获取订单号
	// 2. 更新订单状态为 PENDING_SHIP
	// 3. 将冻结库存转为已扣减
	
	response.Success(c, gin.H{'code': 'SUCCESS', 'message': '成功'})
}

// GetPayStatus 获取支付状态
func (h *PaymentHandler) GetPayStatus(c *gin.Context) {
	userID := c.GetString('user_id')
	orderID := c.Param('order_id')

	order, err := h.orderRepo.FindByIDAndUserID(orderID, userID)
	if err != nil || order == nil {
		response.Error(c, 404, '订单不存在')
		return
	}

	response.Success(c, gin.H{
		'order_id': order.ID,
		'status':   order.Status,
	})
}
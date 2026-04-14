package service

import (
	"fmt"

	// wechatmall "github.com/ethanqian1990/wechat-mall-backend/internal/model"
	wechatrepo "github.com/ethanqian1990/wechat-mall-backend/internal/repository"
	"github.com/ethanqian1990/wechat-mall-backend/pkg/utils"
	// "github.com/google/uuid"
)

type PaymentService struct {
	orderRepo    *wechatrepo.OrderRepository
	stockRepo   *wechatrepo.StockReservationRepository
	paymentUtil *utils.PaymentUtil
}

func NewPaymentService(paymentUtil *utils.PaymentUtil) *PaymentService {
	return &PaymentService{
		orderRepo:  wechatrepo.NewOrderRepository(),
		stockRepo: wechatrepo.NewStockReservationRepository(),
		paymentUtil: paymentUtil,
	}
}

// CreatePayment 创建支付
func (s *PaymentService) CreatePayment(userID, orderID string) (map[string]string, error) {
	// 1. 获取订单
	order, err := s.orderRepo.FindByIDAndUserID(orderID, userID)
	if err != nil || order == nil {
		return nil, fmt.Errorf("订单不存在")
	}

	if order.Status != "PENDING_PAY" {
		return nil, fmt.Errorf("订单状态不允许支付")
	}

	// 2. 调用微信支付统一下单
	amount := int64(order.FinalAmount * 100) // 微信支付单位是分
	description := fmt.Sprintf("订单支付-%s", order.OrderNo)

	result, err := s.paymentUtil.CreateUnifiedOrder(order.OrderNo, amount, description)
	if err != nil {
		return nil, fmt.Errorf("创建支付失败: %v", err)
	}

	if result.ErrCode != "" {
		return nil, fmt.Errorf("微信支付错误: %s - %s", result.ErrCode, result.ErrMsg)
	}

	// 3. 返回支付参数
	payParams := s.paymentUtil.GetJSAPIPayParams(result.PrepayID)
	return payParams, nil
}

// PayCallback 支付回调
func (s *PaymentService) PayCallback(params map[string]string) error {
	// 1. 验证签名
	if !s.paymentUtil.VerifyCallback(params) {
		return fmt.Errorf("签名验证失败")
	}

	// 2. 检查返回码
	if params["return_code"] != "SUCCESS" {
		return fmt.Errorf("微信返回失败: %s", params["return_msg"])
	}

	// 3. 获取订单号和支付状态
	orderNo := params["out_trade_no"]
	transactionID := params["transaction_id"]

	// 4. 查找订单
	orders, err := s.orderRepo.FindByOrderNo(orderNo)
	if err != nil || len(orders) == 0 {
		return fmt.Errorf("订单不存在: %s", orderNo)
	}

	order := orders[0]

	// 5. 检查订单状态，防止重复处理
	if order.Status != "PENDING_PAY" {
		return nil // 已处理过，直接返回成功
	}

	// 6. 更新订单状态
	if err := s.orderRepo.UpdateStatus(order.ID, "PENDING_SHIP"); err != nil {
		return fmt.Errorf("更新订单状态失败: %v", err)
	}

	// 7. 更新订单交易号
	s.orderRepo.UpdateOrder(order.ID, map[string]interface{}{
		"transaction_id": transactionID,
		"pay_at":        "NOW()",
	})

	// 8. 确认库存冻结（从FROZEN转为已扣减）
	s.stockRepo.ConfirmByOrderID(order.ID)

	return nil
}

// GetPayStatus 获取支付状态
func (s *PaymentService) GetPayStatus(orderID, userID string) (string, error) {
	order, err := s.orderRepo.FindByIDAndUserID(orderID, userID)
	if err != nil || order == nil {
		return "", fmt.Errorf("订单不存在")
	}
	return order.Status, nil
}

// ConfirmPayment 确认支付（内部使用）
func (s *PaymentService) ConfirmPayment(orderID string) error {
	return s.orderRepo.UpdateStatus(orderID, "PENDING_SHIP")
}

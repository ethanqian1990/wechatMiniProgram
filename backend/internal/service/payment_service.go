package service

import (
	"fmt"
	"time"

	"github.com/ethanqian1990/wechat-mall-backend/internal/model"
	"github.com/ethanqian1990/wechat-mall-backend/internal/repository"
	"gorm.io/gorm"
)

type PaymentService struct {
	orderRepo *repository.OrderRepository
	stockRepo *repository.StockReservationRepository
}

func NewPaymentService() *PaymentService {
	return &PaymentService{
		orderRepo: repository.NewOrderRepository(),
		stockRepo: repository.NewStockReservationRepository(),
	}
}

type CreatePaymentResp struct {
	PayParams map[string]string `json:"pay_params"`
}

func (s *PaymentService) CreatePayment(userID, orderID string) (*CreatePaymentResp, error) {
	order, err := s.orderRepo.FindByIDAndUser(orderID, userID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("订单不存在")
	}
	if order.Status != "PENDING_PAY" {
		return nil, fmt.Errorf("当前状态不可支付")
	}
	// MVP：返回 mock 的微信支付参数
	return &CreatePaymentResp{
		PayParams: map[string]string{
			"mock":     "1",
			"order_id": order.ID,
			"amount":   fmt.Sprintf("%.2f", order.FinalAmount),
		},
	}, nil
}

func (s *PaymentService) MarkPaid(orderID string) error {
	now := time.Now()
	return repository.DB.Transaction(func(tx *gorm.DB) error {
		var order model.Order
		if err := tx.First(&order, "id = ?", orderID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("订单不存在")
			}
			return err
		}

		// 幂等：已支付/已发货/已完成等重复回调直接返回成功
		if order.Status != "PENDING_PAY" {
			return nil
		}

		if err := tx.Model(&model.Order{}).
			Where("id = ? AND status = ?", orderID, "PENDING_PAY").
			Updates(map[string]interface{}{
				"status":     "PENDING_SHIP",
				"pay_at":     &now,
				"updated_at": now,
			}).Error; err != nil {
			return err
		}

		// 冻结库存转已消耗
		return s.stockRepo.MarkConsumedByOrderIDTx(tx, orderID)
	})
}

func (s *PaymentService) MarkPaidWithTransaction(orderID, transactionID string) error {
	if transactionID == "" {
		return s.MarkPaid(orderID)
	}

	now := time.Now()
	return repository.DB.Transaction(func(tx *gorm.DB) error {
		var order model.Order
		if err := tx.First(&order, "id = ?", orderID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("订单不存在")
			}
			return err
		}

		// 幂等：已处理过
		if order.PaidTransactionID != "" {
			if order.PaidTransactionID == transactionID {
				return nil
			}
			return fmt.Errorf("交易号冲突")
		}
		if order.Status != "PENDING_PAY" {
			// 已经不是待支付，但没有记录交易号：补不上时直接当作已处理
			return nil
		}

		// 原子写入：status=待支付 且 paid_transaction_id 为空时写入交易号并更新状态
		err := tx.Model(&model.Order{}).
			Where("id = ? AND status = ? AND (paid_transaction_id = '' OR paid_transaction_id IS NULL)", orderID, "PENDING_PAY").
			Updates(map[string]interface{}{
				"status":              "PENDING_SHIP",
				"pay_at":              &now,
				"paid_transaction_id": transactionID,
				"updated_at":          now,
			}).Error
		if err != nil {
			return err
		}
		return s.stockRepo.MarkConsumedByOrderIDTx(tx, orderID)
	})
}

func (s *PaymentService) GetPayStatus(userID, orderID string) (string, error) {
	order, err := s.orderRepo.FindByIDAndUser(orderID, userID)
	if err != nil {
		return "", err
	}
	if order == nil {
		return "", fmt.Errorf("订单不存在")
	}
	return order.Status, nil
}

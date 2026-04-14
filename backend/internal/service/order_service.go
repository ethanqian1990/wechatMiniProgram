package service

import (
	"fmt"
	"time"

	"github.com/ethanqian1990/wechat-mall-backend/internal/model"
	"github.com/ethanqian1990/wechat-mall-backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo   *repository.OrderRepository
	orderItemRepo *repository.OrderItemRepository
	addressRepo *repository.AddressRepository
	skuRepo     *repository.SKURepository
	stockRepo   *repository.StockReservationRepository
}

func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo:     repository.NewOrderRepository(),
		orderItemRepo: repository.NewOrderItemRepository(),
		addressRepo:   repository.NewAddressRepository(),
		skuRepo:       repository.NewSKURepository(),
		stockRepo:     repository.NewStockReservationRepository(),
	}
}

type CreateOrderItemReq struct {
	ProductID string `json:"product_id"`
	SkuID     string `json:"sku_id"`
	Quantity  int    `json:"quantity"`
}

type CreateOrderReq struct {
	RegionCode       string               `json:"region_code"`
	AddressID        string               `json:"address_id"`
	ClientOrderToken string               `json:"client_order_token"`
	Items            []CreateOrderItemReq `json:"items"`
}

func (s *OrderService) GetOrders(userID, status string, page, pageSize int) ([]model.Order, int64, error) {
	return s.orderRepo.FindByUser(userID, status, page, pageSize)
}

func (s *OrderService) GetOrderDetail(userID, orderID string) (*model.Order, []model.OrderItem, error) {
	order, err := s.orderRepo.FindByIDAndUser(orderID, userID)
	if err != nil {
		return nil, nil, err
	}
	if order == nil {
		return nil, nil, nil
	}
	items, err := s.orderItemRepo.FindByOrderID(order.ID)
	return order, items, err
}

func (s *OrderService) CreateOrder(userID string, req CreateOrderReq) (*model.Order, error) {
	if req.ClientOrderToken == "" {
		return nil, fmt.Errorf("client_order_token 必填")
	}
	if req.RegionCode == "" {
		return nil, fmt.Errorf("region_code 必填")
	}
	if req.AddressID == "" {
		return nil, fmt.Errorf("address_id 必填")
	}
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("items 不能为空")
	}

	// 幂等：同一 user_id + client_order_token 只允许创建一次
	existing, err := s.orderRepo.FindByUserAndClientToken(userID, req.ClientOrderToken)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	addr, err := s.addressRepo.FindByIDAndUserID(req.AddressID, userID)
	if err != nil {
		return nil, err
	}
	if addr == nil {
		return nil, fmt.Errorf("地址不存在")
	}

	now := time.Now()
	orderID := "ord_" + uuid.New().String()[:12]
	orderNo := fmt.Sprintf("NO%s", now.Format("20060102150405"))

	var created *model.Order
	err = repository.DB.Transaction(func(tx *gorm.DB) error {
		var total float64
		var items []model.OrderItem

		for _, it := range req.Items {
			if it.SkuID == "" {
				return fmt.Errorf("sku_id 必填")
			}
			if it.Quantity <= 0 {
				return fmt.Errorf("quantity 必须 > 0")
			}
			sku, err := s.skuRepo.FindByID(it.SkuID)
			if err != nil {
				return err
			}
			if sku == nil || sku.ProductID != it.ProductID {
				return fmt.Errorf("SKU 不存在")
			}

			// 冻结库存：直接扣减 sku.stock，并写入冻结记录
			res := tx.Model(&model.ProductSKU{}).
				Where("id = ? AND stock >= ?", it.SkuID, it.Quantity).
				Update("stock", gorm.Expr("stock - ?", it.Quantity))
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return fmt.Errorf("库存不足")
			}

			subtotal := sku.Price * float64(it.Quantity)
			total += subtotal

			skuID := it.SkuID
			items = append(items, model.OrderItem{
				ID:           "oi_" + uuid.New().String()[:12],
				OrderID:      orderID,
				ProductID:    it.ProductID,
				SkuID:        &skuID,
				ProductName:  "", // MVP：后续补齐商品快照
				ProductImage: "",
				SkuName:      sku.SkuName,
				Price:        sku.Price,
				Quantity:     it.Quantity,
				Subtotal:     subtotal,
				CreatedAt:    now,
			})

			if err := tx.Create(&model.StockReservation{
				ID:        "sr_" + uuid.New().String()[:12],
				OrderID:   orderID,
				UserID:    userID,
				ProductID: it.ProductID,
				SkuID:     it.SkuID,
				Quantity:  it.Quantity,
				Status:    "FROZEN",
				ExpireAt:  now.Add(15 * time.Minute),
				CreatedAt: now,
				UpdatedAt: now,
			}).Error; err != nil {
				return err
			}
		}

		order := &model.Order{
			ID:               orderID,
			OrderNo:          orderNo,
			UserID:           userID,
			TotalAmount:      total,
			DiscountAmount:   0,
			FreightAmount:    0,
			FinalAmount:      total,
			Status:           "PENDING_PAY",
			Receiver:         addr.Receiver,
			Phone:            addr.Phone,
			Province:         addr.Province,
			City:             addr.City,
			District:         addr.District,
			Detail:           addr.Detail,
			ClientOrderToken: req.ClientOrderToken,
			RegionCode:       req.RegionCode,
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		created = order
		return nil
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *OrderService) CancelOrder(userID, orderID, reason string) error {
	order, err := s.orderRepo.FindByIDAndUser(orderID, userID)
	if err != nil {
		return err
	}
	if order == nil {
		return fmt.Errorf("订单不存在")
	}
	if order.Status != "PENDING_PAY" {
		return fmt.Errorf("当前状态不可取消")
	}
	now := time.Now()
	return repository.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Order{}).Where("id = ?", order.ID).
			Updates(map[string]interface{}{
				"status":        "CANCELED",
				"cancel_at":     &now,
				"cancel_reason": reason,
				"updated_at":    now,
			}).Error; err != nil {
			return err
		}
		// 释放冻结库存（事务内，避免重复返还）
		return s.stockRepo.ReleaseByOrderIDTx(tx, order.ID)
	})
}

func (s *OrderService) ConfirmReceive(userID, orderID string) error {
	order, err := s.orderRepo.FindByIDAndUser(orderID, userID)
	if err != nil {
		return err
	}
	if order == nil {
		return fmt.Errorf("订单不存在")
	}
	if order.Status != "SHIPPED" {
		return fmt.Errorf("当前状态不可确认收货")
	}
	now := time.Now()
	return repository.DB.Model(&model.Order{}).Where("id = ?", order.ID).
		Updates(map[string]interface{}{
			"status":     "FINISHED",
			"receive_at": &now,
			"updated_at": now,
		}).Error
}


package service

import (
	"fmt"
	"time"

	"github.com/ethanqian1990/wechat-mall-backend/internal/config"
	"github.com/ethanqian1990/wechat-mall-backend/internal/model"
	wechatrepo "github.com/ethanqian1990/wechat-mall-backend/internal/repository"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type CronService struct {
	c        *cron.Cron
	orderRepo *wechatrepo.OrderRepository
	skuRepo   *wechatrepo.SKURepository
	stockRepo *wechatrepo.StockReservationRepository

	autoCancelMinutes int
	autoReceiveDays   int
}

func NewCronService(cfg *config.Config) *CronService {
	autoCancel := cfg.Order.AutoCancelMinutes
	if autoCancel <= 0 {
		autoCancel = 15
	}
	autoReceive := cfg.Order.AutoReceiveDays
	if autoReceive <= 0 {
		autoReceive = 7
	}

	return &CronService{
		c:        cron.New(),
		orderRepo: wechatrepo.NewOrderRepository(),
		skuRepo:   wechatrepo.NewSKURepository(),
		stockRepo: wechatrepo.NewStockReservationRepository(),
		autoCancelMinutes: autoCancel,
		autoReceiveDays:   autoReceive,
	}
}

func (s *CronService) Start() {
	// 订单超时取消 - 每分钟执行
	s.c.AddFunc("*/1 * * * *", func() {
		s.CancelExpiredOrders()
	})

	// 自动确认收货 - 每小时执行
	s.c.AddFunc("0 * * * *", func() {
		s.AutoConfirmReceive()
	})

	// 库存释放 - 每5分钟执行
	s.c.AddFunc("*/5 * * * *", func() {
		s.ReleaseExpiredStock()
	})

	s.c.Start()
	fmt.Println("定时任务已启动")
}

// CancelExpiredOrders 取消超时订单（15分钟未支付）
func (s *CronService) CancelExpiredOrders() {
	expireTime := time.Now().Add(-time.Duration(s.autoCancelMinutes) * time.Minute)
	
	orders, err := s.orderRepo.FindExpiredOrders("PENDING_PAY", expireTime)
	if err != nil {
		fmt.Printf("查询超时订单失败: %v\n", err)
		return
	}

	for _, order := range orders {
		// 事务化：先更新订单状态，再释放冻结库存（避免并发/重入导致重复返还）
		err := wechatrepo.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&model.Order{}).Where("id = ? AND status = ?", order.ID, "PENDING_PAY").
				Update("status", "CANCELED").Error; err != nil {
				return err
			}
			return s.stockRepo.ReleaseByOrderIDTx(tx, order.ID)
		})
		if err != nil {
			fmt.Printf("超时取消失败: %v\n", err)
		} else {
			fmt.Printf("订单 %s 已超时取消\n", order.OrderNo)
		}
	}
}

// AutoConfirmReceive 自动确认收货（发货7天后）
func (s *CronService) AutoConfirmReceive() {
	expireTime := time.Now().Add(-time.Duration(s.autoReceiveDays) * 24 * time.Hour)
	
	orders, err := s.orderRepo.FindExpiredOrders("SHIPPED", expireTime)
	if err != nil {
		fmt.Printf("查询待收货订单失败: %v\n", err)
		return
	}

	for _, order := range orders {
		if err := s.orderRepo.UpdateStatus(order.ID, "FINISHED"); err != nil {
			fmt.Printf("确认收货失败: %v\n", err)
		} else {
			fmt.Printf("订单 %s 已自动确认收货\n", order.OrderNo)
		}
	}
}

// ReleaseExpiredStock 释放过期的库存冻结记录
func (s *CronService) ReleaseExpiredStock() {
	now := time.Now()
	
	reservations, err := s.stockRepo.FindExpired(now)
	if err != nil {
		fmt.Printf("查询过期冻结库存失败: %v\n", err)
		return
	}

	for _, res := range reservations {
		// 事务化：原子将 FROZEN -> RELEASED 成功后再返还库存，避免重复加回
		err := wechatrepo.DB.Transaction(func(tx *gorm.DB) error {
			upd := tx.Model(&model.StockReservation{}).
				Where("id = ? AND status = ?", res.ID, "FROZEN").
				Update("status", "RELEASED")
			if upd.Error != nil {
				return upd.Error
			}
			if upd.RowsAffected == 0 {
				return nil
			}
			return tx.Model(&model.ProductSKU{}).Where("id = ?", res.SkuID).
				Update("stock", gorm.Expr("stock + ?", res.Quantity)).Error
		})
		if err != nil {
			fmt.Printf("释放过期冻结库存失败: %v\n", err)
		} else {
			fmt.Printf("已释放过期冻结库存: %s, 数量: %d\n", res.SkuID, res.Quantity)
		}
	}
}

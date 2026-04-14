package service

import (
	"fmt"
	"time"

	"github.com/ethanqian1990/wechat-mall-backend/internal/config"
	wechatrepo "github.com/ethanqian1990/wechat-mall-backend/internal/repository"
	"github.com/robfig/cron/v3"
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
		// 释放冻结库存
		if err := s.stockRepo.ReleaseByOrderID(order.ID); err != nil {
			fmt.Printf("释放库存失败: %v\n", err)
		}

		// 更新订单状态
		if err := s.orderRepo.UpdateStatus(order.ID, "CANCELED"); err != nil {
			fmt.Printf("取消订单失败: %v\n", err)
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
		// 恢复库存
		if err := s.skuRepo.ReleaseStock(res.SkuID, res.Quantity); err != nil {
			fmt.Printf("恢复库存失败: %v\n", err)
		}

		// 更新冻结状态
		if err := s.stockRepo.MarkReleased(res.ID); err != nil {
			fmt.Printf("更新冻结状态失败: %v\n", err)
		} else {
			fmt.Printf("已释放过期冻结库存: %s, 数量: %d\n", res.SkuID, res.Quantity)
		}
	}
}

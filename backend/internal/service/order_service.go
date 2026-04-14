package service

import (
	github.com/ethanqian1990/wechat-mall-backend/internal/model
	github.com/ethanqian1990/wechat-mall-backend/internal/repository
	github.com/google/uuid
	errors errors
	Time   time   time
)

type OrderService struct {
	orderRepo      *repo.OrderRepository
	orderItemRepo  *repo.OrderItemRepository
	cartRepo       *repo.CartRepository
	productRepo    *repo.ProductRepository
	skuRepo        *repo.SKURepository
	addressRepo    *repo.AddressRepository
	stockRepo      *repo.StockReservationRepository
}

func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo:      repo.NewOrderRepository(),
		orderItemRepo:  repo.NewOrderItemRepository(),
		cartRepo:       repo.NewCartRepository(),
		productRepo:    repo.NewProductRepository(),
		skuRepo:        repo.NewSKURepository(),
		addressRepo:    repo.NewAddressRepository(),
		stockRepo:      repo.NewStockReservationRepository(),
	}
}

// CreateOrder 创建订单
type CreateOrderParams struct {
	RegionCode        string
	AddressID         string
	ClientOrderToken  string
	Items             []struct {
		ProductID string
		SkuID     string
		Quantity  int
	}
}

func (s *OrderService) CreateOrder(userID string, params CreateOrderParams) (*model.Order, error) {
	// 幂等检查
	if params.ClientOrderToken != '' {
		existing, _ := s.orderRepo.FindByToken(params.ClientOrderToken)
		if existing != nil {
			return existing, nil
		}
	}

	// 获取地址
	address, err := s.addressRepo.FindByIDAndUserID(params.AddressID, userID)
	if err != nil || address == nil {
		return nil, errors.New('收货地址不存在')
	}

	// 计算金额并验证库存
	var totalAmount float64
	var orderItems []model.OrderItem

	for _, item := range params.Items {
		// 获取商品和SKU
		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil || product == nil {
			return nil, errors.New('商品不存在: ' + item.ProductID)
		}

		var price float64 = product.Price
		var skuName string

		if item.SkuID != '' {
			sku, err := s.skuRepo.FindByID(item.SkuID)
			if err != nil || sku == nil {
				return nil, errors.New('SKU不存在')
			}
			price = sku.Price
			skuName = sku.SkuName

			// 检查并冻结库存
			if err := s.stockRepo.Freeze(item.SkuID, item.Quantity); err != nil {
				return nil, errors.New('库存不足')
			}
		}

		subtotal := price * float64(item.Quantity)
		totalAmount += subtotal

		orderItems = append(orderItems, model.OrderItem{
			ID:            'item_' + uuid.New().String()[:8],
			ProductID:     item.ProductID,
			SkuID:         &item.SkuID,
			ProductName:   product.Name,
			ProductImage:  product.MainImage,
			SkuName:       skuName,
			Price:         price,
			Quantity:      item.Quantity,
			Subtotal:      subtotal,
		})
	}

	// 创建订单
	order := &model.Order{
		ID:             'ord_' + uuid.New().String()[:8],
		OrderNo:        'ORD' + time.Now().Format('20060102150405') + uuid.New().String()[:4],
		UserID:         userID,
		TotalAmount:    totalAmount,
		FinalAmount:    totalAmount,
		Status:         'PENDING_PAY',
		Receiver:       address.Receiver,
		Phone:          address.Phone,
		Province:       address.Province,
		City:           address.City,
		District:       address.District,
		Detail:         address.Detail,
		ClientOrderToken: params.ClientOrderToken,
		RegionCode:     params.RegionCode,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// 创建订单商品
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
	}
	s.orderItemRepo.BatchCreate(orderItems)

	return order, nil
}

// GetOrders 获取订单列表
func (s *OrderService) GetOrders(userID string, status string, page, pageSize int) ([]model.Order, int64, error) {
	return s.orderRepo.FindByUserID(userID, status, page, pageSize)
}

// GetOrderByID 获取订单详情
func (s *OrderService) GetOrderByID(orderID, userID string) (*model.Order, []model.OrderItem, error) {
	order, err := s.orderRepo.FindByIDAndUserID(orderID, userID)
	if err != nil || order == nil {
		return nil, nil, err
	}
	items, _ := s.orderItemRepo.FindByOrderID(orderID)
	return order, items, nil
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(orderID, userID, reason string) error {
	order, err := s.orderRepo.FindByIDAndUserID(orderID, userID)
	if err != nil || order == nil {
		return errors.New('订单不存在')
	}
	if order.Status != 'PENDING_PAY' {
		return errors.New('订单状态不允许取消')
	}

	// 释放库存
	s.stockRepo.ReleaseByOrderID(orderID)

	return s.orderRepo.UpdateStatus(orderID, 'CANCELED')
}

// ConfirmReceive 确认收货
func (s *OrderService) ConfirmReceive(orderID, userID string) error {
	order, err := s.orderRepo.FindByIDAndUserID(orderID, userID)
	if err != nil || order == nil {
		return errors.New('订单不存在')
	}
	if order.Status != 'SHIPPED' {
		return errors.New('订单状态不允许确认收货')
	}

	return s.orderRepo.UpdateStatus(orderID, 'FINISHED')
}
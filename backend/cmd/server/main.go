package main

import (
	fmt   fmt
	github.com/gin-gonic/gin
	github.com/ethanqian1990/wechat-mall-backend/internal/config
	github.com/ethanqian1990/wechat-mall-backend/internal/middleware
)

func main() {
	// 加载配置
	cfg, err := config.Load(./configs/config.yaml)
	if err != nil {
		fmt.Printf(加载配置失败: %v\n, err)
		return
	}

	// 创建 Gin 引擎
	r := gin.Default()
	
	// 添加中间件
	r.Use(middleware.CORS())

	// 健康检查
	r.GET(/health, func(c *gin.Context) {
		c.JSON(200, gin.H{status: ok})
	})

	// TODO: 添加路由
	// userHandler := handler.NewUserHandler()
	// productHandler := handler.NewProductHandler()
	// ...

	api := r.Group(/api/v1)
	{
		// 公开接口
		api.GET(/regions, regionHandler.GetRegions)
		api.GET(/categories, categoryHandler.GetCategories)
		api.GET(/products, productHandler.GetProducts)
		api.GET(/products/:id, productHandler.GetProduct)
		api.GET(/home/banners, homeHandler.GetBanners)
		api.GET(/home/featured, homeHandler.GetFeatured)
		api.GET(/home/products, homeHandler.GetHomeProducts)
		
		// 用户接口（需要登录）
		user := api.Group(/)
		user.Use(middleware.AuthMiddleware())
		{
			user.GET(/user/info, userHandler.GetInfo)
			user.POST(/cart, cartHandler.AddCart)
			user.GET(/cart, cartHandler.GetCart)
			user.PUT(/cart/:id, cartHandler.UpdateCart)
			user.DELETE(/cart/:id, cartHandler.DeleteCart)
			user.GET(/addresses, addressHandler.GetAddresses)
			user.POST(/addresses, addressHandler.CreateAddress)
			user.PUT(/addresses/:id, addressHandler.UpdateAddress)
			user.DELETE(/addresses/:id, addressHandler.DeleteAddress)
			user.PUT(/addresses/:id/default, addressHandler.SetDefault)
			user.GET(/orders, orderHandler.GetOrders)
			user.GET(/orders/:id, orderHandler.GetOrder)
			user.POST(/orders, orderHandler.CreateOrder)
			user.PUT(/orders/:id/cancel, orderHandler.CancelOrder)
			user.PUT(/orders/:id/confirm, orderHandler.ConfirmReceive)
			user.POST(/pay/create, paymentHandler.CreatePayment)
			user.GET(/pay/status/:order_id, paymentHandler.GetPayStatus)
			user.GET(/search/history, searchHandler.GetHistory)
			user.DELETE(/search/history, searchHandler.ClearHistory)
		}
		
		// 后台管理接口
		admin := api.Group(/admin)
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.AdminMiddleware())
		{
			// 商品管理
			admin.GET(/products, productHandler.AdminGetProducts)
			admin.POST(/products, productHandler.CreateProduct)
			admin.PUT(/products/:id, productHandler.UpdateProduct)
			admin.DELETE(/products/:id, productHandler.DeleteProduct)
			admin.PUT(/products/:id/shelf, productHandler.SetShelf)
			admin.PUT(/products/:id/home, productHandler.SetShowOnHome)
			
			// 分类管理
			admin.GET(/categories, categoryHandler.AdminGetCategories)
			admin.POST(/categories, categoryHandler.CreateCategory)
			admin.PUT(/categories/:id, categoryHandler.UpdateCategory)
			admin.DELETE(/categories/:id, categoryHandler.DeleteCategory)
			
			// 订单管理
			admin.GET(/orders, orderHandler.AdminGetOrders)
			admin.GET(/orders/:id, orderHandler.AdminGetOrder)
			admin.PUT(/orders/:id/ship, orderHandler.ShipOrder)
			admin.PUT(/orders/:id/price, orderHandler.UpdatePrice)
			
			// 用户管理
			admin.GET(/users, userHandler.AdminGetUsers)
			admin.GET(/users/:id, userHandler.AdminGetUser)
			
			// 区域管理
			admin.GET(/regions, regionHandler.AdminGetRegions)
			admin.POST(/regions, regionHandler.CreateRegion)
			admin.PUT(/regions/:code, regionHandler.UpdateRegion)
			admin.DELETE(/regions/:code, regionHandler.DeleteRegion)
			
			// 首页运营配置
			admin.GET(/home/config, homeHandler.GetConfig)
			admin.PUT(/home/config, homeHandler.UpdateConfig)
		}
		
		// 微信登录
		api.POST(/user/wx_login, userHandler.WxLogin)
		
		// 支付回调（微信服务器调用）
		api.POST(/pay/callback, paymentHandler.PayCallback)
		
		// 搜索
		api.GET(/search/suggest, searchHandler.GetSuggest)
		api.GET(/search/hot, searchHandler.GetHotSearch)
	}

	// 启动服务
	addr := fmt.Sprintf(:%d, cfg.Server.Port)
	fmt.Printf(服务启动: %s\n, addr)
	r.Run(addr)
}

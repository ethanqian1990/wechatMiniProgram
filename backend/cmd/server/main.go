package main

import (
	"fmt"
	"time"
	"github.com/ethanqian1990/wechat-mall-backend/internal/config"
	"github.com/ethanqian1990/wechat-mall-backend/internal/handler"
	"github.com/ethanqian1990/wechat-mall-backend/internal/middleware"
	"github.com/ethanqian1990/wechat-mall-backend/internal/repository"
	"github.com/ethanqian1990/wechat-mall-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load("./backend/configs/config.yaml")
	if err != nil {
		panic("加载配置失败: " + err.Error())
	}

	if cfg.JWT.Secret == "" {
		panic("jwt.secret 不能为空")
	}
	middleware.SetJWTSecret(cfg.JWT.Secret)
	if cfg.JWT.Expire != "" {
		if d, err := time.ParseDuration(cfg.JWT.Expire); err == nil {
			middleware.SetJWTExpire(d)
		}
	}

	// 初始化数据库
	if err := repository.InitDB(cfg); err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	service.NewCronService(cfg).Start()

	r := gin.Default()
	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	{
		adminAuthHandler := handler.NewAdminAuthHandler(cfg)
		api.POST("/admin/login", adminAuthHandler.Login)

		regionHandler := handler.NewRegionHandler()
		categoryHandler := handler.NewCategoryHandler()
		productHandler := handler.NewProductHandler()
		homeHandler := handler.NewHomeHandler(cfg)
		searchHandler := handler.NewSearchHandler()

		api.GET("/regions", regionHandler.GetRegions)
		api.GET("/categories", categoryHandler.GetCategories)
		api.GET("/products", productHandler.GetProducts)
		api.GET("/products/:id", productHandler.GetProduct)
		api.GET("/home/banners", homeHandler.GetBanners)
		api.GET("/home/featured", homeHandler.GetFeatured)
		api.GET("/home/products", homeHandler.GetHomeProducts)
		api.GET("/search/suggest", searchHandler.GetSuggest)
		api.GET("/search/hot", searchHandler.GetHotSearch)

		userHandler := handler.NewUserHandler()
		cartHandler := handler.NewCartHandler()
		addressHandler := handler.NewAddressHandler()
		orderHandler := handler.NewOrderHandler()
		paymentHandler := handler.NewPaymentHandler()

		api.POST("/user/wx_login", userHandler.WxLogin)

		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware())
		{
			user.GET("/info", userHandler.GetInfo)
			user.POST("/profile", userHandler.UpdateProfile)
		}

		cart := api.Group("/cart")
		cart.Use(middleware.AuthMiddleware())
		{
			cart.GET("", cartHandler.GetCart)
			cart.POST("", cartHandler.AddCart)
			cart.PUT("/:id", cartHandler.UpdateCart)
			cart.DELETE("/:id", cartHandler.DeleteCart)
		}

		addr := api.Group("/addresses")
		addr.Use(middleware.AuthMiddleware())
		{
			addr.GET("", addressHandler.GetAddresses)
			addr.POST("", addressHandler.CreateAddress)
			addr.PUT("/:id", addressHandler.UpdateAddress)
			addr.DELETE("/:id", addressHandler.DeleteAddress)
			addr.PUT("/:id/default", addressHandler.SetDefault)
		}

		order := api.Group("/orders")
		order.Use(middleware.AuthMiddleware())
		{
			order.GET("", orderHandler.GetOrders)
			order.GET("/:id", orderHandler.GetOrder)
			order.POST("", orderHandler.CreateOrder)
			order.PUT("/:id/cancel", orderHandler.CancelOrder)
			order.PUT("/:id/pay", orderHandler.PayOrder)
			order.PUT("/:id/confirm", orderHandler.ConfirmReceive)
		}

		pay := api.Group("/pay")
		pay.Use(middleware.AuthMiddleware())
		{
			pay.POST("/create", paymentHandler.CreatePayment)
			pay.GET("/status/:order_id", paymentHandler.GetPayStatus)
		}
		pay.POST("/callback", paymentHandler.PayCallback)

		search := api.Group("/search")
		search.Use(middleware.AuthMiddleware())
		{
			search.GET("/history", searchHandler.GetHistory)
			search.DELETE("/history", searchHandler.ClearHistory)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.AdminMiddleware())
		{
			adminProductHandler := handler.NewAdminProductHandler()
			adminOrderHandler := handler.NewAdminOrderHandler()
			adminUserHandler := handler.NewAdminUserHandler()
			adminRegionHandler := handler.NewAdminRegionHandler()
			adminCategoryHandler := handler.NewAdminCategoryHandler()

			// 商品
			admin.GET("/products", adminProductHandler.GetProducts)
			admin.POST("/products", adminProductHandler.CreateProduct)
			admin.PUT("/products/:id", adminProductHandler.UpdateProduct)
			admin.DELETE("/products/:id", adminProductHandler.DeleteProduct)
			admin.PUT("/products/:id/shelf", adminProductHandler.SetShelf)
			admin.PUT("/products/:id/home", adminProductHandler.SetHome)

			// 订单
			admin.GET("/orders", adminOrderHandler.GetOrders)
			admin.GET("/orders/:id", adminOrderHandler.GetOrder)
			admin.PUT("/orders/:id/ship", adminOrderHandler.ShipOrder)
			admin.PUT("/orders/:id/price", adminOrderHandler.UpdatePrice)

			// 用户
			admin.GET("/users", adminUserHandler.GetUsers)
			admin.GET("/users/:id", adminUserHandler.GetUser)

			// 区域
			admin.GET("/regions", adminRegionHandler.GetRegions)
			admin.POST("/regions", adminRegionHandler.CreateRegion)
			admin.PUT("/regions/:code", adminRegionHandler.UpdateRegion)
			admin.DELETE("/regions/:code", adminRegionHandler.DeleteRegion)

			// 分类
			admin.GET("/categories", adminCategoryHandler.GetCategories)
			admin.POST("/categories", adminCategoryHandler.CreateCategory)
			admin.PUT("/categories/:id", adminCategoryHandler.UpdateCategory)
			admin.DELETE("/categories/:id", adminCategoryHandler.DeleteCategory)

			admin.GET("/home/config", homeHandler.GetConfig)
			admin.PUT("/home/config", homeHandler.UpdateConfig)
		}
	}

	fmt.Printf("服务启动: :%d\n", cfg.Server.Port)
	r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}

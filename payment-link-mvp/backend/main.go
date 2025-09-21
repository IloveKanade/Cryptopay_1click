package main

import (
	"log"
	"math/rand"
	"payment-link-mvp/internal/config"
	"payment-link-mvp/internal/handler"
	"payment-link-mvp/internal/middleware"
	"payment-link-mvp/internal/model"
	"payment-link-mvp/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())

	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := model.InitDB(model.DatabaseConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移数据库表
	if err := model.AutoMigrate(db); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 初始化服务
	userService := service.NewUserService(db)
	epusdtService := service.NewEpusdtService(cfg.Epusdt.BaseURL, cfg.Epusdt.AuthToken)
	tronService := service.NewTronService("") // 暂时不使用API Key
	callbackService := service.NewCallbackService(cfg)
	paymentLinkService := service.NewPaymentLinkService(db, epusdtService)
	adminService := service.NewAdminService(db, epusdtService)
	middlemanWalletService := service.NewMiddlemanWalletService(db, epusdtService, tronService, callbackService)
	feeConfigService := service.NewFeeConfigService(db)

	// 初始化支付服务配置
	paymentConfig := &service.Config{
		BaseURL:     "http://localhost:" + cfg.Server.Port,
		FrontendURL: "http://localhost:3001",
	}
	paymentService := service.NewPaymentService(db, epusdtService, middlemanWalletService, paymentConfig)

	// 初始化处理器
	userHandler := handler.NewUserHandler(userService)
	paymentLinkHandler := handler.NewPaymentLinkHandler(paymentLinkService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	adminHandler := handler.NewAdminHandler(adminService)
	adminMiddlemanHandler := handler.NewAdminMiddlemanHandler(middlemanWalletService, feeConfigService)
	callbackHandler := handler.NewCallbackHandler(callbackService, cfg)

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// 静态文件服务
	r.Static("/static", "./static")

	// 健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	// API路由组
	api := r.Group("/api")
	{
		// 用户相关
		api.POST("/auth/register", userHandler.Register)
		api.POST("/auth/login", userHandler.Login)

		// 需要认证的路由
		auth := api.Group("/")
		auth.Use(middleware.Auth())
		{
			auth.GET("/user/profile", userHandler.GetProfile)

			// 支付链接管理
			auth.POST("/payment-links", paymentLinkHandler.Create)
			auth.GET("/payment-links", paymentLinkHandler.List)
			auth.GET("/payment-links/:id", paymentLinkHandler.Get)
			auth.DELETE("/payment-links/:id", paymentLinkHandler.Delete)

			// 支付订单管理
			auth.GET("/payment/orders", paymentHandler.GetUserPaymentOrders)
			auth.GET("/payment/orders/:orderId", paymentHandler.GetPaymentOrder)
		}

		// 管理员路由 - 需要管理员权限
		admin := api.Group("/admin")
		admin.Use(middleware.AdminAuth())
		{
			// 仪表板统计
			admin.GET("/dashboard/stats", adminHandler.GetDashboardStats)

			// 用户管理
			admin.GET("/users", adminHandler.GetAllUsers)
			admin.GET("/users/search", adminHandler.SearchUsers)
			admin.GET("/users/:id", adminHandler.GetUserByID)
			admin.PUT("/users/:id", adminHandler.UpdateUser)
			admin.DELETE("/users/:id", adminHandler.DeleteUser)

			// 订单管理
			admin.GET("/orders", adminHandler.GetAllOrders)
			admin.GET("/orders/search", adminHandler.SearchOrders)
			admin.GET("/orders/:id", adminHandler.GetOrderByID)
			admin.PUT("/orders/:id/status", adminHandler.UpdateOrderStatus)
			admin.DELETE("/orders/:id", adminHandler.DeleteOrder)

			// 钱包地址管理
			admin.GET("/wallet-addresses", adminHandler.GetWalletAddresses)
			admin.POST("/wallet-addresses", adminHandler.AddWalletAddress)

			// 中间人钱包管理
			admin.GET("/middleman-wallets", adminMiddlemanHandler.ListMiddlemanWallets)
			admin.POST("/middleman-wallets", adminMiddlemanHandler.AddMiddlemanWallet)
			admin.PUT("/middleman-wallets/:id/status", adminMiddlemanHandler.UpdateMiddlemanWalletStatus)
			admin.DELETE("/middleman-wallets/:id", adminMiddlemanHandler.DeleteMiddlemanWallet)

			// 转账记录管理
			admin.GET("/transfer-records", adminMiddlemanHandler.ListTransferRecords)
			admin.GET("/transfer-records/:id", adminMiddlemanHandler.GetTransferRecord)
			admin.POST("/transfer-records/:id/process", adminMiddlemanHandler.ProcessTransfer)

			// 费率配置管理
			admin.GET("/fee-config", adminMiddlemanHandler.GetFeeConfig)
			admin.PUT("/fee-config", adminMiddlemanHandler.UpdateFeeConfig)
			admin.POST("/fee-config", adminMiddlemanHandler.CreateFeeConfig)
			admin.GET("/fee-config/history", adminMiddlemanHandler.GetFeeConfigHistory)
			admin.GET("/fee-config/calculate", adminMiddlemanHandler.CalculateFee)

			// 回调配置管理
			admin.GET("/callback/config", callbackHandler.GetCallbackConfig)
			admin.PUT("/callback/config", callbackHandler.UpdateCallbackConfig)
			admin.POST("/callback/test", callbackHandler.TestCallback)
		}

		// 支付相关路由（部分需要认证）
		api.POST("/payment/create", paymentHandler.CreatePayment)            // 创建支付订单
		api.GET("/payment/status/:tradeId", paymentHandler.QueryOrderStatus) // 查询订单状态
		api.POST("/payment/notify", paymentHandler.HandleNotify)             // Epusdt回调

		// 直接访问Epusdt付款页面 - 不需要认证
		api.GET("/payment/redirect", paymentHandler.RedirectToEpusdtPayment) // 直接重定向到Epusdt付款页面
		api.GET("/payment/url", paymentHandler.GetEpusdtPaymentURL)          // 获取Epusdt付款页面URL

		// 公开支付页面 - 不需要认证
		api.GET("/payment-links/public/:id", paymentLinkHandler.GetPublicPayment)
		api.GET("/pay/:id", paymentLinkHandler.GetPublicPayment)
		api.GET("/payment/page/:linkId", paymentHandler.GetPaymentPage)
	}

	log.Printf("服务器启动在端口 %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

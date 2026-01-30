package main

import (
	"log"
	"akrick.com/mychat/cache"
	"akrick.com/mychat/database"
	"akrick.com/mychat/handlers"
	"akrick.com/mychat/middleware"
	"akrick.com/mychat/tasks"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "akrick.com/mychat/docs"
)

// @title MyChat API
// @version 1.0
// @description 这是一个基于Gin和GORM的聊天系统API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	// 初始化数据库
	if err := database.InitDB(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	log.Println("数据库连接成功")

	// 初始化Redis
	if err := cache.InitRedis(); err != nil {
		log.Printf("Redis连接失败（将跳过缓存）: %v", err)
	} else {
		log.Println("Redis连接成功")
	}

	// 启动定时任务
	tasks.StartScheduler()

	// 创建Gin路由
	r := gin.Default()

	// 使用中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Recovery())

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		redisStatus := "connected"
		if cache.Rdb == nil {
			redisStatus = "disconnected"
		}
		c.JSON(200, gin.H{
			"status": "ok",
			"redis":  redisStatus,
		})
	})

	// 注册接口
	r.POST("/api/register", handlers.Register)

	// 登录接口
	r.POST("/api/login", handlers.Login)

	// Token刷新接口
	r.POST("/api/token/refresh", handlers.RefreshToken)

	// 用户接口
	r.GET("/api/user/info", middleware.AuthMiddleware(), handlers.GetUserInfo)
	r.PUT("/api/user/profile", middleware.AuthMiddleware(), handlers.UpdateProfile)
	r.POST("/api/user/password", middleware.AuthMiddleware(), handlers.ChangePassword)
	r.POST("/api/upload/avatar", middleware.AuthMiddleware(), handlers.UploadAvatar)
	r.POST("/api/user/recharge", middleware.AuthMiddleware(), handlers.Recharge)
	r.GET("/api/user/transactions", middleware.AuthMiddleware(), handlers.GetTransactions)

	// 咨询师接口（只读）
	r.GET("/api/counselor/list", handlers.GetCounselorList)
	r.GET("/api/counselor/:id", handlers.GetCounselorDetail)
	r.GET("/api/counselor/:id/reviews", handlers.GetCounselorReviews)

	// 订单接口
	r.POST("/api/order/create", middleware.AuthMiddleware(), handlers.CreateOrder)
	r.GET("/api/order/:id", middleware.AuthMiddleware(), handlers.GetOrderDetail)
	r.GET("/api/order/list", middleware.AuthMiddleware(), handlers.GetUserOrders)
	r.PUT("/api/order/:id/status", middleware.AuthMiddleware(), handlers.UpdateOrderStatus)
	r.POST("/api/order/:id/cancel", middleware.AuthMiddleware(), handlers.CancelOrder)
	r.GET("/api/counselor/orders", middleware.AuthMiddleware(), handlers.GetCounselorOrders)
	// 支付接口
	r.POST("/api/payment/create", middleware.AuthMiddleware(), handlers.CreatePayment)
	r.GET("/api/payment/:id", middleware.AuthMiddleware(), handlers.GetPaymentStatus)
	r.GET("/api/payment/list", middleware.AuthMiddleware(), handlers.GetUserPayments)
	r.POST("/api/payment/refund", middleware.AuthMiddleware(), handlers.RefundPayment)
	r.POST("/api/payment/wechat/callback", handlers.WeChatPayCallback)
	r.POST("/api/payment/alipay/callback", handlers.AlipayCallback)

	// 通知接口
	r.GET("/api/notification/list", middleware.AuthMiddleware(), handlers.GetNotifications)
	r.POST("/api/notification/:id/read", middleware.AuthMiddleware(), handlers.MarkNotificationRead)
	r.POST("/api/notification/read-all", middleware.AuthMiddleware(), handlers.MarkAllNotificationsRead)
	r.DELETE("/api/notification/:id", middleware.AuthMiddleware(), handlers.DeleteNotification)

	// 聊天接口（部分功能依赖 WebSocket 服务）
	r.POST("/api/chat/start/:order_id", middleware.AuthMiddleware(), handlers.StartChatSession)
	r.POST("/api/chat/session/:session_id/message", middleware.AuthMiddleware(), handlers.SendMessage)
	r.GET("/api/chat/messages/:session_id", middleware.AuthMiddleware(), handlers.GetMessages)
	r.POST("/api/chat/end/:session_id", middleware.AuthMiddleware(), handlers.EndChatSession)
	r.GET("/api/chat/sessions", middleware.AuthMiddleware(), handlers.GetChatSessions)

	// 文件接口
	r.POST("/api/upload", middleware.AuthMiddleware(), handlers.UploadFile)
	r.GET("/api/file/:id", middleware.AuthMiddleware(), handlers.GetFile)
	r.DELETE("/api/file/:id", middleware.AuthMiddleware(), handlers.DeleteFile)

	// 启动服务
	log.Println("==================================================")
	log.Println("MyChat API 服务启动")
	log.Println("地址: :8080")
	log.Println("PID:", os.Getpid())
	log.Println("Swagger文档: http://localhost:8080/swagger/index.html")
	if cache.Rdb != nil {
		log.Println("缓存状态: Redis + SingleFlight防穿透")
	}
	log.Println("定时任务: 已启动")
	log.Println("==================================================")

	// 优雅退出信号处理
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("接收到退出信号，正在优雅关闭...")

		// 关闭Redis连接
		if cache.Rdb != nil {
			if err := cache.CloseRedis(); err != nil {
				log.Printf("关闭Redis连接失败: %v", err)
			} else {
				log.Println("Redis连接已关闭")
			}
		}

		log.Println("服务已关闭")
		os.Exit(0)
	}()

	// 启动服务
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}

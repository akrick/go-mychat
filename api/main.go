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

	// 个人中心接口
	r.GET("/api/admin/profile", middleware.AuthMiddleware(), handlers.GetProfile)
	r.PUT("/api/admin/profile", middleware.AuthMiddleware(), handlers.UpdateProfile)
	r.POST("/api/admin/user/password", middleware.AuthMiddleware(), handlers.ChangePassword)
	r.POST("/api/admin/upload", middleware.AuthMiddleware(), handlers.UploadAvatar)

	// 管理员接口
	r.POST("/api/admin/login", handlers.AdminLogin)
	r.POST("/api/admin/logout", middleware.AuthMiddleware(), handlers.AdminLogout)
	r.GET("/api/admin/user/info", middleware.AuthMiddleware(), handlers.GetAdminUserInfo)
	r.GET("/api/admin/user/permissions", middleware.AuthMiddleware(), handlers.GetAdminPermissions)
	r.GET("/api/admin/statistics", middleware.AuthMiddleware(), handlers.GetAdminStatistics)
	r.GET("/api/admin/session/stats", middleware.AuthMiddleware(), handlers.GetSessionStats)
	r.GET("/api/admin/online/users", middleware.AuthMiddleware(), handlers.GetOnlineUsers)
	r.POST("/api/admin/broadcast", middleware.AuthMiddleware(), handlers.BroadcastSystemMessage)
	r.POST("/api/admin/withdraw/:id/approve", middleware.AuthMiddleware(), handlers.ApproveWithdraw)
	r.GET("/api/admin/withdraws/pending", middleware.AuthMiddleware(), handlers.GetPendingWithdraws)
	
	// 用户管理接口
	r.GET("/api/admin/users", middleware.AuthMiddleware(), handlers.GetUserList)
	r.POST("/api/admin/users", middleware.AuthMiddleware(), handlers.CreateUser)
	r.PUT("/api/admin/users/:id", middleware.AuthMiddleware(), handlers.UpdateUser)
	r.DELETE("/api/admin/users/:id", middleware.AuthMiddleware(), handlers.DeleteUser)
	r.POST("/api/admin/users/:id/password", middleware.AuthMiddleware(), handlers.ResetUserPassword)
	
	// 角色管理接口
	r.GET("/api/admin/roles", middleware.AuthMiddleware(), handlers.GetRoleList)
	r.POST("/api/admin/roles", middleware.AuthMiddleware(), handlers.CreateRole)
	r.PUT("/api/admin/roles/:id", middleware.AuthMiddleware(), handlers.UpdateRole)
	r.DELETE("/api/admin/roles/:id", middleware.AuthMiddleware(), handlers.DeleteRole)
	r.GET("/api/admin/roles/:id/permissions", middleware.AuthMiddleware(), handlers.GetRolePermissions)
	r.POST("/api/admin/roles/:id/permissions", middleware.AuthMiddleware(), handlers.AssignPermissions)
	r.GET("/api/admin/roles/:id/users", middleware.AuthMiddleware(), handlers.GetRoleUsers)

	// 权限管理接口
	r.GET("/api/admin/permissions/tree", middleware.AuthMiddleware(), handlers.GetPermissionTree)
	r.GET("/api/admin/permissions", middleware.AuthMiddleware(), handlers.GetPermissionList)
	r.POST("/api/admin/permissions", middleware.AuthMiddleware(), handlers.CreatePermission)
	r.PUT("/api/admin/permissions/:id", middleware.AuthMiddleware(), handlers.UpdatePermission)
	r.DELETE("/api/admin/permissions/:id", middleware.AuthMiddleware(), handlers.DeletePermission)

	// 菜单管理接口
	r.GET("/api/admin/menus/tree", middleware.AuthMiddleware(), handlers.GetMenuTree)
	r.GET("/api/admin/menus", middleware.AuthMiddleware(), handlers.GetMenuList)
	r.POST("/api/admin/menus", middleware.AuthMiddleware(), handlers.CreateMenu)
	r.PUT("/api/admin/menus/:id", middleware.AuthMiddleware(), handlers.UpdateMenu)
	r.DELETE("/api/admin/menus/:id", middleware.AuthMiddleware(), handlers.DeleteMenu)

	// 系统管理接口
	r.GET("/api/admin/logs", middleware.AuthMiddleware(), handlers.GetSystemLogList)
	r.GET("/api/admin/online/users", middleware.AuthMiddleware(), handlers.GetOnlineUserList)
	r.GET("/api/admin/configs", middleware.AuthMiddleware(), handlers.GetConfigList)
	r.PUT("/api/admin/configs/:id", middleware.AuthMiddleware(), handlers.UpdateConfig)
	r.GET("/api/admin/dashboard/statistics", middleware.AuthMiddleware(), handlers.GetDashboardStatistics)
	
	// 低代码接口
	r.GET("/api/admin/lowcode/forms", middleware.AuthMiddleware(), handlers.GetFormList)
	r.POST("/api/admin/lowcode/forms", middleware.AuthMiddleware(), handlers.SaveFormDesign)
	r.GET("/api/admin/lowcode/forms/:id", middleware.AuthMiddleware(), handlers.GetFormDesign)
	r.DELETE("/api/admin/lowcode/forms/:id", middleware.AuthMiddleware(), handlers.DeleteForm)
	r.POST("/api/admin/lowcode/forms/:id/data", handlers.SubmitFormData)
	r.GET("/api/admin/lowcode/forms/:id/data", middleware.AuthMiddleware(), handlers.GetFormDataList)
	r.GET("/api/admin/lowcode/pages", middleware.AuthMiddleware(), handlers.GetPageList)
	r.POST("/api/admin/lowcode/pages", middleware.AuthMiddleware(), handlers.SavePageDesign)
	r.GET("/api/admin/lowcode/pages/:id", middleware.AuthMiddleware(), handlers.GetPageDesign)
	r.DELETE("/api/admin/lowcode/pages/:id", middleware.AuthMiddleware(), handlers.DeletePage)
	r.GET("/api/admin/lowcode/pages/:id/preview", middleware.AuthMiddleware(), handlers.PreviewPage)

	// 管理后台-聊天记录接口（需要连接独立的 WebSocket 服务）
	r.GET("/api/admin/chat/sessions", middleware.AuthMiddleware(), handlers.GetAdminChatSessions)
	r.GET("/api/admin/chat/sessions/:session_id/messages", middleware.AuthMiddleware(), handlers.GetAdminChatMessages)
	r.GET("/api/admin/chat/statistics", middleware.AuthMiddleware(), handlers.GetChatStatistics)
	r.GET("/api/admin/chat/messages/search", middleware.AuthMiddleware(), handlers.SearchChatMessages)
	r.DELETE("/api/admin/chat/sessions/:id", middleware.AuthMiddleware(), handlers.DeleteChatSession)

	// 管理后台-订单管理接口
	r.GET("/api/admin/orders", middleware.AuthMiddleware(), handlers.GetOrderList)
	r.GET("/api/admin/orders/statistics", middleware.AuthMiddleware(), handlers.GetOrderStatistics)
	r.PUT("/api/admin/orders/:id/status", middleware.AuthMiddleware(), handlers.AdminUpdateOrderStatus)

	// 咨询师接口
	r.GET("/api/counselor/list", handlers.GetCounselorList)
	r.GET("/api/counselor/:id", handlers.GetCounselorDetail)
	r.POST("/api/counselor/create", middleware.AuthMiddleware(), handlers.CreateCounselor)
	r.PUT("/api/counselor/:id", middleware.AuthMiddleware(), handlers.UpdateCounselor)
	r.DELETE("/api/counselor/:id", middleware.AuthMiddleware(), handlers.DeleteCounselor)

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

	// 评价接口
	r.POST("/api/review/create", middleware.AuthMiddleware(), handlers.CreateReview)
	r.GET("/api/review/counselor/:counselor_id", handlers.GetReviewList)
	r.GET("/api/review/counselor/:counselor_id/statistics", handlers.GetCounselorStatistics)
	r.GET("/api/review/:id", handlers.GetReviewDetail)
	r.POST("/api/review/:id/reply", middleware.AuthMiddleware(), handlers.ReplyReview)
	r.GET("/api/review/my", middleware.AuthMiddleware(), handlers.GetUserReviews)

	// 统计接口
	r.GET("/api/stats/dashboard", middleware.AuthMiddleware(), handlers.DashboardStatistics)
	r.GET("/api/stats/order", middleware.AuthMiddleware(), handlers.OrderStatistics)
	r.GET("/api/stats/counselor/ranking", middleware.AuthMiddleware(), handlers.CounselorRanking)

	// 配置接口
	r.GET("/api/config/payment", middleware.AuthMiddleware(), handlers.GetPaymentConfig)
	r.PUT("/api/config/payment/:id", middleware.AuthMiddleware(), handlers.UpdatePaymentConfig)
	r.POST("/api/config/payment/:id/test", middleware.AuthMiddleware(), handlers.TestPaymentConfig)

	// 订单验证接口
	r.GET("/api/order/:id/validate", middleware.AuthMiddleware(), handlers.ValidateOrderStatus)
	r.GET("/api/order/:id/timeline", middleware.AuthMiddleware(), handlers.GetOrderTimeline)

	// 通知接口
	r.GET("/api/notification/list", middleware.AuthMiddleware(), handlers.GetNotifications)
	r.POST("/api/notification/:id/read", middleware.AuthMiddleware(), handlers.MarkNotificationRead)
	r.POST("/api/notification/read-all", middleware.AuthMiddleware(), handlers.MarkAllNotificationsRead)
	r.DELETE("/api/notification/:id", middleware.AuthMiddleware(), handlers.DeleteNotification)

	// 聊天接口（部分功能依赖 WebSocket 服务）
	r.POST("/api/chat/session/:order_id/start", middleware.AuthMiddleware(), handlers.StartChatSession)
	r.GET("/api/chat/session/:session_id/messages", middleware.AuthMiddleware(), handlers.GetMessages)
	r.POST("/api/chat/session/:session_id/end", middleware.AuthMiddleware(), handlers.EndChatSession)
	r.GET("/api/chat/sessions", middleware.AuthMiddleware(), handlers.GetChatSessions)

	// 聊天账单接口
	r.GET("/api/chat/billings", middleware.AuthMiddleware(), handlers.GetBillingList)
	r.GET("/api/chat/counselor/billings", middleware.AuthMiddleware(), handlers.GetCounselorBillings)
	r.GET("/api/chat/counselor/account", middleware.AuthMiddleware(), handlers.GetCounselorAccount)
	r.POST("/api/chat/counselor/withdraw", middleware.AuthMiddleware(), handlers.CreateWithdraw)
	r.GET("/api/chat/counselor/withdraws", middleware.AuthMiddleware(), handlers.GetWithdrawList)

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

package main

import (
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/handlers"
	middlewarepkg "akrick.com/mychat/admin/backend/middleware"
	"akrick.com/mychat/admin/backend/websocket"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	database.InitDB()

	// 初始化系统配置
	InitSystemConfigs()

	// 初始化 WebSocket Hub
	websocket.InitHub()
	fmt.Println("✅ WebSocket Hub 已初始化")

	// 创建路由
	r := gin.Default()

	// 提供静态文件服务
	r.Static("/uploads", "./uploads")

	// CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// JWT 认证中间件
	authMiddleware := middlewarepkg.AuthMiddleware()

	// 公开路由
	public := r.Group("/api")
	{
		public.POST("/admin/login", handlers.AdminLogin)
		public.POST("/admin2/login", handlers.AdminLogin2)
		fmt.Println("✅ 注册公开路由: POST /api/admin/login")
		public.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"msg": "后端正常运行"})
		})
		public.GET("/counselor/list", handlers.GetCounselorList)
		public.GET("/counselor/:id", handlers.GetCounselorDetail)
	}

	// 需要认证的路由
	auth := r.Group("/api")
	auth.Use(authMiddleware)
	{
		// 管理员路由(使用Administrator表)
		admin2 := auth.Group("/admin2")
		admin2.Use(middlewarepkg.AdminAuthMiddleware())
		{
			// 管理员管理
			admin2.GET("/administrators", handlers.GetAdministratorList)
			admin2.POST("/administrators", handlers.CreateAdministrator)
			admin2.PUT("/administrators/:id", handlers.UpdateAdministrator)
			admin2.DELETE("/administrators/:id", handlers.DeleteAdministrator)
			admin2.POST("/administrators/:id/password", handlers.ResetAdministratorPassword)
			admin2.PUT("/administrators/:id/status", handlers.ToggleAdministratorStatus)

			// 个人信息
			admin2.GET("/info", handlers.GetAdminInfo2)
			admin2.PUT("/profile", handlers.UpdateMyProfile)
			admin2.POST("/password", handlers.ChangeMyPassword)
			admin2.POST("/logout", handlers.AdminLogout2)

			// 权限
			admin2.GET("/permissions", handlers.GetAdminPermissions2)
		}
		// 管理员路由(兼容旧接口，使用Administrator表)
		admin := auth.Group("/admin")
		{
			// 文件上传
			admin.POST("/upload/image", handlers.UploadAvatar)
			admin.POST("/upload/file", handlers.UploadFile)

			// 用户管理(普通用户表)
			admin.GET("/users", handlers.GetUserList)
			admin.POST("/users", handlers.CreateUser)
			admin.PUT("/users/:id", handlers.UpdateUser)
			admin.DELETE("/users/:id", handlers.DeleteUser)
			admin.POST("/users/:id/password", handlers.ResetUserPassword)

			// 咨询师管理
			admin.GET("/counselors", handlers.GetCounselorList)
			admin.POST("/counselors", handlers.CreateCounselor)
			admin.PUT("/counselors/:id", handlers.UpdateCounselor)
			admin.DELETE("/counselors/:id", handlers.DeleteCounselor)

			// 入驻申请管理
			admin.GET("/counselor/applications", handlers.GetApplicationList)
			admin.GET("/counselor/applications/:id", handlers.GetApplicationDetail)
			admin.PUT("/counselor/applications/:id/review", handlers.ReviewApplication)

			// 订单管理
			admin.GET("/orders", handlers.GetOrderList)
			admin.GET("/orders/statistics", handlers.GetOrderStatistics)
			admin.PUT("/orders/:id/status", handlers.AdminUpdateOrderStatus)

			// 统计数据
			admin.GET("/stats/counselor/ranking", handlers.CounselorRanking)
			admin.GET("/stats/order/trend", handlers.OrderTrend)

			// 聊天管理
			admin.GET("/chat/sessions", handlers.GetAdminChatSessions)
			admin.GET("/chat/sessions/:session_id/messages", handlers.GetAdminChatMessages)
			admin.GET("/chat/statistics", handlers.GetChatStatistics)
			admin.GET("/chat/messages/search", handlers.SearchChatMessages)
			admin.DELETE("/chat/sessions/:id", handlers.DeleteChatSession)

			// 财务管理
			admin.GET("/withdraws/pending", handlers.GetPendingWithdraws)
			admin.POST("/withdraw/:id/approve", handlers.ApproveWithdraw)
			admin.POST("/withdraw/:id/transfer", handlers.ConfirmWithdrawTransfer)
			admin.GET("/withdraws", handlers.GetWithdrawList)
			admin.GET("/finance/stats", handlers.GetFinanceStats)
			admin.GET("/finance/revenue", handlers.GetRevenueReport)
			admin.GET("/finance/reports", handlers.GetFinanceReports)
			admin.GET("/finance/accounts", handlers.GetCounselorAccountList)
			admin.GET("/finance/accounts/:id", handlers.GetCounselorAccountDetail)
			admin.GET("/statistics", handlers.GetAdminStatistics)

			// 系统管理
			admin.GET("/user/info", handlers.GetAdminUserInfo)
			admin.GET("/user/permissions", handlers.GetAdminPermissions)
			admin.POST("/logout", handlers.AdminLogout)
			admin.GET("/session/stats", handlers.GetSessionStats)
			admin.GET("/online/users", handlers.GetOnlineUsers)
			admin.GET("/online/users/detailed", handlers.GetOnlineUsersDetailed)
			admin.POST("/online/users/:id/kick", handlers.KickOutUser)
			admin.POST("/online/mute", handlers.MuteUser)
			admin.GET("/online/statistics", handlers.GetOnlineStatistics)
			admin.POST("/online/users/:id/message", handlers.SendToUser)
			admin.POST("/broadcast", handlers.BroadcastSystemMessage)
			admin.GET("/logs", handlers.GetSystemLogs)
			admin.GET("/configs", handlers.GetSystemConfigs)
			admin.POST("/configs", handlers.CreateSystemConfig)
			admin.PUT("/configs/:id", handlers.UpdateSystemConfig)
			admin.POST("/configs/batch", handlers.BatchSaveConfigs)
			admin.DELETE("/configs/:id", handlers.DeleteSystemConfig)

			// RBAC 权限管理
			admin.GET("/roles", handlers.GetRoleList)
			admin.POST("/roles", handlers.CreateRole)
			admin.PUT("/roles/:id", handlers.UpdateRole)
			admin.DELETE("/roles/:id", handlers.DeleteRole)
			admin.GET("/roles/:id/permissions", handlers.GetRolePermissions)
			admin.PUT("/roles/:id/permissions", handlers.AssignPermissions)

			admin.GET("/permissions/tree", handlers.GetPermissionTree)
			admin.GET("/permissions", handlers.GetPermissionList)
			admin.POST("/permissions", handlers.CreatePermission)
			admin.PUT("/permissions/:id", handlers.UpdatePermission)
			admin.DELETE("/permissions/:id", handlers.DeletePermission)

			// 菜单管理
			admin.GET("/menus/tree", handlers.GetMenuTree)
			admin.GET("/menus", handlers.GetMenuList)
			admin.POST("/menus", handlers.CreateMenu)
			admin.PUT("/menus/:id", handlers.UpdateMenu)
			admin.DELETE("/menus/:id", handlers.DeleteMenu)

			// 低代码平台
			admin.GET("/lowcode/forms", handlers.GetFormList)
			admin.POST("/lowcode/forms", handlers.SaveFormDesign)
			admin.PUT("/lowcode/forms/:id", handlers.SaveFormDesign)
			admin.DELETE("/lowcode/forms/:id", handlers.DeleteForm)
			admin.GET("/lowcode/forms/:id", handlers.GetFormDesign)

			admin.GET("/lowcode/pages", handlers.GetPageList)
			admin.POST("/lowcode/pages", handlers.SavePageDesign)
			admin.PUT("/lowcode/pages/:id", handlers.SavePageDesign)
			admin.DELETE("/lowcode/pages/:id", handlers.DeletePage)
			admin.GET("/lowcode/pages/:id", handlers.GetPageDesign)
			admin.GET("/lowcode/pages/:id/preview", handlers.PreviewPage)

			admin.GET("/lowcode/forms/:id/data", handlers.GetFormDataList)
			admin.POST("/lowcode/forms/:id/submit", handlers.SubmitFormData)
		}
	}

	// 启动服务
	fmt.Println("管理后台服务启动在端口 :3003")
	log.Fatal(r.Run(":3003"))
}

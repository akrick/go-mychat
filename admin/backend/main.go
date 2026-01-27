package main

import (
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/handlers"
	"akrick.com/mychat/admin/backend/middleware"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	database.InitDB()

	// 创建路由
	r := gin.Default()

	// CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// JWT 认证中间件
	authMiddleware := middleware.AuthMiddleware()

	// 公开路由
	public := r.Group("/api")
	{
		public.POST("/admin/login", handlers.AdminLogin)
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
		// 管理员路由
		admin := auth.Group("/admin")
		{
			// 用户管理
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

			// 订单管理
			admin.GET("/orders", handlers.GetOrderList)
			admin.GET("/orders/statistics", handlers.GetOrderStatistics)
			admin.PUT("/orders/:id/status", handlers.AdminUpdateOrderStatus)

			// 聊天管理
			admin.GET("/chat/sessions", handlers.GetAdminChatSessions)
			admin.GET("/chat/sessions/:session_id/messages", handlers.GetAdminChatMessages)
			admin.GET("/chat/statistics", handlers.GetChatStatistics)
			admin.GET("/chat/messages/search", handlers.SearchChatMessages)
			admin.DELETE("/chat/sessions/:id", handlers.DeleteChatSession)

			// 财务管理
			admin.GET("/withdraws/pending", handlers.GetPendingWithdraws)
			admin.POST("/withdraw/:id/approve", handlers.ApproveWithdraw)
			admin.GET("/withdraws", handlers.GetWithdrawList)
			admin.GET("/finance/stats", handlers.GetFinanceStats)
			admin.GET("/finance/revenue", handlers.GetRevenueReport)
			admin.GET("/statistics", handlers.GetAdminStatistics)

			// 系统管理
			admin.GET("/user/info", handlers.GetAdminUserInfo)
			admin.GET("/user/permissions", handlers.GetAdminPermissions)
			admin.POST("/logout", handlers.AdminLogout)
			admin.GET("/session/stats", handlers.GetSessionStats)
			admin.GET("/online/users", handlers.GetOnlineUsers)
			admin.POST("/broadcast", handlers.BroadcastSystemMessage)

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
	fmt.Println("管理后台服务启动在端口 :8081")
	log.Fatal(r.Run(":8081"))
}

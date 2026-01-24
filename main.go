package main

import (
	"log"
	"akrick.com/mychat/cache"
	"akrick.com/mychat/database"
	"akrick.com/mychat/handlers"
	"akrick.com/mychat/middleware"
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

	// 创建Gin路由
	r := gin.Default()

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

	// 获取用户信息（需要认证）
	r.GET("/api/user/info", middleware.AuthMiddleware(), handlers.GetUserInfo)

	// 启动服务
	log.Println("==================================================")
	log.Println("MyChat API 服务启动")
	log.Println("地址: :8080")
	log.Println("PID:", os.Getpid())
	log.Println("Swagger文档: http://localhost:8080/swagger/index.html")
	if cache.Rdb != nil {
		log.Println("缓存状态: Redis + SingleFlight防穿透")
	}
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

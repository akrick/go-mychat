package main

import (
	"log"
	"websocket/cache"
	"websocket/database"
	"websocket/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	if err := database.InitDB(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	log.Println("数据库连接成功")

	// 自动迁移数据库
	database.DB.AutoMigrate(&models.ChatSession{}, &models.ChatMessage{}, &models.ChatBilling{})

	// 初始化Redis
	if err := cache.InitRedis(); err != nil {
		log.Printf("Redis连接失败（将跳过缓存）: %v", err)
	} else {
		log.Println("Redis连接成功")
	}

	// 初始化WebSocket Hub
	InitHub()
	log.Println("WebSocket Hub初始化成功")

	// 初始化会话管理器
	InitSessionManager()
	log.Println("会话管理器初始化成功")

	// 创建Gin路由
	r := gin.Default()

	// WebSocket 路由
	r.GET("/ws", HandleWebSocket)
	r.GET("/ws/counselor/:id", HandleCounselorWebSocket)

	// 统计信息路由
	r.GET("/ws/stats", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "获取成功",
			"data": GetSessionStats(),
		})
	})

	// 启动服务
	log.Println("==================================================")
	log.Println("MyChat WebSocket 服务启动")
	log.Println("地址: :8082")
	log.Println("WebSocket: ws://localhost:8082/ws")
	log.Println("咨询师: ws://localhost:8082/ws/counselor/{id}")
	log.Println("==================================================")

	if err := r.Run(":8082"); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}

package tasks

import (
	"fmt"
	"log"
	"time"

	"akrick.com/mychat/database"
)

// StartScheduler 启动定时任务
// 用于订单超时取消、消息清理等后台任务
func StartScheduler() {
	log.Println("定时任务调度器已启动")

	// 启动一个后台goroutine来执行定时任务
	go func() {
		// 订单超时检查 - 每5分钟执行一次
		orderTicker := time.NewTicker(5 * time.Minute)
		defer orderTicker.Stop()

		// 会话清理 - 每1小时执行一次
		sessionTicker := time.NewTicker(1 * time.Hour)
		defer sessionTicker.Stop()

		for {
			select {
			case <-orderTicker.C:
				checkExpiredOrders()
			case <-sessionTicker.C:
				cleanupExpiredSessions()
			}
		}
	}()
}

// checkExpiredOrders 检查并取消超时未支付的订单
func checkExpiredOrders() {
	log.Println("执行订单超时检查...")

	// 创建类型别名以避免循环依赖
	type Order struct {
		ID        uint      `gorm:"primaryKey"`
		UserID    uint      `gorm:"index"`
		Status    int       `gorm:"default:0"` // 0:待支付, 1:已支付, 2:已取消, 3:已完成
		CreatedAt time.Time
	}

	type Notification struct {
		ID        uint   `gorm:"primaryKey"`
		UserID    uint   `gorm:"index"`
		Title     string
		Content   string `gorm:"type:text"`
		IsRead    bool   `gorm:"default:false"`
		CreatedAt time.Time
	}

	// 查询创建超过30分钟且状态为待支付的订单
	timeout := 30 * time.Minute
	var expiredOrders []Order
	err := database.DB.Where("status = ? AND created_at < ?", 0, time.Now().Add(-timeout)).Find(&expiredOrders).Error
	if err != nil {
		log.Printf("查询超时订单失败: %v", err)
		return
	}

	if len(expiredOrders) == 0 {
		log.Println("没有超时订单需要处理")
		return
	}

	log.Printf("发现 %d 个超时订单", len(expiredOrders))

	// 批量取消订单
	tx := database.DB.Begin()
	for _, order := range expiredOrders {
		// 更新订单状态
		if err := tx.Model(&Order{}).Where("id = ?", order.ID).Update("status", 2).Error; err != nil {
			log.Printf("更新订单 %d 状态失败: %v", order.ID, err)
			tx.Rollback()
			return
		}

		// 创建通知
		notification := Notification{
			UserID:  order.UserID,
			Title:   "订单已取消",
			Content: fmt.Sprintf("您的订单 #%d 因超时未支付已自动取消", order.ID),
			IsRead:  false,
		}
		if err := tx.Create(&notification).Error; err != nil {
			log.Printf("创建通知失败: %v", err)
			tx.Rollback()
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("事务提交失败: %v", err)
		return
	}

	log.Printf("成功取消 %d 个超时订单", len(expiredOrders))
}

// cleanupExpiredSessions 清理过期的聊天会话
func cleanupExpiredSessions() {
	log.Println("执行会话清理...")

	// 定义结构体
	type ChatMessage struct {
		ID        uint   `gorm:"primaryKey"`
		SessionID uint   `gorm:"index"`
	}

	// 查询已结束超过7天的会话
	expireDays := 7
	expiredTime := time.Now().AddDate(0, 0, -expireDays)

	type ChatSession struct {
		ID        uint      `gorm:"primaryKey"`
		Status    int       // 0:待开始, 1:进行中, 2:已结束
		EndTime   *time.Time
		CreatedAt time.Time
	}

	var expiredSessions []ChatSession
	err := database.DB.Where("status = ? AND end_time < ?", 2, expiredTime).Find(&expiredSessions).Error
	if err != nil {
		log.Printf("查询过期会话失败: %v", err)
		return
	}

	if len(expiredSessions) == 0 {
		log.Println("没有过期会话需要清理")
		return
	}

	log.Printf("发现 %d 个过期会话", len(expiredSessions))

	// 批量删除过期的会话和消息
	tx := database.DB.Begin()
	for _, session := range expiredSessions {
		// 删除相关消息
		if err := tx.Where("session_id = ?", session.ID).Delete(&ChatMessage{}).Error; err != nil {
			log.Printf("删除会话 %d 的消息失败: %v", session.ID, err)
			tx.Rollback()
			return
		}

		// 删除会话
		if err := tx.Delete(&session).Error; err != nil {
			log.Printf("删除会话 %d 失败: %v", session.ID, err)
			tx.Rollback()
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("事务提交失败: %v", err)
		return
	}

	log.Printf("成功清理 %d 个过期会话及其消息", len(expiredSessions))
}

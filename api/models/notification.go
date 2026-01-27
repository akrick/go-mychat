package models

import (
	"time"
)

// NotificationType 通知类型
const (
	NotificationTypeOrder   = "order"   // 订单通知
	NotificationTypePayment = "payment" // 支付通知
	NotificationTypeReview  = "review"  // 评价通知
	NotificationTypeSystem  = "system"  // 系统通知
	NotificationTypeChat    = "chat"    // 聊天通知
)

// NotificationLevel 通知级别
const (
	NotificationLevelInfo     = "info"     // 信息
	NotificationLevelWarning  = "warning"  // 警告
	NotificationLevelError    = "error"    // 错误
	NotificationLevelSuccess  = "success"  // 成功
)

// Notification 通知表
type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index;comment:接收用户ID" json:"user_id"`
	Type      string    `gorm:"type:varchar(20);not null;index;comment:通知类型" json:"type"`
	Level     string    `gorm:"type:varchar(20);not null;default:info;comment:通知级别" json:"level"`
	Title     string    `gorm:"type:varchar(255);not null;comment:通知标题" json:"title"`
	Content   string    `gorm:"type:text;comment:通知内容" json:"content"`
	ExtraData string    `gorm:"type:text;comment:额外数据(JSON)" json:"extra_data"`
	IsRead    bool      `gorm:"not null;default:false;index;comment:是否已读" json:"is_read"`
	ReadTime  *time.Time `json:"read_time"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

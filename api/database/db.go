package database

import (
	"fmt"
	"akrick.com/mychat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	dsn := "root:123456@tcp(localhost:3306)/mychat?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 自动迁移表结构
	err = DB.AutoMigrate(
		// 用户相关
		&models.User{},
		&models.Counselor{},
		&models.CounselorAccount{},
		&models.CounselorStatistics{},

		// 订单相关
		&models.Order{},
		&models.Payment{},
		&models.PaymentConfig{},

		// 聊天相关
		&models.ChatSession{},
		&models.ChatMessage{},
		&models.ChatBilling{},

		// 评价相关
		&models.Review{},

		// 财务相关
		&models.WithdrawRecord{},

		// RBAC权限相关
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},

		// 系统管理
		&models.SysLog{},
		&models.SysConfig{},
		&models.OnlineUser{},

		// 低代码平台
		&models.FormDesign{},
		&models.FormData{},
		&models.PageDesign{},

		// 文件和通知
		&models.File{},
		&models.Notification{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

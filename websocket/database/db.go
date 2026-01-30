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
	err = DB.AutoMigrate(&models.User{}, &models.Counselor{}, &models.Order{}, &models.Payment{}, &models.PaymentConfig{}, &models.Notification{}, &models.ChatSession{}, &models.ChatMessage{}, &models.File{}, &models.ChatBilling{}, &models.CounselorAccount{}, &models.WithdrawRecord{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

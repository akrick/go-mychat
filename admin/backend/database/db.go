package database

import (
	"fmt"
	"akrick.com/mychat/admin/backend/models"
	"akrick.com/mychat/admin/backend/utils"
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
		&models.User{},
		&models.Counselor{},
		&models.Order{},
		&models.Payment{},
		&models.PaymentConfig{},
		&models.Review{},
		&models.CounselorStatistics{},
		&models.Notification{},
		&models.ChatSession{},
		&models.ChatMessage{},
		&models.File{},
		&models.ChatBilling{},
		&models.CounselorAccount{},
		&models.WithdrawRecord{},
		&models.Role{},
		&models.Permission{},
		&models.Menu{},
		&models.LowcodeForm{},
		&models.LowcodePage{},
		&models.LowcodeFormData{},
		&models.RolePermission{},
		&models.UserRole{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// 创建默认管理员账号
	createDefaultAdmin()

	return nil
}

// createDefaultAdmin 创建默认管理员账号
func createDefaultAdmin() {
	var count int64
	DB.Model(&models.User{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		hashedPassword, _ := utils.HashPassword("admin123")
		admin := models.User{
			Username: "admin",
			Password: hashedPassword,
			Email:    "admin@mychat.com",
			Status:   1,
			IsAdmin:  true,
		}
		DB.Create(&admin)
		fmt.Println("默认管理员账号创建成功: admin / admin123")
	}
}

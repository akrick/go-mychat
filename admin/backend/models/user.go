package models

import (
	"time"
	"gorm.io/gorm"
)

// User 普通用户表（前端用户）
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null;comment:用户名" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null;comment:密码" json:"-"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;comment:邮箱" json:"email"`
	Phone     string         `gorm:"type:varchar(20);comment:手机号" json:"phone"`
	Avatar    string         `gorm:"type:varchar(255);comment:头像" json:"avatar"`
	Status    int            `gorm:"default:1;comment:状态:1-正常,0-禁用" json:"status"`
}

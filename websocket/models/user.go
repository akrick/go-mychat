package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Phone     string         `gorm:"type:varchar(20)" json:"phone"`
	Avatar    string         `gorm:"type:varchar(255)" json:"avatar"`
	Status    int            `gorm:"default:1;comment:1-正常,0-禁用" json:"status"`
	IsAdmin   bool           `gorm:"default:false;comment:是否管理员" json:"is_admin"`
}

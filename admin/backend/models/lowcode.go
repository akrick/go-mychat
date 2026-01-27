package models

import "time"

// LowcodeForm 低代码表单
type LowcodeForm struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Title       string    `gorm:"type:varchar(200);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	FormJSON    string    `gorm:"type:longtext;not null" json:"form_json"`
	Status      int       `gorm:"default:1;comment:1-启用,0-禁用" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// LowcodePage 低代码页面
type LowcodePage struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Title       string    `gorm:"type:varchar(200);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	PageJSON    string    `gorm:"type:longtext;not null" json:"page_json"`
	Status      int       `gorm:"default:1;comment:1-启用,0-禁用" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// LowcodeFormData 低代码表单数据
type LowcodeFormData struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	FormID    uint      `gorm:"not null;index" json:"form_id"`
	FormData  string    `gorm:"type:longtext;not null" json:"form_data"`
	CreatedBy uint      `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           uint `gorm:"primarykey" json:"id"`
	RoleID       uint `gorm:"not null;index" json:"role_id"`
	PermissionID uint `gorm:"not null;index" json:"permission_id"`
}

// UserRole 用户角色关联表
type UserRole struct {
	ID     uint `gorm:"primarykey" json:"id"`
	UserID uint `gorm:"not null;index" json:"user_id"`
	RoleID uint `gorm:"not null;index" json:"role_id"`
}

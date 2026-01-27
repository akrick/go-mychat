package models

import (
	"time"
)

// UserRole 用户角色关联表
type UserRole struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null;uniqueIndex:idx_user_role,priority:1;comment:用户ID" json:"user_id"`
	RoleID uint `gorm:"not null;uniqueIndex:idx_user_role,priority:2;comment:角色ID" json:"role_id"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           uint `gorm:"primaryKey" json:"id"`
	RoleID       uint `gorm:"not null;uniqueIndex:idx_role_permission,priority:1;comment:角色ID" json:"role_id"`
	PermissionID uint `gorm:"not null;uniqueIndex:idx_role_permission,priority:2;comment:权限ID" json:"permission_id"`

	// 关联
	Role       Role       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Permission Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}

// SysLog 系统日志表
type SysLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index;comment:操作用户ID" json:"user_id"`
	Username    string    `gorm:"type:varchar(50);not null;comment:操作用户名" json:"username"`
	Module      string    `gorm:"type:varchar(50);not null;comment:操作模块" json:"module"`
	Action      string    `gorm:"type:varchar(100);not null;comment:操作动作" json:"action"`
	Method      string    `gorm:"type:varchar(10);not null;comment:请求方法" json:"method"`
	IP          string    `gorm:"type:varchar(50);comment:IP地址" json:"ip"`
	URL         string    `gorm:"type:varchar(255);comment:请求URL" json:"url"`
	Params      string    `gorm:"type:text;comment:请求参数" json:"params"`
	Result      string    `gorm:"type:text;comment:返回结果" json:"result"`
	Status      int       `gorm:"default:1;comment:状态:1-成功,0-失败" json:"status"`
	ErrorMsg    string    `gorm:"type:text;comment:错误信息" json:"error_msg"`
	Duration    int       `gorm:"comment:执行时长(毫秒)" json:"duration"`
	CreatedAt   time.Time `json:"created_at"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// OnlineUser 在线用户表
type OnlineUser struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex;index;comment:用户ID" json:"user_id"`
	Token     string    `gorm:"type:varchar(255);not null;comment:Token" json:"token"`
	IP        string    `gorm:"type:varchar(50);comment:IP地址" json:"ip"`
	UserAgent string    `gorm:"type:varchar(500);comment:浏览器UA" json:"user_agent"`
	LoginAt   time.Time `gorm:"not null;comment:登录时间" json:"login_at"`
	UpdatedAt  time.Time `gorm:"not null;comment:最后活动时间" json:"updated_at"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// SysConfig 系统配置表
type SysConfig struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ConfigKey  string    `gorm:"type:varchar(100);uniqueIndex;not null;comment:配置键" json:"config_key"`
	ConfigName string    `gorm:"type:varchar(100);not null;comment:配置名称" json:"config_name"`
	ConfigType string    `gorm:"type:varchar(20);default:string;comment:配置类型:string/number/boolean/json" json:"config_type"`
	ConfigVal string    `gorm:"type:text;comment:配置值" json:"config_val"`
	IsPublic  bool      `gorm:"default:false;comment:是否公开" json:"is_public"`
	Remark     string    `gorm:"type:varchar(500);comment:备注" json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

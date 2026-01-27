package models

import "time"

// Role 角色表
type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(50);uniqueIndex;not null;comment:角色名称" json:"name"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null;comment:角色代码" json:"code"`
	Description string    `gorm:"type:varchar(255);comment:描述" json:"description"`
	Sort        int       `gorm:"default:0;comment:排序" json:"sort"`
	Status      int       `gorm:"default:1;comment:状态:0-禁用,1-启用" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Permission 权限表
type Permission struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	ParentID    uint         `gorm:"default:0;index;comment:父权限ID" json:"parent_id"`
	Name        string       `gorm:"type:varchar(50);not null;comment:权限名称" json:"name"`
	Code        string       `gorm:"type:varchar(100);uniqueIndex;not null;comment:权限代码" json:"code"`
	Type        string       `gorm:"type:varchar(20);default:menu;comment:类型:menu-菜单,button-按钮,api-接口" json:"type"`
	Path        string       `gorm:"type:varchar(255);comment:路由路径" json:"path"`
	Icon        string       `gorm:"type:varchar(50);comment:图标" json:"icon"`
	Component   string       `gorm:"type:varchar(255);comment:组件路径" json:"component"`
	Sort        int          `gorm:"default:0;comment:排序" json:"sort"`
	Status      int          `gorm:"default:1;comment:状态:0-禁用,1-启用" json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Children    []Permission `gorm:"-" json:"children,omitempty"`
}

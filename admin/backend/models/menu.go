package models

import (
	"time"
)

// Menu 菜单模型
type Menu struct {
	ID         int        `json:"id" gorm:"primaryKey;autoIncrement"`
	ParentID   *int       `json:"parent_id" gorm:"index;comment:父菜单ID"`
	Name       string     `json:"name" gorm:"size:100;not null;comment:菜单名称"`
	Type       int        `json:"type" gorm:"not null;default:2;comment:类型:1目录,2菜单,3按钮"`
	Path       string     `json:"path" gorm:"size:200;comment:路由路径"`
	Component  string     `json:"component" gorm:"size:200;comment:组件路径"`
	Permission string     `json:"permission" gorm:"size:100;comment:权限标识"`
	Icon       string     `json:"icon" gorm:"size:50;comment:图标"`
	Sort       int        `json:"sort" gorm:"default:0;comment:排序"`
	Status     int        `json:"status" gorm:"not null;default:1;comment:状态:0禁用,1启用"`
	CreatedAt  time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}

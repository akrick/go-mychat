package models

import (
	"time"
)

// FormDesign 表单设计表
type FormDesign struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null;comment:表单名称" json:"name"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null;comment:表单代码" json:"code"`
	Description string    `gorm:"type:varchar(500);comment:表单描述" json:"description"`
	FormSchema  string    `gorm:"type:longtext;not null;comment:表单配置JSON" json:"form_schema"`
	IsPublished bool      `gorm:"default:false;comment:是否发布" json:"is_published"`
	CreatedBy   uint      `gorm:"comment:创建人ID" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联
	Creator User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Forms   []FormData `gorm:"foreignKey:FormID" json:"forms,omitempty"`
}

// FormData 表单数据表
type FormData struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FormID    uint      `gorm:"not null;index;comment:表单ID" json:"form_id"`
	SubmitBy  uint      `gorm:"not null;index;comment:提交人ID" json:"submit_by"`
	Data      string    `gorm:"type:longtext;not null;comment:表单数据JSON" json:"data"`
	IP        string    `gorm:"type:varchar(50);comment:提交IP" json:"ip"`
	UserAgent string    `gorm:"type:varchar(500);comment:浏览器UA" json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	Form    FormDesign `gorm:"foreignKey:FormID" json:"form,omitempty"`
	Submited User       `gorm:"foreignKey:SubmitBy" json:"submited,omitempty"`
}

// PageDesign 页面设计表
type PageDesign struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null;comment:页面名称" json:"name"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null;comment:页面代码" json:"code"`
	Path        string    `gorm:"type:varchar(255);uniqueIndex;not null;comment:页面路径" json:"path"`
	Description string    `gorm:"type:varchar(500);comment:页面描述" json:"description"`
	PageConfig string    `gorm:"type:longtext;not null;comment:页面配置JSON" json:"page_config"`
	IsPublished bool      `gorm:"default:false;comment:是否发布" json:"is_published"`
	CreatedBy   uint      `gorm:"comment:创建人ID" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联
	Creator User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

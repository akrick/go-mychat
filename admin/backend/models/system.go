package models

import (
	"time"
)

// SystemLog 系统日志表
type SystemLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Operator    string    `gorm:"type:varchar(50);comment:操作人" json:"operator"`
	Action      string    `gorm:"type:varchar(50);comment:操作类型" json:"action"`
	Module      string    `gorm:"type:varchar(50);comment:操作模块" json:"module"`
	Description string    `gorm:"type:text;comment:操作描述" json:"description"`
	IPAddress   string    `gorm:"type:varchar(50);comment:IP地址" json:"ip_address"`
	UserAgent   string    `gorm:"type:varchar(500);comment:浏览器" json:"user_agent"`
	RequestData string    `gorm:"type:text;comment:请求参数" json:"request_data"`
	Response    string    `gorm:"type:text;comment:响应数据" json:"response"`
	Status      int       `gorm:"default:1;comment:状态:1-成功,0-失败" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// SystemConfig 系统配置表
type SystemConfig struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"type:varchar(100);uniqueIndex;not null;comment:配置键" json:"key"`
	Value     string    `gorm:"type:text;comment:配置值" json:"value"`
	Category  string    `gorm:"type:varchar(50);comment:配置分类" json:"category"`
	Label     string    `gorm:"type:varchar(100);comment:配置标签" json:"label"`
	Type      string    `gorm:"type:varchar(20);comment:配置类型" json:"type"`
	IsSystem  bool      `gorm:"default:false;comment:是否系统配置" json:"is_system"`
	Sort      int       `gorm:"default:0;comment:排序" json:"sort"`
	Remark    string    `gorm:"type:varchar(255);comment:备注" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

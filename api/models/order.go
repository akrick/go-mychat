package models

import (
	"time"
)

// 订单状态
const (
	OrderStatusPending   = 0 // 待支付
	OrderStatusPaid      = 1 // 已支付
	OrderStatusCompleted = 2 // 已完成
	OrderStatusCancelled = 3 // 已取消
	OrderStatusRefunded  = 4 // 已退款
)

// Order 订单表
type Order struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OrderNo      string    `gorm:"type:varchar(32);uniqueIndex;not null" json:"order_no"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	CounselorID  uint      `gorm:"not null;index" json:"counselor_id"`
	Duration     int       `gorm:"not null;comment:咨询时长(分钟)" json:"duration"`
	Amount       float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status       int       `gorm:"not null;default:0;index" json:"status"`
	ScheduleTime time.Time `gorm:"not null;comment:预约时间" json:"schedule_time"`
	Notes        string    `gorm:"type:text;comment:备注" json:"notes"`
	PayTime      *time.Time `json:"pay_time"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 关联
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Counselor Counselor `gorm:"foreignKey:CounselorID" json:"counselor,omitempty"`
}

// Counselor 咨询师表
type Counselor struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;comment:用户ID" json:"user_id"`
	Name      string    `gorm:"type:varchar(50);not null" json:"name"`
	Title     string    `gorm:"type:varchar(50);comment:职称" json:"title"`
	Avatar    string    `gorm:"type:varchar(255);comment:头像" json:"avatar"`
	Bio       string    `gorm:"type:text;comment:个人简介" json:"bio"`
	Specialty string    `gorm:"type:varchar(255);comment:擅长领域" json:"specialty"`
	Price     float64   `gorm:"type:decimal(10,2);not null;comment:单价(元/分钟)" json:"price"`
	YearsExp  int       `gorm:"comment:从业年限" json:"years_exp"`
	Rating    float64   `gorm:"type:decimal(3,2);default:5.00;comment:评分" json:"rating"`
	Status    int       `gorm:"not null;default:1;comment:状态:1-启用,0-禁用" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 统计字段（不持久化到数据库）
	ServiceCount int `gorm:"-" json:"service_count"`      // 服务人数
	ReviewCount  int `gorm:"-" json:"review_count"`       // 评价数量
	IsRecommended bool `gorm:"-" json:"is_recommended"`  // 是否推荐

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// CounselorStatistics 咨询师统计表
type CounselorStatistics struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	CounselorID     uint      `gorm:"not null;uniqueIndex;comment:咨询师ID" json:"counselor_id"`
	TotalOrders     int       `gorm:"not null;default:0;comment:总订单数" json:"total_orders"`
	CompletedOrders int       `gorm:"not null;default:0;comment:已完成订单数" json:"completed_orders"`
	CancelledOrders int       `gorm:"not null;default:0;comment:已取消订单数" json:"cancelled_orders"`
	TotalDuration   int       `gorm:"not null;default:0;comment:总咨询时长(分钟)" json:"total_duration"`
	TotalAmount     float64   `gorm:"type:decimal(12,2);not null;default:0.00;comment:总金额" json:"total_amount"`
	ReviewCount     int       `gorm:"not null;default:0;comment:评价数量" json:"review_count"`
	AvgRating       float64   `gorm:"type:decimal(3,2);not null;default:0.00;comment:平均评分" json:"avg_rating"`
	SumRating       int       `gorm:"not null;default:0;comment:总评分" json:"sum_rating"`
	LastOrderTime   *time.Time `gorm:"comment:最后订单时间" json:"last_order_time"`
	UpdatedAt       time.Time `json:"updated_at"`

	// 关联
	Counselor Counselor `gorm:"foreignKey:CounselorID" json:"counselor,omitempty"`
}

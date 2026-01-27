package models

import (
	"time"
)

// ReviewType 评价类型
const (
	ReviewTypeCounselor = "counselor" // 评价咨询师
)

// Review 评价表
type Review struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OrderID     uint      `gorm:"not null;uniqueIndex:idx_reviews_order;comment:订单ID" json:"order_id"`
	OrderNo     string    `gorm:"type:varchar(32);not null;index;comment:订单号" json:"order_no"`
	UserID      uint      `gorm:"not null;index;comment:用户ID" json:"user_id"`
	CounselorID uint      `gorm:"not null;index;comment:咨询师ID" json:"counselor_id"`
	Rating      int       `gorm:"not null;comment:评分(1-5)" json:"rating"`
	ServiceRating int     `gorm:"not null;default:0;comment:服务评分" json:"service_rating"`
	Professionalism int   `gorm:"not null;default:0;comment:专业度评分" json:"professionalism"`
	Effectiveness int    `gorm:"not null;default:0;comment:有效性评分" json:"effectiveness"`
	Content     string    `gorm:"type:text;comment:评价内容" json:"content"`
	IsAnonymous  bool      `gorm:"default:false;comment:是否匿名" json:"is_anonymous"`
	Status       int       `gorm:"not null;default:1;index;comment:状态:1-显示,0-隐藏" json:"status"`
	ReplyContent string    `gorm:"type:text;comment:咨询师回复" json:"reply_content"`
	ReplyTime    *time.Time `json:"reply_time"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 关联
	Order     Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Counselor Counselor `gorm:"foreignKey:CounselorID" json:"counselor,omitempty"`
}

// CounselorStatistics 咨询师统计表
type CounselorStatistics struct {
	ID             uint    `gorm:"primaryKey" json:"id"`
	CounselorID    uint    `gorm:"not null;uniqueIndex:idx_counselor_statistics_counselor;comment:咨询师ID" json:"counselor_id"`
	TotalOrders    int     `gorm:"not null;default:0;comment:总订单数" json:"total_orders"`
	CompletedOrders int    `gorm:"not null;default:0;comment:已完成订单数" json:"completed_orders"`
	CancelledOrders int    `gorm:"not null;default:0;comment:已取消订单数" json:"cancelled_orders"`
	TotalDuration  int     `gorm:"not null;default:0;comment:总咨询时长(分钟)" json:"total_duration"`
	TotalAmount    float64 `gorm:"type:decimal(12,2);not null;default:0;comment:总金额" json:"total_amount"`
	ReviewCount    int     `gorm:"not null;default:0;comment:评价数量" json:"review_count"`
	AvgRating      float64 `gorm:"type:decimal(3,2);not null;default:0;comment:平均评分" json:"avg_rating"`
	SumRating      int     `gorm:"not null;default:0;comment:总评分" json:"sum_rating"`
	LastOrderTime  *time.Time `gorm:"comment:最后订单时间" json:"last_order_time"`
	UpdatedAt      time.Time `json:"updated_at"`

	// 关联
	Counselor Counselor `gorm:"foreignKey:CounselorID" json:"counselor,omitempty"`
}

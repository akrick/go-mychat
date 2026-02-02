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
}

// 入驻申请状态
const (
	ApplicationStatusPending   = 0 // 待审核
	ApplicationStatusApproved = 1 // 审核通过
	ApplicationStatusRejected = 2 // 审核拒绝
)

// CounselorApplication 咨询师入驻申请表
type CounselorApplication struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	UserID           uint       `gorm:"not null;index" json:"user_id"`
	Name             string     `gorm:"type:varchar(50);not null" json:"name"`
	Gender           string     `gorm:"type:varchar(10);comment:性别" json:"gender"`
	Phone            string     `gorm:"type:varchar(20);not null" json:"phone"`
	Email            string     `gorm:"type:varchar(100);comment:邮箱" json:"email"`
	Title            string     `gorm:"type:varchar(50);comment:职称" json:"title"`
	Bio              string     `gorm:"type:text;comment:个人简介" json:"bio"`
	Specialty        string     `gorm:"type:varchar(255);comment:擅长领域" json:"specialty"`
	YearsExp         int        `gorm:"comment:从业年限" json:"years_exp"`
	CertificateImg1  string     `gorm:"type:varchar(255);comment:资质证书1" json:"certificate_img1"`
	CertificateImg2  string     `gorm:"type:varchar(255);comment:资质证书2" json:"certificate_img2"`
	CertificateImg3  string     `gorm:"type:varchar(255);comment:资质证书3" json:"certificate_img3"`
	Status           int        `gorm:"not null;default:0;index;comment:状态:0-待审核,1-审核通过,2-审核拒绝" json:"status"`
	RejectReason     string     `gorm:"type:text;comment:拒绝原因" json:"reject_reason"`
	ReviewerID       *uint      `gorm:"comment:审核人ID" json:"reviewer_id"`
	ReviewedAt       *time.Time `gorm:"comment:审核时间" json:"reviewed_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// 关联
	User     User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Reviewer *User `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
}

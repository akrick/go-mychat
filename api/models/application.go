package models

import (
	"time"
)

// CounselorApplication 咨询师入驻申请表
type CounselorApplication struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"not null;uniqueIndex;comment:用户ID" json:"user_id"`
	Name            string    `gorm:"type:varchar(50);not null;comment:真实姓名" json:"name"`
	Gender          string    `gorm:"type:varchar(10);comment:性别" json:"gender"`
	Phone           string    `gorm:"type:varchar(20);not null;comment:联系电话" json:"phone"`
	Email           string    `gorm:"type:varchar(100);comment:邮箱" json:"email"`
	Title           string    `gorm:"type:varchar(50);comment:职称" json:"title"`
	Specialty       string    `gorm:"type:varchar(500);comment:擅长领域" json:"specialty"`
	YearsExp        int       `gorm:"comment:从业年限" json:"years_exp"`
	Bio             string    `gorm:"type:text;comment:个人简介" json:"bio"`
	CertificateImg1 string    `gorm:"type:varchar(255);comment:资质证书1" json:"certificate_img1"`
	CertificateImg2 string    `gorm:"type:varchar(255);comment:资质证书2" json:"certificate_img2"`
	CertificateImg3 string    `gorm:"type:varchar(255);comment:资质证书3" json:"certificate_img3"`
	Status          int       `gorm:"not null;default:0;comment:状态:0-待审核,1-审核通过,2-审核拒绝" json:"status"`
	RejectReason    string    `gorm:"type:text;comment:拒绝原因" json:"reject_reason"`
	ReviewedBy      uint      `gorm:"comment:审核人ID" json:"reviewed_by"`
	ReviewedAt      *time.Time `json:"reviewed_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// 关联
	User     User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Reviewer *User `gorm:"foreignKey:ReviewedBy" json:"reviewer,omitempty"`
}

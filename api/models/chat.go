package models

import (
	"time"
)

// ChatSession 聊天会话表
type ChatSession struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OrderID     uint      `gorm:"not null;index;comment:关联订单ID" json:"order_id"`
	UserID      uint      `gorm:"not null;index;comment:用户ID" json:"user_id"`
	CounselorID uint      `gorm:"not null;index;comment:咨询师ID" json:"counselor_id"`
	Status      int       `gorm:"not null;default:0;comment:状态:0-待开始,1-进行中,2-已结束,3-已超时" json:"status"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	Duration    int       `gorm:"comment:实际时长(秒)" json:"duration"`
	Price       float64   `gorm:"type:decimal(10,2);comment:单价(元/分钟)" json:"price"`
	TotalAmount float64   `gorm:"type:decimal(10,2);comment:总金额(元)" json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联
	Order     Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Counselor Counselor `gorm:"foreignKey:CounselorID" json:"counselor,omitempty"`
}

// ChatMessage 聊天消息表
type ChatMessage struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	SessionID   uint       `gorm:"not null;index;comment:会话ID" json:"session_id"`
	SenderID    uint       `gorm:"not null;index;comment:发送者ID" json:"sender_id"`
	SenderType  string     `gorm:"type:varchar(20);not null;comment:发送者类型:user/counselor" json:"sender_type"`
	ContentType string     `gorm:"type:varchar(20);default:text;comment:内容类型:text/image/file" json:"content_type"`
	Content     string     `gorm:"type:text;comment:消息内容" json:"content"`
	FileURL     string     `gorm:"type:varchar(255);comment:文件URL" json:"file_url"`
	IsRead      bool       `gorm:"default:false;comment:是否已读" json:"is_read"`
	ReadTime    *time.Time `json:"read_time"`
	CreatedAt   time.Time  `json:"created_at"`

	// 关联
	Session ChatSession `gorm:"foreignKey:SessionID" json:"session,omitempty"`
}

// ChatBilling 聊天计费记录表
type ChatBilling struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	SessionID       uint      `gorm:"not null;uniqueIndex;index;comment:会话ID" json:"session_id"`
	OrderID         uint      `gorm:"not null;index;comment:关联订单ID" json:"order_id"`
	UserID          uint      `gorm:"not null;index;comment:用户ID" json:"user_id"`
	CounselorID     uint      `gorm:"not null;index;comment:咨询师ID" json:"counselor_id"`
	Duration        int       `gorm:"not null;comment:计费时长(秒)" json:"duration"`
	PricePerMinute  float64   `gorm:"type:decimal(10,2);not null;comment:单价(元/分钟)" json:"price_per_minute"`
	TotalAmount     float64   `gorm:"type:decimal(10,2);not null;comment:总金额" json:"total_amount"`
	PlatformFee     float64   `gorm:"type:decimal(10,2);not null;comment:平台费用(30%)" json:"platform_fee"`
	CounselorFee    float64   `gorm:"type:decimal(10,2);not null;comment:咨询师收入(70%)" json:"counselor_fee"`
	Status          int       `gorm:"not null;default:0;comment:状态:0-待结算,1-已结算" json:"status"`
	SettledAt       *time.Time `json:"settled_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// 关联
	Session   ChatSession `gorm:"foreignKey:SessionID" json:"session,omitempty"`
	Order     Order       `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	User      User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Counselor Counselor   `gorm:"foreignKey:CounselorID" json:"counselor,omitempty"`
}

// CounselorAccount 咨询师账户表
type CounselorAccount struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CounselorID  uint      `gorm:"not null;uniqueIndex;index;comment:咨询师ID" json:"counselor_id"`
	TotalIncome  float64   `gorm:"type:decimal(10,2);default:0;comment:总收入" json:"total_income"`
	Withdrawn    float64   `gorm:"type:decimal(10,2);default:0;comment:已提现" json:"withdrawn"`
	Balance      float64   `gorm:"type:decimal(10,2);default:0;comment:可用余额" json:"balance"`
	FrozenAmount float64   `gorm:"type:decimal(10,2);default:0;comment:冻结金额" json:"frozen_amount"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 关联
	Counselor Counselor `gorm:"foreignKey:CounselorID" json:"counselor,omitempty"`
}

// WithdrawRecord 提现记录表
type WithdrawRecord struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CounselorID   uint      `gorm:"not null;index;comment:咨询师ID" json:"counselor_id"`
	Amount        float64   `gorm:"type:decimal(10,2);not null;comment:提现金额" json:"amount"`
	Status        int       `gorm:"not null;default:0;comment:状态:0-待审核,1-已通过,2-已拒绝,3-已打款" json:"status"`
	BankName      string    `gorm:"type:varchar(50);comment:开户行" json:"bank_name"`
	BankAccount   string    `gorm:"type:varchar(50);comment:银行账号" json:"bank_account"`
	AccountName   string    `gorm:"type:varchar(50);comment:账户名" json:"account_name"`
	RejectedReason string   `gorm:"type:varchar(255);comment:拒绝原因" json:"rejected_reason"`
	AuditedAt     *time.Time `json:"audited_at"`
	TransferredAt *time.Time `json:"transferred_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// 关联
	Counselor Counselor `gorm:"foreignKey:CounselorID" json:"counselor,omitempty"`
}

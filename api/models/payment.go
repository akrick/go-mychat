package models

import (
	"time"
)

// 支付方式
const (
	PaymentMethodWeChat = "wechat"
	PaymentMethodAlipay = "alipay"
)

// 支付状态
const (
	PaymentStatusPending   = 0 // 待支付
	PaymentStatusPaid      = 1 // 已支付
	PaymentStatusFailed    = 2 // 支付失败
	PaymentStatusRefunded  = 3 // 已退款
	PaymentStatusCancelled = 4 // 已取消
)

// 支付交易类型
const (
	TradeTypeApp     = "APP"     // APP支付
	TradeTypeJSAPI   = "JSAPI"   // 公众号/小程序支付
	TradeTypeNative  = "NATIVE"  // 扫码支付
	TradeTypeH5      = "H5"      // H5支付
	TradeTypeAlipay  = "ALIPAY"  // 支付宝
)

// Payment 支付记录表
type Payment struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	PaymentNo       string    `gorm:"type:varchar(32);uniqueIndex;not null;comment:支付单号" json:"payment_no"`
	OrderID         uint      `gorm:"not null;index;comment:关联订单ID" json:"order_id"`
	OrderNo         string    `gorm:"type:varchar(32);not null;index;comment:订单号" json:"order_no"`
	UserID          uint      `gorm:"not null;index;comment:用户ID" json:"user_id"`
	PaymentMethod   string    `gorm:"type:varchar(20);not null;comment:支付方式:wechat/alipay" json:"payment_method"`
	TradeType       string    `gorm:"type:varchar(20);comment:交易类型" json:"trade_type"`
	TransactionID   string    `gorm:"type:varchar(64);uniqueIndex;comment:第三方支付交易号" json:"transaction_id"`
	Amount          float64   `gorm:"type:decimal(10,2);not null;comment:支付金额" json:"amount"`
	Status          int       `gorm:"not null;default:0;index;comment:支付状态" json:"status"`
	PayTime         *time.Time `json:"pay_time"`
	NotifyTime      *time.Time `json:"notify_time"`
	NotifyData      string    `gorm:"type:text;comment:支付回调原始数据" json:"notify_data"`
	FailureReason   string    `gorm:"type:varchar(255);comment:失败原因" json:"failure_reason"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// 关联
	Order Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

// PaymentConfig 支付配置
type PaymentConfig struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	PaymentMethod  string `gorm:"type:varchar(20);uniqueIndex;not null;comment:支付方式" json:"payment_method"`
	AppID          string `gorm:"type:varchar(64);comment:应用ID" json:"app_id"`
	MchID          string `gorm:"type:varchar(64);comment:商户号" json:"mch_id"`
	APISecret      string `gorm:"type:varchar(128);comment:API密钥" json:"api_secret"`
	APICertPath    string `gorm:"type:varchar(255);comment:证书路径" json:"api_cert_path"`
	APIKeyPath     string `gorm:"type:varchar(255);comment:密钥路径" json:"api_key_path"`
	NotifyURL      string `gorm:"type:varchar(255);comment:回调地址" json:"notify_url"`
	PrivateKeyPath string `gorm:"type:varchar(255);comment:私钥路径" json:"private_key_path"`
	PublicKeyPath  string `gorm:"type:varchar(255);comment:公钥路径" json:"public_key_path"`
	IsEnabled      bool   `gorm:"default:true;comment:是否启用" json:"is_enabled"`
	IsSandbox      bool   `gorm:"default:false;comment:是否沙箱环境" json:"is_sandbox"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

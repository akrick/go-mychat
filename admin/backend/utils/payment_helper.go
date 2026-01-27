package utils

import (
	"fmt"
	"time"
)

// GenerateNonceStr 生成随机字符串
func GenerateNonceStr() string {
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
	return fmt.Sprintf("%x", timestamp)[:32]
}

// GeneratePaymentNo 生成支付单号
func GeneratePaymentNo() string {
	return fmt.Sprintf("PAY%d", time.Now().UnixNano())
}

// GenerateTradeNo 生成第三方交易号（模拟）
func GenerateTradeNo(prefix string) string {
	return fmt.Sprintf("%s%s", prefix, fmt.Sprintf("%d", time.Now().UnixNano()))
}

// ConvertYuanToFen 元转分
func ConvertYuanToFen(yuan float64) int {
	return int(yuan * 100)
}

// ConvertFenToYuan 分转元
func ConvertFenToYuan(fen int) float64 {
	return float64(fen) / 100
}

// GetTradeTypeText 获取交易类型描述
func GetTradeTypeText(tradeType string) string {
	textMap := map[string]string{
		"APP":     "APP支付",
		"JSAPI":   "公众号/小程序支付",
		"NATIVE":  "扫码支付",
		"H5":      "H5支付",
		"ALIPAY":  "支付宝",
	}
	if text, ok := textMap[tradeType]; ok {
		return text
	}
	return tradeType
}

// GetPaymentStatusText 获取支付状态描述
func GetPaymentStatusText(status int) string {
	textMap := map[int]string{
		0: "待支付",
		1: "已支付",
		2: "支付失败",
		3: "已退款",
		4: "已取消",
	}
	if text, ok := textMap[status]; ok {
		return text
	}
	return "未知状态"
}

// IsPaymentSuccess 判断支付是否成功
func IsPaymentSuccess(status int) bool {
	return status == 1
}

// GetPaymentExpiryTime 获取支付过期时间
func GetPaymentExpiryTime() time.Time {
	// 支付订单默认2小时过期
	return time.Now().Add(2 * time.Hour)
}

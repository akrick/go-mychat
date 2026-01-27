package utils

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Alipay 支付宝支付工具
type Alipay struct {
	AppID            string
	PrivateKey       string
	AlipayPublicKey  string
	NotifyURL        string
	IsSandbox        bool
	Charset          string
	SignType         string
	Format           string
	Version          string
}

// TradeCreateRequest 统一收单下单并支付请求
type TradeCreateRequest struct {
	OutTradeNo     string `json:"out_trade_no"`
	TotalAmount    string `json:"total_amount"`
	Subject        string `json:"subject"`
	Body           string `json:"body,omitempty"`
	TimeoutExpress string `json:"timeout_express,omitempty"`
	NotifyURL      string `json:"notify_url,omitempty"`
	ReturnURL      string `json:"return_url,omitempty"`
	GoodsDetail    string `json:"goods_detail,omitempty"`
	OperateTime    string `json:"operate_time,omitempty"`
	ExtendParams   string `json:"extend_params,omitempty"`
	SellerID       string `json:"seller_id,omitempty"`
	ProductCode    string `json:"product_code"` // QUICK_WAP_WAY(手机网站), FACE_TO_FACE_PAYMENT(当面付), CYCLE_PAY_AUTH(周期扣款)
}

// TradeCreateResponse 统一收单下单并支付响应
type TradeCreateResponse struct {
	Code      string `json:"code"`
	Msg       string `json:"msg"`
	SubCode   string `json:"sub_code,omitempty"`
	SubMsg    string `json:"sub_msg,omitempty"`
	OutTradeNo string `json:"out_trade_no,omitempty"`
	TradeNo   string `json:"trade_no,omitempty"`
	PayURL    string `json:"qr_code,omitempty"` // PC网页支付二维码地址
	OrderInfo string `json:"order_string,omitempty"` // APP支付
}

// AlipayNotifyRequest 支付宝异步通知
type AlipayNotifyRequest struct {
	TradeNo          string `json:"trade_no"`
	OutTradeNo       string `json:"out_trade_no"`
	TotalAmount      string `json:"total_amount"`
	BuyerId          string `json:"buyer_id"`
	SellerId         string `json:"seller_id"`
	TradeStatus      string `json:"trade_status"` // WAIT_BUYER_PAY, TRADE_SUCCESS, TRADE_FINISHED, TRADE_CLOSED
	NotifyTime       string `json:"notify_time"`
	NotifyType       string `json:"notify_type"`
	NotifyID         string `json:"notify_id"`
	Subject          string `json:"subject"`
	Body             string `json:"body"`
	GmtCreate        string `json:"gmt_create"`
	GmtPayment       string `json:"gmt_payment"`
	RefundFee        string `json:"refund_fee"`
	ReceiptBillNo    string `json:"receipt_bill_no"`
	Charset          string `json:"charset"`
	Sign             string `json:"sign"`
	SignType         string `json:"sign_type"`
}

// NewAlipay 创建支付宝支付实例
func NewAlipay(appID, privateKey, alipayPublicKey, notifyURL string, isSandbox bool) *Alipay {
	return &Alipay{
		AppID:           appID,
		PrivateKey:      privateKey,
		AlipayPublicKey: alipayPublicKey,
		NotifyURL:       notifyURL,
		IsSandbox:       isSandbox,
		Charset:         "utf-8",
		SignType:        "RSA2",
		Format:          "JSON",
		Version:         "1.0",
	}
}

// CreateTrade 创建交易（PC网页支付/手机网站支付）
func (a *Alipay) CreateTrade(request TradeCreateRequest) (*TradeCreateResponse, error) {
	// 构建公共请求参数
	bizContent, _ := json.Marshal(request)

	// 构建请求参数
	params := map[string]string{
		"app_id":      a.AppID,
		"method":      "alipay.trade.page.pay",
		"format":      a.Format,
		"charset":     a.Charset,
		"sign_type":   a.SignType,
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     a.Version,
		"biz_content": string(bizContent),
		"notify_url":  a.NotifyURL,
	}

	// 生成签名
	sign := a.Sign(params)
	params["sign"] = sign

	// 模拟API调用（实际项目中需要调用支付宝API）
	// 返回模拟数据
	response := &TradeCreateResponse{
		Code:       "10000",
		Msg:        "Success",
		OutTradeNo: request.OutTradeNo,
		PayURL:     fmt.Sprintf("https://qr.alipay.com/mock_%s", request.OutTradeNo),
	}

	return response, nil
}

// CreateAppTrade 创建APP支付
func (a *Alipay) CreateAppTrade(request TradeCreateRequest) (*TradeCreateResponse, error) {
	// 构建公共请求参数
	bizContent, _ := json.Marshal(request)

	// 构建请求参数
	params := map[string]string{
		"app_id":      a.AppID,
		"method":      "alipay.trade.app.pay",
		"format":      a.Format,
		"charset":     a.Charset,
		"sign_type":   a.SignType,
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     a.Version,
		"biz_content": string(bizContent),
		"notify_url":  a.NotifyURL,
	}

	// 生成签名
	sign := a.Sign(params)
	params["sign"] = sign

	// 生成订单字符串（实际应用需要对参数进行URL编码）
	orderInfo := a.buildOrderString(params)

	// 返回模拟数据
	response := &TradeCreateResponse{
		Code:      "10000",
		Msg:       "Success",
		OrderInfo: orderInfo,
	}

	return response, nil
}

// buildOrderString 构建订单字符串
func (a *Alipay) buildOrderString(params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf strings.Builder
	for _, k := range keys {
		if buf.Len() > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(params[k])
	}
	return buf.String()
}

// Sign 生成签名
func (a *Alipay) Sign(params map[string]string) string {
	// 过滤空值和sign
	var keys []string
	for k, v := range params {
		if v != "" && k != "sign" {
			keys = append(keys, k)
		}
	}

	// 字典序排序
	sort.Strings(keys)

	// 拼接字符串
	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(params[k])
	}

	// RSA2签名
	return a.RSASign(buf.String())
}

// RSASign RSA2签名
func (a *Alipay) RSASign(data string) string {
	// 解析私钥（实际项目中需要从文件或配置加载）
	// 这里使用模拟实现
	hash := sha256.Sum256([]byte(data))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// VerifySign 验证签名
func (a *Alipay) VerifySign(params map[string]string, sign string) bool {
	// 构建待签名字符串
	var keys []string
	for k, v := range params {
		if v != "" && k != "sign" && k != "sign_type" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(params[k])
	}

	// RSA2验签（实际需要使用支付宝公钥验签）
	return sign == a.RSASign(buf.String())
}

// ParseNotify 解析异步通知
func (a *Alipay) ParseNotify(notifyData string) (*AlipayNotifyRequest, error) {
	// 解析参数
	var request AlipayNotifyRequest
	// 简化实现，实际应解析URL编码的表单数据
	json.Unmarshal([]byte(notifyData), &request)

	// 验证签名
	params := map[string]string{
		"trade_no":     request.TradeNo,
		"out_trade_no": request.OutTradeNo,
		"total_amount": request.TotalAmount,
		"buyer_id":     request.BuyerId,
		"seller_id":    request.SellerId,
		"trade_status": request.TradeStatus,
		"notify_time":  request.NotifyTime,
		"notify_type":  request.NotifyType,
		"notify_id":    request.NotifyID,
		"subject":      request.Subject,
		"body":         request.Body,
		"gmt_create":   request.GmtCreate,
		"gmt_payment":  request.GmtPayment,
	}

	if !a.VerifySign(params, request.Sign) {
		return nil, fmt.Errorf("签名验证失败")
	}

	return &request, nil
}

// RSAVerify RSA验签
func (a *Alipay) RSAVerify(data, sign string, publicKey *rsa.PublicKey) bool {
	// 将签名解码
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false
	}

	// 计算数据的哈希
	hashed := sha256.Sum256([]byte(data))

	// 使用公钥验证签名
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signBytes)
	return err == nil
}

// ParsePublicKey 解析公钥
func (a *Alipay) ParsePublicKey(publicKey string) (*rsa.PublicKey, error) {
	// 去掉头尾注释
	publicKey = strings.ReplaceAll(publicKey, "-----BEGIN PUBLIC KEY-----", "")
	publicKey = strings.ReplaceAll(publicKey, "-----END PUBLIC KEY-----", "")
	publicKey = strings.ReplaceAll(publicKey, "\n", "")
	publicKey = strings.ReplaceAll(publicKey, "\r", "")

	// Base64解码
	keyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}

	// 解析公钥
	pubKey, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return nil, err
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("不是RSA公钥")
	}

	return rsaPubKey, nil
}

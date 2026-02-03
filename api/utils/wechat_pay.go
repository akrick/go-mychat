package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

// WeChatPay 微信支付工具
type WeChatPay struct {
	AppID     string
	MchID     string
	APISecret string
	NotifyURL string
	IsSandbox bool
}

// UnifiedOrderRequest 统一下单请求
type UnifiedOrderRequest struct {
	AppID           string `json:"appid"`
	MchID           string `json:"mch_id"`
	NonceStr        string `json:"nonce_str"`
	Body            string `json:"body"`
	OutTradeNo      string `json:"out_trade_no"`
	TotalFee        int    `json:"total_fee"` // 单位：分
	SpbillCreateIP  string `json:"spbill_create_ip"`
	NotifyURL       string `json:"notify_url"`
	TradeType       string `json:"trade_type"`
	OpenID          string `json:"openid,omitempty"`       // JSAPI支付需要
	TimeStart       string `json:"time_start,omitempty"`   // 交易起始时间
	TimeExpire      string `json:"time_expire,omitempty"`  // 交易结束时间
	Attach          string `json:"attach,omitempty"`       // 附加数据
	Detail          string `json:"detail,omitempty"`      // 商品详情
	GoodsTag        string `json:"goods_tag,omitempty"`    // 订单优惠标记
	LimitPay        string `json:"limit_pay,omitempty"`    // 限定支付方式
	ProductID       string `json:"product_id,omitempty"`   // NATIVE支付需要
	SceneInfo       string `json:"scene_info,omitempty"`   // 场景信息
	H5Info          string `json:"h5_info,omitempty"`      // H5支付场景信息
}

// UnifiedOrderResponse 统一下单响应
type UnifiedOrderResponse struct {
	ReturnCode     string `json:"return_code"`
	ReturnMsg      string `json:"return_msg"`
	AppID          string `json:"appid,omitempty"`
	MchID          string `json:"mch_id,omitempty"`
	NonceStr       string `json:"nonce_str,omitempty"`
	Sign           string `json:"sign,omitempty"`
	ResultCode     string `json:"result_code,omitempty"`
	ErrCode        string `json:"err_code,omitempty"`
	ErrCodeDes     string `json:"err_code_des,omitempty"`
	TradeType      string `json:"trade_type,omitempty"`
	PrepayID       string `json:"prepay_id,omitempty"`
	CodeURL        string `json:"code_url,omitempty"`     // NATIVE支付返回
	MWebURL        string `json:"mweb_url,omitempty"`     // H5支付返回
	CodeURL2       string `json:"code_url_2,omitempty"`   // 备用二维码
	AppResponse    string `json:"app_response,omitempty"` // APP支付参数
}

// NewWeChatPay 创建微信支付实例
func NewWeChatPay(appID, mchID, apiSecret, notifyURL string, isSandbox bool) *WeChatPay {
	return &WeChatPay{
		AppID:     appID,
		MchID:     mchID,
		APISecret: apiSecret,
		NotifyURL: notifyURL,
		IsSandbox: isSandbox,
	}
}

// CreateUnifiedOrder 创建统一下单（模拟）
func (w *WeChatPay) CreateUnifiedOrder(order UnifiedOrderRequest) (*UnifiedOrderResponse, error) {
	// 设置参数
	order.AppID = w.AppID
	order.MchID = w.MchID
	order.NotifyURL = w.NotifyURL
	order.NonceStr = GenerateNonceStr()

	// 计算签名
	sign := w.Sign(order)

	// 模拟API调用（实际项目中需要调用微信API）
	// 这里返回模拟数据
	response := &UnifiedOrderResponse{
		ReturnCode: "SUCCESS",
		ReturnMsg:  "OK",
		AppID:      w.AppID,
		MchID:      w.MchID,
		NonceStr:   order.NonceStr,
		ResultCode: "SUCCESS",
		TradeType:  order.TradeType,
		Sign:       sign,
	}

	// 根据不同的交易类型返回不同的参数
	switch order.TradeType {
	case "NATIVE":
		response.CodeURL = "weixin://wxpay/bizpayurl?pr=" + order.NonceStr
	case "H5":
		response.MWebURL = "https://wx.tenpay.com/cgi-bin/mmpayweb-bin/checkmweb?prepay_id=" + order.NonceStr
	case "JSAPI":
		response.PrepayID = "wx" + order.NonceStr
	case "APP":
		response.PrepayID = "wx" + order.NonceStr
	}

	return response, nil
}

// Sign 计算签名
func (w *WeChatPay) Sign(params interface{}) string {
	m := toMap(params)

	// 过滤空值和sign字段
	var keys []string
	for k, v := range m {
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
		buf.WriteString(m[k])
	}
	buf.WriteString("&key=")
	buf.WriteString(w.APISecret)

	// MD5加密并转大写
	hash := md5.Sum([]byte(buf.String()))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// VerifySign 验证签名
func (w *WeChatPay) VerifySign(params interface{}) bool {
	sign := w.Sign(params)
	return sign == getSignFromParams(params)
}

// GetJSAPIPayParams 获取JSAPI支付参数
func (w *WeChatPay) GetJSAPIPayParams(prepayID string) map[string]string {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonceStr := GenerateNonceStr()

	params := map[string]string{
		"appId":     w.AppID,
		"timeStamp": timestamp,
		"nonceStr":  nonceStr,
		"package":   "prepay_id=" + prepayID,
		"signType":  "MD5",
	}

	sign := w.Sign(params)
	params["paySign"] = sign

	return params
}

// GetAppPayParams 获取APP支付参数
func (w *WeChatPay) GetAppPayParams(prepayID string) map[string]string {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonceStr := GenerateNonceStr()

	params := map[string]string{
		"appid":     w.AppID,
		"partnerid": w.MchID,
		"prepayid":  prepayID,
		"package":   "Sign=WXPay",
		"noncestr":  nonceStr,
		"timestamp": timestamp,
	}

	sign := w.Sign(params)
	params["sign"] = sign

	return params
}

// ParseNotify 解析支付回调通知
func (w *WeChatPay) ParseNotify(notifyData string) (map[string]string, error) {
	var data map[string]string
	if err := xmlToMap(notifyData, &data); err != nil {
		return nil, err
	}

	// 验证签名
	if !w.VerifySign(data) {
		return nil, fmt.Errorf("签名验证失败")
	}

	return data, nil
}

// toMap 将结构体转换为map
func toMap(v interface{}) map[string]string {
	m := make(map[string]string)
	b, _ := json.Marshal(v)
	json.Unmarshal(b, &m)
	return m
}

// getSignFromParams 从参数中获取sign
func getSignFromParams(params interface{}) string {
	m := toMap(params)
	return m["sign"]
}

// xmlToMap XML转Map（简化实现，实际项目应使用xml解析）
func xmlToMap(_ string, result *map[string]string) error {
	// 简化实现，实际应使用encoding/xml解析
	*result = map[string]string{
		"return_code": "SUCCESS",
		"result_code": "SUCCESS",
		"openid":      "mock_openid",
		"trade_type":  "NATIVE",
		"bank_type":   "CMC",
		"total_fee":   "100",
		"transaction_id": "mock_transaction_id",
		"out_trade_no":   "mock_order_no",
		"time_end":       "20240125000000",
	}
	return nil
}

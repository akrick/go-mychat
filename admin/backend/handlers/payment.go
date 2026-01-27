package handlers

import (
	"context"
	"fmt"
	"time"
	"akrick.com/mychat/admin/backend/cache"
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"
	"akrick.com/mychat/admin/backend/utils"
	"github.com/gin-gonic/gin"
)

type CreatePaymentRequest struct {
	OrderID       uint   `json:"order_id" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required,oneof=wechat alipay"`
	TradeType     string `json:"trade_type" binding:"required"`
	ClientIP      string `json:"client_ip"`
	ReturnURL     string `json:"return_url"` // 支付成功后的跳转地址
}

type RefundPaymentRequest struct {
	PaymentID     uint    `json:"payment_id" binding:"required"`
	RefundAmount  float64 `json:"refund_amount" binding:"required"`
	RefundReason  string  `json:"refund_reason"`
}

// CreatePayment godoc
// @Summary 创建支付
// @Description 创建支付订单，支持微信支付和支付宝支付
// @Tags 支付
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePaymentRequest true "支付信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:创建成功,data:{payment_no,pay_url,pay_params}"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "订单不存在"
// @Router /api/payment/create [post]
func CreatePayment(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 查询订单信息
	var order models.Order
	if err := database.DB.First(&order, req.OrderID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "订单不存在",
		})
		return
	}

	// 检查订单是否属于当前用户
	if order.UserID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权操作此订单",
		})
		return
	}

	// 检查订单状态
	if order.Status != models.OrderStatusPending {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "订单状态不允许创建支付",
		})
		return
	}

	// 生成支付单号
	paymentNo := utils.GeneratePaymentNo()

	// 创建支付记录
	payment := models.Payment{
		PaymentNo:     paymentNo,
		OrderID:       order.ID,
		OrderNo:       order.OrderNo,
		UserID:        userID.(uint),
		PaymentMethod: req.PaymentMethod,
		TradeType:     req.TradeType,
		Amount:        order.Amount,
		Status:        models.PaymentStatusPending,
	}

	if err := database.DB.Create(&payment).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "创建支付记录失败: " + err.Error(),
		})
		return
	}

	var payURL string
	var payParams map[string]string

	// 根据支付方式发起支付
	if req.PaymentMethod == models.PaymentMethodWeChat {
		payURL, payParams = createWeChatPayment(&order, &payment, req.ClientIP)
	} else if req.PaymentMethod == models.PaymentMethodAlipay {
		payURL, payParams = createAlipayPayment(&order, &payment, req.ReturnURL)
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建支付成功",
		"data": gin.H{
			"payment_id":  payment.ID,
			"payment_no":  payment.PaymentNo,
			"amount":      payment.Amount,
			"pay_url":     payURL,
			"pay_params":  payParams,
			"trade_type":  payment.TradeType,
			"expired_at":  utils.GetPaymentExpiryTime(),
		},
	})
}

// WeChatPayCallback 微信支付回调
// @Summary 微信支付回调
// @Description 处理微信支付异步通知
// @Tags 支付
// @Accept xml
// @Produce xml
// @Param body body string true "微信支付回调数据"
// @Success 200 {string} string "XML格式响应"
// @Router /api/payment/wechat/callback [post]
func WeChatPayCallback(c *gin.Context) {
	ctx := context.Background()
	notifyData, _ := c.GetRawData()

	// 解析回调数据
	notifyMap := make(map[string]string)
	// 简化实现，实际应解析XML
	notifyMap["out_trade_no"] = "mock_order_no"
	notifyMap["transaction_id"] = utils.GenerateTradeNo("WX")
	notifyMap["total_fee"] = "100"
	notifyMap["time_end"] = "20240125000000"

	// 根据支付单号查询支付记录
	var payment models.Payment
	if err := database.DB.Where("payment_no = ?", notifyMap["out_trade_no"]).First(&payment).Error; err != nil {
		c.XML(200, gin.H{
			"return_code": "FAIL",
			"return_msg":  "支付记录不存在",
		})
		return
	}

	// 检查支付状态，避免重复处理
	if payment.Status == models.PaymentStatusPaid {
		c.XML(200, gin.H{
			"return_code": "SUCCESS",
			"return_msg":  "OK",
		})
		return
	}

	// 更新支付状态
	now := time.Now()
	updates := map[string]interface{}{
		"status":        models.PaymentStatusPaid,
		"transaction_id": notifyMap["transaction_id"],
		"pay_time":      &now,
		"notify_time":   &now,
		"notify_data":   string(notifyData),
	}

	if err := database.DB.Model(&payment).Updates(updates).Error; err != nil {
		c.XML(200, gin.H{
			"return_code": "FAIL",
			"return_msg":  "更新支付记录失败",
		})
		return
	}

	// 更新订单状态
	var order models.Order
	if err := database.DB.First(&order, payment.OrderID).Error; err == nil {
		database.DB.Model(&order).Updates(map[string]interface{}{
			"status":   models.OrderStatusPaid,
			"pay_time": &now,
		})

		// 清除缓存
		cache.DeleteOrderCache(ctx, order.ID)
		cache.InvalidateUserOrdersCache(ctx, order.UserID)
		cache.InvalidateCounselorOrdersCache(ctx, order.CounselorID)
	}

	// 清除支付缓存
	cache.DeletePaymentCache(ctx, payment.ID)

	// 返回成功响应
	c.XML(200, gin.H{
		"return_code": "SUCCESS",
		"return_msg":  "OK",
	})
}

// AlipayCallback 支付宝支付回调
// @Summary 支付宝支付回调
// @Description 处理支付宝异步通知
// @Tags 支付
// @Accept json
// @Produce json
// @Param body body string true "支付宝回调数据"
// @Success 200 {string} string "success"
// @Router /api/payment/alipay/callback [post]
func AlipayCallback(c *gin.Context) {
	ctx := context.Background()

	// 解析回调数据
	var notifyData map[string]string
	if err := c.ShouldBindJSON(&notifyData); err != nil {
		c.String(200, "fail")
		return
	}

	// 获取支付单号
	outTradeNo := notifyData["out_trade_no"]
	if outTradeNo == "" {
		c.String(200, "fail")
		return
	}

	// 根据支付单号查询支付记录
	var payment models.Payment
	if err := database.DB.Where("payment_no = ?", outTradeNo).First(&payment).Error; err != nil {
		c.String(200, "fail")
		return
	}

	// 检查支付状态，避免重复处理
	if payment.Status == models.PaymentStatusPaid {
		c.String(200, "success")
		return
	}

	// 更新支付状态
	now := time.Now()
	updates := map[string]interface{}{
		"status":        models.PaymentStatusPaid,
		"transaction_id": notifyData["trade_no"],
		"pay_time":      &now,
		"notify_time":   &now,
	}

	// 序列化通知数据
	if notifyBytes, err := c.GetRawData(); err == nil {
		updates["notify_data"] = string(notifyBytes)
	}

	if err := database.DB.Model(&payment).Updates(updates).Error; err != nil {
		c.String(200, "fail")
		return
	}

	// 更新订单状态
	var order models.Order
	if err := database.DB.First(&order, payment.OrderID).Error; err == nil {
		database.DB.Model(&order).Updates(map[string]interface{}{
			"status":   models.OrderStatusPaid,
			"pay_time": &now,
		})

		// 清除缓存
		cache.DeleteOrderCache(ctx, order.ID)
		cache.InvalidateUserOrdersCache(ctx, order.UserID)
		cache.InvalidateCounselorOrdersCache(ctx, order.CounselorID)
	}

	// 清除支付缓存
	cache.DeletePaymentCache(ctx, payment.ID)

	c.String(200, "success")
}

// GetPaymentStatus godoc
// @Summary 查询支付状态
// @Description 查询支付记录状态
// @Tags 支付
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "支付记录ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{payment}"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "支付记录不存在"
// @Router /api/payment/{id} [get]
func GetPaymentStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	paymentID := c.Param("id")
	ctx := context.Background()

	var payment models.Payment
	var fromCache bool

	// 尝试从缓存获取
	if cache.Rdb != nil {
		id := parseUint(paymentID)
		cachedPayment, err := cache.GetPaymentWithCache(ctx, id)
		if err == nil {
			payment = *cachedPayment
			fromCache = true
		}
	}

	// 缓存未命中，从数据库查询
	if !fromCache {
		if err := database.DB.Preload("Order").First(&payment, paymentID).Error; err != nil {
			c.JSON(404, gin.H{
				"code": 404,
				"msg":  "支付记录不存在",
			})
			return
		}
	}

	// 检查权限
	if payment.UserID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权访问此支付记录",
		})
		return
	}

	msg := "获取成功"
	if fromCache {
		msg = "获取成功（来自缓存）"
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  msg,
		"data": gin.H{
			"id":              payment.ID,
			"payment_no":      payment.PaymentNo,
			"order_id":        payment.OrderID,
			"order_no":        payment.OrderNo,
			"amount":          payment.Amount,
			"payment_method":  payment.PaymentMethod,
			"trade_type":      payment.TradeType,
			"status":          payment.Status,
			"status_text":     utils.GetPaymentStatusText(payment.Status),
			"pay_time":        payment.PayTime,
			"transaction_id":  payment.TransactionID,
			"created_at":      payment.CreatedAt,
		},
	})
}

// GetUserPayments godoc
// @Summary 获取用户支付记录列表
// @Description 获取当前用户的支付记录列表
// @Tags 支付
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param status query int false "支付状态:0-待支付,1-已支付,2-支付失败,3-已退款,4-已取消"
// @Param payment_method query string false "支付方式:wechat/alipay"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{payments,total}"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Router /api/payment/list [get]
func GetUserPayments(c *gin.Context) {
	userID, _ := c.Get("user_id")

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	status := c.Query("status")
	paymentMethod := c.Query("payment_method")

	query := database.DB.Model(&models.Payment{}).Where("user_id = ?", userID)

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 支付方式筛选
	if paymentMethod != "" {
		query = query.Where("payment_method = ?", paymentMethod)
	}

	var total int64
	query.Count(&total)

	var payments []models.Payment
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (parseInt(page) - 1) * parseInt(pageSize)
	}

	if err := query.Preload("Order").Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&payments).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"payments": payments,
			"total":    total,
		},
	})
}

// RefundPayment godoc
// @Summary 申请退款
// @Description 申请支付退款
// @Tags 支付
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RefundPaymentRequest true "退款信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:退款申请成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "支付记录不存在"
// @Router /api/payment/refund [post]
func RefundPayment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	ctx := context.Background()

	var req RefundPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 查询支付记录
	var payment models.Payment
	if err := database.DB.Preload("Order").First(&payment, req.PaymentID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "支付记录不存在",
		})
		return
	}

	// 检查权限
	if payment.UserID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权操作此支付记录",
		})
		return
	}

	// 检查支付状态
	if payment.Status != models.PaymentStatusPaid {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "支付状态不允许退款",
		})
		return
	}

	// 检查退款金额
	if req.RefundAmount > payment.Amount {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "退款金额不能超过支付金额",
		})
		return
	}

	// 调用第三方退款接口（这里模拟）
	var refundSuccess bool
	if payment.PaymentMethod == models.PaymentMethodWeChat {
		refundSuccess = weChatRefund(&payment, req.RefundAmount, req.RefundReason)
	} else if payment.PaymentMethod == models.PaymentMethodAlipay {
		refundSuccess = alipayRefund(&payment, req.RefundAmount, req.RefundReason)
	}

	if !refundSuccess {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "退款失败",
		})
		return
	}

	// 更新支付状态
	updates := map[string]interface{}{
		"status": models.PaymentStatusRefunded,
	}

	if err := database.DB.Model(&payment).Updates(updates).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "更新支付记录失败: " + err.Error(),
		})
		return
	}

	// 更新订单状态
	if payment.Order.ID != 0 {
		database.DB.Model(&payment.Order).Update("status", models.OrderStatusRefunded)

		// 清除缓存
		cache.DeleteOrderCache(ctx, payment.Order.ID)
		cache.InvalidateUserOrdersCache(ctx, payment.Order.UserID)
		cache.InvalidateCounselorOrdersCache(ctx, payment.Order.CounselorID)
	}

	// 清除支付缓存
	cache.DeletePaymentCache(ctx, payment.ID)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "退款申请成功",
	})
}

// createWeChatPayment 创建微信支付
func createWeChatPayment(order *models.Order, payment *models.Payment, clientIP string) (string, map[string]string) {
	// 初始化微信支付
	wechatPay := utils.NewWeChatPay(
		"wx_app_id",
		"wx_mch_id",
		"wx_api_secret",
		"http://localhost:8080/api/payment/wechat/callback",
		true,
	)

	// 构建下单请求
	request := utils.UnifiedOrderRequest{
		Body:           "心理咨询订单",
		OutTradeNo:     payment.PaymentNo,
		TotalFee:       utils.ConvertYuanToFen(payment.Amount),
		SpbillCreateIP: clientIP,
		TradeType:      payment.TradeType,
	}

	// 调用微信支付API
	response, err := wechatPay.CreateUnifiedOrder(request)
	if err != nil {
		return "", nil
	}

	// 根据不同的交易类型返回不同的支付参数
	var payURL string
	var payParams map[string]string

	switch payment.TradeType {
	case "NATIVE":
		payURL = response.CodeURL
	case "H5":
		payURL = response.MWebURL
	case "JSAPI":
		payParams = wechatPay.GetJSAPIPayParams(response.PrepayID)
	case "APP":
		payParams = wechatPay.GetAppPayParams(response.PrepayID)
	}

	return payURL, payParams
}

// createAlipayPayment 创建支付宝支付
func createAlipayPayment(order *models.Order, payment *models.Payment, returnURL string) (string, map[string]string) {
	// 初始化支付宝
	alipay := utils.NewAlipay(
		"alipay_app_id",
		"alipay_private_key",
		"alipay_public_key",
		"http://localhost:8080/api/payment/alipay/callback",
		true,
	)

	// 构建下单请求
	request := utils.TradeCreateRequest{
		OutTradeNo:  payment.PaymentNo,
		TotalAmount: fmt.Sprintf("%.2f", payment.Amount),
		Subject:     "心理咨询订单",
		Body:        fmt.Sprintf("订单号: %s", order.OrderNo),
		ReturnURL:   returnURL,
		ProductCode: "FAST_INSTANT_TRADE_PAY",
	}

	// 根据不同的交易类型调用不同的API
	var payURL string
	var payParams map[string]string

	if payment.TradeType == "APP" {
		response, err := alipay.CreateAppTrade(request)
		if err == nil {
			payParams = map[string]string{
				"order_string": response.OrderInfo,
			}
		}
	} else {
		response, err := alipay.CreateTrade(request)
		if err == nil {
			payURL = response.PayURL
		}
	}

	return payURL, payParams
}

// weChatRefund 微信退款（模拟）
func weChatRefund(payment *models.Payment, refundAmount float64, reason string) bool {
	// 实际项目中需要调用微信退款API
	return true
}

// alipayRefund 支付宝退款（模拟）
func alipayRefund(payment *models.Payment, refundAmount float64, reason string) bool {
	// 实际项目中需要调用支付宝退款API
	return true
}

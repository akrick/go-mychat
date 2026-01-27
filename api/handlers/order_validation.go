package handlers

import (
	"fmt"
	"akrick.com/mychat/database"
	"akrick.com/mychat/models"
	"akrick.com/mychat/utils"
	"github.com/gin-gonic/gin"
)

// OrderValidation 订单验证工具
type OrderValidation struct {
	Order    *models.Order
	Counselor *models.Counselor
	Errors   []string
}

// NewOrderValidation 创建订单验证器
func NewOrderValidation(orderID uint) *OrderValidation {
	v := &OrderValidation{
		Errors: make([]string, 0),
	}

	// 查询订单
	if err := database.DB.First(&v.Order, orderID).Error; err != nil {
		v.Errors = append(v.Errors, "订单不存在")
		return v
	}

	// 查询咨询师信息
	database.DB.First(&v.Counselor, v.Order.CounselorID)

	return v
}

// CanPay 检查订单是否可以支付
func (v *OrderValidation) CanPay() bool {
	if v.Order == nil {
		v.Errors = append(v.Errors, "订单不存在")
		return false
	}

	if v.Order.Status != models.OrderStatusPending {
		v.Errors = append(v.Errors, "订单状态不允许支付")
		return false
	}

	if v.Counselor != nil && v.Counselor.Status != 1 {
		v.Errors = append(v.Errors, "咨询师暂不可用")
		return false
	}

	return true
}

// CanCancel 检查订单是否可以取消
func (v *OrderValidation) CanCancel() bool {
	if v.Order == nil {
		v.Errors = append(v.Errors, "订单不存在")
		return false
	}

	if v.Order.Status != models.OrderStatusPending && v.Order.Status != models.OrderStatusPaid {
		v.Errors = append(v.Errors, "订单状态不允许取消")
		return false
	}

	// 检查预约时间，如果预约时间在24小时内，不允许取消
	if v.Order.ScheduleTime.Before(utils.GetPaymentExpiryTime()) {
		v.Errors = append(v.Errors, "预约时间过近，无法取消")
		return false
	}

	return true
}

// CanRefund 检查订单是否可以退款
func (v *OrderValidation) CanRefund() bool {
	if v.Order == nil {
		v.Errors = append(v.Errors, "订单不存在")
		return false
	}

	if v.Order.Status != models.OrderStatusPaid && v.Order.Status != models.OrderStatusCompleted {
		v.Errors = append(v.Errors, "订单状态不允许退款")
		return false
	}

	return true
}

// CanReview 检查订单是否可以评价
func (v *OrderValidation) CanReview() bool {
	if v.Order == nil {
		v.Errors = append(v.Errors, "订单不存在")
		return false
	}

	if v.Order.Status != models.OrderStatusCompleted {
		v.Errors = append(v.Errors, "订单未完成，无法评价")
		return false
	}

	// 检查是否已评价
	var existingReview models.Review
	if err := database.DB.Where("order_id = ?", v.Order.ID).First(&existingReview).Error; err == nil {
		v.Errors = append(v.Errors, "该订单已评价")
		return false
	}

	return true
}

// CanComplete 检查订单是否可以完成
func (v *OrderValidation) CanComplete() bool {
	if v.Order == nil {
		v.Errors = append(v.Errors, "订单不存在")
		return false
	}

	if v.Order.Status != models.OrderStatusPaid {
		v.Errors = append(v.Errors, "订单未支付，无法完成")
		return false
	}

	return true
}

// ValidateOrderStatus godoc
// @Summary 验证订单状态
// @Description 验证订单是否可以执行指定操作
// @Tags 订单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Param action query string true "操作类型:pay,cancel,refund,review,complete"
// @Success 200 {object} map[string]interface{} "code:200,msg:验证成功,data:{valid,errors}"
// @Router /api/order/{id}/validate [get]
func ValidateOrderStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	orderID := c.Param("id")
	action := c.Query("action")

	// 查询订单
	var order models.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "订单不存在",
		})
		return
	}

	// 检查权限
	if order.UserID != userID.(uint) && order.CounselorID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权访问此订单",
		})
		return
	}

	// 创建验证器
	validation := NewOrderValidation(order.ID)

	var valid bool
	switch action {
	case "pay":
		valid = validation.CanPay()
	case "cancel":
		valid = validation.CanCancel()
	case "refund":
		valid = validation.CanRefund()
	case "review":
		valid = validation.CanReview()
	case "complete":
		valid = validation.CanComplete()
	default:
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "不支持的操作类型",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "验证成功",
		"data": gin.H{
			"valid":  valid,
			"errors": validation.Errors,
			"status": order.Status,
			"status_text": map[int]string{
				0: "待支付",
				1: "已支付",
				2: "已完成",
				3: "已取消",
				4: "已退款",
			}[order.Status],
		},
	})
}

// GetOrderTimeline godoc
// @Summary 获取订单时间线
// @Description 获取订单的状态变化时间线
// @Tags 订单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{timeline}"
// @Router /api/order/{id}/timeline [get]
func GetOrderTimeline(c *gin.Context) {
	userID, _ := c.Get("user_id")
	orderID := c.Param("id")

	// 查询订单
	var order models.Order
	if err := database.DB.Preload("User").Preload("Counselor").First(&order, orderID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "订单不存在",
		})
		return
	}

	// 检查权限
	if order.UserID != userID.(uint) && order.CounselorID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权访问此订单",
		})
		return
	}

	// 构建时间线
	timeline := make([]map[string]interface{}, 0)

	// 创建订单
	timeline = append(timeline, map[string]interface{}{
		"event":   "创建订单",
		"time":    order.CreatedAt,
		"message": fmt.Sprintf("用户 %s 创建了订单", order.User.Username),
	})

	// 支付订单
	if order.PayTime != nil {
		timeline = append(timeline, map[string]interface{}{
			"event":   "支付订单",
			"time":    *order.PayTime,
			"message": fmt.Sprintf("订单已支付，金额：%.2f元", order.Amount),
		})
	}

	// 预约时间
	timeline = append(timeline, map[string]interface{}{
		"event":   "预约咨询",
		"time":    order.ScheduleTime,
		"message": fmt.Sprintf("预约咨询时间：%s", order.ScheduleTime.Format("2006-01-02 15:04:05")),
	})

	// 根据订单状态添加时间线
	switch order.Status {
	case models.OrderStatusCompleted:
		timeline = append(timeline, map[string]interface{}{
			"event":   "完成咨询",
			"time":    order.UpdatedAt,
			"message": fmt.Sprintf("咨询已完成，咨询师：%s", order.Counselor.Name),
		})
	case models.OrderStatusCancelled:
		timeline = append(timeline, map[string]interface{}{
			"event":   "取消订单",
			"time":    order.UpdatedAt,
			"message": "订单已取消",
		})
	case models.OrderStatusRefunded:
		timeline = append(timeline, map[string]interface{}{
			"event":   "退款成功",
			"time":    order.UpdatedAt,
			"message": "订单已退款",
		})
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"timeline": timeline,
		},
	})
}

package handlers

import (
	"akrick.com/mychat/database"
	"akrick.com/mychat/models"
	"github.com/gin-gonic/gin"
)

// GetOrderList godoc
// @Summary 获取订单列表
// @Description 获取订单列表（管理员接口）
// @Tags 订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param status query int false "订单状态:0-待支付,1-已支付,2-已完成,3-已取消,4-已退款"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/orders [get]
func GetOrderList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	status := c.Query("status")
	keyword := c.Query("keyword")

	query := database.DB.Model(&models.Order{})

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 搜索
	if keyword != "" {
		query = query.Where("order_no LIKE ?", "%"+keyword+"%").
			Or("notes LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var orders []models.Order
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (parseInt(page) - 1) * parseInt(pageSize)
	}

	if err := query.Preload("User").Preload("Counselor").
		Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&orders).Error; err != nil {
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
			"orders": orders,
			"total":  total,
		},
	})
}

// GetOrderStatistics godoc
// @Summary 获取订单统计
// @Description 获取订单统计数据（管理员接口）
// @Tags 订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/orders/statistics [get]
func GetOrderStatistics(c *gin.Context) {
	var totalOrders int64
	var pendingOrders int64
	var paidOrders int64
	var completedOrders int64
	var cancelledOrders int64
	var totalAmount float64
	var todayAmount float64
	var thisMonthAmount float64

	database.DB.Model(&models.Order{}).Count(&totalOrders)
	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderStatusPending).Count(&pendingOrders)
	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderStatusPaid).Count(&paidOrders)
	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderStatusCompleted).Count(&completedOrders)
	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderStatusCancelled).Count(&cancelledOrders)

	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderStatusPaid).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount)

	database.DB.Model(&models.Order{}).Where("status = ? AND DATE(created_at) = CURDATE()",
		models.OrderStatusPaid).
		Select("COALESCE(SUM(amount), 0)").Scan(&todayAmount)

	database.DB.Model(&models.Order{}).Where("status = ? AND YEAR(created_at) = YEAR(NOW()) AND MONTH(created_at) = MONTH(NOW())",
		models.OrderStatusPaid).
		Select("COALESCE(SUM(amount), 0)").Scan(&thisMonthAmount)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"total_orders":       totalOrders,
			"pending_orders":     pendingOrders,
			"paid_orders":        paidOrders,
			"completed_orders":   completedOrders,
			"cancelled_orders":   cancelledOrders,
			"total_amount":       totalAmount,
			"today_amount":       todayAmount,
			"this_month_amount":  thisMonthAmount,
		},
	})
}

// AdminUpdateOrderStatus godoc
// @Summary 更新订单状态
// @Description 管理员更新订单状态
// @Tags 订单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Param request body map[string]interface{} true "状态信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:更新成功"
// @Router /api/admin/orders/{id}/status [put]
func AdminUpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")

	var req struct {
		Status int `json:"status" binding:"required,oneof=0 1 2 3 4"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	var order models.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "订单不存在",
		})
		return
	}

	if err := database.DB.Model(&order).Update("status", req.Status).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "更新失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

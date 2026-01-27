package handlers

import (
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"
	"github.com/gin-gonic/gin"
	"time"
)

// DashboardStatistics godoc
// @Summary 获取仪表盘统计数据
// @Description 获取系统整体统计数据（管理员）
// @Tags 统计
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{statistics}"
// @Router /api/stats/dashboard [get]
func DashboardStatistics(c *gin.Context) {
	var statistics map[string]interface{}
	statistics = make(map[string]interface{})

	// 用户统计
	var userCount int64
	database.DB.Model(&models.User{}).Count(&userCount)
	statistics["user_count"] = userCount

	// 咨询师统计
	var counselorCount int64
	database.DB.Model(&models.Counselor{}).Count(&counselorCount)
	statistics["counselor_count"] = counselorCount

	// 订单统计
	var orderCount int64
	database.DB.Model(&models.Order{}).Count(&orderCount)
	statistics["order_count"] = orderCount

	// 今日订单
	var todayOrderCount int64
	today := time.Now().Format("2006-01-02")
	database.DB.Model(&models.Order{}).Where("DATE(created_at) = ?", today).Count(&todayOrderCount)
	statistics["today_order_count"] = todayOrderCount

	// 支付订单统计
	var paidOrderCount int64
	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderStatusPaid).Count(&paidOrderCount)
	statistics["paid_order_count"] = paidOrderCount

	// 已完成订单统计
	var completedOrderCount int64
	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderStatusCompleted).Count(&completedOrderCount)
	statistics["completed_order_count"] = completedOrderCount

	// 评价统计
	var reviewCount int64
	database.DB.Model(&models.Review{}).Count(&reviewCount)
	statistics["review_count"] = reviewCount

	// 今日收入
	var todayIncome float64
	database.DB.Model(&models.Order{}).
		Where("DATE(pay_time) = ? AND status = ?", today, models.OrderStatusPaid).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&todayIncome)
	statistics["today_income"] = todayIncome

	// 总收入
	var totalIncome float64
	database.DB.Model(&models.Order{}).
		Where("status >= ?", models.OrderStatusPaid).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome)
	statistics["total_income"] = totalIncome

	// 近7天订单趋势
	var orderTrend []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var count int64
		database.DB.Model(&models.Order{}).Where("DATE(created_at) = ?", date).Count(&count)
		orderTrend = append(orderTrend, struct {
			Date  string `json:"date"`
			Count int64  `json:"count"`
		}{
			Date:  date,
			Count: count,
		})
	}
	statistics["order_trend"] = orderTrend

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": statistics,
	})
}

// OrderStatistics godoc
// @Summary 获取订单统计数据
// @Description 获取订单相关的统计数据
// @Tags 统计
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "开始日期" default("2024-01-01")
// @Param end_date query string false "结束日期" default()
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{statistics}"
// @Router /api/stats/order [get]
func OrderStatistics(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", "2024-01-01")
	endDate := c.Query("end_date")

	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	query := database.DB.Model(&models.Order{}).Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate)

	var statistics map[string]interface{}
	statistics = make(map[string]interface{})

	// 总订单数
	var totalOrders int64
	query.Count(&totalOrders)
	statistics["total_orders"] = totalOrders

	// 待支付订单
	var pendingOrders int64
	query.Where("status = ?", models.OrderStatusPending).Count(&pendingOrders)
	statistics["pending_orders"] = pendingOrders

	// 已支付订单
	var paidOrders int64
	database.DB.Model(&models.Order{}).
		Where("DATE(created_at) BETWEEN ? AND ? AND status = ?", startDate, endDate, models.OrderStatusPaid).
		Count(&paidOrders)
	statistics["paid_orders"] = paidOrders

	// 已完成订单
	var completedOrders int64
	database.DB.Model(&models.Order{}).
		Where("DATE(created_at) BETWEEN ? AND ? AND status = ?", startDate, endDate, models.OrderStatusCompleted).
		Count(&completedOrders)
	statistics["completed_orders"] = completedOrders

	// 已取消订单
	var cancelledOrders int64
	database.DB.Model(&models.Order{}).
		Where("DATE(created_at) BETWEEN ? AND ? AND status = ?", startDate, endDate, models.OrderStatusCancelled).
		Count(&cancelledOrders)
	statistics["cancelled_orders"] = cancelledOrders

	// 退款订单
	var refundedOrders int64
	database.DB.Model(&models.Order{}).
		Where("DATE(created_at) BETWEEN ? AND ? AND status = ?", startDate, endDate, models.OrderStatusRefunded).
		Count(&refundedOrders)
	statistics["refunded_orders"] = refundedOrders

	// 总金额
	var totalAmount float64
	database.DB.Model(&models.Order{}).
		Where("DATE(created_at) BETWEEN ? AND ? AND status >= ?", startDate, endDate, models.OrderStatusPaid).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalAmount)
	statistics["total_amount"] = totalAmount

	// 总咨询时长
	var totalDuration int
	database.DB.Model(&models.Order{}).
		Where("DATE(created_at) BETWEEN ? AND ? AND status >= ?", startDate, endDate, models.OrderStatusPaid).
		Select("COALESCE(SUM(duration), 0)").
		Scan(&totalDuration)
	statistics["total_duration"] = totalDuration

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": statistics,
	})
}

// CounselorRanking godoc
// @Summary 获取咨询师排行榜
// @Description 按不同维度获取咨询师排行榜
// @Tags 统计
// @Accept json
// @Produce json
// @Param type query string false "排行榜类型:orders-订单数,income-收入,rating-评分" default("orders")
// @Param limit query int false "返回数量" default(10)
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{ranking}"
// @Router /api/stats/counselor/ranking [get]
func CounselorRanking(c *gin.Context) {
	rankType := c.DefaultQuery("type", "orders")
	limit := c.DefaultQuery("limit", "10")

	var counselors []models.Counselor

	switch rankType {
	case "orders":
		// 按订单数排名
		database.DB.
			Select("counselors.*, COUNT(orders.id) as order_count").
			Joins("LEFT JOIN orders ON counselors.id = orders.counselor_id AND orders.status >= ?", models.OrderStatusPaid).
			Group("counselors.id").
			Order("order_count DESC").
			Limit(parseInt(limit)).
			Find(&counselors)

	case "income":
		// 按收入排名
		database.DB.
			Select("counselors.*, COALESCE(SUM(orders.amount), 0) as total_income").
			Joins("LEFT JOIN orders ON counselors.id = orders.counselor_id AND orders.status >= ?", models.OrderStatusPaid).
			Group("counselors.id").
			Order("total_income DESC").
			Limit(parseInt(limit)).
			Find(&counselors)

	case "rating":
		// 按评分排名
		database.DB.
			Select("counselors.*").
			Where("status = ?", 1).
			Order("rating DESC").
			Limit(parseInt(limit)).
			Find(&counselors)

	default:
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "不支持的排行榜类型",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"type":      rankType,
			"counselors": counselors,
		},
	})
}

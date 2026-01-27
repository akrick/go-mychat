package handlers

import (
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

// GetFinanceStats godoc
// @Summary 获取财务统计数据
// @Description 获取财务相关的统计数据
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/finance/stats [get]
func GetFinanceStats(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")

	var totalRevenue float64
	var totalWithdrawn float64
	var totalCommission float64
	var todayRevenue float64
	var todayOrders int64
	var pendingWithdraws int64
	var approvedWithdraws int64

	query := database.DB.Model(&models.Order{}).Where("status = ?", models.OrderStatusPaid)
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	// 总营收
	query.Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue)

	// 今日数据
	today := time.Now().Format("2006-01-02")
	database.DB.Model(&models.Order{}).
		Where("status = ? AND DATE(created_at) = ?", models.OrderStatusPaid, today).
		Select("COALESCE(SUM(amount), 0)").Scan(&todayRevenue)
	database.DB.Model(&models.Order{}).
		Where("status = ? AND DATE(created_at) = ?", models.OrderStatusPaid, today).
		Count(&todayOrders)

	// 提现统计
	database.DB.Model(&models.WithdrawRecord{}).
		Where("status = ?", 3).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalWithdrawn)

	database.DB.Model(&models.WithdrawRecord{}).
		Where("status = ?", 0).
		Count(&pendingWithdraws)

	database.DB.Model(&models.WithdrawRecord{}).
		Where("status = ?", 3).
		Count(&approvedWithdraws)

	// 平台佣金 (假设佣金率为20%)
	totalCommission = totalRevenue * 0.2

	// 咨询师总收益
	totalCounselorEarnings := totalRevenue - totalCommission

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"total_revenue":          totalRevenue,
			"total_withdrawn":        totalWithdrawn,
			"total_commission":       totalCommission,
			"total_counselor_earning": totalCounselorEarnings,
			"today_revenue":          todayRevenue,
			"today_orders":           todayOrders,
			"pending_withdraws":      pendingWithdraws,
			"approved_withdraws":     approvedWithdraws,
		},
	})
}

// GetRevenueReport godoc
// @Summary 获取营收报表
// @Description 获取指定时间段内的营收报表数据
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Param group_by query string false "分组方式:day,month,year" Enums(day,month,year)
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/finance/revenue [get]
func GetRevenueReport(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")
	groupBy := c.DefaultQuery("group_by", "day")

	var dateFormat string
	switch groupBy {
	case "month":
		dateFormat = "%Y-%m"
	case "year":
		dateFormat = "%Y"
	default:
		dateFormat = "%Y-%m-%d"
	}

	type RevenueData struct {
		Date   string  `json:"date"`
		Amount float64 `json:"amount"`
		Count  int64   `json:"count"`
	}

	var revenueData []RevenueData

	query := database.DB.Model(&models.Order{}).
		Select(dateFormat+" as date, COALESCE(SUM(amount), 0) as amount, COUNT(*) as count").
		Where("status = ?", models.OrderStatusPaid)

	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	query.Group(dateFormat).Order("date ASC").Scan(&revenueData)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"revenue_data": revenueData,
			"group_by":     groupBy,
		},
	})
}

// GetWithdrawList godoc
// @Summary 获取提现记录列表
// @Description 获取所有提现记录
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param status query int false "状态"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/withdraws [get]
func GetWithdrawList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	status := c.Query("status")

	query := database.DB.Model(&models.WithdrawRecord{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var withdraws []models.WithdrawRecord
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (parseInt(page) - 1) * parseInt(pageSize)
	}

	if err := query.Preload("Counselor").
		Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&withdraws).Error; err != nil {
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
			"withdraws": withdraws,
			"total":     total,
		},
	})
}

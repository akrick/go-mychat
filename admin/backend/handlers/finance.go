package handlers

import (
	"context"
	"akrick.com/mychat/admin/backend/cache"
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"
	"akrick.com/mychat/admin/backend/utils"
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
	var pendingWithdrawAmount float64

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
		Where("status = ?", 0).
		Select("COALESCE(SUM(amount), 0)").Scan(&pendingWithdrawAmount)

	database.DB.Model(&models.WithdrawRecord{}).
		Where("status = ?", 3).
		Count(&approvedWithdraws)

	// 平台佣金 (假设佣金率为20%)
	totalCommission = totalRevenue * 0.2

	// 咨询师总收益
	totalCounselorEarnings := totalRevenue - totalCommission

	// 获取所有咨询师账户总余额
	var totalCounselorBalance float64
	database.DB.Model(&models.CounselorAccount{}).
		Select("COALESCE(SUM(balance), 0)").Scan(&totalCounselorBalance)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"total_revenue":            totalRevenue,
			"total_withdrawn":          totalWithdrawn,
			"total_commission":         totalCommission,
			"total_counselor_earning":   totalCounselorEarnings,
			"total_counselor_balance":  totalCounselorBalance,
			"pending_withdraws":        pendingWithdraws,
			"pending_withdraw_amount":   pendingWithdrawAmount,
			"approved_withdraws":       approvedWithdraws,
			"today_revenue":            todayRevenue,
			"today_orders":             todayOrders,
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
// @Param counselor_id query int false "咨询师ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/withdraws [get]
func GetWithdrawList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	status := c.Query("status")
	counselorID := c.Query("counselor_id")

	query := database.DB.Model(&models.WithdrawRecord{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if counselorID != "" {
		query = query.Where("counselor_id = ?", counselorID)
	}

	var total int64
	query.Count(&total)

	var withdraws []models.WithdrawRecord
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (utils.ParseInt(page) - 1) * utils.ParseInt(pageSize)
	}

	if err := query.Preload("Counselor").
		Offset(offset).Limit(utils.ParseInt(pageSize)).Order("created_at DESC").Find(&withdraws).Error; err != nil {
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

// GetCounselorAccountList godoc
// @Summary 获取咨询师账户列表
// @Description 获取所有咨询师账户信息
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param counselor_id query int false "咨询师ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/finance/accounts [get]
func GetCounselorAccountList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	counselorID := c.Query("counselor_id")

	query := database.DB.Model(&models.CounselorAccount{})
	if counselorID != "" {
		query = query.Where("counselor_id = ?", counselorID)
	}

	var total int64
	query.Count(&total)

	var accounts []models.CounselorAccount
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (utils.ParseInt(page) - 1) * utils.ParseInt(pageSize)
	}

	if err := query.Preload("Counselor").
		Offset(offset).Limit(utils.ParseInt(pageSize)).Order("created_at DESC").Find(&accounts).Error; err != nil {
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
			"accounts": accounts,
			"total":    total,
		},
	})
}

// GetCounselorAccountDetail godoc
// @Summary 获取咨询师账户详情
// @Description 获取指定咨询师的账户详情
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "咨询师ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/finance/accounts/{id} [get]
func GetCounselorAccountDetail(c *gin.Context) {
	id := c.Param("id")

	var account models.CounselorAccount
	if err := database.DB.Preload("Counselor").Where("counselor_id = ?", id).First(&account).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "账户不存在",
		})
		return
	}

	// 获取提现记录统计
	var withdrawCount int64
	var withdrawTotal float64
	database.DB.Model(&models.WithdrawRecord{}).
		Where("counselor_id = ? AND status IN (1,3)", id).
		Count(&withdrawCount)
	database.DB.Model(&models.WithdrawRecord{}).
		Where("counselor_id = ? AND status IN (1,3)", id).
		Select("COALESCE(SUM(amount), 0)").Scan(&withdrawTotal)

	// 获取待审核提现
	var pendingWithdraws []models.WithdrawRecord
	database.DB.Where("counselor_id = ? AND status = ?", id, 0).
		Order("created_at DESC").Find(&pendingWithdraws)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"account":           account,
			"withdraw_count":    withdrawCount,
			"withdraw_total":    withdrawTotal,
			"pending_withdraws": pendingWithdraws,
		},
	})
}

// ConfirmWithdrawTransfer godoc
// @Summary 确认提现打款
// @Description 管理员确认提现已完成打款
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "提现记录ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:确认成功"
// @Router /api/admin/withdraw/{id}/transfer [post]
func ConfirmWithdrawTransfer(c *gin.Context) {
	withdrawID := c.Param("id")

	var withdraw models.WithdrawRecord
	if err := database.DB.First(&withdraw, withdrawID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "提现记录不存在",
		})
		return
	}

	// 检查状态
	if withdraw.Status != 1 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "提现记录状态不正确",
		})
		return
	}

	now := time.Now()
	
	// 更新提现记录状态为已打款
	database.DB.Model(&withdraw).Updates(map[string]interface{}{
		"status":         3,
		"transferred_at":  &now,
	})

	// 更新咨询师账户 - 减少已提现金额
	var account models.CounselorAccount
	database.DB.Where("counselor_id = ?", withdraw.CounselorID).First(&account)
	database.DB.Model(&account).Updates(map[string]interface{}{
		"withdrawn": account.Withdrawn + withdraw.Amount,
	})

	// 清除缓存
	ctx := context.Background()
	cache.DeleteCounselorAccountCache(ctx, withdraw.CounselorID)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "确认打款成功",
		"data": withdraw,
	})
}

// GetFinanceReports godoc
// @Summary 获取财务报表
// @Description 获取综合财务报表
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Param type query string false "报表类型:income,withdraw,account" Enums(income,withdraw,account)
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/finance/reports [get]
func GetFinanceReports(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")
	reportType := c.DefaultQuery("type", "income")

	var data interface{}

	switch reportType {
	case "income":
		// 收入报表
		type IncomeData struct {
			Date        string  `json:"date"`
			Amount      float64 `json:"amount"`
			OrderCount  int64   `json:"order_count"`
			Commission  float64 `json:"commission"`
			CounselorEarnings float64 `json:"counselor_earnings"`
		}
		var incomeData []IncomeData

		query := database.DB.Model(&models.Order{}).
			Select("DATE(created_at) as date, COALESCE(SUM(amount), 0) as amount, COUNT(*) as order_count").
			Where("status = ?", models.OrderStatusPaid)

		if startDate != "" {
			query = query.Where("created_at >= ?", startDate)
		}
		if endDate != "" {
			query = query.Where("created_at <= ?", endDate)
		}

		query.Group("DATE(created_at)").Order("date ASC").Scan(&incomeData)

		// 计算佣金和咨询师收益
		for i := range incomeData {
			incomeData[i].Commission = incomeData[i].Amount * 0.2
			incomeData[i].CounselorEarnings = incomeData[i].Amount - incomeData[i].Commission
		}

		data = incomeData

	case "withdraw":
		// 提现报表
		type WithdrawData struct {
			Date          string  `json:"date"`
			Amount        float64 `json:"amount"`
			Count         int64   `json:"count"`
			ApprovedCount int64   `json:"approved_count"`
			RejectedCount int64   `json:"rejected_count"`
		}
		var withdrawData []WithdrawData

		query := database.DB.Model(&models.WithdrawRecord{}).
			Select("DATE(created_at) as date, COALESCE(SUM(amount), 0) as amount, COUNT(*) as count, "+
				"SUM(CASE WHEN status IN (1,3) THEN 1 ELSE 0 END) as approved_count, "+
				"SUM(CASE WHEN status = 2 THEN 1 ELSE 0 END) as rejected_count")

		if startDate != "" {
			query = query.Where("created_at >= ?", startDate)
		}
		if endDate != "" {
			query = query.Where("created_at <= ?", endDate)
		}

		query.Group("DATE(created_at)").Order("date ASC").Scan(&withdrawData)

		data = withdrawData

	case "account":
		// 账户报表
		type AccountData struct {
			CounselorID   uint    `json:"counselor_id"`
			CounselorName string  `json:"counselor_name"`
			TotalIncome   float64 `json:"total_income"`
			Balance       float64 `json:"balance"`
			Withdrawn     float64 `json:"withdrawn"`
			FrozenAmount  float64 `json:"frozen_amount"`
		}
		var accountData []AccountData

		database.DB.Model(&models.CounselorAccount{}).
			Select("counselor_accounts.counselor_id, counselors.name as counselor_name, "+
				"counselor_accounts.total_income, counselor_accounts.balance, "+
				"counselor_accounts.withdrawn, counselor_accounts.frozen_amount").
			Joins("LEFT JOIN counselors ON counselors.id = counselor_accounts.counselor_id").
			Order("counselor_accounts.total_income DESC").
			Scan(&accountData)

		data = accountData
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"report_type": reportType,
			"start_date":  startDate,
			"end_date":    endDate,
			"records":     data,
		},
	})
}


package handlers

import (
	"akrick.com/mychat/database"
	"akrick.com/mychat/models"
	"github.com/gin-gonic/gin"
)

// GetSystemLogList 获取系统日志列表
// @Summary 获取系统日志列表
// @Description 获取系统操作日志列表（管理员接口）
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param module query string false "模块名称"
// @Param username query string false "用户名"
// @Param action query string false "操作类型"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/logs [get]
func GetSystemLogList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	module := c.Query("module")
	username := c.Query("username")
	action := c.Query("action")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	query := database.DB.Model(&models.SysLog{})

	if module != "" {
		query = query.Where("module = ?", module)
	}
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if startDate != "" {
		query = query.Where("DATE(created_at) >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("DATE(created_at) <= ?", endDate)
	}

	var total int64
	query.Count(&total)

	var logs []models.SysLog
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (parseInt(page) - 1) * parseInt(pageSize)
	}

	if err := query.Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&logs).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":  logs,
			"total": total,
		},
	})
}

// GetOnlineUserList 获取在线用户列表
// @Summary 获取在线用户列表
// @Description 获取所有在线用户列表（管理员接口）
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/online/users [get]
func GetOnlineUserList(c *gin.Context) {
	var onlineUsers []models.OnlineUser
	if err := database.DB.Preload("User").Find(&onlineUsers).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"users": onlineUsers,
			"total": len(onlineUsers),
		},
	})
}

// GetConfigList 获取系统配置列表
// @Summary 获取系统配置列表
// @Description 获取系统配置列表（管理员接口）
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param is_public query bool false "是否公开"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/configs [get]
func GetConfigList(c *gin.Context) {
	isPublic := c.Query("is_public")
	query := database.DB.Model(&models.SysConfig{})

	if isPublic == "true" {
		query = query.Where("is_public = ?", true)
	}

	var configs []models.SysConfig
	if err := query.Order("config_key ASC").Find(&configs).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": configs,
	})
}

// UpdateConfig 更新系统配置
// @Summary 更新系统配置
// @Description 更新系统配置（管理员接口）
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "配置ID"
// @Param request body map[string]interface{} true "配置信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:更新成功"
// @Router /api/admin/configs/{id} [put]
func UpdateConfig(c *gin.Context) {
	configID := c.Param("id")

	var req struct {
		ConfigVal string `json:"config_val" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	if err := database.DB.Model(&models.SysConfig{}).Where("id = ?", configID).Update("config_val", req.ConfigVal).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "更新失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// GetDashboardStatistics 获取Dashboard统计数据
// @Summary 获取Dashboard统计数据
// @Description 获取Dashboard统计数据（管理员接口）
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/dashboard/statistics [get]
func GetDashboardStatistics(c *gin.Context) {
	stats := gin.H{}

	// 用户统计
	var userCount int64
	database.DB.Model(&models.User{}).Count(&userCount)
	stats["user_count"] = userCount

	var todayUserCount int64
	database.DB.Model(&models.User{}).Where("DATE(created_at) = CURDATE()").Count(&todayUserCount)
	stats["today_user_count"] = todayUserCount

	// 咨询师统计
	var counselorCount int64
	database.DB.Model(&models.Counselor{}).Where("status = ?", 1).Count(&counselorCount)
	stats["counselor_count"] = counselorCount

	// 订单统计
	var orderCount int64
	database.DB.Model(&models.Order{}).Count(&orderCount)
	stats["order_count"] = orderCount

	var todayOrderCount int64
	database.DB.Model(&models.Order{}).Where("DATE(created_at) = CURDATE()").Count(&todayOrderCount)
	stats["today_order_count"] = todayOrderCount

	var completedOrderCount int64
	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderStatusCompleted).Count(&completedOrderCount)
	stats["completed_order_count"] = completedOrderCount

	// 交易额统计
	var totalAmount float64
	database.DB.Model(&models.Payment{}).Where("status = ?", models.PaymentStatusSuccess).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount)
	stats["total_amount"] = totalAmount

	var todayAmount float64
	database.DB.Model(&models.Payment{}).Where("status = ? AND DATE(created_at) = CURDATE()", models.PaymentStatusSuccess).
		Select("COALESCE(SUM(amount), 0)").Scan(&todayAmount)
	stats["today_amount"] = todayAmount

	// 会话统计
	var sessionCount int64
	database.DB.Model(&models.ChatSession{}).Count(&sessionCount)
	stats["session_count"] = sessionCount

	var activeSessionCount int64
	database.DB.Model(&models.ChatSession{}).Where("status = ?", models.SessionStatusActive).Count(&activeSessionCount)
	stats["active_session_count"] = activeSessionCount

	// 消息统计
	var messageCount int64
	database.DB.Model(&models.ChatMessage{}).Count(&messageCount)
	stats["message_count"] = messageCount

	var todayMessageCount int64
	database.DB.Model(&models.ChatMessage{}).Where("DATE(created_at) = CURDATE()").Count(&todayMessageCount)
	stats["today_message_count"] = todayMessageCount

	// 提现统计
	var pendingWithdrawCount int64
	database.DB.Model(&models.WithdrawRecord{}).Where("status = ?", models.WithdrawStatusPending).Count(&pendingWithdrawCount)
	stats["pending_withdraw_count"] = pendingWithdrawCount

	var totalWithdrawAmount float64
	database.DB.Model(&models.WithdrawRecord{}).Where("status = ?", models.WithdrawStatusCompleted).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalWithdrawAmount)
	stats["total_withdraw_amount"] = totalWithdrawAmount

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": stats,
	})
}

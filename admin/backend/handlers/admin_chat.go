package handlers

import (
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"
	"github.com/gin-gonic/gin"
)

// GetAdminChatSessions godoc
// @Summary 获取聊天会话列表
// @Description 获取聊天会话列表（管理员接口）
// @Tags 聊天管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param status query int false "状态:0-待开始,1-进行中,2-已结束,3-已超时"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/chat/sessions [get]
func GetAdminChatSessions(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	status := c.Query("status")
	keyword := c.Query("keyword")

	query := database.DB.Model(&models.ChatSession{})

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 搜索
	if keyword != "" {
		query = query.Joins("LEFT JOIN users u ON user_id = u.id").
			Joins("LEFT JOIN counselors c ON counselor_id = c.id").
			Where("u.username LIKE ? OR c.name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var sessions []models.ChatSession
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (parseInt(page) - 1) * parseInt(pageSize)
	}

	if err := query.Preload("User").Preload("Counselor").Preload("Order").
		Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&sessions).Error; err != nil {
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
			"sessions": sessions,
			"total":    total,
		},
	})
}

// GetAdminChatMessages godoc
// @Summary 获取聊天消息列表
// @Description 获取指定会话的聊天消息列表（管理员接口）
// @Tags 聊天管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param session_id path int true "会话ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(50)
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/chat/sessions/{session_id}/messages [get]
func GetAdminChatMessages(c *gin.Context) {
	sessionID := c.Param("session_id")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "50")

	// 检查会话是否存在
	var session models.ChatSession
	if err := database.DB.First(&session, sessionID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "会话不存在",
		})
		return
	}

	query := database.DB.Model(&models.ChatMessage{}).Where("session_id = ?", sessionID)

	var total int64
	query.Count(&total)

	var messages []models.ChatMessage
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (parseInt(page) - 1) * parseInt(pageSize)
	}

	if err := query.Offset(offset).Limit(parseInt(pageSize)).Order("created_at ASC").Find(&messages).Error; err != nil {
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
			"messages": messages,
			"total":    total,
			"session":  session,
		},
	})
}

// GetChatStatistics godoc
// @Summary 获取聊天统计
// @Description 获取聊天统计数据（管理员接口）
// @Tags 聊天管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/chat/statistics [get]
func GetChatStatistics(c *gin.Context) {
	var totalSessions int64
	var activeSessions int64
	var totalMessages int64
	var todayMessages int64
	var totalBilling float64

	database.DB.Model(&models.ChatSession{}).Count(&totalSessions)
	database.DB.Model(&models.ChatSession{}).Where("status = ?", 1).Count(&activeSessions)
	database.DB.Model(&models.ChatMessage{}).Count(&totalMessages)

	today := "NOW()"
	database.DB.Model(&models.ChatMessage{}).Where("DATE(created_at) = ?", today).Count(&todayMessages)

	database.DB.Model(&models.ChatBilling{}).Select("COALESCE(SUM(total_amount), 0)").Scan(&totalBilling)

	// 获取最近的消息
	var recentMessages []models.ChatMessage
	database.DB.Preload("Session").
		Order("created_at DESC").
		Limit(10).
		Find(&recentMessages)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"total_sessions":   totalSessions,
			"active_sessions":  activeSessions,
			"total_messages":   totalMessages,
			"today_messages":   todayMessages,
			"total_billing":    totalBilling,
			"recent_messages":  recentMessages,
		},
	})
}

// SearchChatMessages godoc
// @Summary 搜索聊天消息
// @Description 搜索聊天消息（管理员接口）
// @Tags 聊天管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/chat/messages/search [get]
func SearchChatMessages(c *gin.Context) {
	keyword := c.Query("keyword")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	if keyword == "" {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "搜索关键词不能为空",
		})
		return
	}

	query := database.DB.Model(&models.ChatMessage{}).
		Where("content LIKE ?", "%"+keyword+"%")

	var total int64
	query.Count(&total)

	var messages []models.ChatMessage
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (parseInt(page) - 1) * parseInt(pageSize)
	}

	if err := query.Preload("Session").Preload("Session.User").Preload("Session.Counselor").
		Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&messages).Error; err != nil {
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
			"messages": messages,
			"total":    total,
		},
	})
}

// DeleteChatSession godoc
// @Summary 删除聊天会话
// @Description 删除聊天会话（管理员接口）
// @Tags 聊天管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "会话ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:删除成功"
// @Router /api/admin/chat/sessions/{id} [delete]
func DeleteChatSession(c *gin.Context) {
	sessionID := c.Param("id")

	var session models.ChatSession
	if err := database.DB.First(&session, sessionID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "会话不存在",
		})
		return
	}

	// 开始事务
	tx := database.DB.Begin()

	// 删除会话的消息
	if err := tx.Where("session_id = ?", sessionID).Delete(&models.ChatMessage{}).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "删除消息失败: " + err.Error(),
		})
		return
	}

	// 删除会话
	if err := tx.Delete(&session).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "删除会话失败: " + err.Error(),
		})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

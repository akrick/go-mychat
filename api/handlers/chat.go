package handlers

import (
	"akrick.com/mychat/database"
	"akrick.com/mychat/models"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

// StartChatSession godoc
// @Summary 开始聊天会话
// @Description 创建聊天会话（咨询师发起）
// @Tags 聊天
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order_id path int true "订单ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:创建成功,data:{session_id}"
// @Router /api/chat/session/{order_id}/start [post]
func StartChatSession(c *gin.Context) {
	userID, _ := c.Get("user_id")
	orderID := c.Param("order_id")

	// 查询订单
	var order models.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "订单不存在",
		})
		return
	}

	// 检查是否为咨询师
	if order.CounselorID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "只有咨询师可以发起会话",
		})
		return
	}

	// 检查订单状态
	if order.Status != models.OrderStatusPaid {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "订单状态不允许开始会话",
		})
		return
	}

	// 检查是否已存在会话
	var existingSession models.ChatSession
	if err := database.DB.Where("order_id = ?", orderID).First(&existingSession).Error; err == nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "会话已存在",
			"data": gin.H{
				"session_id": existingSession.ID,
			},
		})
		return
	}

	// 创建会话
	now := time.Now()
	session := models.ChatSession{
		OrderID:     order.ID,
		UserID:      order.UserID,
		CounselorID: order.CounselorID,
		Status:      1, // 进行中
		StartTime:   &now,
	}

	if err := database.DB.Create(&session).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "创建会话失败: " + err.Error(),
		})
		return
	}

	// 发送通知给用户
	go CreateNotification(order.UserID, models.NotificationTypeChat, models.NotificationLevelInfo, "咨询会话已开始", "您的咨询会话已经开始，请及时参与", "")

	// 更新订单状态
	database.DB.Model(&order).Update("status", models.OrderStatusCompleted)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": gin.H{
			"session_id": session.ID,
		},
	})
}

// SendMessage godoc
// @Summary 发送消息
// @Description 发送聊天消息
// @Tags 聊天
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param session_id path int true "会话ID"
// @Param request body map[string]interface{} true "消息内容:content,sender_type,content_type,file_url"
// @Success 200 {object} map[string]interface{} "code:200,msg:发送成功,data:{message}"
// @Router /api/chat/session/{session_id}/message [post]
func SendMessage(c *gin.Context) {
	userID, _ := c.Get("user_id")
	sessionID := c.Param("session_id")

	var req struct {
		Content     string `json:"content" binding:"required"`
		SenderType  string `json:"sender_type" binding:"required,oneof=user counselor"`
		ContentType string `json:"content_type"`
		FileURL     string `json:"file_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 查询会话
	var session models.ChatSession
	if err := database.DB.First(&session, sessionID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "会话不存在",
		})
		return
	}

	// 检查会话状态
	if session.Status != 1 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "会话未进行中",
		})
		return
	}

	// 检查发送者权限
	if req.SenderType == "user" && session.UserID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权发送消息",
		})
		return
	}

	if req.SenderType == "counselor" && session.CounselorID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权发送消息",
		})
		return
	}

	// 设置默认内容类型
	if req.ContentType == "" {
		req.ContentType = "text"
	}

	// 创建消息
	message := models.ChatMessage{
		SessionID:   session.ID,
		SenderID:    userID.(uint),
		SenderType:  req.SenderType,
		ContentType: req.ContentType,
		Content:     req.Content,
		FileURL:     req.FileURL,
		IsRead:      false,
	}

	if err := database.DB.Create(&message).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "发送失败: " + err.Error(),
		})
		return
	}

	// 发送通知给接收者
	var receiverID uint
	if req.SenderType == "user" {
		receiverID = session.CounselorID
	} else {
		receiverID = session.UserID
	}

	go CreateNotification(receiverID, models.NotificationTypeChat, models.NotificationLevelInfo, "新消息", "您收到一条新消息", "")

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "发送成功",
		"data": message,
	})
}

// GetMessages godoc
// @Summary 获取消息列表
// @Description 获取会话的消息列表
// @Tags 聊天
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param session_id path int true "会话ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{messages,total}"
// @Router /api/chat/session/{session_id}/messages [get]
func GetMessages(c *gin.Context) {
	userID, _ := c.Get("user_id")
	sessionID := c.Param("session_id")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	// 查询会话
	var session models.ChatSession
	if err := database.DB.First(&session, sessionID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "会话不存在",
		})
		return
	}

	// 检查权限
	if session.UserID != userID.(uint) && session.CounselorID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权访问此会话",
		})
		return
	}

	// 查询消息
	query := database.DB.Model(&models.ChatMessage{}).Where("session_id = ?", sessionID)

	var total int64
	query.Count(&total)

	var messages []models.ChatMessage
	offset := 0
	if page != "1" {
		p, _ := strconv.Atoi(page)
		ps, _ := strconv.Atoi(pageSize)
		offset = (p - 1) * ps
	}

	ps, _ := strconv.Atoi(pageSize)
	if err := query.Preload("Sender").Offset(offset).Limit(ps).Order("created_at ASC").Find(&messages).Error; err != nil {
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

// EndChatSession godoc
// @Summary 结束聊天会话
// @Description 结束聊天会话
// @Tags 聊天
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param session_id path int true "会话ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:结束成功"
// @Router /api/chat/session/{session_id}/end [post]
func EndChatSession(c *gin.Context) {
	userID, _ := c.Get("user_id")
	sessionID := c.Param("session_id")

	// 查询会话
	var session models.ChatSession
	if err := database.DB.First(&session, sessionID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "会话不存在",
		})
		return
	}

	// 检查权限（只有咨询师可以结束会话）
	if session.CounselorID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "只有咨询师可以结束会话",
		})
		return
	}

	// 检查会话状态
	if session.Status != 1 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "会话已结束",
		})
		return
	}

	// 结束会话
	now := time.Now()
	duration := int(now.Sub(*session.StartTime).Seconds())

	if err := database.DB.Model(&session).Updates(map[string]interface{}{
		"status":   2,
		"end_time": &now,
		"duration": duration,
	}).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "结束失败: " + err.Error(),
		})
		return
	}

	// 发送通知给用户
	go CreateNotification(session.UserID, models.NotificationTypeChat, models.NotificationLevelSuccess, "咨询已结束", "您的咨询会话已结束", "")

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "结束成功",
		"data": gin.H{
			"duration": duration,
		},
	})
}

// GetChatSessions godoc
// @Summary 获取聊天会话列表
// @Description 获取当前用户的聊天会话列表
// @Tags 聊天
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{sessions,total}"
// @Router /api/chat/sessions [get]
func GetChatSessions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	query := database.DB.Model(&models.ChatSession{}).Where("user_id = ? OR counselor_id = ?", userID, userID)

	var total int64
	query.Count(&total)

	var sessions []models.ChatSession
	offset := 0
	if page != "1" {
		p, _ := strconv.Atoi(page)
		ps, _ := strconv.Atoi(pageSize)
		offset = (p - 1) * ps
	}

	ps, _ := strconv.Atoi(pageSize)
	if err := query.Preload("Order").Preload("User").Preload("Counselor").Offset(offset).Limit(ps).Order("created_at DESC").Find(&sessions).Error; err != nil {
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

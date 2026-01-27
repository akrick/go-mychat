package handlers

import (
	"akrick.com/mychat/database"
	"akrick.com/mychat/models"
	"akrick.com/mychat/websocket"
	"time"

	"github.com/gin-gonic/gin"
)

// WSChatHandler WebSocket聊天处理器
func WSChatHandler(c *gin.Context) {
	websocket.HandleWebSocket(c)
}

// GetSessionInfo godoc
// @Summary 获取会话信息
// @Description 获取聊天会话详细信息
// @Tags WebSocket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param session_id path int true "会话ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{session}"
// @Router /api/ws/session/{session_id} [get]
func GetSessionInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")
	sessionID := c.Param("session_id")

	var session models.ChatSession
	if err := database.DB.Preload("Order").Preload("User").Preload("Counselor").First(&session, sessionID).Error; err != nil {
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

	// 获取在线状态
	onlineUsers := websocket.GetSessionParticipants(session.ID)
	online := false
	for _, uid := range onlineUsers {
		if uid == userID.(uint) {
			online = true
			break
		}
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"session": session,
			"online":  online,
		},
	})
}

// GetBillingInfo godoc
// @Summary 获取计费信息
// @Description 获取聊天会话的计费信息
// @Tags WebSocket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param session_id path int true "会话ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{billing}"
// @Router /api/ws/session/{session_id}/billing [get]
func GetBillingInfo(c *gin.Context) {
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

	// 检查权限
	if session.UserID != userID.(uint) && session.CounselorID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权访问此会话",
		})
		return
	}

	// 查询计费记录
	var billing models.ChatBilling
	if err := database.DB.Where("session_id = ?", sessionID).First(&billing).Error; err != nil {
		// 如果没有计费记录，说明会话还未结束
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "会话进行中",
			"data": gin.H{
				"billing": nil,
				"session": session,
			},
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"billing": billing,
		},
	})
}

// GetBillingList godoc
// @Summary 获取计费列表
// @Description 获取当前用户的计费记录列表
// @Tags WebSocket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{billings,total}"
// @Router /api/ws/billings [get]
func GetBillingList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	query := database.DB.Model(&models.ChatBilling{}).Where("user_id = ?", userID)

	var total int64
	query.Count(&total)

	var billings []models.ChatBilling
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (parseInt(page) - 1) * parseInt(pageSize)
	}

	if err := query.Preload("Session").Preload("Order").Preload("User").Preload("Counselor").
		Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&billings).Error; err != nil {
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
			"billings": billings,
			"total":    total,
		},
	})
}

// GetCounselorBillings godoc
// @Summary 获取咨询师计费列表
// @Description 获取咨询师的所有计费记录
// @Tags WebSocket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{billings,total}"
// @Router /api/ws/counselor/billings [get]
func GetCounselorBillings(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	// 检查是否为咨询师
	var counselor models.Counselor
	if err := database.DB.Where("user_id = ?", userID).First(&counselor).Error; err != nil {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "不是咨询师",
		})
		return
	}

	query := database.DB.Model(&models.ChatBilling{}).Where("counselor_id = ?", counselor.ID)

	var total int64
	query.Count(&total)

	var billings []models.ChatBilling
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (parseInt(page) - 1) * parseInt(pageSize)
	}

	if err := query.Preload("Session").Preload("Order").Preload("User").Preload("Counselor").
		Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&billings).Error; err != nil {
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
			"billings": billings,
			"total":    total,
		},
	})
}

// GetCounselorAccount godoc
// @Summary 获取咨询师账户信息
// @Description 获取当前咨询师的账户信息
// @Tags WebSocket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{account}"
// @Router /api/ws/counselor/account [get]
func GetCounselorAccount(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// 检查是否为咨询师
	var counselor models.Counselor
	if err := database.DB.Where("user_id = ?", userID).First(&counselor).Error; err != nil {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "不是咨询师",
		})
		return
	}

	// 查询账户
	var account models.CounselorAccount
	if err := database.DB.Where("counselor_id = ?", counselor.ID).First(&account).Error; err != nil {
		// 如果没有账户，创建一个
		account = models.CounselorAccount{
			CounselorID:  counselor.ID,
			TotalIncome:  0,
			Withdrawn:    0,
			Balance:      0,
			FrozenAmount: 0,
		}
		database.DB.Create(&account)
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"account": account,
		},
	})
}

// CreateWithdraw godoc
// @Summary 创建提现申请
// @Description 咨询师创建提现申请
// @Tags WebSocket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]interface{} true "提现信息:amount,bank_name,bank_account,account_name"
// @Success 200 {object} map[string]interface{} "code:200,msg:申请成功,data:{withdraw}"
// @Router /api/ws/counselor/withdraw [post]
func CreateWithdraw(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		BankName    string  `json:"bank_name" binding:"required"`
		BankAccount string  `json:"bank_account" binding:"required"`
		AccountName string  `json:"account_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 检查是否为咨询师
	var counselor models.Counselor
	if err := database.DB.Where("user_id = ?", userID).First(&counselor).Error; err != nil {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "不是咨询师",
		})
		return
	}

	// 查询账户
	var account models.CounselorAccount
	if err := database.DB.Where("counselor_id = ?", counselor.ID).First(&account).Error; err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "账户不存在",
		})
		return
	}

	// 检查余额是否足够
	if account.Balance < req.Amount {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "余额不足",
		})
		return
	}

	// 创建提现记录
	withdraw := models.WithdrawRecord{
		CounselorID: counselor.ID,
		Amount:      req.Amount,
		Status:      0, // 待审核
		BankName:    req.BankName,
		BankAccount: req.BankAccount,
		AccountName: req.AccountName,
	}

	if err := database.DB.Create(&withdraw).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "创建失败: " + err.Error(),
		})
		return
	}

	// 冻结金额
	database.DB.Model(&account).Update("frozen_amount", account.FrozenAmount+req.Amount)
	database.DB.Model(&account).Update("balance", account.Balance-req.Amount)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "申请成功",
		"data": gin.H{
			"withdraw": withdraw,
		},
	})
}

// GetWithdrawList godoc
// @Summary 获取提现记录列表
// @Description 获取当前咨询师的提现记录列表
// @Tags WebSocket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{withdraws,total}"
// @Router /api/ws/counselor/withdraws [get]
func GetWithdrawList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	// 检查是否为咨询师
	var counselor models.Counselor
	if err := database.DB.Where("user_id = ?", userID).First(&counselor).Error; err != nil {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "不是咨询师",
		})
		return
	}

	query := database.DB.Model(&models.WithdrawRecord{}).Where("counselor_id = ?", counselor.ID)

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
			"total":    total,
		},
	})
}

// GetMessageHistory godoc
// @Summary 获取消息历史
// @Description 分页获取聊天消息历史
// @Tags WebSocket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param session_id path int true "会话ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(50)
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{messages,total}"
// @Router /api/ws/session/:session_id/messages [get]
func GetMessageHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	sessionID := c.Param("session_id")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "50")

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
	if page == "1" {
		offset = 0
	} else {
		offset = (parseInt(page) - 1) * parseInt(pageSize)
	}

	if err := query.Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&messages).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败: " + err.Error(),
		})
		return
	}

	// 反转消息顺序，使最新的在最后
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
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

// GetUnreadCount godoc
// @Summary 获取未读消息数
// @Description 获取用户在各会话中的未读消息数
// @Tags WebSocket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{unread_count,sessions}"
// @Router /api/ws/unread/count [get]
func GetUnreadCount(c *gin.Context) {
	userID, _ := c.Get("user_id")

	type UnreadSession struct {
		SessionID   uint `json:"session_id"`
		UnreadCount int  `json:"unread_count"`
	}

	// 查询用户参与的会话
	var sessions []models.ChatSession
	database.DB.Where("user_id = ? OR counselor_id = ?", userID, userID).Find(&sessions)

	sessionIDs := make([]uint, len(sessions))
	for i, s := range sessions {
		sessionIDs[i] = s.ID
	}

	// 统计未读消息
	var unreadSessions []UnreadSession
	if len(sessionIDs) > 0 {
		rows, _ := database.DB.Model(&models.ChatMessage{}).
			Select("session_id, count(*) as unread_count").
			Where("session_id IN ? AND is_read = false AND sender_id != ?", sessionIDs, userID).
			Group("session_id").
			Rows()

		for rows.Next() {
			var us UnreadSession
			database.DB.ScanRows(rows, &us)
			unreadSessions = append(unreadSessions, us)
		}
		rows.Close()
	}

	totalUnread := 0
	for _, us := range unreadSessions {
		totalUnread += us.UnreadCount
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"unread_count": totalUnread,
			"sessions":     unreadSessions,
		},
	})
}

// MarkMessagesRead godoc
// @Summary 标记会话消息为已读
// @Description 标记指定会话中的所有消息为已读
// @Tags WebSocket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param session_id path int true "会话ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:标记成功"
// @Router /api/ws/session/:session_id/read [post]
func MarkMessagesRead(c *gin.Context) {
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

	// 检查权限
	if session.UserID != userID.(uint) && session.CounselorID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权访问此会话",
		})
		return
	}

	// 更新未读消息为已读
	nowTime := time.Now()
	database.DB.Model(&models.ChatMessage{}).
		Where("session_id = ? AND sender_id != ? AND is_read = false", sessionID, userID).
		Updates(map[string]interface{}{
			"is_read":   true,
			"read_time": &nowTime,
		})

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "标记成功",
	})
}

// SearchMessages godoc
// @Summary 搜索消息
// @Description 在会话中搜索包含关键词的消息
// @Tags WebSocket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param session_id path int true "会话ID"
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{messages,total}"
// @Router /api/ws/session/:session_id/search [get]
func SearchMessages(c *gin.Context) {
	userID, _ := c.Get("user_id")
	sessionID := c.Param("session_id")
	keyword := c.Query("keyword")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	if keyword == "" {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "请输入搜索关键词",
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

	// 检查权限
	if session.UserID != userID.(uint) && session.CounselorID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权访问此会话",
		})
		return
	}

	// 搜索消息
	query := database.DB.Model(&models.ChatMessage{}).
		Where("session_id = ? AND content LIKE ?", sessionID, "%"+keyword+"%")

	var total int64
	query.Count(&total)

	var messages []models.ChatMessage
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (parseInt(page) - 1) * parseInt(pageSize)
	}

	if err := query.Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&messages).Error; err != nil {
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

// RevokeMessage godoc
// @Summary 撤回消息
// @Description 撤回已发送的消息（发送后2分钟内）
// @Tags WebSocket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message_id path int true "消息ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:撤回成功"
// @Router /api/ws/message/:message_id/revoke [post]
func RevokeMessage(c *gin.Context) {
	userID, _ := c.Get("user_id")
	messageID := c.Param("message_id")

	// 查询消息
	var message models.ChatMessage
	if err := database.DB.First(&message, messageID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "消息不存在",
		})
		return
	}

	// 检查是否为发送者
	if message.SenderID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "只能撤回自己的消息",
		})
		return
	}

	// 检查时间（2分钟内可撤回）
	if time.Since(message.CreatedAt) > 2*time.Minute {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "超过2分钟，无法撤回",
		})
		return
	}

	// 标记为已撤回
	message.Content = "[消息已撤回]"
	message.ContentType = "system"
	database.DB.Save(&message)

	// 通知会话内所有客户端
	websocket.BroadcastToSession(message.SessionID, websocket.BuildRevokeMessage(messageID))

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "撤回成功",
	})
}

// GetCounselorOnlineStatus godoc
// @Summary 获取咨询师在线状态
// @Description 获取指定咨询师的在线状态
// @Tags WebSocket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param counselor_id path int true "咨询师ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{online,last_active}"
// @Router /api/ws/counselor/:counselor_id/status [get]
func GetCounselorOnlineStatus(c *gin.Context) {
	counselorID := c.Param("counselor_id")

	// 查询咨询师
	var counselor models.Counselor
	if err := database.DB.First(&counselor, counselorID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "咨询师不存在",
		})
		return
	}

	// 检查是否在线（咨询师ID也作为用户ID使用）
	online := websocket.IsUserOnline(counselor.ID)

	// 获取最后活跃时间
	var lastActive interface{}
	if online {
		lastActive = time.Now()
	} else {
		lastActive = nil
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"online":      online,
			"last_active": lastActive,
			"counselor_id": counselor.ID,
		},
	})
}

// GetOnlineCounselors godoc
// @Summary 获取在线咨询师列表
// @Description 获取当前所有在线的咨询师
// @Tags WebSocket
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{counselors}"
// @Router /api/ws/counselors/online [get]
func GetOnlineCounselors(c *gin.Context) {
	// 获取所有在线用户
	onlineUserIDs := websocket.GetOnlineUsers()

	// 查询在线咨询师
	var counselors []models.Counselor
	if len(onlineUserIDs) > 0 {
		database.DB.Where("user_id IN ? AND status = 1", onlineUserIDs).
			Preload("User").
			Find(&counselors)
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"counselors": counselors,
			"count":     len(counselors),
		},
	})
}

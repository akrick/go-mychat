package handlers

import (
	"encoding/json"
	"akrick.com/mychat/database"
	"akrick.com/mychat/models"
	"akrick.com/mychat/websocket"
	"akrick.com/mychat/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// GetSessionStats godoc
// @Summary 获取会话统计信息
// @Description 获取WebSocket会话统计信息
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{stats}"
// @Router /api/admin/session/stats [get]
func GetSessionStats(c *gin.Context) {
	stats := GetSessionStatsInternal()

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": stats,
	})
}

// GetSessionStatsInternal 获取会话统计内部方法
func GetSessionStatsInternal() map[string]interface{} {
	return websocket.GetSessionStats()
}

// GetOnlineUsers godoc
// @Summary 获取在线用户列表
// @Description 获取当前在线的所有用户
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{user_ids}"
// @Router /api/admin/online/users [get]
func GetOnlineUsers(c *gin.Context) {
	userIDs := websocket.GetOnlineUsers()

	// 查询用户详细信息
	var users []models.User
	if len(userIDs) > 0 {
		database.DB.Where("id IN ?", userIDs).Find(&users)
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"user_ids": userIDs,
			"users":    users,
			"count":    len(userIDs),
		},
	})
}

// BroadcastSystemMessage godoc
// @Summary 广播系统消息
// @Description 向所有在线用户广播系统消息
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]interface{} true "消息内容"
// @Success 200 {object} map[string]interface{} "code:200,msg:广播成功"
// @Router /api/admin/broadcast [post]
func BroadcastSystemMessage(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	wsBroadcastSystemMessage(req.Content)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "广播成功",
	})
}

// wsBroadcastSystemMessage WebSocket广播系统消息
func wsBroadcastSystemMessage(content string) {
	msg, _ := json.Marshal(map[string]interface{}{
		"type": "system_message",
		"data": gin.H{
			"content":    content,
			"created_at": time.Now(),
		},
	})

	websocket.BroadcastToAll(msg)
}

// ApproveWithdraw godoc
// @Summary 审核提现申请
// @Description 管理员审核提现申请
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "提现记录ID"
// @Param request body map[string]interface{} true "审核结果:approved(true/false),rejected_reason"
// @Success 200 {object} map[string]interface{} "code:200,msg:审核成功"
// @Router /api/admin/withdraw/:id/approve [post]
func ApproveWithdraw(c *gin.Context) {
	var req struct {
		Approved       bool   `json:"approved" binding:"required"`
		RejectedReason string `json:"rejected_reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

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
	if withdraw.Status != 0 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "该申请已处理",
		})
		return
	}

	now := time.Now()

	if req.Approved {
		// 通过审核
		database.DB.Model(&withdraw).Updates(map[string]interface{}{
			"status":       3,
			"audited_at":   &now,
			"transferred_at": &now,
		})

		// 解冻并扣除余额
		var account models.CounselorAccount
		database.DB.Where("counselor_id = ?", withdraw.CounselorID).First(&account)
		database.DB.Model(&account).Updates(map[string]interface{}{
			"frozen_amount": account.FrozenAmount - withdraw.Amount,
			"withdrawn":     account.Withdrawn + withdraw.Amount,
		})
	} else {
		// 拒绝审核
		database.DB.Model(&withdraw).Updates(map[string]interface{}{
			"status":          2,
			"audited_at":      &now,
			"rejected_reason": req.RejectedReason,
		})

		// 解冻并退回余额
		var account models.CounselorAccount
		database.DB.Where("counselor_id = ?", withdraw.CounselorID).First(&account)
		database.DB.Model(&account).Updates(map[string]interface{}{
			"frozen_amount": account.FrozenAmount - withdraw.Amount,
			"balance":       account.Balance + withdraw.Amount,
		})
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "审核成功",
		"data": withdraw,
	})
}

// GetPendingWithdraws godoc
// @Summary 获取待审核提现列表
// @Description 获取所有待审核的提现申请
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/withdraws/pending [get]
func GetPendingWithdraws(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	query := database.DB.Model(&models.WithdrawRecord{}).Where("status = ?", 0)

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

// GetAdminUserInfo godoc
// @Summary 获取管理员信息
// @Description 获取当前登录的管理员信息
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/user/info [get]
func GetAdminUserInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": user,
	})
}

// GetAdminPermissions godoc
// @Summary 获取管理员权限列表
// @Description 获取当前管理员的权限列表
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/user/permissions [get]
func GetAdminPermissions(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// 获取用户角色
	var user models.User
	database.DB.First(&user, userID)

	// 获取角色权限
	var permissions []models.Permission
	database.DB.Raw(`
		SELECT p.* FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN roles r ON rp.role_id = r.id
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = ? AND p.status = 1 AND r.status = 1
	`, userID).Scan(&permissions)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": permissions,
	})
}

// GetAdminStatistics godoc
// @Summary 获取管理员统计数据
// @Description 获取管理后台的统计数据
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/statistics [get]
func GetAdminStatistics(c *gin.Context) {
	var userCount int64
	var onlineUserCount int64
	var counselorCount int64
	var orderCount int64
	var todayOrderCount int64
	var totalAmount float64
	var todayAmount float64
	var sessionCount int64
	var activeSessionCount int64

	database.DB.Model(&models.User{}).Count(&userCount)
	database.DB.Model(&models.User{}).Where("status = ?", 1).Count(&onlineUserCount)
	database.DB.Model(&models.Counselor{}).Where("status = ?", 1).Count(&counselorCount)
	database.DB.Model(&models.Order{}).Count(&orderCount)
	
	today := time.Now().Format("2006-01-02")
	database.DB.Model(&models.Order{}).Where("DATE(created_at) = ?", today).Count(&todayOrderCount)
	
	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderStatusPaid).Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount)
	database.DB.Model(&models.Order{}).Where("status = ? AND DATE(created_at) = ?", models.OrderStatusPaid, today).Select("COALESCE(SUM(amount), 0)").Scan(&todayAmount)
	
	database.DB.Model(&models.ChatSession{}).Count(&sessionCount)
	database.DB.Model(&models.ChatSession{}).Where("status = ?", 1).Count(&activeSessionCount)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"user_count":           userCount,
			"online_user_count":    onlineUserCount,
			"counselor_count":      counselorCount,
			"order_count":          orderCount,
			"today_order_count":    todayOrderCount,
			"total_amount":         totalAmount,
			"today_amount":         todayAmount,
			"session_count":        sessionCount,
			"active_session_count": activeSessionCount,
		},
	})
}

// AdminLogin godoc
// @Summary 管理员登录
// @Description 管理员登录接口
// @Tags 管理员
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "登录信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:登录成功,data:{token,user}"
// @Router /api/admin/login [post]
func AdminLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 检查是否为管理员
	if !user.IsAdmin {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "非管理员账号",
		})
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 检查状态
	if user.Status != 1 {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "账号已被禁用",
		})
		return
	}

	// 生成Token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "生成Token失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"avatar":   user.Avatar,
				"is_admin": user.IsAdmin,
			},
		},
	})
}

// AdminLogout godoc
// @Summary 管理员退出登录
// @Description 管理员退出登录接口
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:退出成功"
// @Router /api/admin/logout [post]
func AdminLogout(c *gin.Context) {
	// 可以在这里添加将token加入黑名单的逻辑
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "退出成功",
	})
}

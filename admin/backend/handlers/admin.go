package handlers

import (
	"encoding/json"
	"fmt"
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"
	"akrick.com/mychat/admin/backend/websocket"
	"akrick.com/mychat/admin/backend/utils"
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
		database.DB.Where("id IN (?)", userIDs).Find(&users)
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

// KickOutUser godoc
// @Summary 强制下线用户
// @Description 强制指定用户下线
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:下线成功"
// @Router /api/admin/online/users/:id/kick [post]
func KickOutUser(c *gin.Context) {
	userID := c.Param("id")

	// 将字符串转换为uint
	var uid uint
	fmt.Sscanf(userID, "%d", &uid)

	// 检查用户是否在线
	if !websocket.IsUserOnline(uid) {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "用户不在线",
		})
		return
	}

	// 发送强制下线消息给用户
	msg, _ := json.Marshal(map[string]interface{}{
		"type": "force_logout",
		"data": gin.H{
			"reason":     "管理员强制下线",
			"created_at": time.Now(),
		},
	})

	websocket.SendToUser(uid, msg)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "用户已强制下线",
	})
}

// SendToUser godoc
// @Summary 发送消息给指定用户
// @Description 向指定在线用户发送消息
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param request body map[string]interface{} true "消息内容"
// @Success 200 {object} map[string]interface{} "code:200,msg:发送成功"
// @Router /api/admin/online/users/:id/message [post]
func SendToUser(c *gin.Context) {
	userID := c.Param("id")

	// 将字符串转换为uint
	var uid uint
	fmt.Sscanf(userID, "%d", &uid)

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

	// 检查用户是否在线
	if !websocket.IsUserOnline(uid) {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "用户不在线",
		})
		return
	}

	// 发送消息给指定用户
	msg, _ := json.Marshal(map[string]interface{}{
		"type": "admin_message",
		"data": gin.H{
			"content":    req.Content,
			"created_at": time.Now(),
		},
	})

	websocket.SendToUser(uid, msg)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "消息发送成功",
	})
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
		"data": gin.H{
			"permissions": permissions,
			"roles":       []string{"admin"},
		},
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

	// 趋势数据
	var yesterdayOrderCount int64
	var yesterdayAmount float64

	// 订单状态统计
	var pendingOrders int64
	var paidOrders int64
	var completedOrders int64
	var cancelledOrders int64

	database.DB.Model(&models.User{}).Count(&userCount)
	database.DB.Model(&models.User{}).Where("status = ?", 1).Count(&onlineUserCount)
	database.DB.Model(&models.Counselor{}).Where("status = ?", 1).Count(&counselorCount)
	database.DB.Model(&models.Order{}).Count(&orderCount)

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	database.DB.Model(&models.Order{}).Where("DATE(created_at) = ?", today).Count(&todayOrderCount)
	database.DB.Model(&models.Order{}).Where("DATE(created_at) = ?", yesterday).Count(&yesterdayOrderCount)

	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderStatusPaid).Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount)
	database.DB.Model(&models.Order{}).Where("status = ? AND DATE(created_at) = ?", models.OrderStatusPaid, today).Select("COALESCE(SUM(amount), 0)").Scan(&todayAmount)
	database.DB.Model(&models.Order{}).Where("status = ? AND DATE(created_at) = ?", models.OrderStatusPaid, yesterday).Select("COALESCE(SUM(amount), 0)").Scan(&yesterdayAmount)

	// 订单状态统计
	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderStatusPending).Count(&pendingOrders)
	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderStatusPaid).Count(&paidOrders)
	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderStatusCompleted).Count(&completedOrders)
	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderStatusCancelled).Count(&cancelledOrders)

	database.DB.Model(&models.ChatSession{}).Count(&sessionCount)
	database.DB.Model(&models.ChatSession{}).Where("status = ?", 1).Count(&activeSessionCount)

	// 计算趋势百分比
	var userTrend float64
	var counselorTrend float64
	var orderTrend float64
	var revenueTrend float64

	if yesterdayOrderCount > 0 {
		orderTrend = float64(todayOrderCount-yesterdayOrderCount) / float64(yesterdayOrderCount) * 100
	}

	if yesterdayAmount > 0 {
		revenueTrend = (todayAmount - yesterdayAmount) / yesterdayAmount * 100
	}

	// 用户和咨询师趋势使用月度对比
	var lastMonthUserCount int64
	var thisMonthUserCount int64
	lastMonth := time.Now().AddDate(0, -1, 0).Format("2006-01")
	thisMonth := time.Now().Format("2006-01")
	database.DB.Model(&models.User{}).Where("DATE_FORMAT(created_at, '%Y-%m') = ?", lastMonth).Count(&lastMonthUserCount)
	database.DB.Model(&models.User{}).Where("DATE_FORMAT(created_at, '%Y-%m') = ?", thisMonth).Count(&thisMonthUserCount)

	if lastMonthUserCount > 0 {
		userTrend = float64(thisMonthUserCount-lastMonthUserCount) / float64(lastMonthUserCount) * 100
	}

	var lastMonthCounselorCount int64
	var thisMonthCounselorCount int64
	database.DB.Model(&models.Counselor{}).Where("DATE_FORMAT(created_at, '%Y-%m') = ?", lastMonth).Count(&lastMonthCounselorCount)
	database.DB.Model(&models.Counselor{}).Where("DATE_FORMAT(created_at, '%Y-%m') = ?", thisMonth).Count(&thisMonthCounselorCount)

	if lastMonthCounselorCount > 0 {
		counselorTrend = float64(thisMonthCounselorCount-lastMonthCounselorCount) / float64(lastMonthCounselorCount) * 100
	}

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
			"user_trend":           userTrend,
			"counselor_trend":      counselorTrend,
			"order_trend":          orderTrend,
			"revenue_trend":        revenueTrend,
			"pending_orders":       pendingOrders,
			"paid_orders":          paidOrders,
			"completed_orders":     completedOrders,
			"cancelled_orders":      cancelledOrders,
		},
	})
}

// AdminLogin godoc
// @Summary 管理员登录（兼容旧接口，使用Administrator表）
// @Description 管理员登录接口（公开路由，无需鉴权）
// @Tags 管理员
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "登录信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:登录成功,data:{token,user}"
// @Router /api/admin/login [post]
func AdminLogin(c *gin.Context) {
	fmt.Println("\n========== 收到登录请求 ==========")
	fmt.Println("请求方法:", c.Request.Method)
	fmt.Println("请求路径:", c.Request.URL.Path)
	fmt.Println("请求头:", c.Request.Header)
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	fmt.Println("=== 登录请求开始 ===")
	fmt.Printf("接收到的请求体: %+v\n", req)

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("参数绑定失败:", err)
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	fmt.Printf("用户名: %s, 密码: %s\n", req.Username, req.Password)

	var admin models.Administrator
	if err := database.DB.Where("username = ?", req.Username).First(&admin).Error; err != nil {
		fmt.Println("管理员查询失败:", err)
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}

	fmt.Printf("找到管理员: ID=%d, Username=%s, Role=%s, Status=%d\n", admin.ID, admin.Username, admin.Role, admin.Status)

	// 验证密码
	passwordValid := utils.CheckPassword(req.Password, admin.Password)
	fmt.Printf("密码验证结果: %v\n", passwordValid)

	if !passwordValid {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 检查状态
	if admin.Status != 1 {
		fmt.Println("管理员状态异常:", admin.Status)
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "账号已被禁用",
		})
		return
	}

	// 生成Token
	token, err := utils.GenerateToken(admin.ID, admin.Username)
	if err != nil {
		fmt.Println("生成Token失败:", err)
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "生成Token失败",
		})
		return
	}

	// 更新最后登录时间
	now := time.Now()
	database.DB.Model(&admin).Update("last_login", now)

	fmt.Println("登录成功，生成的Token:", token)
	fmt.Println("=== 登录请求结束 ===")

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":        admin.ID,
				"username":  admin.Username,
				"real_name": admin.RealName,
				"email":     admin.Email,
				"avatar":    admin.Avatar,
				"role":     admin.Role,
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

// MuteUser godoc
// @Summary 禁言/解禁用户
// @Description 禁言或解禁指定咨询师（预留功能，当前数据库未支持）
// @Tags 在线用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]interface{} true "禁言信息:user_id,is_muted"
// @Success 200 {object} map[string]interface{} "code:200,msg:操作成功"
// @Router /api/admin/online/mute [post]
func MuteUser(c *gin.Context) {
	// 当前数据库模型没有is_muted字段，此接口暂时返回成功
	// 如需启用，请先在CounselorApplication表中添加is_muted字段
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "操作成功（功能待启用）",
	})
}

// GetOnlineStatistics godoc
// @Summary 获取在线统计
// @Description 获取在线用户统计数据
// @Tags 在线用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{statistics}"
// @Router /api/admin/online/statistics [get]
func GetOnlineStatistics(c *gin.Context) {
	// 获取在线用户
	onlineUserIDs := websocket.GetOnlineUsers()
	totalOnline := len(onlineUserIDs)

	// 统计咨询师在线数（通过CounselorApplication状态为1的用户）
	var onlineCounselorCount int64
	database.DB.Model(&models.CounselorApplication{}).
		Where("user_id IN (?) AND status = 1", onlineUserIDs).
		Count(&onlineCounselorCount)

	// 普通用户在线数
	regularUserCount := totalOnline - int(onlineCounselorCount)

	c.JSON(200, gin.H{
		"code": 200,
		"msg": "获取成功",
		"data": gin.H{
			"total_online":       totalOnline,
			"counselor_online":   int(onlineCounselorCount),
			"user_online":        regularUserCount,
			"muted_counselor":    0, // 数据库暂不支持
			"muted_online_count": 0,
		},
	})
}

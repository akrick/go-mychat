package handlers

import (
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"
	"akrick.com/mychat/admin/backend/websocket"

	"github.com/gin-gonic/gin"
)

// OnlineUserInfo 在线用户信息
type OnlineUserInfo struct {
	UserID     uint   `json:"user_id"`
	Username    string `json:"username"`
	Role       string `json:"role"`       // user, counselor, admin
	OnlineTime string `json:"online_time"` // 在线时长
	SessionID  *uint  `json:"session_id"`
	IsMuted    bool   `json:"is_muted"`   // 是否被禁言
}

// GetOnlineUsersDetailed godoc
// @Summary 获取在线用户详细信息
// @Description 获取所有在线用户的详细信息
// @Tags 在线用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{users,total}"
// @Router /api/admin/online/users/detailed [get]
func GetOnlineUsersDetailed(c *gin.Context) {
	// 获取在线用户ID列表和会话信息
	onlineUserIDs := websocket.GetOnlineUsers()
	userSessions := websocket.GetOnlineUsersWithSessions()

	if len(onlineUserIDs) == 0 {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "获取成功",
			"data": gin.H{
				"users": []OnlineUserInfo{},
				"total": 0,
			},
		})
		return
	}

	// 查询用户基本信息
	var users []models.User
	if err := database.DB.Where("id IN ?", onlineUserIDs).Find(&users).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg": "查询用户信息失败: " + err.Error(),
		})
		return
	}

	// 查询咨询师信息（通过ID关联User）
	var counselors []models.CounselorApplication
	if err := database.DB.Where("user_id IN ? AND status = 1", onlineUserIDs).Find(&counselors).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg": "查询咨询师信息失败: " + err.Error(),
		})
		return
	}

	// 构建咨询师映射
	counselorUserMap := make(map[uint]bool)
	for i := range counselors {
		counselorUserMap[counselors[i].UserID] = true
	}

	onlineUsers := make([]OnlineUserInfo, 0, len(users))
	for _, user := range users {
		role := "user"
		var sessionID *uint = nil

		// 判断角色
		if counselorUserMap[user.ID] {
			role = "counselor"
		}

		// 获取会话ID
		if sessionInfo, ok := userSessions[user.ID]; ok {
			if sessID, hasSession := sessionInfo.(map[string]interface{})["session_id"]; hasSession {
				if sessID != nil {
					if sid, ok := sessID.(uint); ok {
						sessionID = &sid
					}
				}
			}
		}

		onlineUsers = append(onlineUsers, OnlineUserInfo{
			UserID:     user.ID,
			Username:    user.Username,
			Role:       role,
			OnlineTime: "未知",
			SessionID:  sessionID,
			IsMuted:    false, // 数据库中没有is_muted字段，暂时设为false
		})
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

package handlers

import (
	"akrick.com/mychat/database"
	"akrick.com/mychat/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

// CreateNotification 创建通知（内部使用）
func CreateNotification(userID uint, notificationType, level, title, content, extraData string) error {
	notification := models.Notification{
		UserID:    userID,
		Type:      notificationType,
		Level:     level,
		Title:     title,
		Content:   content,
		ExtraData: extraData,
		IsRead:    false,
	}

	return database.DB.Create(&notification).Error
}

// GetNotifications godoc
// @Summary 获取通知列表
// @Description 获取当前用户的通知列表
// @Tags 通知
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param type query string false "通知类型"
// @Param is_read query bool false "是否已读"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{notifications,total,unread_count}"
// @Router /api/notification/list [get]
func GetNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	notificationType := c.Query("type")
	isRead := c.Query("is_read")

	query := database.DB.Model(&models.Notification{}).Where("user_id = ?", userID)

	// 类型筛选
	if notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}

	// 已读筛选
	if isRead != "" {
		if isRead == "true" {
			query = query.Where("is_read = ?", true)
		} else {
			query = query.Where("is_read = ?", false)
		}
	}

	var total int64
	var unreadCount int64

	query.Count(&total)
	database.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&unreadCount)

	var notifications []models.Notification
	offset := 0
	if page != "1" {
		p, _ := strconv.Atoi(page)
		ps, _ := strconv.Atoi(pageSize)
		offset = (p - 1) * ps
	}

	ps, _ := strconv.Atoi(pageSize)
	if err := query.Preload("User").Offset(offset).Limit(ps).Order("created_at DESC").Find(&notifications).Error; err != nil {
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
			"notifications": notifications,
			"total":        total,
			"unread_count": unreadCount,
		},
	})
}

// MarkNotificationRead godoc
// @Summary 标记通知为已读
// @Description 标记单个通知为已读
// @Tags 通知
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "通知ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:标记成功"
// @Router /api/notification/{id}/read [post]
func MarkNotificationRead(c *gin.Context) {
	userID, _ := c.Get("user_id")
	notificationID := c.Param("id")

	// 查询通知
	var notification models.Notification
	if err := database.DB.First(&notification, notificationID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "通知不存在",
		})
		return
	}

	// 检查权限
	if notification.UserID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权操作此通知",
		})
		return
	}

	// 标记为已读
	now := database.DB.NowFunc()
	if err := database.DB.Model(&notification).Updates(map[string]interface{}{
		"is_read":  true,
		"read_time": &now,
	}).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "标记失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "标记成功",
	})
}

// MarkAllNotificationsRead godoc
// @Summary 标记所有通知为已读
// @Description 标记当前用户的所有通知为已读
// @Tags 通知
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:标记成功"
// @Router /api/notification/read-all [post]
func MarkAllNotificationsRead(c *gin.Context) {
	userID, _ := c.Get("user_id")

	now := database.DB.NowFunc()
	if err := database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read":   true,
			"read_time": &now,
		}).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "标记失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "标记成功",
	})
}

// DeleteNotification godoc
// @Summary 删除通知
// @Description 删除通知
// @Tags 通知
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "通知ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:删除成功"
// @Router /api/notification/{id} [delete]
func DeleteNotification(c *gin.Context) {
	userID, _ := c.Get("user_id")
	notificationID := c.Param("id")

	// 查询通知
	var notification models.Notification
	if err := database.DB.First(&notification, notificationID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "通知不存在",
		})
		return
	}

	// 检查权限
	if notification.UserID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权操作此通知",
		})
		return
	}

	// 删除通知
	if err := database.DB.Delete(&notification).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "删除失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

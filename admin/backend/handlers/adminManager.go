package handlers

import (
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"
	"akrick.com/mychat/admin/backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

// GetAdminList 获取管理员列表
// @Summary 获取管理员列表
// @Description 获取所有管理员用户列表
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Param status query int false "状态:0-禁用,1-正常"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/managers [get]
func GetAdminList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	keyword := c.Query("keyword")
	status := c.Query("status")

	query := database.DB.Model(&models.User{}).Where("is_admin = ?", true)

	// 搜索
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var admins []models.User
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (utils.ParseInt(page) - 1) * utils.ParseInt(pageSize)
	}

	if err := query.Offset(offset).Limit(utils.ParseInt(pageSize)).Order("created_at DESC").Find(&admins).Error; err != nil {
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
			"admins": admins,
			"total":  total,
		},
	})
}

// CreateAdmin 创建管理员
// @Summary 创建管理员
// @Description 创建新的管理员账号
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateAdminRequest true "管理员信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:创建成功"
// @Router /api/admin/managers [post]
func CreateAdmin(c *gin.Context) {
	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 检查用户名是否存在
	var count int64
	database.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "用户名已存在",
		})
		return
	}

	// 检查邮箱是否存在
	if req.Email != "" {
		database.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
		if count > 0 {
			c.JSON(400, gin.H{
				"code": 400,
				"msg":  "邮箱已存在",
			})
			return
		}
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "密码加密失败",
		})
		return
	}

	admin := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Status:   req.Status,
		IsAdmin:  true, // 强制设为管理员
	}

	if err := database.DB.Create(&admin).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "创建失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": admin,
	})
}

// UpdateAdmin 更新管理员
// @Summary 更新管理员
// @Description 更新管理员信息
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "管理员ID"
// @Param request body UpdateAdminRequest true "更新信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:更新成功"
// @Router /api/admin/managers/{id} [put]
func UpdateAdmin(c *gin.Context) {
	adminID := c.Param("id")

	var req UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	var admin models.User
	if err := database.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "管理员不存在",
		})
		return
	}

	// 检查是否为管理员
	if !admin.IsAdmin {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "该用户不是管理员",
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := database.DB.Model(&admin).Updates(updates).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "更新失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// DeleteAdmin 删除管理员
// @Summary 删除管理员
// @Description 删除管理员账号
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "管理员ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:删除成功"
// @Router /api/admin/managers/{id} [delete]
func DeleteAdmin(c *gin.Context) {
	adminID := c.Param("id")

	// 获取当前登录用户ID
	currentUserID, _ := c.Get("user_id")

	// 不允许删除自己
	if currentUserID != nil && adminID == fmt.Sprintf("%v", currentUserID) {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "不能删除自己",
		})
		return
	}

	var admin models.User
	if err := database.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "管理员不存在",
		})
		return
	}

	// 检查是否为管理员
	if !admin.IsAdmin {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "该用户不是管理员",
		})
		return
	}

	if err := database.DB.Delete(&admin).Error; err != nil {
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

// ResetAdminPassword 重置管理员密码
// @Summary 重置管理员密码
// @Description 重置管理员密码
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "管理员ID"
// @Param request body ResetAdminPasswordRequest true "密码信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:重置成功"
// @Router /api/admin/managers/{id}/password [post]
func ResetAdminPassword(c *gin.Context) {
	adminID := c.Param("id")

	var req ResetAdminPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	var admin models.User
	if err := database.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "管理员不存在",
		})
		return
	}

	// 检查是否为管理员
	if !admin.IsAdmin {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "该用户不是管理员",
		})
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "密码加密失败",
		})
		return
	}

	if err := database.DB.Model(&admin).Update("password", hashedPassword).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "重置失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "重置成功",
	})
}

// ToggleAdminStatus 切换管理员状态
// @Summary 切换管理员状态
// @Description 禁用或启用管理员账号
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "管理员ID"
// @Param request body ToggleAdminStatusRequest true "状态信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:操作成功"
// @Router /api/admin/managers/{id}/status [put]
func ToggleAdminStatus(c *gin.Context) {
	adminID := c.Param("id")

	// 获取当前登录用户ID
	currentUserID, _ := c.Get("user_id")

	// 不允许禁用自己
	if currentUserID != nil && adminID == fmt.Sprintf("%v", currentUserID) {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "不能禁用自己",
		})
		return
	}

	var req ToggleAdminStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	var admin models.User
	if err := database.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "管理员不存在",
		})
		return
	}

	// 检查是否为管理员
	if !admin.IsAdmin {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "该用户不是管理员",
		})
		return
	}

	if err := database.DB.Model(&admin).Update("status", req.Status).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "操作失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "操作成功",
	})
}

// Request structs
type CreateAdminRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
	Avatar   string `json:"avatar"`
	Status   int    `json:"status" binding:"omitempty,oneof=0 1"`
}

type UpdateAdminRequest struct {
	Email  string `json:"email" binding:"omitempty,email"`
	Phone  string `json:"phone" binding:"omitempty,len=11"`
	Avatar string `json:"avatar"`
	Status *int   `json:"status" binding:"omitempty,oneof=0 1"`
}

type ResetAdminPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6"`
}

type ToggleAdminStatusRequest struct {
	Status int `json:"status" binding:"required,oneof=0 1"`
}

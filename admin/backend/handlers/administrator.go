package handlers

import (
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"
	"akrick.com/mychat/admin/backend/utils"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// AdminLogin2 管理员登录(使用Administrator表)
// @Summary 管理员登录
// @Description 管理员登录接口
// @Tags 管理员
// @Accept json
// @Produce json
// @Param request body AdminLoginRequest true "登录信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:登录成功,data:{token,admin}"
// @Router /api/admin2/login [post]
func AdminLogin2(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	var admin models.Administrator
	if err := database.DB.Where("username = ?", req.Username).First(&admin).Error; err != nil {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, admin.Password) {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 检查状态
	if admin.Status != 1 {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "账号已被禁用",
		})
		return
	}

	// 更新最后登录时间
	database.DB.Model(&admin).Update("last_login", time.Now())

	// 生成Token
	token, err := utils.GenerateToken(uint(admin.ID), admin.Username)
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
			"admin": gin.H{
				"id":        admin.ID,
				"username":  admin.Username,
				"real_name": admin.RealName,
				"email":     admin.Email,
				"phone":     admin.Phone,
				"avatar":    admin.Avatar,
				"role":      admin.Role,
			},
		},
	})
}

// GetAdminInfo2 获取当前管理员信息
// @Summary 获取当前管理员信息
// @Description 获取当前登录的管理员信息
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin2/info [get]
func GetAdminInfo2(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	var admin models.Administrator
	if err := database.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "管理员不存在",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": admin,
	})
}

// GetAdministratorList 获取管理员列表
// @Summary 获取管理员列表
// @Description 获取所有管理员列表
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Param status query int false "状态:0-禁用,1-正常"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin2/administrators [get]
func GetAdministratorList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	keyword := c.Query("keyword")
	status := c.Query("status")

	query := database.DB.Model(&models.Administrator{})

	// 搜索
	if keyword != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ? OR email LIKE ? OR phone LIKE ?", 
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var admins []models.Administrator
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (utils.ParseInt(page) - 1) * utils.ParseInt(pageSize)
	}

	if err := query.Offset(offset).Limit(utils.ParseInt(pageSize)).
		Order("created_at DESC").Find(&admins).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败: " + err.Error(),
		})
		return
	}

	// 隐藏密码
	for i := range admins {
		admins[i].Password = ""
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"administrators": admins,
			"total":          total,
		},
	})
}

// CreateAdministrator 创建管理员
// @Summary 创建管理员
// @Description 创建新的管理员账号
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateAdministratorRequest true "管理员信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:创建成功"
// @Router /api/admin2/administrators [post]
func CreateAdministrator(c *gin.Context) {
	var req CreateAdministratorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 检查用户名是否存在
	var count int64
	database.DB.Model(&models.Administrator{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "用户名已存在",
		})
		return
	}

	// 检查邮箱是否存在
	if req.Email != "" {
		database.DB.Model(&models.Administrator{}).Where("email = ?", req.Email).Count(&count)
		if count > 0 {
			c.JSON(400, gin.H{
				"code": 400,
				"msg":  "邮箱已存在",
			})
			return
		}
	}

	// 检查手机号是否存在
	if req.Phone != "" {
		database.DB.Model(&models.Administrator{}).Where("phone = ?", req.Phone).Count(&count)
		if count > 0 {
			c.JSON(400, gin.H{
				"code": 400,
				"msg":  "手机号已存在",
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

	admin := models.Administrator{
		Username: req.Username,
		Password: hashedPassword,
		RealName: req.RealName,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Role:     req.Role,
		Status:   req.Status,
	}

	if err := database.DB.Create(&admin).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "创建失败: " + err.Error(),
		})
		return
	}

	// 隐藏密码
	admin.Password = ""

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": admin,
	})
}

// UpdateAdministrator 更新管理员
// @Summary 更新管理员
// @Description 更新管理员信息
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "管理员ID"
// @Param request body UpdateAdministratorRequest true "更新信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:更新成功"
// @Router /api/admin2/administrators/{id} [put]
func UpdateAdministrator(c *gin.Context) {
	adminID := c.Param("id")

	// 获取当前登录管理员ID
	currentAdminID, _ := c.Get("admin_id")

	// 不允许修改自己
	if currentAdminID != nil && adminID == fmt.Sprintf("%v", currentAdminID) {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "不能修改自己的信息",
		})
		return
	}

	var req UpdateAdministratorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	var admin models.Administrator
	if err := database.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "管理员不存在",
		})
		return
	}

	// 检查用户名是否被其他管理员使用
	if req.Username != "" && req.Username != admin.Username {
		var count int64
		database.DB.Model(&models.Administrator{}).Where("username = ? AND id != ?", req.Username, adminID).Count(&count)
		if count > 0 {
			c.JSON(400, gin.H{
				"code": 400,
				"msg":  "用户名已被其他管理员使用",
			})
			return
		}
	}

	updates := make(map[string]interface{})
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Role != "" {
		updates["role"] = req.Role
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

// DeleteAdministrator 删除管理员
// @Summary 删除管理员
// @Description 删除管理员账号
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "管理员ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:删除成功"
// @Router /api/admin2/administrators/{id} [delete]
func DeleteAdministrator(c *gin.Context) {
	adminID := c.Param("id")

	// 获取当前登录管理员ID
	currentAdminID, _ := c.Get("admin_id")

	// 不允许删除自己
	if currentAdminID != nil && adminID == fmt.Sprintf("%v", currentAdminID) {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "不能删除自己",
		})
		return
	}

	var admin models.Administrator
	if err := database.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "管理员不存在",
		})
		return
	}

	// 不允许删除超级管理员
	if admin.Role == "super_admin" {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "不能删除超级管理员",
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

// ResetAdministratorPassword 重置管理员密码
// @Summary 重置管理员密码
// @Description 重置管理员密码
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "管理员ID"
// @Param request body ResetPasswordRequest true "密码信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:重置成功"
// @Router /api/admin2/administrators/{id}/password [post]
func ResetAdministratorPassword(c *gin.Context) {
	adminID := c.Param("id")

	// 不允许重置超级管理员密码
	var targetAdmin models.Administrator
	if err := database.DB.First(&targetAdmin, adminID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "管理员不存在",
		})
		return
	}

	if targetAdmin.Role == "super_admin" {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "不能重置超级管理员密码",
		})
		return
	}

	var req AdminResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
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

	if err := database.DB.Model(&targetAdmin).Update("password", hashedPassword).Error; err != nil {
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

// ToggleAdministratorStatus 切换管理员状态
// @Summary 切换管理员状态
// @Description 禁用或启用管理员账号
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "管理员ID"
// @Param request body ToggleStatusRequest true "状态信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:操作成功"
// @Router /api/admin2/administrators/{id}/status [put]
func ToggleAdministratorStatus(c *gin.Context) {
	adminID := c.Param("id")

	// 获取当前登录管理员ID
	currentAdminID, _ := c.Get("admin_id")

	// 不允许禁用自己
	if currentAdminID != nil && adminID == fmt.Sprintf("%v", currentAdminID) {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "不能禁用自己",
		})
		return
	}

	var admin models.Administrator
	if err := database.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "管理员不存在",
		})
		return
	}

	// 不允许禁用超级管理员
	if admin.Role == "super_admin" {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "不能禁用超级管理员",
		})
		return
	}

	var req ToggleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
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

// UpdateMyProfile 更新个人资料
// @Summary 更新个人资料
// @Description 管理员更新自己的资料
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AdminUpdateProfileRequest true "个人资料"
// @Success 200 {object} map[string]interface{} "code:200,msg:更新成功"
// @Router /api/admin2/profile [put]
func UpdateMyProfile(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	var req AdminUpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	var admin models.Administrator
	if err := database.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "管理员不存在",
		})
		return
	}

	updates := make(map[string]interface{})
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
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

// ChangeMyPassword 修改自己的密码
// @Summary 修改自己的密码
// @Description 管理员修改自己的密码
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AdminChangePasswordRequest true "密码信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:修改成功"
// @Router /api/admin2/password [post]
func ChangeMyPassword(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	var req AdminChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	var admin models.Administrator
	if err := database.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "管理员不存在",
		})
		return
	}

	// 验证原密码
	if !utils.CheckPassword(req.OldPassword, admin.Password) {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "原密码错误",
		})
		return
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
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
			"msg":  "修改失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})
}

// AdminLogout2 管理员退出登录
// @Summary 管理员退出登录
// @Description 管理员退出登录接口
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:退出成功"
// @Router /api/admin2/logout [post]
func AdminLogout2(c *gin.Context) {
	// 可以在这里添加将token加入黑名单的逻辑
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "退出成功",
	})
}

// GetAdminPermissions2 获取管理员权限列表(使用Administrator表)
// @Summary 获取管理员权限列表
// @Description 获取当前管理员的权限列表
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin2/permissions [get]
func GetAdminPermissions2(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	// 获取管理员
	var admin models.Administrator
	database.DB.First(&admin, adminID)

	// 根据角色返回权限
	// 简单实现：根据角色返回固定权限
	permissions := []string{}
	roles := []string{admin.Role}

	if admin.Role == "super_admin" {
		// 超级管理员拥有所有权限
		permissions = []string{
			"dashboard:view",
			"user:view", "user:add", "user:edit", "user:delete", "user:reset_password",
			"counselor:view", "counselor:add", "counselor:edit", "counselor:delete",
			"order:view", "order:edit_status", "order:view_statistics",
			"chat:view", "chat:manage",
			"finance:view", "finance:approve_withdraw",
			"system:view", "system:edit_config", "system:view_logs", "system:manage_online",
			"admin:view", "admin:add", "admin:edit", "admin:delete", "admin:reset_password",
			"role:view", "role:add", "role:edit", "role:delete",
			"permission:view", "permission:add", "permission:edit", "permission:delete",
			"menu:view", "menu:add", "menu:edit", "menu:delete",
		}
	} else if admin.Role == "admin" {
		// 普通管理员拥有部分权限
		permissions = []string{
			"dashboard:view",
			"user:view",
			"counselor:view", "counselor:edit",
			"order:view", "order:edit_status",
			"chat:view",
			"finance:view",
			"system:view",
		}
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"permissions": permissions,
			"roles":       roles,
		},
	})
}

// Request structs for Administrator management
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateAdministratorRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	RealName string `json:"real_name" binding:"required,max=50"`
	Email    string `json:"email" binding:"omitempty,email,max=100"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role" binding:"required,oneof=admin super_admin"`
	Status   int    `json:"status" binding:"omitempty,oneof=0 1"`
}

type UpdateAdministratorRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=50"`
	RealName string `json:"real_name" binding:"omitempty,max=50"`
	Email    string `json:"email" binding:"omitempty,email,max=100"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role" binding:"omitempty,oneof=admin super_admin"`
	Status   *int   `json:"status" binding:"omitempty,oneof=0 1"`
}

type AdminResetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6"`
}

type ToggleStatusRequest struct {
	Status int `json:"status" binding:"required,oneof=0 1"`
}

type AdminUpdateProfileRequest struct {
	RealName string `json:"real_name" binding:"omitempty,max=50"`
	Email    string `json:"email" binding:"omitempty,email,max=100"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
	Avatar   string `json:"avatar"`
}

type AdminChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

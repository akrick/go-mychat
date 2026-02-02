package handlers

import (
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"
	"akrick.com/mychat/admin/backend/utils"
	"github.com/gin-gonic/gin"
)

// UpdateProfile godoc
// @Summary 更新用户资料
// @Description 更新当前用户的资料信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UserUpdateProfileRequest true "更新信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:更新成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "更新失败"
// @Router /api/user/profile [put]
func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req UserUpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
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

	if len(updates) == 0 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "没有需要更新的字段",
		})
		return
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
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

// ChangePassword godoc
// @Summary 修改密码
// @Description 修改当前用户的密码
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UserChangePasswordRequest true "密码信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:修改成功"
// @Failure 400 {object} map[string]interface{} "参数错误或旧密码错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "修改失败"
// @Router /api/user/password [post]
func ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req UserChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "旧密码错误",
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

	if err := database.DB.Model(&user).Update("password", hashedPassword).Error; err != nil {
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

// GetStatistics godoc
// @Summary 获取统计数据
// @Description 获取系统统计数据（管理员接口）
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{statistics}"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Router /api/admin/statistics [get]
func GetStatistics(c *gin.Context) {
	var userCount int64
	var counselorCount int64
	var orderCount int64
	var totalAmount float64

	database.DB.Model(&models.User{}).Count(&userCount)
	database.DB.Model(&models.Counselor{}).Count(&counselorCount)
	database.DB.Model(&models.Order{}).Count(&orderCount)
	database.DB.Model(&models.Order{}).Where("status = ?", models.OrderStatusPaid).Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"user_count":       userCount,
			"counselor_count":  counselorCount,
			"order_count":      orderCount,
			"total_amount":     totalAmount,
		},
	})
}

// GetUserList godoc
// @Summary 获取用户列表
// @Description 获取普通用户列表（管理员接口）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Param status query int false "状态:0-禁用,1-正常"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/users [get]
func GetUserList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	keyword := c.Query("keyword")
	status := c.Query("status")

	query := database.DB.Model(&models.User{})

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

	var users []models.User
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (utils.ParseInt(page) - 1) * utils.ParseInt(pageSize)
	}

	if err := query.Offset(offset).Limit(utils.ParseInt(pageSize)).Order("created_at DESC").Find(&users).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败: " + err.Error(),
		})
		return
	}

	// 隐藏密码
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"users": users,
			"total": total,
		},
	})
}

// CreateUser godoc
// @Summary 创建用户
// @Description 创建普通用户（管理员接口）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UserCreateRequest true "用户信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:创建成功"
// @Router /api/admin/users [post]
func CreateUser(c *gin.Context) {
	var req UserCreateRequest
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

	// 检查手机号是否存在
	if req.Phone != "" {
		database.DB.Model(&models.User{}).Where("phone = ?", req.Phone).Count(&count)
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

	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Status:   req.Status,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "创建失败: " + err.Error(),
		})
		return
	}

	// 隐藏密码
	user.Password = ""

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": user,
	})
}

// UpdateUser godoc
// @Summary 更新用户
// @Description 更新普通用户信息（管理员接口）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param request body UserUpdateRequest true "更新信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:更新成功"
// @Router /api/admin/users/{id} [put]
func UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var req UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "用户不存在",
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

	if len(updates) == 0 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "没有需要更新的字段",
		})
		return
	}

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
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

// DeleteUser godoc
// @Summary 删除用户
// @Description 删除普通用户（管理员接口）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:删除成功"
// @Router /api/admin/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	// 软删除
	if err := database.DB.Delete(&user).Error; err != nil {
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

// ResetUserPassword godoc
// @Summary 重置用户密码
// @Description 重置普通用户密码（管理员接口）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param request body UserResetPasswordRequest true "密码信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:重置成功"
// @Router /api/admin/users/{id}/password [post]
func ResetUserPassword(c *gin.Context) {
	userID := c.Param("id")

	var req UserResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "用户不存在",
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

	if err := database.DB.Model(&user).Update("password", hashedPassword).Error; err != nil {
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

// Request structs for User management
type UserCreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"omitempty,email,max=100"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
	Avatar   string `json:"avatar"`
	Status   int    `json:"status" binding:"omitempty,oneof=0 1"`
}

type UserUpdateRequest struct {
	Email  string `json:"email" binding:"omitempty,email,max=100"`
	Phone  string `json:"phone" binding:"omitempty,max=20"`
	Avatar string `json:"avatar"`
	Status *int   `json:"status" binding:"omitempty,oneof=0 1"`
}

type UserUpdateProfileRequest struct {
	Email  string `json:"email" binding:"omitempty,email"`
	Phone  string `json:"phone" binding:"omitempty,len=11"`
	Avatar string `json:"avatar"`
}

type UserChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type UserResetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6"`
}

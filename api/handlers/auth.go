package handlers

import (
	"context"
	"akrick.com/mychat/cache"
	"akrick.com/mychat/database"
	"akrick.com/mychat/models"
	"akrick.com/mychat/utils"
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" example:"testuser"`
	Password string `json:"password" binding:"required,min=6" example:"123456"`
	Email    string `json:"email" example:"test@example.com"`
	Phone    string `json:"phone" example:"13800138000"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"testuser"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required" example:"your-jwt-token"`
}

// Register godoc
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:注册成功,data:{user_id,username}"
// @Failure 400 {object} map[string]interface{} "参数错误或用户名已存在"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "用户名已存在",
		})
		return
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "密码加密失败",
		})
		return
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   1,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "注册失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "注册成功",
		"data": gin.H{
			"user_id":  user.ID,
			"username": user.Username,
		},
	})
}

// Login godoc
// @Summary 用户登录
// @Description 使用用户名和密码登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:登录成功,data:{token,user}"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "用户名或密码错误"
// @Failure 403 {object} map[string]interface{} "账户已被禁用"
// @Router /api/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 查找用户
	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "账户已被禁用",
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

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "生成token失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"user_id":  user.ID,
				"username": user.Username,
				"email":    user.Email,
			},
		},
	})
}

// RefreshToken godoc
// @Summary 刷新Token
// @Description 使用现有token刷新获取新的token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "旧token"
// @Success 200 {object} map[string]interface{} "code:200,msg:刷新成功,data:{token}"
// @Failure 400 {object} map[string]interface{} "参数错误或token不需要刷新"
// @Router /api/token/refresh [post]
func RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	newToken, err := utils.RefreshToken(req.Token)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "刷新token失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "刷新成功",
		"data": gin.H{
			"token": newToken,
		},
	})
}

// GetUserInfo godoc
// @Summary 获取用户信息
// @Description 获取当前登录用户的信息（需要认证，使用Redis缓存和SingleFlight防穿透）
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{user_info}"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "用户不存在"
// @Router /api/user/info [get]
func GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")

	ctx := context.Background()
	
	// 尝试从缓存获取
	var user models.User
	var fromCache bool
	
	if cache.Rdb != nil {
		userInfo, err := cache.GetUserInfoWithCache(ctx, userID.(uint))
		if err == nil {
			user = models.User{
				ID:       userInfo.UserID,
				Username: userInfo.Username,
				Email:    userInfo.Email,
				Phone:    userInfo.Phone,
				Status:   userInfo.Status,
			}
			fromCache = true
		}
	}

	// 缓存未命中或Redis未连接，从数据库查询
	if !fromCache {
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(404, gin.H{
				"code": 404,
				"msg":  "用户不存在",
			})
			return
		}
	}

	msg := "获取成功"
	if fromCache {
		msg = "获取成功（来自缓存）"
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  msg,
		"data": gin.H{
			"id":       user.ID,
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
			"phone":    user.Phone,
			"avatar":   user.Avatar,
			"balance":  user.Balance,
		},
	})
}

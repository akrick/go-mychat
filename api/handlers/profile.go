package handlers

import (
	"fmt"
	"net/http"

	"akrick.com/mychat/database"
	"akrick.com/mychat/models"
	"akrick.com/mychat/utils"

	"github.com/gin-gonic/gin"
)

// GetProfile godoc
// @Summary 获取个人信息
// @Description 获取当前登录用户的信息
// @Tags 个人中心
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/profile [get]
func GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	// 清除敏感信息
	user.Password = ""

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": user,
	})
}

// UploadAvatar godoc
// @Summary 上传头像
// @Description 上传用户头像
// @Tags 个人中心
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "头像文件"
// @Success 200 {object} map[string]interface{} "code:200,msg:上传成功"
// @Router /api/admin/upload [post]
func UploadAvatar(c *gin.Context) {
	userID := c.GetUint("user_id")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "文件上传失败: " + err.Error(),
		})
		return
	}

	// 验证文件类型
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/jpg":  true,
	}

	fileHeader, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "文件打开失败",
		})
		return
	}
	defer fileHeader.Close()

	// 读取文件头部以检查类型
	buffer := make([]byte, 512)
	_, err = fileHeader.Read(buffer)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "文件读取失败",
		})
		return
	}

	contentType := http.DetectContentType(buffer)
	if !allowedTypes[contentType] {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "只支持JPG、PNG格式的图片",
		})
		return
	}

	// 验证文件大小（2MB）
	if file.Size > 2*1024*1024 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "图片大小不能超过2MB",
		})
		return
	}

	// 生成文件名
	filename := fmt.Sprintf("avatar_%d_%s", userID, file.Filename)

	// 保存文件路径
	uploadDir := "./uploads/avatars"
	if err := utils.EnsureDir(uploadDir); err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "创建上传目录失败",
		})
		return
	}

	// 保存文件
	filepath := fmt.Sprintf("%s/%s", uploadDir, filename)
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "文件保存失败",
		})
		return
	}

	// 返回访问URL
	url := fmt.Sprintf("/uploads/avatars/%s", filename)

	// 更新用户头像
	database.DB.Model(&models.User{}).Where("id = ?", userID).Update("avatar", url)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "上传成功",
		"data": gin.H{
			"url": url,
		},
	})
}

// RechargeRequest 充值请求
type RechargeRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// Recharge godoc
// @Summary 账户充值
// @Description 充值用户账户余额
// @Tags 个人中心
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RechargeRequest true "充值金额"
// @Success 200 {object} map[string]interface{} "code:200,msg:充值成功"
// @Router /api/user/recharge [post]
func Recharge(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 获取用户
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	// 更新余额
	newBalance := user.Balance + req.Amount
	if err := database.DB.Model(&user).Update("balance", newBalance).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "充值失败: " + err.Error(),
		})
		return
	}

	// 创建交易记录
	transaction := models.UserTransaction{
		UserID:      userID,
		Type:        "recharge",
		Amount:      req.Amount,
		Description: "账户充值",
		Balance:     newBalance,
	}
	if err := database.DB.Create(&transaction).Error; err != nil {
		// 交易记录创建失败不影响充值
		fmt.Printf("创建交易记录失败: %v\n", err)
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "充值成功",
		"data": gin.H{
			"balance": newBalance,
		},
	})
}

// GetTransactions godoc
// @Summary 获取交易记录
// @Description 获取当前用户的交易记录
// @Tags 个人中心
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param type query string false "交易类型:recharge/consume/refund"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/user/transactions [get]
func GetTransactions(c *gin.Context) {
	userID := c.GetUint("user_id")

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	transactionType := c.Query("type")

	query := database.DB.Model(&models.UserTransaction{}).Where("user_id = ?", userID)

	// 类型筛选
	if transactionType != "" {
		query = query.Where("type = ?", transactionType)
	}

	var total int64
	query.Count(&total)

	var transactions []models.UserTransaction
	offset := 0
	if page != "1" {
		p := utils.ParseInt(page)
		ps := utils.ParseInt(pageSize)
		offset = (p - 1) * ps
	}

	ps := utils.ParseInt(pageSize)
	if err := query.Offset(offset).Limit(ps).Order("created_at DESC").Find(&transactions).Error; err != nil {
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
			"transactions": transactions,
			"total":        total,
		},
	})
}

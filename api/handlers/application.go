package handlers

import (
	"akrick.com/mychat/database"
	"akrick.com/mychat/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateCounselorApplication 提交咨询师入驻申请
// @Summary 提交咨询师入驻申请
// @Tags 咨询师入驻
// @Accept json
// @Produce json
// @Param request body models.CounselorApplication true "申请信息"
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "已有申请或已是咨询师"
// @Router /api/counselor/application [post]
func CreateCounselorApplication(c *gin.Context) {
	var req models.CounselorApplication

	// 绑定JSON参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"error": err.Error(),
		})
		return
	}

	// 从JWT获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未登录",
		})
		return
	}
	req.UserID = uint(userID.(float64))

	// 检查是否已有申请
	var existingApp models.CounselorApplication
	result := database.DB.Where("user_id = ?", req.UserID).First(&existingApp)
	if result.Error == nil {
		if existingApp.Status == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "您已提交入驻申请，请等待审核",
			})
			return
		} else if existingApp.Status == 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "您的申请已审核通过",
			})
			return
		}
		// 如果是拒绝状态，允许重新申请
	}

	// 检查用户是否已经是咨询师
	var counselor models.Counselor
	result = database.DB.Where("user_id = ?", req.UserID).First(&counselor)
	if result.Error == nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "您已是认证咨询师",
		})
		return
	}

	// 验证必填字段
	if req.Name == "" || req.Phone == "" || req.Title == "" || req.Specialty == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "姓名、电话、职称和擅长领域为必填项",
		})
		return
	}

	// 至少需要一张证书
	if req.CertificateImg1 == "" && req.CertificateImg2 == "" && req.CertificateImg3 == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "至少需要上传一张资质证书",
		})
		return
	}

	// 设置状态为待审核
	req.Status = 0

	// 创建申请
	if err := database.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "提交申请失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "申请提交成功，请等待审核",
		"data": gin.H{
			"id": req.ID,
		},
	})
}

// GetMyApplication 获取我的入驻申请
// @Summary 获取我的入驻申请
// @Tags 咨询师入驻
// @Produce json
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Router /api/counselor/my-application [get]
func GetMyApplication(c *gin.Context) {
	// 从JWT获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未登录",
		})
		return
	}

	var application models.CounselorApplication
	result := database.DB.Where("user_id = ?", userID).Order("id desc").First(&application)

	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "暂无入驻申请",
			"data": nil,
		})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询申请失败",
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": application,
	})
}

// GetAllApplications 获取所有入驻申请（管理员）
// @Summary 获取所有入驻申请
// @Tags 咨询师入驻
// @Produce json
// @Param status query int false "状态筛选:0-待审核,1-审核通过,2-审核拒绝"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Router /api/counselor/applications [get]
func GetAllApplications(c *gin.Context) {
	// 从JWT获取用户ID和是否管理员
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未登录",
		})
		return
	}

	// 检查是否管理员
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户不存在",
		})
		return
	}

	if !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "无权限访问",
		})
		return
	}

	// 获取查询参数
	statusStr := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 构建查询
	query := database.DB.Model(&models.CounselorApplication{})
	if statusStr != "" {
		status, err := strconv.Atoi(statusStr)
		if err == nil {
			query = query.Where("status = ?", status)
		}
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var applications []models.CounselorApplication
	offset := (page - 1) * pageSize
	if err := query.Preload("User").Preload("Reviewer").Offset(offset).Limit(pageSize).Order("created_at desc").Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": gin.H{
			"list":      applications,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ReviewApplication 审核入驻申请（管理员）
// @Summary 审核入驻申请
// @Tags 咨询师入驻
// @Accept json
// @Produce json
// @Param id path int true "申请ID"
// @Param request body object true "审核信息: status(1-通过,2-拒绝), reject_reason(拒绝时必填)"
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Router /api/counselor/application/:id/review [put]
func ReviewApplication(c *gin.Context) {
	// 从JWT获取用户ID和是否管理员
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未登录",
		})
		return
	}

	// 检查是否管理员
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户不存在",
		})
		return
	}

	if !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "无权限访问",
		})
		return
	}

	// 获取申请ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "申请ID错误",
		})
		return
	}

	// 绑定请求参数
	var req struct {
		Status       int    `json:"status" binding:"required"`
		RejectReason string `json:"reject_reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"error": err.Error(),
		})
		return
	}

	// 验证状态
	if req.Status != 1 && req.Status != 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "状态错误，1-通过，2-拒绝",
		})
		return
	}

	// 拒绝时必须有拒绝原因
	if req.Status == 2 && req.RejectReason == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "拒绝时必须填写拒绝原因",
		})
		return
	}

	// 查询申请
	var application models.CounselorApplication
	if err := database.DB.First(&application, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "申请不存在",
		})
		return
	}

	// 检查状态
	if application.Status != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "该申请已审核",
		})
		return
	}

	// 开启事务
	tx := database.DB.Begin()

	// 更新申请状态
	now := database.DB.NowFunc()
	application.Status = req.Status
	application.RejectReason = req.RejectReason
	application.ReviewedBy = uint(userID.(float64))
	application.ReviewedAt = &now
	if err := tx.Save(&application).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新申请失败",
			"error": err.Error(),
		})
		return
	}

	// 如果审核通过，创建咨询师账户
	if req.Status == 1 {
		// 创建咨询师记录
		counselor := models.Counselor{
			UserID:    application.UserID,
			Name:      application.Name,
			Title:     application.Title,
			Avatar:    application.User.Avatar,
			Bio:       application.Bio,
			Specialty: application.Specialty,
			Price:     5.0, // 默认价格5元/分钟
			YearsExp:  application.YearsExp,
			Rating:    5.0,
			Status:    1,
		}

		if err := tx.Create(&counselor).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "创建咨询师失败",
				"error": err.Error(),
			})
			return
		}

		// 创建统计记录
		statistics := models.CounselorStatistics{
			CounselorID: counselor.ID,
		}
		if err := tx.Create(&statistics).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "创建统计记录失败",
				"error": err.Error(),
			})
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "审核失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "审核成功",
	})
}

// UploadCertificate 上传证书图片
// @Summary 上传证书图片
// @Tags 咨询师入驻
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "证书图片"
// @Success 200 {object} map[string]interface{} "成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Router /api/counselor/upload-certificate [post]
func UploadCertificate(c *gin.Context) {
	// 从JWT获取用户ID
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未登录",
		})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请选择文件",
		})
		return
	}

	// 验证文件类型
	allowedTypes := []string{"image/jpeg", "image/jpg", "image/png"}
	isValidType := false
	for _, t := range allowedTypes {
		if file.Header.Get("Content-Type") == t {
			isValidType = true
			break
		}
	}

	if !isValidType {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "仅支持JPG、JPEG、PNG格式的图片",
		})
		return
	}

	// 限制文件大小（5MB）
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "图片大小不能超过5MB",
		})
		return
	}

	// 生成文件名
	ext := ".jpg"
	if file.Header.Get("Content-Type") == "image/png" {
		ext = ".png"
	}
	filename := fmt.Sprintf("certificate_%d_%d%s", uint(c.GetFloat64("user_id")), time.Now().Unix(), ext)

	// 保存文件路径
	uploadPath := "./uploads/certificates/"
	filePath := uploadPath + filename

	// 创建目录
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "文件上传失败",
			"error": err.Error(),
		})
		return
	}

	// 返回文件URL
	fileURL := "/uploads/certificates/" + filename

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "上传成功",
		"data": gin.H{
			"url": fileURL,
		},
	})
}

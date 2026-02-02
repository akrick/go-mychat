package handlers

import (
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"
	"akrick.com/mychat/admin/backend/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type ReviewApplicationRequest struct {
	Status        int    `json:"status" binding:"required,oneof=1 2"`
	RejectReason string `json:"reject_reason"`
}

// GetApplicationList 获取入驻申请列表
// @Summary 获取入驻申请列表
// @Description 获取所有入驻申请列表
// @Tags 入驻申请
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param status query int false "状态:0-待审核,1-审核通过,2-审核拒绝"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/counselor/applications [get]
func GetApplicationList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	status := c.Query("status")
	keyword := c.Query("keyword")

	query := database.DB.Model(&models.CounselorApplication{})

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 搜索
	if keyword != "" {
		query = query.Where("name LIKE ? OR phone LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var applications []models.CounselorApplication
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (utils.ParseInt(page) - 1) * utils.ParseInt(pageSize)
	}

	if err := query.Preload("User").Preload("Reviewer").
		Offset(offset).Limit(utils.ParseInt(pageSize)).
		Order("created_at DESC").Find(&applications).Error; err != nil {
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
			"list":  applications,
			"total": total,
		},
	})
}

// GetApplicationDetail 获取入驻申请详情
// @Summary 获取入驻申请详情
// @Description 根据ID获取入驻申请详情
// @Tags 入驻申请
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "申请ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/counselor/applications/{id} [get]
func GetApplicationDetail(c *gin.Context) {
	applicationID := c.Param("id")

	var application models.CounselorApplication
	if err := database.DB.Preload("User").Preload("Reviewer").
		First(&application, applicationID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "申请不存在",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": application,
	})
}

// ReviewApplication 审核入驻申请
// @Summary 审核入驻申请
// @Description 审核入驻申请，通过后自动创建咨询师账户
// @Tags 入驻申请
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "申请ID"
// @Param request body ReviewApplicationRequest true "审核信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:审核成功"
// @Router /api/admin/counselor/applications/{id}/review [put]
func ReviewApplication(c *gin.Context) {
	applicationID := c.Param("id")
	reviewerID, _ := c.Get("admin_id")

	var application models.CounselorApplication
	if err := database.DB.First(&application, applicationID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "申请不存在",
		})
		return
	}

	// 检查申请状态
	if application.Status != models.ApplicationStatusPending {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "该申请已处理",
		})
		return
	}

	var req ReviewApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 检查拒绝原因
	if req.Status == models.ApplicationStatusRejected && req.RejectReason == "" {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "请填写拒绝原因",
		})
		return
	}

	now := time.Now()

	// 更新申请状态
	updates := map[string]interface{}{
		"status":     req.Status,
		"reviewer_id": reviewerID,
		"reviewed_at": &now,
	}

	if req.RejectReason != "" {
		updates["reject_reason"] = req.RejectReason
	}

	// 如果审核通过，创建咨询师账户
	if req.Status == models.ApplicationStatusApproved {
		// 创建咨询师
		counselor := models.Counselor{
			Name:      application.Name,
			Title:     application.Title,
			Bio:       application.Bio,
			Specialty: application.Specialty,
			YearsExp:  application.YearsExp,
			Price:     1.0, // 默认单价
			Rating:    5.0,
			Status:    1,
		}

		if err := database.DB.Create(&counselor).Error; err != nil {
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "创建咨询师失败: " + err.Error(),
			})
			return
		}
	}

	// 更新申请状态
	if err := database.DB.Model(&application).Updates(updates).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "审核失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "审核成功",
	})
}

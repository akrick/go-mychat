package handlers

import (
	"context"
	"akrick.com/mychat/admin/backend/cache"
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"
	"akrick.com/mychat/admin/backend/utils"
	"github.com/gin-gonic/gin"
)

type CreateCounselorRequest struct {
	Name      string  `json:"name" binding:"required,min=2,max=50"`
	Title     string  `json:"title"`
	Avatar    string  `json:"avatar"`
	Bio       string  `json:"bio"`
	Specialty string  `json:"specialty"`
	Price     float64 `json:"price" binding:"required,min=0"`
	YearsExp  int     `json:"years_exp"`
	Rating    float64 `json:"rating"`
}

type UpdateCounselorRequest struct {
	Name      string  `json:"name" binding:"min=2,max=50"`
	Title     string  `json:"title"`
	Avatar    string  `json:"avatar"`
	Bio       string  `json:"bio"`
	Specialty string  `json:"specialty"`
	Price     float64 `json:"price" binding:"min=0"`
	YearsExp  int     `json:"years_exp"`
	Rating    float64 `json:"rating"`
	Status    *int    `json:"status" binding:"omitempty,oneof=0 1"`
}

// CreateCounselor godoc
// @Summary 创建咨询师
// @Description 创建新的咨询师（管理员接口）
// @Tags 咨询师
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCounselorRequest true "咨询师信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:创建成功,data:{counselor}"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Router /api/counselor/create [post]
func CreateCounselor(c *gin.Context) {
	var req CreateCounselorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	counselor := models.Counselor{
		Name:      req.Name,
		Title:     req.Title,
		Avatar:    req.Avatar,
		Bio:       req.Bio,
		Specialty: req.Specialty,
		Price:     req.Price,
		YearsExp:  req.YearsExp,
		Rating:    req.Rating,
		Status:    1,
	}

	if err := database.DB.Create(&counselor).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "创建失败: " + err.Error(),
		})
		return
	}

	// 删除咨询师列表缓存（创建新咨询师后列表会变化）
	if cache.Rdb != nil {
		// 这里可以添加删除列表缓存的逻辑
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": counselor,
	})
}

// GetCounselorList godoc
// @Summary 获取咨询师列表
// @Description 获取所有启用的咨询师列表
// @Tags 咨询师
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{counselors,total}"
// @Router /api/counselor/list [get]
func GetCounselorList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	var total int64
	query := database.DB.Model(&models.Counselor{}).Where("status = ?", 1)
	query.Count(&total)

	var counselors []models.Counselor
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (utils.ParseInt(page) - 1) * utils.ParseInt(pageSize)
	}

	if err := query.Offset(offset).Limit(utils.ParseInt(pageSize)).Order("rating DESC, created_at DESC").Find(&counselors).Error; err != nil {
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
			"counselors": counselors,
			"total":      total,
		},
	})
}

// GetCounselorDetail godoc
// @Summary 获取咨询师详情
// @Description 根据ID获取咨询师详细信息（使用Redis缓存和SingleFlight防穿透）
// @Tags 咨询师
// @Accept json
// @Produce json
// @Param id path int true "咨询师ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{counselor}"
// @Failure 404 {object} map[string]interface{} "咨询师不存在"
// @Router /api/counselor/{id} [get]
func GetCounselorDetail(c *gin.Context) {
	counselorID := c.Param("id")
	ctx := context.Background()

	var counselor models.Counselor
	var fromCache bool

	// 尝试从缓存获取
	if cache.Rdb != nil {
		id := parseUint(counselorID)
		counselorData, err := cache.GetCounselorWithCache(ctx, id)
		if err == nil {
			counselor = models.Counselor{
				ID:        counselorData.ID,
				Name:      counselorData.Name,
				Title:     counselorData.Title,
				Avatar:    counselorData.Avatar,
				Bio:       counselorData.Bio,
				Specialty: counselorData.Specialty,
				Price:     counselorData.Price,
				YearsExp:  counselorData.YearsExp,
				Rating:    counselorData.Rating,
				Status:    counselorData.Status,
				CreatedAt: counselorData.CreatedAt,
				UpdatedAt: counselorData.UpdatedAt,
			}
			fromCache = true
		}
	}

	// 缓存未命中，从数据库查询
	if !fromCache {
		if err := database.DB.First(&counselor, counselorID).Error; err != nil {
			c.JSON(404, gin.H{
				"code": 404,
				"msg":  "咨询师不存在",
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
		"data": counselor,
	})
}

// UpdateCounselor godoc
// @Summary 更新咨询师信息
// @Description 更新咨询师信息（仅咨询师本人或管理员）
// @Tags 咨询师
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "咨询师ID"
// @Param request body UpdateCounselorRequest true "更新信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:更新成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 403 {object} map[string]interface{} "无权操作"
// @Failure 404 {object} map[string]interface{} "咨询师不存在"
// @Router /api/counselor/{id} [put]
func UpdateCounselor(c *gin.Context) {
	counselorID := c.Param("id")

	var req UpdateCounselorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	var counselor models.Counselor
	if err := database.DB.First(&counselor, counselorID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "咨询师不存在",
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}
	if req.Specialty != "" {
		updates["specialty"] = req.Specialty
	}
	if req.Price > 0 {
		updates["price"] = req.Price
	}
	if req.YearsExp > 0 {
		updates["years_exp"] = req.YearsExp
	}
	if req.Rating > 0 {
		updates["rating"] = req.Rating
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := database.DB.Model(&counselor).Updates(updates).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "更新失败: " + err.Error(),
		})
		return
	}

	// 删除缓存
	if cache.Rdb != nil {
		cache.DeleteCounselorCache(context.Background(), counselor.ID)
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// DeleteCounselor godoc
// @Summary 删除咨询师
// @Description 删除咨询师（仅管理员可操作）
// @Tags 咨询师
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "咨询师ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:删除成功"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 403 {object} map[string]interface{} "无权操作"
// @Failure 404 {object} map[string]interface{} "咨询师不存在"
// @Router /api/counselor/{id} [delete]
func DeleteCounselor(c *gin.Context) {
	counselorID := c.Param("id")
	ctx := context.Background()

	var counselor models.Counselor
	if err := database.DB.First(&counselor, counselorID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "咨询师不存在",
		})
		return
	}

	if err := database.DB.Delete(&counselor).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "删除失败: " + err.Error(),
		})
		return
	}

	// 删除缓存
	if cache.Rdb != nil {
		cache.DeleteCounselorCache(ctx, counselor.ID)
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

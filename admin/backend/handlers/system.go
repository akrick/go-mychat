package handlers

import (
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"
	"akrick.com/mychat/admin/backend/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// GetSystemLogs godoc
// @Summary 获取系统日志列表
// @Description 获取系统操作日志
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(20)
// @Param operator query string false "操作人"
// @Param action query string false "操作类型"
// @Param startDate query string false "开始日期"
// @Param endDate query string false "结束日期"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/logs [get]
func GetSystemLogs(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "20")
	operator := c.Query("operator")
	action := c.Query("action")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	query := database.DB.Model(&models.SystemLog{})

	if operator != "" {
		query = query.Where("operator LIKE ?", "%"+operator+"%")
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate+" 23:59:59")
	}

	var total int64
	query.Count(&total)

	var logs []models.SystemLog
	offset := 0
	if page != "1" {
		offset = (utils.ParseInt(page) - 1) * utils.ParseInt(pageSize)
	}

	if err := query.Offset(offset).Limit(utils.ParseInt(pageSize)).
		Order("created_at DESC").Find(&logs).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"list":  logs,
		"total": total,
	})
}

// GetSystemConfigs godoc
// @Summary 获取系统配置列表
// @Description 获取系统配置
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(20)
// @Param category query string false "配置分类"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/configs [get]
func GetSystemConfigs(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "20")
	category := c.Query("category")

	query := database.DB.Model(&models.SystemConfig{})

	if category != "" {
		query = query.Where("category = ?", category)
	}

	var total int64
	query.Count(&total)

	var configs []models.SystemConfig
	offset := 0
	if page != "1" {
		offset = (utils.ParseInt(page) - 1) * utils.ParseInt(pageSize)
	}

	if err := query.Offset(offset).Limit(utils.ParseInt(pageSize)).
		Order("sort ASC, id ASC").Find(&configs).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"list":  configs,
		"total": total,
	})
}

// UpdateSystemConfig godoc
// @Summary 更新系统配置
// @Description 更新系统配置
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "配置ID"
// @Param request body models.SystemConfig true "配置信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:更新成功"
// @Router /api/admin/configs/{id} [put]
func UpdateSystemConfig(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Key      string `json:"key"`
		Value    string `json:"value"`
		Category string `json:"category"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	var config models.SystemConfig
	if err := database.DB.First(&config, id).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "配置不存在",
		})
		return
	}

	// 系统配置不允许修改键
	if config.IsSystem && config.Key != req.Key {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "系统配置不允许修改键",
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Value != "" {
		updates["value"] = req.Value
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}

	if err := database.DB.Model(&config).Updates(updates).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "更新失败: " + err.Error(),
		})
		return
	}

	// 记录日志
	go logSystemAction(c, "更新配置", "系统管理", "更新配置: "+config.Key, nil)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功",
		"data": config,
	})
}

// CreateSystemConfig godoc
// @Summary 创建系统配置
// @Description 创建系统配置
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.SystemConfig true "配置信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:创建成功"
// @Router /api/admin/configs [post]
func CreateSystemConfig(c *gin.Context) {
	var req models.SystemConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 检查键是否已存在
	var count int64
	database.DB.Model(&models.SystemConfig{}).Where("key = ?", req.Key).Count(&count)
	if count > 0 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "配置键已存在",
		})
		return
	}

	if err := database.DB.Create(&req).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "创建失败: " + err.Error(),
		})
		return
	}

	// 记录日志
	go logSystemAction(c, "创建配置", "系统管理", "创建配置: "+req.Key, nil)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": req,
	})
}

// DeleteSystemConfig godoc
// @Summary 删除系统配置
// @Description 删除系统配置
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "配置ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:删除成功"
// @Router /api/admin/configs/{id} [delete]
func DeleteSystemConfig(c *gin.Context) {
	id := c.Param("id")

	var config models.SystemConfig
	if err := database.DB.First(&config, id).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "配置不存在",
		})
		return
	}

	// 系统配置不允许删除
	if config.IsSystem {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "系统配置不允许删除",
		})
		return
	}

	if err := database.DB.Delete(&config).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "删除失败: " + err.Error(),
		})
		return
	}

	// 记录日志
	go logSystemAction(c, "删除配置", "系统管理", "删除配置: "+config.Key, nil)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// BatchSaveConfigs godoc
// @Summary 批量保存系统配置
// @Description 批量保存或更新系统配置
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "配置数组"
// @Success 200 {object} map[string]interface{} "code:200,msg:保存成功"
// @Router /api/admin/configs/batch [post]
func BatchSaveConfigs(c *gin.Context) {
	var req struct {
		Configs []models.SystemConfig `json:"configs"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 使用事务批量保存
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, config := range req.Configs {
		var existingConfig models.SystemConfig
		err := tx.Where("key = ?", config.Key).First(&existingConfig).Error

		if err != nil {
			// 创建新配置
			if err := tx.Create(&config).Error; err != nil {
				tx.Rollback()
				c.JSON(500, gin.H{
					"code": 500,
					"msg":  "保存失败: " + err.Error(),
				})
				return
			}
		} else {
			// 更新现有配置
			updates := map[string]interface{}{
				"value":    config.Value,
				"category": config.Category,
			}
			if config.Label != "" {
				updates["label"] = config.Label
			}
			if config.Type != "" {
				updates["type"] = config.Type
			}
			if err := tx.Model(&existingConfig).Updates(updates).Error; err != nil {
				tx.Rollback()
				c.JSON(500, gin.H{
					"code": 500,
					"msg":  "保存失败: " + err.Error(),
				})
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "保存失败: " + err.Error(),
		})
		return
	}

	// 记录日志
	go logSystemAction(c, "批量保存配置", "系统管理", "批量保存配置", nil)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "保存成功",
	})
}

// logSystemAction 记录系统日志
func logSystemAction(c *gin.Context, action, module, description string, requestData interface{}) {
	// 从上下文获取用户信息
	userID, exists := c.Get("user_id")
	username := ""
	if exists {
		username = userID.(string)
	}

	log := models.SystemLog{
		Operator:    username,
		Action:      action,
		Module:      module,
		Description: description,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
		RequestData: "",
		Status:      1,
		CreatedAt:   time.Now(),
	}

	if requestData != nil {
		// 简化处理，实际应该序列化
		log.RequestData = "操作数据"
	}

	database.DB.Create(&log)
}

package handlers

import (
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"
	"akrick.com/mychat/admin/backend/utils"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// GetFormList 获取表单列表
func GetFormList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	title := c.Query("title")

	query := database.DB.Model(&models.LowcodeForm{})

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	var total int64
	query.Count(&total)

	var forms []models.LowcodeForm
	offset := 0
	if page != "1" {
		offset = (utils.ParseInt(page) - 1) * utils.ParseInt(pageSize)
	}

	if err := query.Offset(offset).Limit(utils.ParseInt(pageSize)).Order("created_at DESC").Find(&forms).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":  forms,
			"total": total,
		},
	})
}

// SaveFormDesign 保存表单设计
func SaveFormDesign(c *gin.Context) {
	var req struct {
		ID          uint                   `json:"id"`
		Title       string                 `json:"title" binding:"required"`
		Description string                 `json:"description"`
		Items       map[string]interface{} `json:"items" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 序列化表单配置
	formJSON, err := json.Marshal(req.Items)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "序列化失败",
		})
		return
	}

	var form models.LowcodeForm

	// 如果有ID，更新；否则创建
	if req.ID > 0 {
		// 更新
		if err := database.DB.First(&form, req.ID).Error; err != nil {
			c.JSON(404, gin.H{
				"code": 404,
				"msg":  "表单不存在",
			})
			return
		}

		form.Title = req.Title
		form.Description = req.Description
		form.FormJSON = string(formJSON)

		if err := database.DB.Save(&form).Error; err != nil {
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "保存失败",
			})
			return
		}
	} else {
		// 创建
		form = models.LowcodeForm{
			Title:       req.Title,
			Description: req.Description,
			FormJSON:    string(formJSON),
		}

		if err := database.DB.Create(&form).Error; err != nil {
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "保存失败",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "保存成功",
		"data": form,
	})
}

// GetFormDesign 获取表单设计
func GetFormDesign(c *gin.Context) {
	formID := c.Param("id")

	var form models.LowcodeForm
	if err := database.DB.First(&form, formID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "表单不存在",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": form,
	})
}

// DeleteForm 删除表单
func DeleteForm(c *gin.Context) {
	formID := c.Param("id")

	if err := database.DB.Delete(&models.LowcodeForm{}, formID).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "删除失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// SubmitFormData 提交表单数据
func SubmitFormData(c *gin.Context) {
	formID := c.Param("id")

	var formData map[string]interface{}
	if err := c.ShouldBindJSON(&formData); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 检查表单是否存在
	var form models.LowcodeForm
	if err := database.DB.First(&form, formID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "表单不存在",
		})
		return
	}

	// 序列化表单数据
	dataJSON, err := json.Marshal(formData)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "序列化失败",
		})
		return
	}

	// 保存表单数据
	formDataRecord := map[string]interface{}{
		"form_id":   formID,
		"form_data": string(dataJSON),
	}

	if err := database.DB.Table("form_data_records").Create(&formDataRecord).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "保存失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "提交成功",
	})
}

// GetFormDataList 获取表单数据列表
func GetFormDataList(c *gin.Context) {
	formID := c.Param("id")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	query := database.DB.Table("form_data_records").Where("form_id = ?", formID)

	var total int64
	query.Count(&total)

	var records []map[string]interface{}
	offset := 0
	if page != "1" {
		offset = (utils.ParseInt(page) - 1) * utils.ParseInt(pageSize)
	}

	if err := query.Offset(offset).Limit(utils.ParseInt(pageSize)).Order("created_at DESC").Find(&records).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":  records,
			"total": total,
		},
	})
}



// GetPageList 获取页面列表
func GetPageList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	title := c.Query("title")

	query := database.DB.Model(&models.LowcodePage{})

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	var total int64
	query.Count(&total)

	var pages []models.LowcodePage
	offset := 0
	if page != "1" {
		offset = (utils.ParseInt(page) - 1) * utils.ParseInt(pageSize)
	}

	if err := query.Offset(offset).Limit(utils.ParseInt(pageSize)).Order("created_at DESC").Find(&pages).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":  pages,
			"total": total,
		},
	})
}

// SavePageDesign 保存页面设计
func SavePageDesign(c *gin.Context) {
	var req struct {
		ID    uint                   `json:"id"`
		Title string                 `json:"title" binding:"required"`
		Path  string                 `json:"path" binding:"required"`
		Items map[string]interface{} `json:"items" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	pageJSON, err := json.Marshal(req.Items)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "序列化失败",
		})
		return
	}

	var page models.LowcodePage

	// 如果有ID，更新；否则创建
	if req.ID > 0 {
		// 更新
		if err := database.DB.First(&page, req.ID).Error; err != nil {
			c.JSON(404, gin.H{
				"code": 404,
				"msg":  "页面不存在",
			})
			return
		}

		page.Title = req.Title
		page.Description = req.Path
		page.PageJSON = string(pageJSON)

		if err := database.DB.Save(&page).Error; err != nil {
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "保存失败",
			})
			return
		}
	} else {
		// 创建
		page = models.LowcodePage{
			Title:       req.Title,
			Description: req.Path,
			PageJSON:    string(pageJSON),
		}

		if err := database.DB.Create(&page).Error; err != nil {
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "保存失败",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "保存成功",
		"data": page,
	})
}

// GetPageDesign 获取页面设计
func GetPageDesign(c *gin.Context) {
	pageID := c.Param("id")

	var page models.LowcodePage
	if err := database.DB.First(&page, pageID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "页面不存在",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": page,
	})
}

// DeletePage 删除页面
func DeletePage(c *gin.Context) {
	pageID := c.Param("id")

	if err := database.DB.Delete(&models.LowcodePage{}, pageID).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "删除失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// PreviewPage 预览页面
func PreviewPage(c *gin.Context) {
	pageID := c.Param("id")

	var page models.LowcodePage
	if err := database.DB.First(&page, pageID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "页面不存在",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"title":     page.Title,
			"page_json": page.PageJSON,
		},
	})
}

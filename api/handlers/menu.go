package handlers

import (
	"akrick.com/mychat/database"
	"akrick.com/mychat/models"

	"github.com/gin-gonic/gin"
)

// GetMenuTree 获取菜单树
// @Summary 获取菜单树
// @Description 获取菜单树结构
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/menus/tree [get]
func GetMenuTree(c *gin.Context) {
	// 获取类型为menu的权限，作为菜单
	var permissions []models.Permission
	database.DB.Where("type = ?", "menu").Order("sort ASC").Find(&permissions)

	tree := buildPermissionTree(permissions, 0)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": tree,
	})
}

// GetMenuList 获取菜单列表
// @Summary 获取菜单列表
// @Description 获取所有菜单列表
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/menus [get]
func GetMenuList(c *gin.Context) {
	// 获取类型为menu的权限
	var permissions []models.Permission
	database.DB.Where("type = ?", "menu").Order("sort ASC").Find(&permissions)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": permissions,
	})
}

// CreateMenu 创建菜单
// @Summary 创建菜单
// @Description 创建新菜单
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Permission true "菜单信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:创建成功"
// @Router /api/admin/menus [post]
func CreateMenu(c *gin.Context) {
	var req models.Permission
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 设置类型为菜单
	req.Type = "menu"

	// 检查菜单代码是否存在
	var count int64
	database.DB.Model(&models.Permission{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "菜单代码已存在",
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

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": req,
	})
}

// UpdateMenu 更新菜单
// @Summary 更新菜单
// @Description 更新菜单信息
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "菜单ID"
// @Param request body models.Permission true "菜单信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:更新成功"
// @Router /api/admin/menus/{id} [put]
func UpdateMenu(c *gin.Context) {
	menuID := c.Param("id")

	var req models.Permission
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 设置类型为菜单
	req.Type = "menu"

	if err := database.DB.Model(&models.Permission{}).Where("id = ?", menuID).Updates(&req).Error; err != nil {
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

// DeleteMenu 删除菜单
// @Summary 删除菜单
// @Description 删除指定菜单
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "菜单ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:删除成功"
// @Router /api/admin/menus/{id} [delete]
func DeleteMenu(c *gin.Context) {
	menuID := c.Param("id")

	// 检查是否有子菜单
	var count int64
	database.DB.Model(&models.Permission{}).Where("parent_id = ? AND type = ?", menuID, "menu").Count(&count)
	if count > 0 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "该菜单下有子菜单，无法删除",
		})
		return
	}

	if err := database.DB.Delete(&models.Permission{}, menuID).Error; err != nil {
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

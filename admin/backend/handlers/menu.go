package handlers

import (
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

// MenuNode 菜单树节点
type MenuNode struct {
	ID       int        `json:"id"`
	ParentID *int       `json:"parent_id"`
	Name     string     `json:"name"`
	Type     int        `json:"type"`
	Path     string     `json:"path"`
	Component string    `json:"component"`
	Permission string   `json:"permission"`
	Icon     string     `json:"icon"`
	Sort     int        `json:"sort"`
	Status   int        `json:"status"`
	Children []MenuNode `json:"children,omitempty"`
}

// GetMenuTree godoc
// @Summary 获取菜单树
// @Description 获取菜单树结构
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/menus/tree [get]
func GetMenuTree(c *gin.Context) {
	var menus []models.Menu
	if err := database.DB.Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败: " + err.Error(),
		})
		return
	}

	menuTree := buildMenuTree(menus, 0)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": menuTree,
	})
}

// GetMenuList godoc
// @Summary 获取菜单列表
// @Description 获取所有菜单列表(带分页)
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/menus [get]
func GetMenuList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	name := c.Query("name")

	query := database.DB.Model(&models.Menu{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	var total int64
	query.Count(&total)

	var menus []models.Menu
	offset := 0
	if page != "1" {
		offset = (parseInt(page) - 1) * parseInt(pageSize)
	}

	if err := query.Offset(offset).Limit(parseInt(pageSize)).Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败: " + err.Error(),
		})
		return
	}

	// 构建树形结构
	menuTree := buildMenuTree(menus, 0)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":  menuTree,
			"total": total,
		},
	})
}

// GetMenus godoc
// @Summary 获取菜单列表（兼容旧API）
// @Description 获取所有菜单列表
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/menus [get]
func GetMenus(c *gin.Context) {
	var menus []models.Menu
	if err := database.DB.Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败: " + err.Error(),
		})
		return
	}

	menuTree := buildMenuTree(menus, 0)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": menuTree,
	})
}

// buildMenuTree 构建菜单树
func buildMenuTree(menus []models.Menu, parentID int) []MenuNode {
	var tree []MenuNode

	for _, menu := range menus {
		var pid *int
		if menu.ParentID != nil && *menu.ParentID == parentID {
			pid = menu.ParentID
			node := MenuNode{
				ID:         menu.ID,
				ParentID:   pid,
				Name:       menu.Name,
				Type:       menu.Type,
				Path:       menu.Path,
				Component:  menu.Component,
				Permission: menu.Permission,
				Icon:       menu.Icon,
				Sort:       menu.Sort,
				Status:     menu.Status,
			}

			// 递归查找子菜单
			children := buildMenuTree(menus, menu.ID)
			if len(children) > 0 {
				node.Children = children
			}

			tree = append(tree, node)
		}
	}

	return tree
}

// CreateMenu godoc
// @Summary 创建菜单
// @Description 创建新菜单
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Menu true "菜单信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:创建成功"
// @Router /api/admin/menus [post]
func CreateMenu(c *gin.Context) {
	var req models.Menu

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 验证必填字段
	if req.Name == "" {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "菜单名称不能为空",
		})
		return
	}

	if req.Type == 0 {
		req.Type = 2 // 默认为菜单
	}

	// 如果是顶级菜单,parent_id为null
	if req.ParentID != nil && *req.ParentID == 0 {
		req.ParentID = nil
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

// UpdateMenu godoc
// @Summary 更新菜单
// @Description 更新菜单信息
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "菜单ID"
// @Param request body models.Menu true "菜单信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:更新成功"
// @Router /api/admin/menus/{id} [put]
func UpdateMenu(c *gin.Context) {
	id := c.Param("id")

	var req models.Menu
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 检查菜单是否存在
	var menu models.Menu
	if err := database.DB.First(&menu, id).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "菜单不存在",
		})
		return
	}

	// 验证parent_id不能是自己
	if req.ParentID != nil && *req.ParentID == menu.ID {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "上级菜单不能是自己",
		})
		return
	}

	// 更新字段
	updates := map[string]interface{}{
		"name":       req.Name,
		"type":       req.Type,
		"path":       req.Path,
		"component":  req.Component,
		"permission": req.Permission,
		"icon":       req.Icon,
		"sort":       req.Sort,
		"status":     req.Status,
	}

	if req.ParentID != nil && *req.ParentID == 0 {
		updates["parent_id"] = nil
	} else {
		updates["parent_id"] = req.ParentID
	}

	if err := database.DB.Model(&menu).Updates(updates).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "更新失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功",
		"data": menu,
	})
}

// DeleteMenu godoc
// @Summary 删除菜单
// @Description 删除指定菜单
// @Tags 管理员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "菜单ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:删除成功"
// @Router /api/admin/menus/{id} [delete]
func DeleteMenu(c *gin.Context) {
	id := c.Param("id")

	// 检查菜单是否存在
	var menu models.Menu
	if err := database.DB.First(&menu, id).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "菜单不存在",
		})
		return
	}

	// 检查是否有子菜单
	var count int64
	database.DB.Model(&models.Menu{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  fmt.Sprintf("该菜单下还有 %d 个子菜单，无法删除", count),
		})
		return
	}

	if err := database.DB.Delete(&menu).Error; err != nil {
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

package handlers

import (
	"akrick.com/mychat/database"
	"akrick.com/mychat/models"

	"github.com/gin-gonic/gin"
)

// GetRoleList 获取角色列表
func GetRoleList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	name := c.Query("name")

	query := database.DB.Model(&models.Role{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	var total int64
	query.Count(&total)

	var roles []models.Role
	offset := 0
	if page != "1" {
		offset = (parseInt(page) - 1) * parseInt(pageSize)
	}

	if err := query.Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&roles).Error; err != nil {
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
			"list":  roles,
			"total": total,
		},
	})
}

// CreateRole 创建角色
func CreateRole(c *gin.Context) {
	var req models.Role

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 检查角色名是否存在
	var count int64
	database.DB.Model(&models.Role{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "角色名已存在",
		})
		return
	}

	if err := database.DB.Create(&req).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "创建失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": req,
	})
}

// UpdateRole 更新角色
func UpdateRole(c *gin.Context) {
	roleID := c.Param("id")

	var req models.Role
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	if err := database.DB.Model(&models.Role{}).Where("id = ?", roleID).Updates(&req).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "更新失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// GetRoleUsers 获取角色用户列表
// @Summary 获取角色用户列表
// @Description 获取指定角色的所有用户
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功"
// @Router /api/admin/roles/:id/users [get]
func GetRoleUsers(c *gin.Context) {
	roleID := c.Param("id")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	var users []models.User
	var total int64

	query := database.DB.Table("users").
		Joins("INNER JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.role_id = ?", roleID)

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := 0
	if page != "1" {
		offset = (parseInt(page) - 1) * parseInt(pageSize)
	}

	if err := query.Offset(offset).Limit(parseInt(pageSize)).
		Find(&users).Error; err != nil {
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
			"users": users,
			"total": total,
		},
	})
}

// GetRolePermissions 获取角色权限
func GetRolePermissions(c *gin.Context) {
	roleID := c.Param("id")

	var permissions []models.Permission
	database.DB.Where("id IN (SELECT permission_id FROM role_permissions WHERE role_id = ?)", roleID).Find(&permissions)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": permissions,
	})
}

// AssignPermissions 分配权限
func AssignPermissions(c *gin.Context) {
	roleID := c.Param("id")

	var req struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 删除旧权限
	database.DB.Exec("DELETE FROM role_permissions WHERE role_id = ?", roleID)

	// 添加新权限
	for _, permissionID := range req.PermissionIDs {
		database.DB.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", roleID, permissionID)
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "分配成功",
	})
}

// GetPermissionTree 获取权限树
func GetPermissionTree(c *gin.Context) {
	var permissions []models.Permission
	database.DB.Order("sort ASC").Find(&permissions)

	tree := buildPermissionTree(permissions, 0)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": tree,
	})
}

// GetPermissionList 获取权限列表
func GetPermissionList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	name := c.Query("name")

	query := database.DB.Model(&models.Permission{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	var total int64
	query.Count(&total)

	var permissions []models.Permission
	offset := 0
	if page != "1" {
		offset = (parseInt(page) - 1) * parseInt(pageSize)
	}

	if err := query.Offset(offset).Limit(parseInt(pageSize)).Order("sort ASC").Find(&permissions).Error; err != nil {
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
			"list":  permissions,
			"total": total,
		},
	})
}

// CreatePermission 创建权限
func CreatePermission(c *gin.Context) {
	var req models.Permission

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	if err := database.DB.Create(&req).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "创建失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": req,
	})
}

// UpdatePermission 更新权限
func UpdatePermission(c *gin.Context) {
	permissionID := c.Param("id")

	var req models.Permission
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	if err := database.DB.Model(&models.Permission{}).Where("id = ?", permissionID).Updates(&req).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "更新失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// DeletePermission 删除权限
func DeletePermission(c *gin.Context) {
	permissionID := c.Param("id")

	if err := database.DB.Delete(&models.Permission{}, permissionID).Error; err != nil {
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

// buildPermissionTree 构建权限树
func buildPermissionTree(permissions []models.Permission, parentID uint) []models.Permission {
	var tree []models.Permission
	for _, p := range permissions {
		if p.ParentID == parentID {
			children := buildPermissionTree(permissions, p.ID)
			if len(children) > 0 {
				p.Children = children
			}
			tree = append(tree, p)
		}
	}
	return tree
}

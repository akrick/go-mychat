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

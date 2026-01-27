package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
	"akrick.com/mychat/database"
	"akrick.com/mychat/models"
	"github.com/gin-gonic/gin"
)

// UploadFile godoc
// @Summary 上传文件
// @Description 上传文件
// @Tags 文件
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "文件"
// @Param file_type formData string false "文件类型:avatar/document/image/audio/video" default(other)
// @Param relation_id formData int false "关联ID"
// @Param relation_type formData string false "关联类型"
// @Success 200 {object} map[string]interface{} "code:200,msg:上传成功,data:{file}"
// @Router /api/upload [post]
func UploadFile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	fileType := c.DefaultPostForm("file_type", models.FileTypeOther)
	relationID := c.PostForm("relation_id")
	relationType := c.PostForm("relation_type")

	// 获取文件
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "文件不存在: " + err.Error(),
		})
		return
	}

	// 检查文件大小（限制10MB）
	if fileHeader.Size > 10*1024*1024 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "文件大小不能超过10MB",
		})
		return
	}

	// 获取文件扩展名
	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "文件类型不支持",
		})
		return
	}

	// 创建上传目录
	uploadDir := fmt.Sprintf("./uploads/%s/%d", fileType, userID.(uint))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "创建目录失败: " + err.Error(),
		})
		return
	}

	// 生成新文件名
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, newFileName)

	// 打开文件
	src, err := fileHeader.Open()
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "打开文件失败: " + err.Error(),
		})
		return
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "创建文件失败: " + err.Error(),
		})
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "保存文件失败: " + err.Error(),
		})
		return
	}

	// 获取MIME类型
	mimeType := detectMimeType(fileHeader)

	// 构建文件URL
	fileURL := fmt.Sprintf("/uploads/%s/%d/%s", fileType, userID.(uint), newFileName)

	// 解析关联ID
	var relationIDUint uint
	if relationID != "" {
		fmt.Sscanf(relationID, "%d", &relationIDUint)
	}

	// 保存文件记录到数据库
	file := models.File{
		FileName:     newFileName,
		OriginalName: fileHeader.Filename,
		FileType:     fileType,
		FilePath:     filePath,
		FileURL:      fileURL,
		FileSize:     fileHeader.Size,
		MimeType:     mimeType,
		UploaderID:   userID.(uint),
		RelationID:   relationIDUint,
		RelationType: relationType,
		Status:       1,
	}

	if err := database.DB.Create(&file).Error; err != nil {
		// 删除已上传的文件
		os.Remove(filePath)
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "保存文件记录失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "上传成功",
		"data": gin.H{
			"id":            file.ID,
			"file_name":     file.FileName,
			"original_name": file.OriginalName,
			"file_url":      file.FileURL,
			"file_size":     file.FileSize,
			"mime_type":     file.MimeType,
		},
	})
}

// detectMimeType 检测MIME类型
func detectMimeType(fileHeader *multipart.FileHeader) string {
	// 简单实现，根据文件扩展名判断
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".mp3":  "audio/mpeg",
		".mp4":  "video/mp4",
		".wav":  "audio/wav",
		".avi":  "video/x-msvideo",
		".mov":  "video/quicktime",
		".txt":  "text/plain",
		".zip":  "application/zip",
		".rar":  "application/x-rar-compressed",
	}

	if mime, ok := mimeTypes[ext]; ok {
		return mime
	}

	return "application/octet-stream"
}

// GetFile godoc
// @Summary 获取文件信息
// @Description 获取文件信息
// @Tags 文件
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文件ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{file}"
// @Router /api/file/{id} [get]
func GetFile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	fileID := c.Param("id")

	var file models.File
	if err := database.DB.Preload("Uploader").First(&file, fileID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "文件不存在",
		})
		return
	}

	// 检查权限（上传者可以访问）
	if file.UploaderID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权访问此文件",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": file,
	})
}

// DeleteFile godoc
// @Summary 删除文件
// @Description 删除文件
// @Tags 文件
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文件ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:删除成功"
// @Router /api/file/{id} [delete]
func DeleteFile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	fileID := c.Param("id")

	var file models.File
	if err := database.DB.First(&file, fileID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "文件不存在",
		})
		return
	}

	// 检查权限
	if file.UploaderID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权删除此文件",
		})
		return
	}

	// 删除物理文件
	if err := os.Remove(file.FilePath); err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "删除文件失败: " + err.Error(),
		})
		return
	}

	// 更新数据库状态
	if err := database.DB.Model(&file).Update("status", 0).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "更新记录失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

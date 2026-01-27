package middleware

import (
	"github.com/gin-gonic/gin"
)

// AdminOnly 仅管理员可访问
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里可以根据实际需求判断用户是否为管理员
		// 例如：检查用户角色、特殊标识等
		// 暂时注释掉，实际使用时需要实现
		/*
		userID, _ := c.Get("user_id")
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(401, gin.H{
				"code": 401,
				"msg":  "未授权",
			})
			c.Abort()
			return
		}

		if user.Role != "admin" {
			c.JSON(403, gin.H{
				"code": 403,
				"msg":  "仅管理员可访问",
			})
			c.Abort()
			return
		}
		*/
		c.Next()
	}
}

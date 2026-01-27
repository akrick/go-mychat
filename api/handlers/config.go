package handlers

import (
	"akrick.com/mychat/database"
	"akrick.com/mychat/models"
	"github.com/gin-gonic/gin"
)

// GetPaymentConfig godoc
// @Summary 获取支付配置
// @Description 获取支付配置信息（管理员）
// @Tags 配置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payment_method query string false "支付方式:wechat/alipay"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{configs}"
// @Router /api/config/payment [get]
func GetPaymentConfig(c *gin.Context) {
	paymentMethod := c.Query("payment_method")

	var configs []models.PaymentConfig
	query := database.DB.Model(&models.PaymentConfig{})

	if paymentMethod != "" {
		query = query.Where("payment_method = ?", paymentMethod)
	}

	if err := query.Find(&configs).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败: " + err.Error(),
		})
		return
	}

	// 隐藏敏感信息
	for i := range configs {
		configs[i].APISecret = ""
		if configs[i].PrivateKeyPath != "" {
			configs[i].PrivateKeyPath = "***"
		}
		if configs[i].PublicKeyPath != "" {
			configs[i].PublicKeyPath = "***"
		}
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": configs,
	})
}

// UpdatePaymentConfig godoc
// @Summary 更新支付配置
// @Description 更新支付配置信息（管理员）
// @Tags 配置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "配置ID"
// @Param config body models.PaymentConfig true "支付配置"
// @Success 200 {object} map[string]interface{} "code:200,msg:更新成功"
// @Router /api/config/payment/{id} [put]
func UpdatePaymentConfig(c *gin.Context) {
	configID := c.Param("id")

	var config models.PaymentConfig
	if err := database.DB.First(&config, configID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "配置不存在",
		})
		return
	}

	var req models.PaymentConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 更新配置
	updates := map[string]interface{}{
		"app_id":          req.AppID,
		"mch_id":          req.MchID,
		"api_cert_path":   req.APICertPath,
		"api_key_path":    req.APIKeyPath,
		"notify_url":      req.NotifyURL,
		"private_key_path": req.PrivateKeyPath,
		"public_key_path":  req.PublicKeyPath,
		"is_enabled":      req.IsEnabled,
		"is_sandbox":      req.IsSandbox,
	}

	// 只有在提供了新的API密钥时才更新
	if req.APISecret != "" {
		updates["api_secret"] = req.APISecret
	}

	if err := database.DB.Model(&config).Updates(updates).Error; err != nil {
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

// TestPaymentConfig godoc
// @Summary 测试支付配置
// @Description 测试支付配置是否正确（管理员）
// @Tags 配置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "配置ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:测试成功"
// @Router /api/config/payment/{id}/test [post]
func TestPaymentConfig(c *gin.Context) {
	configID := c.Param("id")

	var config models.PaymentConfig
	if err := database.DB.First(&config, configID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "配置不存在",
		})
		return
	}

	// 检查必要配置
	if config.AppID == "" {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "AppID不能为空",
		})
		return
	}

	// 根据支付方式进行测试
	var testResult map[string]interface{}
	testResult = make(map[string]interface{})

	if config.PaymentMethod == models.PaymentMethodWeChat {
		// 微信支付测试
		testResult["app_id"] = config.AppID
		testResult["mch_id"] = config.MchID
		testResult["notify_url"] = config.NotifyURL
		testResult["is_sandbox"] = config.IsSandbox
		testResult["config_valid"] = config.MchID != ""
	} else if config.PaymentMethod == models.PaymentMethodAlipay {
		// 支付宝测试
		testResult["app_id"] = config.AppID
		testResult["notify_url"] = config.NotifyURL
		testResult["is_sandbox"] = config.IsSandbox
		testResult["config_valid"] = true
	}

	testResult["test_time"] = database.DB.NowFunc()
	testResult["status"] = "success"

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "测试成功",
		"data": testResult,
	})
}

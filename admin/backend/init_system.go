package main

import (
	"encoding/json"
	"akrick.com/mychat/admin/backend/database"
	"akrick.com/mychat/admin/backend/models"
	"log"
)

// InitSystemConfigs 初始化系统配置
func InitSystemConfigs() {
	configs := []models.SystemConfig{
		// 基本配置
		{
			Key:      "system_name",
			Value:     `"MyChat 心理咨询平台"`,
			Category:  "basic",
			Label:     "系统名称",
			Type:      "string",
			IsSystem:  true,
			Sort:      1,
		},
		{
			Key:      "system_description",
			Value:     `"专业的在线心理咨询服务平台"`,
			Category:  "basic",
			Label:     "系统简介",
			Type:      "string",
			IsSystem:  true,
			Sort:      2,
		},
		{
			Key:      "site_domain",
			Value:     `"https://mychat.com"`,
			Category:  "basic",
			Label:     "网站域名",
			Type:      "string",
			IsSystem:  false,
			Sort:      3,
		},
		{
			Key:      "contact_phone",
			Value:     `"400-888-8888"`,
			Category:  "basic",
			Label:     "联系电话",
			Type:      "string",
			IsSystem:  false,
			Sort:      4,
		},
		{
			Key:      "contact_email",
			Value:     `"support@mychat.com"`,
			Category:  "basic",
			Label:     "联系邮箱",
			Type:      "string",
			IsSystem:  false,
			Sort:      5,
		},
		{
			Key:      "icp_code",
			Value:     `""`,
			Category:  "basic",
			Label:     "ICP备案号",
			Type:      "string",
			IsSystem:  false,
			Sort:      6,
		},

		// 用户配置
		{
			Key:      "allow_register",
			Value:     `true`,
			Category:  "user",
			Label:     "允许注册",
			Type:      "boolean",
			IsSystem:  true,
			Sort:      1,
		},
		{
			Key:      "default_avatar",
			Value:     `"/assets/default-avatar.png"`,
			Category:  "user",
			Label:     "默认头像",
			Type:      "string",
			IsSystem:  false,
			Sort:      2,
		},
		{
			Key:      "password_min_length",
			Value:     `6`,
			Category:  "user",
			Label:     "密码最小长度",
			Type:      "number",
			IsSystem:  true,
			Sort:      3,
		},
		{
			Key:      "login_max_attempts",
			Value:     `5`,
			Category:  "user",
			Label:     "账户锁定阈值",
			Type:      "number",
			IsSystem:  true,
			Sort:      4,
		},
		{
			Key:      "login_lock_time",
			Value:     `30`,
			Category:  "user",
			Label:     "锁定时长(分钟)",
			Type:      "number",
			IsSystem:  true,
			Sort:      5,
		},

		// 聊天配置
		{
			Key:      "free_chat_duration",
			Value:     `5`,
			Category:  "chat",
			Label:     "免费聊天时长(分钟)",
			Type:      "number",
			IsSystem:  true,
			Sort:      1,
		},
		{
			Key:      "billing_cycle",
			Value:     `60`,
			Category:  "chat",
			Label:     "计费周期(秒)",
			Type:      "number",
			IsSystem:  true,
			Sort:      2,
		},
		{
			Key:      "session_timeout",
			Value:     `30`,
			Category:  "chat",
			Label:     "会话超时(分钟)",
			Type:      "number",
			IsSystem:  true,
			Sort:      3,
		},
		{
			Key:      "message_retention_days",
			Value:     `90`,
			Category:  "chat",
			Label:     "消息保留天数",
			Type:      "number",
			IsSystem:  true,
			Sort:      4,
		},
		{
			Key:      "allow_image_message",
			Value:     `true`,
			Category:  "chat",
			Label:     "允许发图",
			Type:      "boolean",
			IsSystem:  true,
			Sort:      5,
		},
		{
			Key:      "allow_voice_message",
			Value:     `true`,
			Category:  "chat",
			Label:     "允许发语音",
			Type:      "boolean",
			IsSystem:  true,
			Sort:      6,
		},

		// 支付配置
		{
			Key:      "alipay_app_id",
			Value:     `""`,
			Category:  "payment",
			Label:     "支付宝AppID",
			Type:      "string",
			IsSystem:  false,
			Sort:      1,
		},
		{
			Key:      "alipay_public_key",
			Value:     `""`,
			Category:  "payment",
			Label:     "支付宝公钥",
			Type:      "textarea",
			IsSystem:  false,
			Sort:      2,
		},
		{
			Key:      "wechat_app_id",
			Value:     `""`,
			Category:  "payment",
			Label:     "微信AppID",
			Type:      "string",
			IsSystem:  false,
			Sort:      3,
		},
		{
			Key:      "wechat_mch_id",
			Value:     `""`,
			Category:  "payment",
			Label:     "微信商户号",
			Type:      "string",
			IsSystem:  false,
			Sort:      4,
		},
		{
			Key:      "wechat_api_key",
			Value:     `""`,
			Category:  "payment",
			Label:     "微信API密钥",
			Type:      "password",
			IsSystem:  false,
			Sort:      5,
		},
		{
			Key:      "min_withdraw_amount",
			Value:     `100`,
			Category:  "payment",
			Label:     "提现最小金额(元)",
			Type:      "number",
			IsSystem:  true,
			Sort:      6,
		},

		// 通知配置
		{
			Key:      "email_enabled",
			Value:     `false`,
			Category:  "notification",
			Label:     "启用邮件通知",
			Type:      "boolean",
			IsSystem:  true,
			Sort:      1,
		},
		{
			Key:      "smtp_host",
			Value:     `"smtp.gmail.com"`,
			Category:  "notification",
			Label:     "SMTP服务器",
			Type:      "string",
			IsSystem:  false,
			Sort:      2,
		},
		{
			Key:      "smtp_port",
			Value:     `465`,
			Category:  "notification",
			Label:     "SMTP端口",
			Type:      "number",
			IsSystem:  false,
			Sort:      3,
		},
		{
			Key:      "smtp_username",
			Value:     `""`,
			Category:  "notification",
			Label:     "SMTP用户名",
			Type:      "string",
			IsSystem:  false,
			Sort:      4,
		},
		{
			Key:      "smtp_password",
			Value:     `""`,
			Category:  "notification",
			Label:     "SMTP密码",
			Type:      "password",
			IsSystem:  false,
			Sort:      5,
		},
		{
			Key:      "sms_enabled",
			Value:     `false`,
			Category:  "notification",
			Label:     "启用短信通知",
			Type:      "boolean",
			IsSystem:  true,
			Sort:      6,
		},
		{
			Key:      "sms_provider",
			Value:     `"aliyun"`,
			Category:  "notification",
			Label:     "短信服务商",
			Type:      "select",
			IsSystem:  false,
			Sort:      7,
		},

		// 存储配置
		{
			Key:      "storage_type",
			Value:     `"local"`,
			Category:  "storage",
			Label:     "存储方式",
			Type:      "radio",
			IsSystem:  true,
			Sort:      1,
		},
		{
			Key:      "upload_path",
			Value:     `"/uploads"`,
			Category:  "storage",
			Label:     "上传路径",
			Type:      "string",
			IsSystem:  false,
			Sort:      2,
		},
		{
			Key:      "allowed_file_types",
			Value:     `"jpg,jpeg,png,gif,pdf,doc,docx,xls,xlsx"`,
			Category:  "storage",
			Label:     "允许文件类型",
			Type:      "string",
			IsSystem:  true,
			Sort:      3,
		},
		{
			Key:      "max_file_size",
			Value:     `10`,
			Category:  "storage",
			Label:     "最大文件大小(MB)",
			Type:      "number",
			IsSystem:  true,
			Sort:      4,
		},
		{
			Key:      "image_compress",
			Value:     `true`,
			Category:  "storage",
			Label:     "图片压缩",
			Type:      "boolean",
			IsSystem:  true,
			Sort:      5,
		},
		{
			Key:      "compress_quality",
			Value:     `80`,
			Category:  "storage",
			Label:     "压缩质量",
			Type:      "slider",
			IsSystem:  true,
			Sort:      6,
		},
	}

	for _, config := range configs {
		var count int64
		database.DB.Model(&models.SystemConfig{}).Where("key = ?", config.Key).Count(&count)

		if count == 0 {
			if err := database.DB.Create(&config).Error; err != nil {
				log.Printf("创建配置失败: %s, 错误: %v", config.Key, err)
			} else {
				log.Printf("创建配置成功: %s", config.Key)
			}
		}
	}

	// 验证配置是否为有效的JSON
	for i := range configs {
		var result interface{}
		if err := json.Unmarshal([]byte(configs[i].Value), &result); err != nil {
			log.Printf("配置值不是有效JSON: %s, 值: %s", configs[i].Key, configs[i].Value)
		}
	}
}

package main

// GetSessionStats 获取会话统计
func GetSessionStats() map[string]interface{} {
	return sessionManager.GetSessionStats()
}

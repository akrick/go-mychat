package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"akrick.com/mychat/cache"
	"akrick.com/mychat/database"
	"akrick.com/mychat/models"
)

// SessionManager 会话管理器
type SessionManager struct {
	mu               sync.RWMutex
	activeSessions   map[uint]*ActiveSession // sessionID -> ActiveSession
	sessionTimers    map[uint]*time.Timer   // sessionID -> Timer
}

type ActiveSession struct {
	SessionID   uint
	UserID      uint
	CounselorID uint
	StartTime  time.Time
	PricePerMin float64
	LastPing    time.Time
}

var sessionManager *SessionManager

// InitSessionManager 初始化会话管理器
func InitSessionManager() {
	sessionManager = &SessionManager{
		activeSessions: make(map[uint]*ActiveSession),
		sessionTimers:  make(map[uint]*time.Timer),
	}
	
	// 启动定时检查超时会话
	go sessionManager.checkTimeoutSessions()
	
	log.Println("会话管理器初始化成功")
}

// StartSession 开始会话
func (sm *SessionManager) StartSession(sessionID, userID, counselorID uint, price float64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	now := time.Now()
	activeSession := &ActiveSession{
		SessionID:   sessionID,
		UserID:      userID,
		CounselorID: counselorID,
		StartTime:  now,
		PricePerMin: price,
		LastPing:   now,
	}
	sm.activeSessions[sessionID] = activeSession

	// 设置30分钟无操作超时
	if timer, ok := sm.sessionTimers[sessionID]; ok {
		timer.Stop()
	}
	sm.sessionTimers[sessionID] = time.AfterFunc(30*time.Minute, func() {
		sm.handleSessionTimeout(sessionID)
	})

	log.Printf("会话开始: sessionID=%d, userID=%d, counselorID=%d", sessionID, userID, counselorID)
}

// EndSession 结束会话
func (sm *SessionManager) EndSession(sessionID uint) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if timer, ok := sm.sessionTimers[sessionID]; ok {
		timer.Stop()
		delete(sm.sessionTimers, sessionID)
	}
	delete(sm.activeSessions, sessionID)

	log.Printf("会话结束: sessionID=%d", sessionID)
}

// UpdateLastPing 更新最后ping时间
func (sm *SessionManager) UpdateLastPing(sessionID uint) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if session, ok := sm.activeSessions[sessionID]; ok {
		session.LastPing = time.Now()
		
		// 重置超时计时器
		if timer, ok := sm.sessionTimers[sessionID]; ok {
			timer.Stop()
		}
		sm.sessionTimers[sessionID] = time.AfterFunc(30*time.Minute, func() {
			sm.handleSessionTimeout(sessionID)
		})
	}
}

// GetActiveSession 获取活跃会话
func (sm *SessionManager) GetActiveSession(sessionID uint) (*ActiveSession, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, ok := sm.activeSessions[sessionID]
	return session, ok
}

// GetActiveSessions 获取所有活跃会话
func (sm *SessionManager) GetActiveSessions() []*ActiveSession {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sessions := make([]*ActiveSession, 0, len(sm.activeSessions))
	for _, session := range sm.activeSessions {
		sessions = append(sessions, session)
	}
	return sessions
}

// handleSessionTimeout 处理会话超时
func (sm *SessionManager) handleSessionTimeout(sessionID uint) {
	sm.mu.Lock()
	session, ok := sm.activeSessions[sessionID]
	if !ok {
		sm.mu.Unlock()
		return
	}
	
	userID := session.UserID
	counselorID := session.CounselorID
	sm.mu.Unlock()

	log.Printf("会话超时: sessionID=%d, 开始自动结束", sessionID)

	// 自动结束会话
	sm.EndSession(sessionID)

	// 查询会话并计费
	var sessionModel models.ChatSession
	if err := database.DB.First(&sessionModel, sessionID).Error; err == nil {
		if sessionModel.Status == 1 {
			// 使用client的endSession方法
			if client, ok := globalHub.clients[counselorID]; ok {
				client.endSession(sessionID, sessionModel)
			}
		}
	}

	// 发送通知
	var notification models.Notification
	notification.UserID = userID
	notification.Type = models.NotificationTypeSystem
	notification.Level = models.NotificationLevelWarning
	notification.Title = "会话超时结束"
	notification.Content = "由于长时间无操作，您的咨询会话已自动结束"
	database.DB.Create(&notification)

	ctx := context.Background()
	cache.ClearOnlineStatus(ctx, userID)
}

// checkTimeoutSessions 定时检查超时会话
func (sm *SessionManager) checkTimeoutSessions() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sm.mu.RLock()
		now := time.Now()
		timeoutSessions := make([]uint, 0)
		
		for sessionID, session := range sm.activeSessions {
			// 检查是否超过30分钟无操作
			if now.Sub(session.LastPing) > 30*time.Minute {
				timeoutSessions = append(timeoutSessions, sessionID)
			}
		}
		sm.mu.RUnlock()

		// 处理超时会话
		for _, sessionID := range timeoutSessions {
			sm.handleSessionTimeout(sessionID)
		}
	}
}

// GetSessionStats 获取会话统计
func (sm *SessionManager) GetSessionStats() map[string]interface{} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	stats := map[string]interface{}{
		"active_sessions": len(sm.activeSessions),
		"total_duration":  0,
		"total_amount":    0.0,
	}

	now := time.Now()
	for _, session := range sm.activeSessions {
		duration := now.Sub(session.StartTime)
		durationMinutes := int(duration.Minutes())
		if durationMinutes < 1 {
			durationMinutes = 1
		}
		
		stats["total_duration"] = stats["total_duration"].(int) + durationMinutes
		stats["total_amount"] = stats["total_amount"].(float64) + float64(durationMinutes)*session.PricePerMin
	}

	return stats
}

// BroadcastSessionStats 广播会话统计
func BroadcastSessionStats() {
	stats := sessionManager.GetSessionStats()
	statsData, _ := json.Marshal(map[string]interface{}{
		"type": "session_stats",
		"data": stats,
	})

	globalHub.broadcast <- statsData
}

// GetCounselorEarnings 获取咨询师收益统计
func GetCounselorEarnings(counselorID uint) map[string]interface{} {
	var account models.CounselorAccount
	database.DB.Where("counselor_id = ?", counselorID).First(&account)

	var totalSessions int64
	var totalDuration int
	database.DB.Model(&models.ChatSession{}).
		Where("counselor_id = ? AND status = 2", counselorID).
		Count(&totalSessions)
	database.DB.Model(&models.ChatSession{}).
		Select("COALESCE(SUM(duration), 0)").
		Where("counselor_id = ? AND status = 2", counselorID).
		Scan(&totalDuration)

	return map[string]interface{}{
		"total_income":   account.TotalIncome,
		"balance":        account.Balance,
		"withdrawn":      account.Withdrawn,
		"frozen_amount":  account.FrozenAmount,
		"total_sessions": totalSessions,
		"total_duration": totalDuration / 60, // 转换为分钟
	}
}
